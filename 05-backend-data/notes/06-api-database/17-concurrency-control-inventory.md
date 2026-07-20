# 并发控制：乐观锁、悲观锁与库存扣减

并发控制要保护业务不变量，而不是追求“没有冲突”。库存场景的核心不变量是可售数量不能为负、同一业务意图不能重复扣减、订单取消只能归还一次，并且数据库记录与外部动作可恢复。

## 1. 从不变量选择机制

```text
available >= 0
reserved >= 0
on_hand = available + reserved + damaged + ...
unique(tenant_id, reservation_id)
```

前端显示“还剩 1 件”和提交前校验都不是不变量保护。两个请求能同时看到 1。最终判断必须在受控数据库的原子写、约束或锁内完成。

## 2. 原子条件更新

单 SKU 扣减最直接的方式：

```sql
UPDATE inventory
SET available = available - $2,
    version = version + 1,
    updated_at = now()
WHERE tenant_id = $1
  AND sku = $3
  AND available >= $2
RETURNING available, version;
```

影响一行表示成功，零行表示不存在或库存不足；为避免泄露可按 API 契约进一步查询分类。条件判断和写入由同一语句原子完成，避免“SELECT 后 UPDATE”的竞态。数据库仍应有 `CHECK (available >= 0)` 作为纵深防御。

这通常优于先显式 `FOR UPDATE` 再扣减：往返更少、锁持有更短。但多 SKU、复杂定价或需要读取多列决策时，显式事务更清晰。

## 3. 乐观并发控制

乐观锁不先阻止别人写，而是在提交更新时验证版本未变：

```sql
UPDATE products
SET title = $1, version = version + 1
WHERE tenant_id = $2 AND product_id = $3 AND version = $4
RETURNING version;
```

零行代表版本冲突或资源不存在。HTTP 可用 ETag/`If-Match` 表达，冲突返回 412。适合读取多、冲突少、用户编辑需要发现覆盖的场景。

版本必须在数据库条件写中比较。仅在 Go 中先比较两个整数再无条件 UPDATE 仍会丢失更新。`updated_at` 可作为版本但时间精度、相同时间和格式转换更复杂，递增整数或不可伪造 ETag 更稳定。

冲突后不能总是自动重试：用户编辑的两个不同标题需要合并或提示；可交换的计数增量可以重新读后重算。策略取决于操作语义。

## 4. 悲观锁

悲观锁在决策前锁住行：

```sql
BEGIN;
SELECT available, price_cents
FROM inventory
WHERE tenant_id = $1 AND sku = $2
FOR UPDATE;
-- 校验后 UPDATE，并写 reservation
COMMIT;
```

适合冲突频繁、事务很短、必须基于当前行做多步决策的场景。成本是等待、死锁和吞吐下降。多 SKU 必须按稳定顺序（例如 sku 升序）锁定，避免 A→B 与 B→A 死锁。

`NOWAIT` 适合宁可快速失败；`SKIP LOCKED` 适合 worker 领取任务，不适合购物库存，因为跳过被锁 SKU 会错误表示不存在/缺货。

## 5. 预留而非立即扣最终库存

下单到付款存在时间。常见模型创建 reservation：

```sql
CREATE TABLE inventory_reservations (
    tenant_id uuid NOT NULL,
    reservation_id uuid NOT NULL,
    order_id uuid NOT NULL,
    sku text NOT NULL,
    quantity integer NOT NULL CHECK (quantity > 0),
    state text NOT NULL CHECK (state IN ('active','confirmed','released','expired')),
    expires_at timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, reservation_id),
    UNIQUE (tenant_id, order_id, sku)
);
```

创建 active reservation 与减少 available 在同一事务。付款成功转 confirmed；取消或超时转 released/expired 并归还一次。状态迁移用条件 UPDATE：

```sql
UPDATE inventory_reservations
SET state = 'released'
WHERE tenant_id = $1 AND reservation_id = $2 AND state = 'active'
RETURNING sku, quantity;
```

只有返回行时才归还库存，重复取消不重复加。定时过期 worker 可能重复执行，仍依赖状态条件和唯一约束。

## 6. 多 SKU 订单

一单包含多个 SKU 时必须定义原子性：要么全部预留，要么全部失败。步骤：

1. 对输入按 `(tenant_id, sku)` 排序并合并重复 SKU 数量。
2. 开启短事务。
3. 按排序顺序锁定全部库存行，或逐条执行条件扣减。
4. 任一不足则回滚全部。
5. 插入每项 reservation 和订单状态。
6. 提交后再发送事件。

逐条条件 UPDATE 失败时整个事务回滚可保持原子性，但仍应固定顺序减少死锁。单条 SQL 也可用输入 CTE 做批量校验/更新，复杂度更高，必须验证影响行数和结果集合完整。

## 7. 热点 SKU 与吞吐边界

单行强一致库存天然串行化竞争，增加应用副本不会提高该行写吞吐。优化方向：缩短事务、避免事务内外部调用、按仓库/批次拆分库存、请求排队、预分配令牌，或在可接受业务语义下采用分段计数。

分片库存会增加防超卖和汇总复杂度，不能只把一个计数拆成十行再任意扣。需要分配算法和总量不变量，迁移/回收也要安全。

“缓存先扣、异步写数据库”只有在缓存成为正确性系统并具备持久、复制、恢复和对账设计时才成立。普通缓存丢失或双写失败会造成库存漂移。

## 8. 队列能解决与不能解决的事

按 SKU 串行队列可减少数据库锁竞争，但队列通常至少一次投递，消费者仍要幂等；分区重平衡、重试和绕过队列的管理接口仍可能并发。数据库约束是最后防线。

排队还改变产品语义：请求可能先 accepted、后失败。若客户端需要即时确认，必须返回同步预留结果或把订单建为 pending 并提供任务状态，不能在 UI 假装成功。

## 9. 跨服务与 Saga

库存、支付和订单分属数据库时，没有普通本地事务覆盖全部。可采用 saga：

1. 订单服务创建 pending 订单与 outbox。
2. 库存服务幂等创建 reservation。
3. 支付服务幂等创建 payment intent。
4. 全部成功后订单 confirmed。
5. 某步永久失败，执行补偿：释放 reservation、取消支付。

补偿不是数据库回滚：可能失败、延迟或只能做反向业务操作。每步需要状态机、幂等 ID、重试、超时和人工对账。不要在持有库存行锁时跨网络调用支付服务。

## 10. Go 1.26 事务示例要点

`*sql.DB` 是并发安全的连接池句柄，不是单连接。事务必须通过 `*sql.Tx` 的方法执行，不能在事务中误用 `db.ExecContext`：

```go
func reserve(ctx context.Context, db *sql.DB, tenant, order, sku string, qty int) error {
    tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
    if err != nil { return err }
    defer tx.Rollback()

    var remaining int
    err = tx.QueryRowContext(ctx, `
        UPDATE inventory SET available = available - $1
        WHERE tenant_id = $2 AND sku = $3 AND available >= $1
        RETURNING available`, qty, tenant, sku).Scan(&remaining)
    if errors.Is(err, sql.ErrNoRows) { return ErrInsufficientStock }
    if err != nil { return err }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO inventory_reservations
          (tenant_id, reservation_id, order_id, sku, quantity, state, expires_at)
        VALUES ($1, gen_random_uuid(), $2, $3, $4, 'active', now()+interval '15 minutes')`,
        tenant, order, sku, qty)
    if err != nil { return err }
    return tx.Commit()
}
```

延迟 Rollback 在 Commit 后返回 `sql.ErrTxDone`，可忽略；真正错误仍由前面的返回处理。context deadline 应覆盖事务，且数据库设置 statement/lock timeout 防御。

## 11. 完整案例：最后一件商品

### 输入

`available=1`，两个不同订单并发各买 1 件；同一订单因超时还会重试三次。预留 15 分钟过期。

### 步骤

1. 每个订单有唯一 `order_id`，reservation 唯一约束为 `(tenant, order, sku)`。
2. 两事务执行 `UPDATE ... available >= 1 RETURNING`；行锁使它们串行，只有一个返回行。
3. 成功事务插入 active reservation 后提交；失败事务返回库存不足。
4. 同订单重试命中 reservation 唯一键，查询并返回已有预留，不再扣减。
5. 到期 worker 用 `state='active'` 条件迁移 expired，仅成功者归还 1。

### 输出

最终一个 active reservation、一个订单缺货，`available=0`。若预留过期，恰好归还一次变为 1。

### 验证

- 100 并发请求下 `available` 永不为负。
- 相同订单重试只存在一条 reservation。
- 付款确认与过期 worker 并发时只允许一个状态迁移，不能既确认又归还。
- 故障注入在 UPDATE 后、COMMIT 前终止连接，事务全部回滚。
- 锁等待和死锁重试率可观测。

### 失败分支

若逻辑是 `SELECT available` 后在事务外 `UPDATE available=0`，两个请求都可能成功。若取消接口每次都 `available=available+1`，重复调用会凭空增加库存。分别修正为原子条件更新，以及 reservation 状态条件迁移后才归还。

## 12. 方案选择

| 场景 | 优先机制 | 原因 |
|---|---|---|
| 单行计数扣减 | 原子条件 UPDATE + CHECK | 最短正确边界 |
| 低冲突用户编辑 | version/ETag 乐观锁 | 不长时间持锁，可发现覆盖 |
| 高频冲突、多步行内决策 | `FOR UPDATE` 悲观锁 | 基于当前状态串行决策 |
| 跨多行不变量 | 约束/共同锁/Serializable | 单行版本不足 |
| 长业务流程 | 状态机 + reservation + saga | 不持有长数据库事务 |

## 13. 常见错误与练习

错误包括：只做前端校验；库存缓存与数据库双写无恢复；事务中调用支付；多 SKU 锁顺序不一致；自动覆盖乐观冲突；取消/过期不幂等；把队列当 exactly-once。

练习：实现多 SKU 预留、确认、取消和过期 worker，加入 50 并发测试与崩溃注入。

完成标准：全部或全部不预留；库存不为负；重复请求/消息不重复扣加；确认与过期竞争只有一个终态；锁顺序固定；外部事件通过 outbox；指标能观测冲突、等待、重试和库存对账差异。

## 来源

- [PostgreSQL 18: Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Transaction Isolation](https://www.postgresql.org/docs/18/transaction-iso.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Constraints](https://www.postgresql.org/docs/18/ddl-constraints.html)（访问日期：2026-07-17）
- [Go 1.26: Executing transactions](https://go.dev/doc/database/execute-transactions)（访问日期：2026-07-17）
