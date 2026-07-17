# PostgreSQL B-tree、复合索引与覆盖索引

PostgreSQL B-tree 索引按操作符类定义的顺序保存键，支持等值、范围、排序和部分模式前缀查询。复合索引把多个键按确定顺序组织；`INCLUDE` 把只用于返回的列作为 payload 保存，使查询在满足 MVCC 可见性条件时可能使用 index-only scan。

Skip scan、`INCLUDE`、`CREATE INDEX CONCURRENTLY` 和执行计划示例以 PostgreSQL 18.4 为基线。

## 索引和 heap 的责任边界

PostgreSQL 表的主要数据区称为 heap。B-tree 索引是独立结构，索引项包含键和指向表行版本的 TID。普通 index scan 的路径是：

```mermaid
flowchart LR
    Q["查询谓词"] --> R["B-tree 根/内部页定位范围"]
    R --> L["叶子页读取键与 TID"]
    L --> H["heap 取行并检查 MVCC 可见性"]
    H --> O["返回所需列"]
```

树保持有序和平衡，查找从高层页面缩小到叶子范围；叶子页面通过相邻关系支持范围扫描。索引提高定位速度，但增加磁盘、缓存、WAL 和写入维护，不会保存查询结果。

PostgreSQL 文档和 SQL 语法称该访问方法为 **B-tree**。不要把其他数据库产品关于聚簇索引、主键组织表或“B+Tree 固定实现”的结论直接套用到 PostgreSQL。

## B-tree 能支持哪些条件

默认操作符类通常支持：

- `<`、`<=`、`=`、`>=`、`>` 范围和等值。
- `BETWEEN`、`IN` 等可转化为相应比较的条件。
- `IS NULL`、`IS NOT NULL`。
- `ORDER BY`，因为 B-tree 可以按键顺序向前或向后扫描。
- 在排序规则允许时，锚定开头且不是不区分大小写的模式，例如 `LIKE 'foo%'`；非 C locale 下可能需要合适的 pattern operator class。

普通 B-tree 不能直接支持：

- `LIKE '%foo%'` 的任意包含搜索。
- 对未建立表达式索引的 `lower(email)` 条件。
- JSON 包含、全文检索、几何邻近等其他访问方法擅长的操作。

表达式必须与索引表达式匹配：

```sql
CREATE INDEX app_users_lower_email_idx ON app_users (lower(email));

SELECT id FROM app_users WHERE lower(email) = lower($1);
```

表达式索引让读取可定位，但每次插入和相关更新都要计算并维护表达式。

## 排序方向和空值顺序

单列 B-tree 默认按 `ASC NULLS LAST` 保存；向后扫描可满足 `DESC NULLS FIRST`。多列索引需要同时考虑各列方向：

```sql
CREATE INDEX events_priority_time_idx
ON events (priority DESC, created_at ASC, id ASC);
```

它可直接产生 `priority DESC, created_at ASC, id ASC`，向后则产生全部方向和空值顺序的反向。若目标排序是一列升序、一列降序，定义方向有实际意义。最终是否使用索引排序取决于过滤行数、成本和所需列，规划器也可能选择顺序扫描后排序。

## 复合 B-tree 的准确规则

对索引 `(a, b, c)`，PostgreSQL 可以使用涉及任意子集键列的条件，但扫描效率不同。

在不考虑 skip scan 时：

1. 前导列上的等值条件用于缩小扫描区间。
2. 第一列没有等值的键上的不等式条件继续限定扫描边界。
3. 更右侧条件仍可在索引中检查并减少 heap 访问，但不一定减少需要扫描的索引范围。

例如：

```sql
CREATE INDEX orders_abc_idx ON orders (tenant_id, status, created_at);
```

| 查询条件 | 一般用途 |
|---|---|
| `tenant_id = ?` | 使用首列范围 |
| `tenant_id = ? AND status = ?` | 使用前两列等值范围 |
| `tenant_id = ? AND status = ? AND created_at >= ?` | 前两列等值加第三列范围，匹配良好 |
| `tenant_id = ? AND created_at >= ?` | status 缺失；可能扫描该租户多个 status 范围或使用 skip scan |
| `status = ?` | 缺少首列；是否值得用取决于 skip scan、基数和成本 |

“只有满足最左前缀才能使用索引”是过度简化。PostgreSQL 可能扫描整个索引、组合多个索引，或者执行 skip scan。正确结论是：前导等值通常最能限制连续扫描范围，右侧列仍可能参与，最终由成本模型选择。

## PostgreSQL 18 的 skip scan

Skip scan 在缺少某个前导列等值时，内部为该列的可能值反复生成动态等值查找，以利用更右列条件。

索引 `(region, customer_id)`，查询只有：

```sql
SELECT * FROM orders WHERE customer_id = 7700;
```

若 `region` 只有很少不同值，规划器可能执行类似：

```text
region = 'apac' AND customer_id = 7700
region = 'emea' AND customer_id = 7700
region = 'amer' AND customer_id = 7700
```

它跳过叶子页中不可能匹配的区间，比完整索引扫描更少读取。若 `region` 有大量不同值，重复定位成本可能高于顺序扫描或独立的 `customer_id` 索引，规划器通常不会选择 skip scan。

因此：

- skip scan 是成本驱动的计划选择，不是 DDL 开关或每次保证。
- 前导列 distinct 数、条件选择性、表大小、缓存和统计信息都影响选择。
- 频繁且延迟敏感的后导列查询不能只因“可能 skip scan”就放弃专用索引。
- 使用 `EXPLAIN (ANALYZE, BUFFERS)` 在代表性数据上验证读取范围和实际时间。

## 列顺序怎样决定

“等值列在前、范围列在后”是有用起点，不是完整算法。索引设计应建立查询矩阵：

| 查询 | 频率 | 过滤 | 排序 | 返回列 |
|---|---:|---|---|---|
| 租户已支付订单时间线 | 高 | tenant、status | created_at/id 降序 | id、total |
| 租户全部订单时间线 | 高 | tenant | created_at/id 降序 | id、status、total |
| 按外部支付号查订单 | 中 | tenant、payment_ref | 无 | id、status |
| 全局按状态批处理 | 低 | status | id 升序 | id、tenant |

一条 `(tenant_id, status, created_at DESC, id DESC)` 很适合第一条，却不能直接以同样效率服务第二条，因为 `status` 形成缺口。可能方案是再建 `(tenant_id, created_at DESC, id DESC)`，或根据 status distinct 数和读写比例验证 skip scan。索引数量是工作负载取舍，不是越多越完整。

复合 B-tree 索引键加 `INCLUDE` 列默认最多 32 列，此编译期上限可以被修改。实际设计远早于上限就会遇到索引元组过大、缓存效率下降和写入成本；官方建议多列索引谨慎使用，超过三键通常只适合高度固定的查询形态。

## 覆盖索引与 `INCLUDE`

```sql
CREATE INDEX orders_tenant_status_created_idx
ON orders (tenant_id, status, created_at DESC, id DESC)
INCLUDE (total, currency);
```

键列和 include 列的区别：

| 属性 | 键列 | `INCLUDE` 列 |
|---|---|---|
| 参与树的搜索顺序 | 是 | 否 |
| 可用于索引扫描谓词 | 是 | 只作为 payload，不能作为该索引搜索键 |
| 决定索引排序 | 是 | 否 |
| 唯一索引的唯一性范围 | 是 | 否 |
| 可直接返回值 | 支持的访问方法下可以 | 支持的访问方法下可以 |

唯一覆盖索引：

```sql
CREATE UNIQUE INDEX users_tenant_email_uq
ON app_users (tenant_id, email)
INCLUDE (display_name);
```

只保证 `(tenant_id, email)` 唯一，`display_name` 不参与唯一性。

### Index-only scan 的三个条件

1. 索引访问方法能返回原值；B-tree 总是支持，其他类型取决于操作符类。
2. 查询所需列都在索引中，或部分索引谓词已让某个条件无需运行时重检。
3. 对应 heap 页的 all-visible 位足够多，才能避免实际 heap fetch。

索引项不保存行版本对当前 MVCC 快照是否可见的信息。index-only scan 先查看 visibility map；all-visible 为真时直接返回，否则仍访问 heap 检查。因此刚经历频繁更新的表即使“覆盖”所有列，也可能显示大量 `Heap Fetches`。

`VACUUM` 维护 visibility map，但不能为了强迫 index-only scan 而在生产频繁手工 vacuum。应根据表变化模式、autovacuum 和查询收益整体判断。

### INCLUDE 的代价

- 索引更大，缓存同样内存能容纳的页面更少。
- 插入和被包含列更新需要维护索引。
- 宽值可能超过索引元组大小上限，使写入失败。
- 频繁变化表难以从 index-only scan 获得持续收益。
- 覆盖某个查询不代表覆盖工作负载；返回列变化会使索引快速膨胀。

只 include 小、稳定、频繁返回且能显著减少 heap I/O 的列。

## 单列索引组合与复合索引

PostgreSQL 可以把多个索引组合为 bitmap scan：

```sql
CREATE INDEX orders_tenant_idx ON orders (tenant_id);
CREATE INDEX orders_status_idx ON orders (status);
```

对 `tenant_id = ? AND status = ?`，规划器可能对两个 bitmap 做 `AND`。优点是单列查询也可分别使用；缺点是组合会丢失索引的原始顺序，通常还需要排序。固定的过滤加排序查询往往由复合索引更直接支持。

是否保留单列、复合或两者都保留，应比较：

- 查询频率和延迟目标；
- 更新频率与写入吞吐；
- 索引大小和缓存命中；
- 是否服务外键检查、唯一约束或低频关键任务；
- 实际计划而非索引名称猜测。

## 完整案例：订单时间线索引

### 输入工作负载

`orders` 有 2000 万行、5000 个租户。典型租户查询：

```sql
SELECT id, created_at, total, currency
FROM orders
WHERE tenant_id = $1
  AND status = 'paid'
  AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT 50;
```

写入持续发生，订单创建后状态会从 `pending` 更新到 `paid`，金额不再变化。

### 第一步：从谓词和排序生成候选

- `tenant_id`、`status` 是等值。
- `(created_at, id)` 是范围边界和排序。
- `total`、`currency` 只用于返回。

候选：

```sql
CREATE INDEX CONCURRENTLY orders_paid_timeline_idx
ON orders (tenant_id, status, created_at DESC, id DESC)
INCLUDE (total, currency);
```

`CONCURRENTLY` 降低阻塞普通写入的时间，但需要更长构建、更多工作，并且不能放在事务块中；失败后可能留下 invalid index，需要检查并处理。

### 第二步：获取基线和计划

先在测试或可控环境：

```sql
ANALYZE orders;

EXPLAIN (ANALYZE, BUFFERS, WAL, SETTINGS)
SELECT id, created_at, total, currency
FROM orders
WHERE tenant_id = 42
  AND status = 'paid'
  AND (created_at, id) < ('2026-07-17 08:00:00+00', 900000)
ORDER BY created_at DESC, id DESC
LIMIT 50;
```

`EXPLAIN ANALYZE` 会执行查询；对修改语句必须包在事务中回滚，不能在生产随意运行。观察：

- 节点是否为 `Index Scan` 或 `Index Only Scan`。
- `Index Cond` 是否包含租户、状态和复合范围。
- `actual rows` 与估算 `rows` 偏差。
- shared hit/read blocks。
- index-only scan 的 `Heap Fetches`。
- planning 与 execution time。

不能预先承诺具体节点。若租户几乎全表、缓存状态不同或统计信息异常，规划器可能选择其他路径。

### 第三步：验证 index-only 收益边界

在只读较稳定的历史分区，vacuum 后 all-visible 页比例高，覆盖索引可能显著减少 heap fetch。在当前活跃分区，状态更新会清除页面 all-visible 位，index-only scan 仍可能访问 heap。

因此可比较两个候选：

```sql
-- A：较窄，写入和缓存成本低
CREATE INDEX orders_timeline_narrow_idx
ON orders (tenant_id, status, created_at DESC, id DESC);

-- B：覆盖返回列，历史只读数据可能减少 heap I/O
CREATE INDEX orders_timeline_covering_idx
ON orders (tenant_id, status, created_at DESC, id DESC)
INCLUDE (total, currency);
```

用相同数据和并发负载比较 p50/p95 延迟、缓冲区、索引大小和写入吞吐，而不是只比较单次 warm-cache 查询。

### 输出与验收

验收不是“计划显示索引名”，而是：

1. 返回结果与无索引基准 SQL 完全相同且顺序稳定。
2. 深游标仍只访问与 50 行同量级的候选范围。
3. 延迟目标在代表性冷/热缓存和并发下满足。
4. 写入吞吐、WAL 和存储增长在预算内。
5. 索引构建、失败恢复和删除旧索引有运维步骤。

### 失败分支

- 把 `created_at` 放在 `status` 前，查询 status 等值加时间范围时仍可能可用，但键顺序与主要查询的等值/范围结构不匹配，需要读取更大区间。
- 只建 `(tenant_id, status)`，数据库还要排序或 top-N sort，深范围效率不足。
- 把所有返回列都 `INCLUDE`，索引过宽，写吞吐和缓存恶化。
- 看到 `Index Only Scan` 就认为无 heap I/O，却忽略 `Heap Fetches` 很高。
- 因为 PostgreSQL 18 支持 skip scan，就让高频 `WHERE payment_ref = ?` 依赖低效偶然计划；应验证并可能建专用索引。

## 检查索引定义和使用

```sql
SELECT
  indexrelname,
  idx_scan,
  last_idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_indexes
WHERE relname = 'orders'
ORDER BY indexrelname;
```

统计计数自统计重置以来累计，且低频关键索引、约束索引和灾备查询索引不能仅凭 `idx_scan = 0` 删除。还要核对：

```sql
SELECT pg_size_pretty(pg_relation_size('orders_paid_timeline_idx'));
```

删除疑似重复索引前，检查唯一/主键约束依赖、所有查询形态、分区索引、写入维护和一个完整业务周期。

## 练习：工单队列索引

工单查询按租户和负责人筛选，未完成状态按 `priority DESC, created_at ASC, id ASC` 取前 100 条；后台还有只按负责人跨租户处理的授权任务。

完成标准：

- 列出两个查询的频率、过滤、排序和返回列。
- 分别提出复合索引、部分索引和专用后台索引候选。
- 解释哪些条件限制扫描范围，哪些只在索引内过滤。
- 构造前导列 distinct 少和多两组数据，观察 skip scan 是否被选择。
- 比较普通 index scan、bitmap scan、index-only scan 的 heap 与排序代价。
- 给出写入吞吐、索引大小、WAL、构建和回滚验收标准。

## 来源

- [PostgreSQL 18：B-Tree Indexes](https://www.postgresql.org/docs/18/btree.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Multicolumn Indexes](https://www.postgresql.org/docs/18/indexes-multicolumn.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Indexes and ORDER BY](https://www.postgresql.org/docs/18/indexes-ordering.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Index-Only Scans and Covering Indexes](https://www.postgresql.org/docs/18/indexes-index-only-scans.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Examining Index Usage](https://www.postgresql.org/docs/18/indexes-examine.html)（访问日期：2026-07-17）
