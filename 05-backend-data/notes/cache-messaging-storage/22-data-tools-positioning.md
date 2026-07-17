---
stage: intermediate
direction: backend-data
topic: data-tools-positioning
---

# ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位

数据工具应按责任选择：Kafka负责持久事件日志，Flink/Spark处理数据，ClickHouse服务分析查询，dbt管理SQL变换，Airflow编排有界任务。产品能力会扩展，但责任边界不能因为“也能做”而消失。

## 1. 数据流与决策边界

```mermaid
flowchart LR
 A["源事实/事件"] --> B["传输与处理"]
 B --> C["模型/存储"]
 C --> D["查询与指标"]
 D --> E["质量、对账与重建"]
```

每个输出必须能追溯输入水位、代码或模型版本、质量结果和权限范围；低延迟不能替代可恢复性。

## 2. ClickHouse

机制：列式OLAP、MergeTree、向量化和数据跳过。

实际用途：实时分析、日志和大聚合。

失败方式：承担账户事务/逐行强约束。

验证：scan、merge、parts和query p99。

取舍：查询快但写模型和运维不同。

ClickHouse 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 3. dbt

机制：编译SQL模型、依赖图、测试、文档和增量构建。

实际用途：仓库内ELT和语义治理。

失败方式：充当通用流处理/传文件。

验证：manifest、run results、tests。

取舍：SQL工程化但依赖目标仓库。

dbt 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 4. Airflow

机制：DAG调度有界任务、重试和依赖。

实际用途：批处理、回填、跨系统编排。

失败方式：每条实时事件或长期业务事务。

验证：scheduler lag、task SLA和idempotence。

取舍：可观察编排但不是数据传输层。

Airflow 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 5. Kafka

机制：分区持久日志、consumer groups和保留。

实际用途：事件骨干、CDC缓冲、重放。

失败方式：任意优先级任务/长期对象存储。

验证：partition lag、ISR、disk。

取舍：高吞吐但schema与consumer治理。

Kafka 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 6. Flink

机制：有状态流处理、event time、watermark、checkpoint。

实际用途：低延迟窗口、join、实时规则。

失败方式：简单日报也强行常驻集群。

验证：checkpoint、backpressure、late events。

取舍：强流语义但运维/状态复杂。

Flink 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 7. Spark

机制：批处理与Structured Streaming统一计算API。

实际用途：大规模ETL、ML、微批/流。

失败方式：低延迟单行事务服务。

验证：stage/shuffle/state/stream progress。

取舍：生态广但启动/资源成本。

Spark 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 8. Connector

机制：在系统间搬运并维护offset/schema。

实际用途：标准source/sink。

失败方式：把业务补偿隐藏在connector。

验证：offset、DLQ、version和exact语义。

取舍：减少代码但专属保证需核验。

Connector 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 9. Orchestration

机制：控制任务何时/依赖/重试。

实际用途：Airflow调dbt/Spark/质量检查。

失败方式：作为高频事件数据面。

验证：DAG run、SLA、backfill。

取舍：控制面清楚但跨系统原子不存在。

Orchestration 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 10. Transformation

机制：把原始数据转换为可用模型。

实际用途：dbt SQL、Flink/Spark代码。

失败方式：无版本无测试脚本。

验证：code/version/input/output lineage。

取舍：灵活但质量需门禁。

Transformation 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 11. Serving

机制：为用户/BI提供受控查询。

实际用途：ClickHouse/warehouse/semantic layer。

失败方式：直接暴露raw lake/Kafka。

验证：权限、并发、成本、freshness。

取舍：优化消费但复制数据。

Serving 的生产契约还要定义输入schema、业务/事件时间、幂等键、水位、迟到/删除和重跑行为；成功状态必须对应已发布且通过质量门禁的数据。

## 12. 方案比较

|方案|主要能力|边界|
|---|---|---|
|Kafka|传输/保留|不做任意BI|
|Flink|低延迟有状态流|长期状态运维|
|Spark|大批/统一计算|交互延迟较高|
|dbt|仓库SQL模型|不搬运实时事件|
|Airflow|编排|不是队列|
|ClickHouse|OLAP serving|不是OLTP事实|

## 13. 完整案例：实时产品指标平台

### 输入与约束

事件10万/s，5分钟看板，保留一年，每日财务校正。

### 处理步骤

1. Kafka作为短期事件日志与回放缓冲。
2. Flink按event time去重/窗口，checkpoint到可靠存储。
3. 对象湖保存原始Parquet供重放。
4. ClickHouse接实时聚合和明细服务看板。
5. Airflow每天编排dbt/Spark批次对账并发布修正。

### 输出

流路径低延迟，批路径校正，事实和水位可追踪。

### 验证

停止Flink从checkpoint恢复；重放不重复；实时与日批差异受阈值控制。

### 失败分支

只在Flink保留一年状态会放大恢复和成本；原始历史进入对象存储。

### 恢复与重跑

作业以batch/run ID和输入水位幂等重跑，已发布版本不被未验证结果原地覆盖；修复先写隔离版本，对账通过后再切换读指针。

## 14. 完整案例：月度财务关账

### 输入与约束

多源文件/数据库，必须可重跑、审批和锁定版本，不要求秒级。

### 处理步骤

1. Airflow以业务月/batch ID编排抽取与质量门禁。
2. Spark处理大文件标准化。
3. dbt构建仓库事实维度和测试。
4. ClickHouse可作交互分析副本但财务发布表有明确权威。
5. 审批后发布immutable close version并记录lineage。

### 输出

每月版本可追溯、重跑不覆盖已审批结果。

### 验证

同batch重跑幂等；输入checksum/代码版本/测试结果完整。

### 失败分支

使用实时流结果直接关账且无截止/修正会让数字持续变化；关账需有界版本。

### 恢复与重跑

作业以batch/run ID和输入水位幂等重跑，已发布版本不被未验证结果原地覆盖；修复先写隔离版本，对账通过后再切换读指针。

## 15. 失败注入矩阵

|注入|预期信号与恢复|禁止结果|
|---|---|---|
|重复输入|`kafka_lag` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|乱序和迟到|`flink_checkpoint` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|源schema变化|`flink_backpressure` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|任务中途崩溃|`spark_shuffle` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|下游429/503|`airflow_scheduler_lag` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|checkpoint损坏|`dbt_test_fail` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|质量规则失败|`clickhouse_merge` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|回填与实时并发|`serving_p95` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|权限撤销|`lineage_coverage` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|
|成本超预算|`cost_per_pipeline` 变化可解释，按水位/版本恢复|静默丢行、重复计量、越权或覆盖已发布版本|

## 16. 数据质量与对账

1. 跨工具输入输出的行数与唯一业务键：定义通过阈值、严重级别、quarantine和是否阻断发布；规则本身进入版本控制。
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

1. `kafka_lag`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
2. `flink_checkpoint`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
3. `flink_backpressure`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
4. `spark_shuffle`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
5. `airflow_scheduler_lag`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
6. `dbt_test_fail`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
7. `clickhouse_merge`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
8. `serving_p95`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
9. `lineage_coverage`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。
10. `cost_per_pipeline`：明确单位、采样点、聚合窗口和低基数维度，并与run ID、水位和代码版本关联。

排障从一个可复现业务分片开始，沿源记录、传输offset、处理checkpoint、目标版本和指标SQL逐跳核对；只看job success不能证明数据正确。

## 18. 安全、成本与运维边界

1. 源凭据最小权限；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
2. PII分层访问和脱敏；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
3. 原始层不可被普通BI任意下载；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
4. 重跑/回填有资源配额；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
5. 流批共享sink有容量仲裁；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
6. 删除请求传播到派生层；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
7. schema发布兼容门禁；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
8. checkpoint/manifest备份；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
9. 灾备恢复演练；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。
10. 成本按pipeline/dataset/tenant归集；ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位 的实现要提供owner、runbook、停止阈值和审计记录。

## 19. 综合练习与验收

实现“实时产品指标平台”，再以“月度财务关账”验证另一类时效/治理约束。提交数据样例、模型、质量测试、故障注入、lineage和成本面板。

- [ ] ClickHouse 的定义、应用、失败和验证均能用真实数据复现。
- [ ] dbt 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Airflow 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Kafka 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Flink 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Spark 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Connector 的定义、应用、失败和验证均能用真实数据复现。
- [ ] Orchestration 的定义、应用、失败和验证均能用真实数据复现。
- [ ] 两个案例包含输入、步骤、输出、验证、失败与重跑。
- [ ] 源与目标按业务分片完成count/sum/version/hash对账。
- [ ] 历史发布版本可回退，回填不压垮在线事实系统。

## 来源

- [ClickHouse docs](https://clickhouse.com/docs/)（访问日期：2026-07-17）
- [dbt docs](https://docs.getdbt.com/docs/introduction)（访问日期：2026-07-17）
- [Apache Airflow stable docs](https://airflow.apache.org/docs/apache-airflow/stable/)（访问日期：2026-07-17）
- [Apache Flink stable docs](https://nightlies.apache.org/flink/flink-docs-stable/)（访问日期：2026-07-17）
- [Apache Spark docs](https://spark.apache.org/docs/latest/)（访问日期：2026-07-17）
