# 在线回填、零停机发布与可执行回滚

“零停机”不是没有风险或没有延迟，而是在迁移期间持续满足明确的服务目标：旧/新版本共存、请求正确、延迟和错误率在预算内、数据可验证、异常时能暂停或回退。

## 1. 定义发布不变量

迁移前写出可测条件：

- 旧版与新版应用都能读写 expand 后 schema。
- 新写入不会遗漏目标列。
- 历史回填不覆盖用户更新。
- 主库 p95、锁等待、复制延迟、WAL 和磁盘不超过停止阈值。
- 回填可暂停、续跑和重复执行。
- 读切换可在分钟级回退，contract 前不删除回退所需数据。

这些条件比“脚本执行成功”更接近完成。

## 2. 在线回填的数据竞态

简单脚本“扫描旧列→计算→更新新列”会与在线写入竞争。若用户在扫描和 UPDATE 之间修改，回填可能把新值覆盖为旧派生值。

使用条件更新保护：

```sql
UPDATE accounts
SET normalized_email = lower(email),
    normalization_version = 2
WHERE account_id > $1
  AND account_id <= $2
  AND normalized_email IS NULL;
```

若计算必须基于读到的旧版本，带版本条件：

```sql
UPDATE accounts
SET target = $computed, migration_version = 2
WHERE account_id = $id
  AND source_version = $observed_version
  AND target IS NULL;
```

零影响行表示在线更新或其他 worker 竞争，重新读取而不是覆盖。

## 3. 稳定分批与 checkpoint

选择不可变或单调主键 keyset，不用不断变大的 OFFSET：

```sql
SELECT account_id
FROM accounts
WHERE account_id > $last
  AND target IS NULL
ORDER BY account_id
LIMIT 1000;
```

UUID 无自然顺序也可按 `(created_at, id)` 或哈希分片处理，但游标和边界必须稳定。多个 worker 可预分配不重叠范围，或用任务表/`SKIP LOCKED` 领取；不能各自从头扫造成放大。

checkpoint 记录 job/version、range、last key、处理/跳过/错误数、开始结束时间。提交数据批次后再更新 checkpoint，或让二者在同一事务；即使 checkpoint 落后，幂等条件允许安全重跑。

## 4. 批大小与节流

批越大吞吐高但事务、锁、WAL 峰值和回滚时间更大。采用反馈控制而非固定全速：

```text
if replica_lag > 10s or db_cpu > 75% or lock_wait_p95 > 100ms:
    pause or reduce batch/rate
else if system has headroom:
    slowly increase rate
```

停止阈值按实际 SLO 定义。监控：rows/s、batch duration、dead tuples、WAL bytes/s、replication lag、checkpoint I/O、autovacuum、锁等待、应用 p95/p99、错误率和 pool wait。

回填 UPDATE 会创建新行版本并增加 WAL/表膨胀。即使每批很短，总量仍可能触发存储和复制压力；演练必须接近生产数据规模和分布。

## 5. 双写策略

### 应用双写

优点：逻辑可见、容易逐步发布；缺点：所有写入口必须升级，包括后台任务、管理脚本、导入和旧客户端。两列应在同一事务写，不能先写旧列后异步写新列而无对账。

### 数据库 trigger

优点：覆盖所有数据库写入者；缺点：隐藏逻辑、额外写成本、递归/顺序与移除风险。触发器函数要版本化、测试 bulk load 和复制行为。

### Change Data Capture

适合跨存储迁移，但通常异步并有延迟/重复/乱序。需要事件位置、幂等、删除语义、schema 演进和全量+增量切换水位。不能宣称瞬时一致。

## 6. 影子读与差异验证

新读路径先不影响用户结果，后台读取新模型并与旧结果比较。影子请求必须有资源预算，不能把生产流量翻倍压垮数据库。

差异分类：等价格式差异、已知舍入、真实数据遗漏、顺序不稳定、权限差异。指标 label 不放资源 ID；详细差异进入受限抽样日志，脱敏并设保留期。

验证方法组合：

- 精确计数：`target IS NULL`、非法状态、孤儿外键。
- 聚合对账：按 tenant/day 的 count/sum。
- 哈希：按稳定排序和规范化格式做分片 hash。
- 随机与风险抽样：大租户、边界金额、多语言字符。
- 业务不变量：总额守恒、状态迁移合法、跨租户为零。

只比较总行数会漏掉“一行丢、一行重复”。

## 7. Canary 与读切换

切换层次可为：内部用户→1%租户→10%→50%→100%。租户级切换比请求随机切换更易保持一致缓存和排障。每阶段比较错误、延迟、差异、数据库负载和工单。

feature flag 本身要有审计、默认值和依赖失效策略。若配置服务不可用，明确保持上次值或回到安全旧路径；不能随机变化。

缓存 key/version 必须区分新旧表示，切读时清理或自然过期。异步任务也可能在数小时后运行旧代码/旧 payload，不能只看在线副本已升级。

## 8. 回滚的层级

### 8.1 暂停

第一动作往往是暂停回填/扩容，保留已处理数据，降低负载。脚本必须响应取消并在当前小批提交后安全退出。

### 8.2 流量回退

feature flag 切回旧读。前提是双写仍保持旧结构可用，新数据能由旧模型表达。切回后继续监控，不能立即执行 destructive down migration。

### 8.3 应用版本回退

旧二进制必须兼容 expand schema 和新写入。若新应用已经产生旧版无法理解的枚举/状态，回退会失败；应先禁用新功能、迁移/隔离这些数据。

### 8.4 数据修复/恢复

错误回填可用记录的 job ID/version 选择性修复。如果旧值未保留且变换不可逆，只能从备份/PITR 恢复到旁路实例，再提取受影响数据；这不是快速 `down.sql`。

回滚脚本本身也要先演练、分批、带条件和审计。

## 9. 部署顺序示例

目标将 `orders.customer_email` 转为 `customer_contact_id`：

1. Expand 新 contact 表和可空 FK 列，旧应用不受影响。
2. 部署双写版本：创建/更新订单时在同一事务 upsert contact 并写 ID，同时保留 email。
3. 运行限速回填，按 order ID 分批，条件为 contact_id IS NULL。
4. 验证每租户 count、孤儿、email 映射和跨租户不变量。
5. 影子读 contact，与旧 email 比较。
6. canary 切新读；逐步扩大。
7. 停止旧写，观察旧列变化归零。
8. 添加并验证 FK/NOT NULL。
9. 经过回退窗口后在独立发布 contract 删除旧列。

## 10. 完整案例：一亿行邮箱规范化

### 输入

`accounts` 一亿行，需要新增 `email_lookup_key` 用于精确匹配。业务仍在线；用户可能修改邮箱；读副本延迟 SLO < 10 s；主库 p95 增幅不得超过 10%。

### 步骤

1. 定义规范化算法版本和 Unicode/域名规则；新增可空 lookup key 与 version。
2. 应用写路径在同一事务按原始 email 计算 key；唯一约束范围包含 tenant，并先处理历史冲突。
3. 回填 worker 按主键 keyset 每批 500，`WHERE key IS NULL` 条件更新，记录 checkpoint。
4. 控制器每批读取 replica lag、DB CPU、lock wait 和应用 p95；超阈值暂停 60 秒并降低速率。
5. 每百万行做 tenant 分组计数和重复 key 验证；完成后全量 NULL=0。
6. 新查询先影子运行，差异分类；canary 按租户切读。
7. 观察旧路径 14 天后再验证/收紧约束，暂不删除原始 email（它仍是展示数据）。

### 输出

回填可中断续跑，用户在线修改不被覆盖；查询逐步切换，数据库和复制指标保持阈值内。算法版本让未来规则更新可再迁移。

### 验证

- 随机杀死 worker，重启从 checkpoint 继续且无错误覆盖。
- 人工制造 replica lag，节流器暂停。
- 回填与用户修改同一行并发，最终 key 对应最新 email。
- 影子读在大小租户和国际化邮箱样本中无未知差异。
- 关闭新读 flag 后旧路径立即恢复且新写仍双写。

### 失败分支

若回填先 SELECT email、长时间计算、再无条件 UPDATE，用户在间隔修改会得到旧 key。修正为在单条条件 UPDATE 中计算，或比较源版本。若新版本立刻停止旧写，应用回退后新订单缺旧字段；应保持双写直到回退窗口结束。

## 11. 故障演练

演练至少包括：DDL lock timeout；回填进程崩溃；重复启动两个 worker；主库 failover；复制延迟；磁盘接近阈值；feature flag 服务故障；新旧读结果不一致；回滚时存在新版状态。

每次演练记录检测信号、停止动作、数据状态、恢复步骤和耗时。没有实际演练的 runbook 只是文档假设。

## 12. 常见错误与练习

错误包括：单次 UPDATE 全表；OFFSET 分页；回填无条件覆盖；只监控 worker rows/s；完成后立即 drop 旧列；down migration 直接逆转不可逆数据；影子流量无预算；contract 前未检查后台任务和旧客户端。

练习：设计 5000 万订单的货币单位迁移，从 decimal 元到 bigint 分，包含舍入、异常数据、双写、回填、对账、canary 和回滚。

完成标准：规则版本化；每批短事务且 keyset；并发更新不覆盖；节流阈值明确；sum/count/异常对账；旧新应用共存；流量回退分钟级；不可逆数据有旁路恢复方案；contract 独立审批。

## 来源

- [PostgreSQL 18: ALTER TABLE](https://www.postgresql.org/docs/18/sql-altertable.html)（访问日期：2026-07-17）
- [PostgreSQL 18: CREATE INDEX](https://www.postgresql.org/docs/18/sql-createindex.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Monitoring Database Activity](https://www.postgresql.org/docs/18/monitoring-stats.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Routine Vacuuming](https://www.postgresql.org/docs/18/routine-vacuuming.html)（访问日期：2026-07-17）
- [PostgreSQL 18.4 Release Notes](https://www.postgresql.org/docs/release/18.4/)（访问日期：2026-07-17）
