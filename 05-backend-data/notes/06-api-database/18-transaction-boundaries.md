# 事务边界、Go `database/sql` 连接池与可靠副作用

事务边界应覆盖必须一起成立的数据库不变量，并尽可能短。HTTP 请求、数据库事务、消息投递和外部 API 分属不同可靠性边界；把它们写在同一个函数中不会自动获得原子性。

## 1. 从业务不变量划边界

“创建订单”可能要求：订单主表、明细、库存预留和 outbox 事件一起提交。邮件发送、支付供应商调用和对象存储上传无法加入普通 PostgreSQL 本地事务，应在提交后通过可靠状态机执行。

```mermaid
flowchart LR
    A["校验请求/授权"] --> B["BEGIN"]
    B --> C["读取并锁定当前状态"]
    C --> D["写业务数据 + outbox"]
    D --> E["COMMIT"]
    E --> F["异步投递外部副作用"]
```

不要在事务开始前读取决定性状态、事务中只写结果，否则读写之间仍有竞态。静态格式校验可在事务外；依赖数据库当前状态的校验要与写入处于同一正确性边界。

## 2. `sql.DB` 的真实含义

Go `*sql.DB` 是数据库连接池句柄，可供多个 goroutine 并发使用；它不是一条连接，也不代表一个事务。第一次查询时才可能建立连接。

池配置：

- `SetMaxOpenConns(n)`：打开连接上限；过大压垮数据库，过小产生排队。
- `SetMaxIdleConns(n)`：保留空闲连接上限；过低增加握手，不能高于 max open 的实际效果。
- `SetConnMaxLifetime(d)`：连接最大复用寿命，用于负载均衡/服务端回收，不应所有实例同步抖动到期。
- `SetConnMaxIdleTime(d)`：连接最大闲置时间。
- `DB.Stats()`：观察 OpenConnections、InUse、Idle、WaitCount、WaitDuration、MaxIdleClosed、MaxLifetimeClosed 等。

池大小要从数据库总连接预算反推：数据库上限减去运维/迁移/后台保留，再除应用实例和任务类型。应用水平扩容会乘以每实例上限。

## 3. 连接池排队与请求 deadline

`QueryContext` 在池中等待连接时也受 context 控制。请求慢可能没有任何 SQL 执行，因为全部时间花在 pool wait。需要同时观测应用池等待、数据库活动连接和语句耗时。

`db.PingContext` 可验证连接，但健康检查不能每次新建无限连接，也不能把暂时数据库慢导致所有实例同时摘除。readiness 与 liveness 语义分开。

## 4. 事务 API 的正确使用

```go
func withTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn func(*sql.Tx) error) (err error) {
    tx, err := db.BeginTx(ctx, opts)
    if err != nil { return err }
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
        if err != nil { _ = tx.Rollback() }
    }()
    if err = fn(tx); err != nil { return err }
    err = tx.Commit()
    return err
}
```

必须检查 `Commit` 错误；最后一条 Exec 成功不等于 commit 已持久确认。事务内所有操作用 `tx.Query/Exec`，不能混用 `db`，后者可能在另一连接、事务外执行。

`defer tx.Rollback()` 的简单模式可靠，Commit 后 Rollback 返回 `ErrTxDone`。封装要处理 panic，但不要吞掉原 panic。事务 closure 若支持重试，必须没有不可重复外部副作用。

## 5. 查询资源生命周期

`Rows` 必须关闭，并在遍历后检查 `rows.Err()`：

```go
func readRows(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
    rows, err := tx.QueryContext(ctx, query, args...)
    if err != nil { return err }
    defer rows.Close()
    for rows.Next() {
        var id string
        var amount int64
        if err := rows.Scan(&id, &amount); err != nil { return err }
    }
    return rows.Err()
}
```

不关闭 Rows 会长期占用连接；在一个事务中未消费完结果就发下一查询，对某些驱动/协议会失败。`QueryRowContext` 的错误在 `Scan` 返回。`sql.ErrNoRows` 是明确不存在，不应当 500。

参数必须使用占位符，列名/排序方向等标识符不能作为普通参数，应白名单映射。不要用 `fmt.Sprintf` 拼用户输入。

## 6. NULL、类型与时间

数据库 NULL 与 Go 零值不同。使用 `sql.NullString` 等 nullable 类型、指针或驱动类型明确区分。扫描 decimal 金额到 float64 会产生精度风险，货币可用最小单位整数或精确 decimal 库/数据库 numeric 并定义舍入。

时间统一明确时区；PostgreSQL `timestamptz` 表示瞬时时刻、按会话时区显示，`timestamp without time zone` 不带时区语义。API 输出 RFC 3339 UTC 通常更稳，日历日期使用 `date` 而非午夜时间戳。

## 7. 超时与数据库设置

应用 context 是第一层，数据库 session/transaction 设置是第二层：

```sql
SET LOCAL statement_timeout = '800ms';
SET LOCAL lock_timeout = '200ms';
SET LOCAL idle_in_transaction_session_timeout = '2s';
```

使用 `SET LOCAL` 限于当前事务，避免池连接状态泄漏。具体值按业务预算，不复制示例到所有查询。statement timeout 也包括锁等待时间；lock timeout 只控制取锁等待并且不应大于 statement timeout。

context 取消后驱动会尝试取消查询，但连接状态和提交结果仍需处理。网络在 COMMIT 时断开可能产生“提交结果未知”；高风险写入用业务幂等 ID/查询接口对账，而不是假定回滚。

## 8. Outbox 解决提交与消息双写

错误模式：先提交订单再发消息，进程崩溃会漏消息；先发消息再提交，消费者可能看到不存在的订单。

Transactional outbox 在同一数据库事务写业务行和待发送事件：

```sql
CREATE TABLE outbox_events (
    event_id uuid PRIMARY KEY,
    aggregate_type text NOT NULL,
    aggregate_id uuid NOT NULL,
    event_type text NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    published_at timestamptz
);
```

worker 读取未发布事件，投递后标记。崩溃可能导致重复投递，因此消费者按 event ID 幂等。把 `published_at` 更新和外部 broker ack 做不到本地原子，outbox 提供至少一次而非神奇 exactly-once。

多个 worker 用 `FOR UPDATE SKIP LOCKED` 短事务领取批次，但不要持数据库事务跨 broker 网络调用。可先标租约/attempt，提交，再发送；失败后租约过期重试。

## 9. 跨系统事务与补偿

两阶段提交可协调支持该协议的资源，但增加阻塞、运维和恢复复杂度，也不覆盖普通 SaaS API。产品流程常用 saga/状态机：每步本地提交，失败执行显式补偿。

补偿要考虑：

- 反向动作是否真实可行，例如已发送邮件无法撤回；
- 补偿本身可能重复/失败，需要幂等；
- 中间状态是否对用户可见；
- 超时后结果未知时如何对账；
- 哪些状态需要人工介入。

## 10. 事务重试封装边界

仅对已识别的暂时错误重试，例如 PostgreSQL `40001` 与 `40P01`。重试从新事务开始，重新读取状态。最大次数、总 deadline、退避和指标必需。

闭包中允许：数据库读取、确定性计算、以数据库唯一 ID 写 outbox。禁止直接：发送邮件、请求支付、发布消息、修改外部文件。随机业务 ID 在闭包外生成并由唯一约束稳定，或首次写入后复用。

## 11. 完整案例：创建报销并通知

### 输入

- 报销、明细和预算占用必须一起成功。
- 成功后发消息给审批系统；broker 可能不可用。
- 客户端超时会用同一 idempotency key 重试。
- API 总预算 2 s，锁等待最多 200 ms。

### 步骤

1. handler 在事务外完成 JSON 结构和认证；从 principal 得到 tenant/user。
2. `BeginTx` 后 `SET LOCAL` 超时，插入/锁定幂等记录。
3. 查询当前预算并用条件 UPDATE 占用；插入 expense 和 items。
4. 同事务插入 `expense.submitted` outbox 和最终幂等结果。
5. 检查 Commit；成功返回 201。broker 故障不回滚已提交业务。
6. outbox worker 投递；审批系统按 event ID 幂等消费。

### 输出

数据库中 expense、items、预算变化、outbox 和幂等结果要么全部存在，要么都不存在。通知可能延迟或重复，但最终可恢复且不会重复创建审批对象。

### 验证

- 在每条 SQL 后注入错误，事务回滚后不存在半成品。
- 在 COMMIT 响应丢失处注入网络故障，同 key 查询得到原 expense。
- broker 停机一小时，outbox 积压可观测；恢复后全部投递。
- 消费同一 event 三次只创建一次审批记录。
- 池 WaitDuration、InUse、事务时长和 outbox oldest age 有指标。

### 失败分支

若事务中调用 broker 且等待 30 秒，数据库锁和连接被长期占用；若 broker 成功后数据库回滚，还会出现幽灵事件。修正为事务内写 outbox，事务外投递。若函数在 tx 内意外调用 `db.ExecContext` 写 outbox，订单回滚时事件仍可能提交；代码审查和接口设计应让 repository 只接收 tx executor。

## 12. 常见错误

- 把 `sql.DB` 当单连接，在其上设置会话变量后假定下一请求沿用。
- 每请求创建/关闭一个 DB，失去池复用。
- 不设 MaxOpenConns，实例扩容耗尽数据库连接。
- 忽略 Commit、Rows.Err 或 Rows.Close。
- 事务内混用 db 与 tx。
- 在事务中做长 HTTP 调用或等待用户输入。
- 将 context 取消当作确定回滚。
- outbox 消费者假定不会重复。

## 13. 练习

用 Go 1.26 `database/sql` 实现转账与 outbox。编写并发、context 超时、commit 结果未知和 broker 重复投递测试。

完成标准：池上限从数据库预算推导；事务内只使用 tx；所有 Rows 关闭并检查错误；SQL 参数化；不变量和 outbox 原子提交；消费者幂等；重试 closure 无外部副作用；故障注入后能通过业务 ID 对账最终结果。

## 来源

- [Go 1.26 package database/sql](https://pkg.go.dev/database/sql)（访问日期：2026-07-17）
- [Go: Managing connections](https://go.dev/doc/database/manage-connections)（访问日期：2026-07-17）
- [Go: Executing transactions](https://go.dev/doc/database/execute-transactions)（访问日期：2026-07-17）
- [PostgreSQL 18: Client Connection Defaults](https://www.postgresql.org/docs/18/runtime-config-client.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html)（访问日期：2026-07-17）
