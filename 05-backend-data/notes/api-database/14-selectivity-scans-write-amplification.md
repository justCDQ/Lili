# 选择性、扫描方式与索引写放大

选择性是一个谓词预计匹配的行比例；规划器用行数、页面数、统计信息和成本参数比较顺序扫描、索引扫描、bitmap 扫描等方案。索引减少某些读取的定位成本，却让写入维护更多数据结构并产生额外 I/O 和 WAL，这就是索引设计中的写放大。

统计对象、计划选项、HOT 和监控字段采用 PostgreSQL 18.4 的语法与行为。

## 选择性不是“不同值数量”

若表有 100 万行，条件匹配 1000 行：

```text
selectivity = matched_rows / total_rows = 1,000 / 1,000,000 = 0.001 = 0.1%
```

这是高过滤能力、低匹配比例的条件。工程讨论中有人把“高选择性”用于“能排除很多行”，也有人用于“匹配比例高”，容易产生歧义。直接说“预计匹配 0.1%”最准确。

**基数（cardinality）** 可指关系行数或列的不同值数量，必须指明语境。列有 100 个不同值并不表示每个等值条件都匹配 1%；数据可能高度倾斜。

### 选择性会组合，但不能总按独立相乘

若：

```text
P(country = 'CN') = 20%
P(city = 'Shanghai') = 3%
```

在独立假设下组合为 `0.2 × 0.03 = 0.6%`。但 city 与 country 强相关，`Shanghai` 基本决定 `CN`，独立相乘会严重低估。规划器可用扩展统计描述函数依赖、不同值组合和多列 most-common-values。

## PostgreSQL 统计信息

`ANALYZE` 抽样并把统计写入系统目录。普通用户应查看 `pg_stats` 视图，而不是直接依赖内部 `pg_statistic`。

```sql
SELECT
  attname,
  null_frac,
  n_distinct,
  most_common_vals,
  most_common_freqs,
  histogram_bounds,
  correlation
FROM pg_stats
WHERE schemaname = 'public'
  AND tablename = 'orders';
```

字段含义：

| 字段 | 含义 | 使用边界 |
|---|---|---|
| `null_frac` | 样本中空值比例 | 影响 `IS NULL` 估算 |
| `n_distinct` | 正数为估算不同值数；负数表示不同值数约为表行数乘其绝对值 | 是估算，不是唯一性约束 |
| `most_common_vals` | 高频值列表 | 只覆盖统计目标允许保存的部分值 |
| `most_common_freqs` | 高频值对应频率 | 处理倾斜分布 |
| `histogram_bounds` | 排除 MCV 后的近似等频边界 | 用于范围选择性估算 |
| `correlation` | 列逻辑顺序与物理行顺序的相关程度，范围为 -1 到 1 | 影响索引访问随机程度估算，不是列间相关性 |

`pg_class.reltuples` 和 `relpages` 是近似行数、页面数，vacuum/analyze 可增量更新，不是实时精确计数。

### 提高单列统计目标

默认样本不能描述极端倾斜时，可对关键列提高统计目标：

```sql
ALTER TABLE orders ALTER COLUMN status SET STATISTICS 1000;
ANALYZE orders (status);
```

更高目标增加 ANALYZE 时间和统计存储，也增加规划时间，不应全库盲目调高。先确认估算偏差与该列统计有关。

### 扩展统计

```sql
CREATE STATISTICS orders_country_city_stats
  (dependencies, ndistinct, mcv)
ON country_code, city
FROM orders;

ANALYZE orders;
```

- `dependencies` 描述列间函数依赖程度，帮助组合等值条件。
- `ndistinct` 改善多列分组或组合不同值数估算。
- `mcv` 保存多列高频组合，处理相关和倾斜。

扩展统计不会创建索引，也不保证任何查询变快；它改善行数估算，使规划器更可能选对已有路径。当前统计能力有适用表达式和查询形态边界，应通过实际计划验证。

PostgreSQL 18 的扩展统计当前不用于表与表之间 JOIN 的选择性估算。它能改善同一基础表上的相关条件与分组估算，但不能据此承诺多表连接基数一定准确。

## 四类常见扫描

### Sequential Scan

顺序扫描按表页面读取并逐行检查条件。它适合：

- 小表；
- 条件匹配表中较大比例；
- 查询需要大部分列和行；
- heap 物理连续读取成本低，而索引随机 heap 访问成本高；
- 没有可用索引或条件不能用于索引。

看到 `Seq Scan` 不等于缺索引。返回 60% 行时，通过索引定位再随机访问 heap 往往比顺序读完整表更贵。

### Index Scan

索引扫描先定位索引项，再按 TID 访问 heap 并检查 MVCC 可见性及未被索引处理的过滤条件。适合少量分散行、需要索引顺序、快速 `LIMIT` 等情况。匹配行很多且物理分散时，大量随机 heap 访问会变贵。

### Index Only Scan

查询所需值都能由索引返回时物理上可行；但仍需 visibility map 判断对应 heap 页是否 all-visible。若不是，执行计划中的 `Heap Fetches` 会增加。覆盖索引只创造可能，不保证无 heap I/O。

### Bitmap Index/Heap Scan

bitmap index scan 收集匹配 TID，bitmap heap scan 按 heap 页面顺序访问，折中处理“匹配不少但仍远小于全表”的范围。它还能对多个索引 bitmap 做 `AND`/`OR` 组合。

bitmap 可能因内存限制变成 lossy 页面粒度，需在 heap 上重新检查条件。它也不能像普通 B-tree index scan 那样直接保持索引顺序，常需额外排序。

## 规划器比较的是估算成本

成本单位不是毫秒。主要因素包括：

- 预计输入和输出行数；
- 表与索引页面数；
- 顺序页面和随机页面读取成本；
- CPU 元组、操作符和索引处理成本；
- 缓存规模估计 `effective_cache_size`；
- 并行启动和处理成本；
- 排序、哈希和聚合的内存与溢写可能性。

```sql
SHOW seq_page_cost;
SHOW random_page_cost;
SHOW cpu_tuple_cost;
SHOW effective_cache_size;
SHOW work_mem;
```

这些参数是成本模型输入，不应为强迫某条查询使用索引而随意修改全局值。先修复统计、SQL 和索引，再以整机 I/O 特征和完整工作负载校准参数。

### `EXPLAIN` 与 `EXPLAIN ANALYZE`

```sql
EXPLAIN (COSTS, VERBOSE, SETTINGS)
SELECT id FROM orders WHERE tenant_id = 42 AND status = 'paid';
```

普通 `EXPLAIN` 不执行查询，显示估算计划。实际验证：

```sql
EXPLAIN (ANALYZE, BUFFERS, WAL, SETTINGS, SUMMARY)
SELECT id FROM orders WHERE tenant_id = 42 AND status = 'paid';
```

`ANALYZE` 会真正执行语句。对写操作可在测试库运行，或在事务中执行后 `ROLLBACK`，但序列值、外部副作用和某些函数效果未必完全回滚；生产环境必须谨慎。

阅读重点：

- `cost=启动..总成本 rows=估算行数 width=估算宽度`。
- `actual time=首行..全部 rows=每次循环行数 loops=循环次数`。
- 总处理行数需结合 `rows × loops`。
- `Rows Removed by Filter` 是否很大。
- shared/local/temp blocks 的 hit、read、dirtied、written。
- sort method、内存和 disk 使用。
- 估算行数与实际行数的数量级偏差。

单次执行时间受缓存、并发、检查点和系统负载影响。比较方案要多次采样，并同时看 I/O、CPU、锁和尾延迟。

## 索引写放大从哪里来

### INSERT

每插入一行，数据库写 heap，并为每个适用索引插入索引项。可能发生页面分裂、WAL 增加、缓存页面被挤出。唯一索引还要检查冲突；部分索引只为满足谓词的行维护项。

### UPDATE

PostgreSQL MVCC 更新通常创建新行版本。若更新了索引键或 included 列，相应索引必须写新项。即使值未改变，具体执行和存储行为也应从实际 WAL 与统计观察。

HOT（Heap-Only Tuple）更新可避免为新版本建立新的索引项，基本前提包括：

- 更新没有改变任何索引所引用的列；
- 同一 heap 页面有足够空间放新版本。

核心 PostgreSQL 中的 BRIN 属于 summarizing index，不受第一项同样限制，但仍可能需要更新摘要。其他索引的表达式、部分索引谓词和 `INCLUDE` 列会使相关更新失去 HOT 资格。降低表 `fillfactor` 可预留页面空间提高 HOT 机会，但增加表大小和扫描页数，需要测量。

### DELETE 与软删除

删除产生 dead tuple，索引项不会在事务提交时立即物理消失，由 vacuum 后续清理。软删除是 UPDATE；若大量部分索引使用 `WHERE deleted_at IS NULL`，改变 `deleted_at` 会同时改变这些索引的成员资格。

### 复制和备份链路

更多 WAL 会增加：

- WAL 存储与归档流量；
- 流复制网络和回放压力；
- 只读副本延迟风险；
- 增量备份与恢复处理量；
- 检查点写入压力。

因此读查询加速不能只报告延迟下降，还要报告写吞吐、WAL bytes、索引大小、缓存和复制延迟变化。

## 测量索引成本

### 大小和使用

```sql
SELECT
  indexrelname,
  pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
  idx_scan,
  last_idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_indexes
WHERE relname = 'orders'
ORDER BY pg_relation_size(indexrelid) DESC;
```

`idx_scan = 0` 只表示自统计重置以来未记录扫描，不能单独作为删除依据。约束索引、低频月末任务、故障处理和季节性查询都可能重要。

### 表写入和 HOT

```sql
SELECT
  relname,
  n_tup_ins,
  n_tup_upd,
  n_tup_hot_upd,
  n_tup_del,
  n_dead_tup,
  last_autovacuum,
  last_autoanalyze
FROM pg_stat_user_tables
WHERE relname = 'orders';
```

可计算观察窗口内 `n_tup_hot_upd / n_tup_upd`，但计数自统计重置累计且可能存在采集延迟。必须记录时间窗、重置点和负载变化。

### WAL

`EXPLAIN (ANALYZE, WAL)` 可展示单语句计划节点相关 WAL 记录、full-page images 和 bytes。系统层可观察 `pg_stat_wal`，基准前后记录差值。不要把单次小样本线性外推为生产日写入量。

## 完整案例：是否为低区分度状态建索引

### 输入

`jobs` 表 1000 万行：

- `done` 约 97%。
- `pending` 约 2%。
- `failed` 约 1%。
- worker 每秒插入和更新大量任务。
- 队列查询只取某租户最早 100 个 pending 任务：

```sql
SELECT id, payload
FROM jobs
WHERE tenant_id = $1 AND status = 'pending'
ORDER BY available_at, id
LIMIT 100;
```

只看 `status` 的不同值数量会误判“只有三个值，索引无用”。真正条件还包含租户、稀少的 pending 谓词、排序和 LIMIT。

### 方案 A：无队列索引

可能顺序扫描或使用仅 tenant 的索引，再过滤 status、排序。优点是写入维护少；缺点是每次取 100 个任务可能读取大量 done 行。

### 方案 B：完整复合索引

```sql
CREATE INDEX jobs_tenant_status_available_idx
ON jobs (tenant_id, status, available_at, id)
INCLUDE (payload);
```

支持多状态查询，但索引覆盖所有 1000 万行，`payload` 若很宽会造成巨大索引与写成本，且状态更新会维护键。

### 方案 C：pending 部分索引

```sql
CREATE INDEX jobs_pending_queue_idx
ON jobs (tenant_id, available_at, id)
WHERE status = 'pending';
```

索引只保存 pending 行，较小且直接支持队列顺序。任务从 pending 变 done 时仍要移除其索引成员；但 done 行无需留在该索引。查询必须使用能让规划器证明 `status = 'pending'` 的谓词。

不把宽 `payload` include，可在定位 100 个 ID 后访问 100 个 heap 行，通常比让每次 payload 更新扩大索引更可控。是否 include 要实测。

### 步骤一：收集分布和估算偏差

```sql
ANALYZE jobs;

SELECT status, count(*)
FROM jobs
GROUP BY status;

SELECT most_common_vals, most_common_freqs
FROM pg_stats
WHERE tablename = 'jobs' AND attname = 'status';
```

精确 group 统计应在可承受环境执行；生产大表可从受控分析副本获取。再分别对大租户、小租户采样计划，避免只测平均租户。

### 步骤二：比较读取

对 A/B/C 使用相同数据快照和参数，记录：

```sql
EXPLAIN (ANALYZE, BUFFERS, SETTINGS)
SELECT id, payload
FROM jobs
WHERE tenant_id = 42 AND status = 'pending'
ORDER BY available_at, id
LIMIT 100;
```

验收观察点：返回相同 ID 顺序、估算与实际行数、读取页面、被过滤行、是否排序、p50/p95 延迟。

### 步骤三：比较写入

回放代表性生命周期：批量插入 pending、更新可用时间、领取后改 running、完成后改 done。每个方案记录：

- transactions/s 与 p95 提交延迟；
- WAL bytes；
- heap 和索引大小；
- `n_tup_upd` 与 HOT 比例；
- autovacuum、dead tuples 和副本延迟。

### 输出决策

若工作负载证实部分索引显著降低队列读取，且写入/WAL 在预算内，选择方案 C。决策输出应包含索引 DDL、目标查询、基准数据、构建方式、监控阈值和回滚步骤，而不是“低选择性列不该建索引”的通则。

### 失败分支

- 测试数据中 status 均匀分布，生产却 97% done，计划结论无效。
- 加完整索引并 include 大 payload，读快少量但写吞吐和索引大小恶化。
- 看到 Seq Scan 就设置 `enable_seqscan = off` 作为永久修复，掩盖统计或模型问题。
- `ANALYZE` 过期导致 pending 估算错误；先更新统计并查分布。
- 删除 `idx_scan = 0` 的唯一索引，破坏约束或低频关键路径。

## 安全删除冗余索引

1. 确认索引不支撑主键、唯一或排他约束。
2. 比较定义，注意排序、opclass、collation、谓词和 include 列差异。
3. 覆盖完整业务周期观察使用，并查询应用 SQL/监控。
4. 在预发布回放读写负载，记录计划变化。
5. 使用适合生产的并发删除方式并了解其事务限制。
6. 保留重建 DDL、预计构建时长、磁盘预算和回滚条件。

## 练习：事件表索引预算

事件表每天写入 5000 万行，查询包括按租户查最近 100 条、按低频错误码查 24 小时事件、按时间全量归档。设计索引并证明每个索引值得维护。

完成标准：

- 给出每个条件预计匹配比例，避免只写“高/低选择性”。
- 为三条查询比较 seq、index、bitmap、BRIN/B-tree 候选的边界。
- 用 `EXPLAIN (ANALYZE, BUFFERS, WAL)` 和统计视图定义验证步骤。
- 记录索引大小、写吞吐、WAL、HOT、vacuum 和副本延迟。
- 注入陈旧统计、倾斜租户和相关列三种失败情形。
- 给出新增、观察、回滚和删除索引的生产流程。

## 来源

- [PostgreSQL 18：Using EXPLAIN](https://www.postgresql.org/docs/18/using-explain.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Statistics Used by the Planner](https://www.postgresql.org/docs/18/planner-stats.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Planner Cost Constants](https://www.postgresql.org/docs/18/runtime-config-query.html#RUNTIME-CONFIG-QUERY-CONSTANTS)（访问日期：2026-07-17）
- [PostgreSQL 18：Heap-Only Tuples](https://www.postgresql.org/docs/18/storage-hot.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Cumulative Statistics Views](https://www.postgresql.org/docs/18/monitoring-stats.html)（访问日期：2026-07-17）
