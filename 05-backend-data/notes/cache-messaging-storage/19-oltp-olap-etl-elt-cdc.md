---
stage: intermediate
direction: backend-data
topic: oltp-olap-etl-elt-cdc
---

# OLTP、OLAP、ETL、ELT 与 CDC

OLTP 优化小范围事务读写，OLAP 优化大范围扫描聚合。ETL/ELT 描述变换位置，CDC 捕获已提交变化；这些是工作负载和数据流模式，不是互斥产品标签。

## 1. 数据流与决策边界

```mermaid
flowchart LR
 A["源事实/事件"] --> B["传输与处理"]
 B --> C["模型/存储"]
 C --> D["查询与指标"]
 D --> E["质量、对账与重建"]
```

每个输出必须能追溯输入水位、代码或模型版本、质量结果和权限范围；低延迟不能替代可恢复性。

## 2. OLTP

机制：短事务、点查、小批写、约束和并发控制。

实际用途：订单、余额、库存事实。

失败方式：在主库跑无界全表BI。

验证：事务p99、锁、WAL、连接。

取舍：强一致写与分析扫描资源冲突。

OLTP 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 3. OLAP

机制：列式读取、压缩、向量化和数据跳过。

实际用途：跨亿行group/filter分析。

失败方式：用其替代高并发逐行事务。

验证：scan bytes、rows/s、query p95。

取舍：吞吐高但更新/事务语义不同。

OLAP 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 4. ETL

机制：抽取后在进入目标前转换。

实际用途：强清洗、固定仓库模型。

失败方式：黑盒脚本无原始重放。

验证：输入输出行数/hash和版本。

取舍：目标干净但变更需重跑管道。

ETL 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 5. ELT

机制：先加载原始/近原始数据，再在目标计算。

实际用途：云仓库/lake保留可重算历史。

失败方式：把未治理原始表直接开放。

验证：分层权限、模型测试和成本。

取舍：灵活但计算/治理转移到目标。

ELT 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 6. CDC

机制：读取WAL/binlog变化与事务顺序。

实际用途：低延迟同步、缓存/搜索/OLAP。

失败方式：把CDC事件当完整领域语义。

验证：LSN/offset、snapshot水位、lag。

取舍：覆盖写者但耦合物理schema。

CDC 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 7. Snapshot

机制：提供历史基线并与增量衔接。

实际用途：首次加载/重建。

失败方式：snapshot期间直接接当前流量产生缺口。

验证：一致快照位置与count/hash。

取舍：成本高但恢复必需。

Snapshot 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 8. Delete/Tombstone

机制：传递源删除和版本。

实际用途：下游不复活数据。

失败方式：忽略delete只做upsert。

验证：源下游删除count和版本。

取舍：隐私/历史保留需分别定义。

Delete/Tombstone 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 9. Schema evolution

机制：源列增删/类型变化传播到pipeline。

实际用途：expand-contract和模型版本。

失败方式：自动接受破坏类型导致停流。

验证：兼容检查、DLQ和canary。

取舍：治理增加发布协调。

Schema evolution 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 10. Watermark

机制：表示已处理的源位置/事件时间进度。

实际用途：恢复、延迟和窗口完成。

失败方式：只用墙钟推断完整性。

验证：按partition记录并检测gap。

取舍：提供进度但需持久一致。

Watermark 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 11. Backfill

机制：历史数据按稳定分片重算。

实际用途：修复逻辑和新增指标。

失败方式：与增量双写无水位导致重复漏失。

验证：batch ID、版本、对账、限速。

取舍：可修复但资源和合并复杂。

Backfill 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 12. 方案比较

|方案|主要能力|边界|
|---|---|---|
|OLTP行存|点查事务|大扫描干扰|
|列式OLAP|聚合扫描|逐行更新成本|
|ETL|目标前转换|原始重算能力弱|
|ELT|目标内转换|成本治理复杂|
|CDC|低延迟增量|不是完整备份|

## 13. 完整案例：PostgreSQL到ClickHouse实时订单分析

### 输入与约束

订单主库不可被BI拖慢；分析5分钟内可见；删除和金额版本必须正确。

### 处理步骤

1. PostgreSQL事务写事实并保留版本。
2. Debezium一致快照记录LSN后读取WAL。
3. Kafka按表/订单key传递，schema受控。
4. ClickHouse staging幂等写入并用版本模型生成查询表。
5. 按tenant/day对账count、sum、max version与delete。

### 输出

交易主库承担正确性，ClickHouse承担扫描聚合，freshness可量化。

### 验证

snapshot+stream无gap；乱序最终最高版本；主库BI连接被禁止/限额。

### 失败分支

把更新当普通append会重复金额；必须保留变更语义并按版本构建事实。

### 恢复与重跑

作业以batch/run ID和输入水位幂等重跑，已发布版本不被未验证结果原地覆盖；修复先写隔离版本，对账通过后再切换读指针。

## 14. 完整案例：批量财务ETL

### 输入与约束

每日对象存储CSV、金额需精确、错误行不能静默丢，次日8点报表。

### 处理步骤

1. 保存原始对象checksum和manifest。
2. 解析到隔离staging，金额用decimal/整数分。
3. 质量规则输出accepted/rejected与reason。
4. 单batch事务/分区发布curated表。
5. 对账文件总额、行数、重复业务键并记录lineage。

### 输出

每个报表可追溯到输入对象、代码版本和质量结果。

### 验证

重复运行同batch幂等；坏行有隔离；总额守恒。

### 失败分支

解析失败后跳过行仍显示job成功会少计；质量阈值和发布门禁必须阻止。

### 恢复与重跑

作业以batch/run ID和输入水位幂等重跑，已发布版本不被未验证结果原地覆盖；修复先写隔离版本，对账通过后再切换读指针。

## 15. 失败注入矩阵

|注入|预期信号与恢复|禁止结果|
|---|---|---|
|重复输入|`oltp_p99` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|乱序和迟到|`olap_scan_bytes` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|源schema变化|`cdc_lag` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|任务中途崩溃|`snapshot_progress` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|下游429/503|`schema_errors` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|checkpoint损坏|`delete_lag` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|质量规则失败|`backfill_rate` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|回填与实时并发|`quality_reject` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|权限撤销|`reconciliation_diff` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|成本超预算|`pipeline_cost` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|

## 16. 数据质量与对账

1. CDC 源目标的行数与唯一业务键：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
2. 金额/数量守恒：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
3. not null与范围：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
4. 引用完整性：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
5. 源/目标最大版本：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
6. 删除和tombstone：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
7. freshness/迟到：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
8. 分区完整性：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
9. schema版本：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
10. 抽样hash：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。

## 17. 调试与观测

1. `oltp_p99`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
2. `olap_scan_bytes`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
3. `cdc_lag`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
4. `snapshot_progress`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
5. `schema_errors`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
6. `delete_lag`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
7. `backfill_rate`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
8. `quality_reject`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
9. `reconciliation_diff`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
10. `pipeline_cost`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。

排障从一个可复现业务分片开始，沿源记录、传输offset、处理checkpoint、目标版本和指标SQL逐跳核对；只看job success不能证明数据正确。

## 18. 安全、成本与运维边界

1. 源凭据最小权限；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
2. PII分层访问和脱敏；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
3. 原始层不可被普通BI任意下载；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
4. 重跑/回填有资源配额；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
5. 流批共享sink有容量仲裁；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
6. 删除请求传播到派生层；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
7. schema发布兼容门禁；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
8. checkpoint/manifest备份；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
9. 灾备恢复演练；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。
10. 成本按pipeline/dataset/tenant归集；OLTP、OLAP、ETL、ELT 与 CDC 的实现要提供owner、runbook、停止阈值和审计记录。

## 19. 综合练习与验收

实现“PostgreSQL到ClickHouse实时订单分析”，再以“批量财务ETL”验证另一类时效/治理约束。提交数据样例、模型、质量测试、故障注入、lineage和成本面板。

- [ ] OLTP 的定义、应用、失败和验证均能用真实数据复现。
- [ ] OLAP 的定义、应用、失败和验证均能用真实数据复现。
- [ ] ETL 的定义、应用、失败和验证均能用真实数据复现。
- [ ] ELT 的定义、应用、失败和验证均能用真实数据复现。
- [ ] CDC 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Snapshot 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Delete/Tombstone 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Schema evolution 的定义、应用、失败和验证均能用真实数据复现。
- [ ] 两个案例包含输入、步骤、输出、验证、失败与重跑。
- [ ] 源与目标按业务分片完成count/sum/version/hash对账。
- [ ] 历史发布版本可回退，回填不压垮在线事实系统。

## 来源

- [PostgreSQL 18 logical decoding](https://www.postgresql.org/docs/18/logicaldecoding.html)（访问日期：2026-07-17）
- [Debezium PostgreSQL connector](https://debezium.io/documentation/reference/stable/connectors/postgresql.html)（访问日期：2026-07-17）
- [ClickHouse architecture](https://clickhouse.com/docs/architecture/introduction)（访问日期：2026-07-17）
- [Apache Kafka 4.3 documentation](https://kafka.apache.org/documentation/)（访问日期：2026-07-17）
- [dbt docs](https://docs.getdbt.com/docs/introduction)（访问日期：2026-07-17）
