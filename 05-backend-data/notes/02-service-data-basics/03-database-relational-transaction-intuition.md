# 数据库、关系表、约束、查询与事务

## 学习目标

本文建立关系数据库入门模型：表、行、列、键与约束如何表达数据不变量，SQL 查询如何形成结果，事务如何把多条语句组成提交或回滚单元。示例以 PostgreSQL 语法为主，并明确方言边界。

## 1. 数据库管理系统

数据库管理系统负责持久存储、查询、并发控制、恢复、权限和协议访问。应用进程通常通过客户端驱动连接数据库服务；连接失败、SQL 错误、事务冲突和约束违反都是不同失败层。

“数据库”可指一个逻辑数据集合，也可能在具体产品中指命名隔离单元。实例、cluster、database、schema 的层级由实现定义，不能跨 PostgreSQL、MySQL、SQLite 直接类推。

SQLite 是嵌入式数据库库，直接操作文件；PostgreSQL 是客户端/服务器数据库。本文 SQL 验证若用 SQLite，只能覆盖双方兼容子集，PostgreSQL identity、并发和隔离语义必须以 PostgreSQL 文档与实例验证。

## 2. 关系、表、行与列

关系模型以关系、属性和元组描述数据；SQL 表是其工程实现，但 SQL 还包含 NULL、重复行和实现类型等特性，不能把所有数学关系性质直接当 SQL 行为。

表定义列名、数据类型、默认值和约束。行是每列一个值的记录。类型决定允许表示和运算，但具体范围、精度、排序规则、时区行为由数据库类型定义。

```sql
CREATE TABLE task (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title TEXT NOT NULL,
  priority SMALLINT NOT NULL,
  done BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT task_priority_range CHECK (priority BETWEEN 1 AND 3)
);
```

`GENERATED ... AS IDENTITY` 与 `TIMESTAMPTZ` 是 PostgreSQL 语法/类型。SQLite 验证兼容版本应使用它支持的自增与时间表示，不能声称两者行为相同。

## 3. NULL 与三值逻辑

NULL 表示缺失或未知，不等于零、空字符串或 false。SQL 条件可能得到 TRUE、FALSE、UNKNOWN；WHERE 只保留 TRUE。

```sql
SELECT * FROM task WHERE completed_at = NULL;  -- 不会按预期匹配
SELECT * FROM task WHERE completed_at IS NULL;
```

`NULL = NULL` 的结果是 UNKNOWN。聚合函数对 NULL 的处理也各有定义，例如 `COUNT(column)` 只计非 NULL，`COUNT(*)` 计行。设计 nullable 列前要定义缺失含义；如果业务要求必填，用 NOT NULL 让数据库保护。

## 4. 主键、候选键与外键

主键唯一标识一行，并隐含唯一与非 NULL。一个表只有一个主键，但可有其他 UNIQUE 候选键。主键应稳定，频繁变化的业务字段通常用唯一约束而不是主键。

```sql
CREATE TABLE project (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE
);

CREATE TABLE task (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id BIGINT NOT NULL REFERENCES project(id),
  title TEXT NOT NULL
);
```

外键要求引用值对应目标唯一/主键行，保护引用完整性。删除/更新目标时可配置 RESTRICT/NO ACTION、CASCADE、SET NULL 等动作；选择必须符合业务生命周期。CASCADE 可能影响大量行，不能为了省代码默认开启。

应用先 SELECT 检查唯一再 INSERT 存在并发竞态；最终由数据库 UNIQUE 约束裁决，应用捕获约束错误映射为业务冲突。

## 5. CHECK、DEFAULT 与约束边界

NOT NULL、UNIQUE、PRIMARY KEY、FOREIGN KEY 和 CHECK 把不变量放到所有写入路径共享的边界。DEFAULT 只在插入未提供值时生成，不验证显式值。

CHECK 适合当前行可判断的稳定规则，例如 priority 范围。跨行总额、外部 HTTP 状态或随时间变化的规则不能简单写成可移植 CHECK；具体数据库对表达式还有额外限制。

数据库约束不替代 API 友好验证。应用可提前返回字段错误，数据库约束防止竞态、脚本、迁移或其他服务绕过。

## 6. SQL 的四类基本操作

INSERT 增加行，SELECT 产生查询结果，UPDATE 修改匹配行，DELETE 删除匹配行。

```sql
INSERT INTO task (title, priority)
VALUES ('read transaction docs', 2)
RETURNING id, title, priority, done;

SELECT id, title, priority
FROM task
WHERE done = FALSE
ORDER BY priority DESC, id ASC
LIMIT 20;

UPDATE task
SET done = TRUE
WHERE id = 42 AND done = FALSE;

DELETE FROM task
WHERE id = 42;
```

`RETURNING` 是 PostgreSQL 支持的扩展，其他数据库支持情况不同。UPDATE/DELETE 没有 WHERE 会作用于所有可见匹配行；维护操作先在事务中 SELECT 目标并核对数量。

查询没有 ORDER BY 时结果顺序不保证。即使当前执行计划总按主键返回，也可因索引、并发和版本改变。分页必须使用确定性排序，排序键相同再补唯一键。

## 7. 查询、连接与聚合

JOIN 按条件组合表。INNER JOIN 只保留双方匹配，LEFT JOIN 保留左侧所有行，右侧不匹配列为 NULL。

```sql
SELECT p.slug, COUNT(t.id) AS open_tasks
FROM project AS p
LEFT JOIN task AS t
  ON t.project_id = p.id AND t.done = FALSE
GROUP BY p.id, p.slug
ORDER BY p.slug;
```

把 `t.done = FALSE` 写在 ON 保留没有开放任务的项目；写到 WHERE 会过滤右侧为 NULL 的行，使结果接近 INNER JOIN。聚合前要确认 join 是否造成一对多重复计数。

参数必须通过驱动参数绑定，不拼接 SQL 字符串。占位符语法由驱动/数据库定义，PostgreSQL 常见 `$1`，database/sql 驱动也各有约定。

## 8. 索引的直觉

索引是额外数据结构，把键组织成更快定位行的形式。它可能帮助 WHERE、JOIN、ORDER BY 和唯一性检查，但占磁盘、缓存，并增加 INSERT/UPDATE/DELETE 成本。

主键和唯一约束通常创建支持索引，但具体实现看数据库文档。外键引用列是否自动建索引是实现相关；PostgreSQL 不自动为引用列创建索引，常需按查询/删除模式评估。

不能用“有索引”推断查询一定使用。优化器依据统计、选择性、数据量和成本选择扫描。用目标数据库 `EXPLAIN` 验证，而不是把某产品计划节点名称当 SQL 标准。

## 9. 事务

事务把一组数据库操作作为提交或回滚单元。PostgreSQL 用 `BEGIN` 开始，`COMMIT` 使修改可见并完成事务，`ROLLBACK` 放弃未提交修改。

```sql
BEGIN;
UPDATE account SET balance_cents = balance_cents - 500 WHERE id = 1;
UPDATE account SET balance_cents = balance_cents + 500 WHERE id = 2;
COMMIT;
```

若第二条失败而没有事务，第一条可能已提交形成资金丢失。事务中失败要回滚；PostgreSQL 中事务内语句错误后，事务进入 aborted 状态，必须 ROLLBACK 或回滚到 savepoint 后才能继续。

ACID 常作四个目标：原子性表示事务修改整体提交/回滚；一致性表示约束与应用不变量在正确事务后成立；隔离性描述并发事务可见关系；持久性表示已提交结果按数据库保证抵抗故障。具体保证依隔离级别、同步提交、存储和配置。

## 10. 事务隔离的直觉

事务不是“把数据库锁住”。多个事务可并发，隔离级别决定允许的可见现象。SQL 标准定义级别与现象，PostgreSQL 以 MVCC 和锁实现自己的精确行为。

在默认 Read Committed 中，一条语句看到语句开始时的已提交快照等，连续两条 SELECT 可能看到其他事务新提交。Repeatable Read、Serializable 提供更强保证但可能需要处理 serialization failure 重试。

任何“先读后写”业务不变量都要分析并发。例如两个请求都读库存 1，再各自扣 1，应用检查可能超卖。可用条件 UPDATE、行锁、Serializable 或原子数据库语句，并检查受影响行数。

事务保持短小：不要在持有事务/锁时等待用户、调用慢 HTTP 或执行无界计算。外部系统不参加本地数据库原子提交，需要 outbox、幂等和补偿等更高层设计。

## 11. 完整案例：原子领取任务

### 11.1 表与不变量

```sql
CREATE TABLE job (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  state TEXT NOT NULL CHECK (state IN ('ready', 'claimed', 'done')),
  worker_id TEXT,
  CONSTRAINT claimed_worker_required CHECK (
    (state = 'ready' AND worker_id IS NULL) OR
    (state IN ('claimed', 'done') AND worker_id IS NOT NULL)
  )
);

INSERT INTO job (state) VALUES ('ready');
```

不变量是 ready 没有 worker，claimed/done 必须有 worker。两个 worker 不能同时成功领取同一 ready job。

### 11.2 原子条件更新

```sql
UPDATE job
SET state = 'claimed', worker_id = $1
WHERE id = $2 AND state = 'ready'
RETURNING id, state, worker_id;
```

输入 `$1='worker-a'`、`$2=1`。数据库在单条 UPDATE 中检查当前状态并修改。第一次返回一行 `(1, claimed, worker-a)`；第二个 worker 对同一 id 执行时更新 0 行，应用返回 409 或“已被领取”。

这比应用先 `SELECT state` 再无条件 UPDATE 更安全，因为判断与写入在同一数据库语句中竞争。仍要在目标 PostgreSQL 并发测试。

### 11.3 Go 事务示意

```go
func claim(ctx context.Context, db *sql.DB, jobID int64, worker string) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return fmt.Errorf("begin claim: %w", err) }
    defer tx.Rollback()

    var gotID int64
    err = tx.QueryRowContext(ctx, `
        UPDATE job
        SET state = 'claimed', worker_id = $1
        WHERE id = $2 AND state = 'ready'
        RETURNING id`, worker, jobID).Scan(&gotID)
    if errors.Is(err, sql.ErrNoRows) { return ErrAlreadyClaimed }
    if err != nil { return fmt.Errorf("update job: %w", err) }
    if err := tx.Commit(); err != nil { return fmt.Errorf("commit claim: %w", err) }
    return nil
}
```

`defer tx.Rollback()` 在已 Commit 后通常返回事务已完成错误，可忽略；它保护提前返回路径。Commit 失败不能当成功，调用者需要知道结果可能不确定并按操作幂等性处理。

### 11.4 验证与失败分支

验证一：同一 job 并发启动 20 次 claim，断言仅一次返回成功，19 次 AlreadyClaimed，最终 worker_id 为成功者。验证二：worker 空字符串触发应用验证；若数据库只检查非 NULL，空串仍合法，说明约束还需 `worker_id <> ''`。

事务中在 UPDATE 后主动返回错误，defer rollback 应使状态仍 ready。杀死数据库连接等故障下，不能仅凭客户端错误断言数据库一定没提交；需要通过幂等读取/重试协议确认。

## 12. 调试清单

- 查询顺序漂移：补完整 ORDER BY 与唯一 tie-breaker。
- 唯一检查偶发重复：把唯一性放数据库约束并处理冲突。
- LEFT JOIN 少行：检查右表过滤是否误放 WHERE。
- 更新影响过多：记录并断言 rows affected，维护操作先预览。
- 事务空闲很久：查事务边界内外部调用和未关闭 rows。
- 死锁/serialization failure：按数据库错误码有限重试整个事务，不只最后语句。
- SQLite 测试通过、PostgreSQL 失败：核对方言、类型、隔离与并发实现。

## 13. 练习

1. 为项目/任务表加入唯一 `(project_id,title)`，模拟两个并发插入并处理冲突。
2. 写 LEFT JOIN 查询同时返回零任务项目，比较过滤条件放 ON 与 WHERE。
3. 为 claim 加 worker 非空 CHECK，并验证数据库直接写也受保护。
4. 用 PostgreSQL 两个会话执行条件 UPDATE，记录一方成功、一方零行的步骤。
5. 设计转账表和 ledger，不在事务内调用外部 HTTP，说明通知如何异步发送。

## 来源

- [PostgreSQL 当前文档：Table Basics](https://www.postgresql.org/docs/current/ddl-basics.html)（访问日期：2026-07-17）
- [PostgreSQL 当前文档：Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-17）
- [PostgreSQL 当前文档：Transactions](https://www.postgresql.org/docs/current/tutorial-transactions.html)（访问日期：2026-07-17）
- [PostgreSQL 当前文档：Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)（访问日期：2026-07-17）
- [PostgreSQL 当前文档：SELECT](https://www.postgresql.org/docs/current/sql-select.html)（访问日期：2026-07-17）
