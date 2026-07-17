---
stage: intermediate
direction: backend-data
topic: retry-dlq-poison-message
---

# Retry、指数退避、Dead Letter Queue 与 Poison Message

重试只适用于暂时失败，并且必须受次数、总时限、容量和幂等约束。DLQ 是隔离与人工/自动修复入口，不是把失败消息移走后不再处理的终点。

## 1. 先分类失败

| 类别 | 例子 | 原样重试 | 动作 |
|---|---|---:|---|
| 暂时依赖失败 | timeout、503、短锁冲突 | 是 | 有界退避 |
| 限流/过载 | 429、consumer 下游饱和 | 是 | Respect Retry-After/降速 |
| 确定性输入错误 | schema 缺字段、非法枚举 | 否 | 隔离/DLQ |
| 业务终态冲突 | 对象已取消、权限撤销 | 通常否 | 记录终态 |
| 代码 bug | panic、nil dereference | 原样会重复 | 停止/发布修复 |
| 未知外部结果 | 请求 timeout 但可能执行 | 先对账 | 不盲重试 |

分类必须基于稳定错误码/类型，不解析错误文案。5xx 也可能是确定性 bug；4xx 中 408/429 可能可重试。

## 2. 重试预算

单条消息的预算包括最大 attempts、elapsed time 和业务 deadline。系统预算还包括全局 retry QPS、并发和队列容量。

```text
max_attempts = 8
max_elapsed = 2 hours
per_attempt_timeout = min(10s, remaining_budget)
```

达到任一上限进入终态/DLQ。无限重试会阻塞 partition、积压并持续攻击故障依赖。

## 3. 指数退避与 jitter

基础指数：`delay = min(cap, base × 2^attempt)`。所有 consumer 同时恢复会同步重试，必须 jitter。

常见 full jitter：

```text
delay = random(0, min(cap, base × 2^attempt))
```

decorrelated jitter 可基于上次延迟随机增长。选型通过负载仿真，不需要不同服务各自发明。`Retry-After` 是下游明确提示时的下限/依据，但仍受业务总 deadline。

## 4. 重试放在哪里

### Consumer 内短重试

适合几十毫秒瞬态网络错误，次数少。期间 partition/visibility 被占用；长 sleep 会降低吞吐和触发 group timeout。

### 延迟/重试 Queue 或 Topic

失败消息发布到按延迟级别的 retry topic，原 consumer 继续处理后续。需保留原 event ID、attempt、first_failed_at、错误类别。跨 topic 会打乱顺序。

### Broker 延迟能力

部分工作队列有延迟/TTL/visibility。语义和最大延迟依产品，不把 RabbitMQ 插件、云队列和 Kafka 当同一 API。

### 调度数据库

对需取消、优先级、查询状态的长任务，用 job 表 `next_attempt_at` 和 worker lease 更清晰。Kafka 不原生提供任意单消息延时队列。

## 5. Kafka Retry Topic

一种结构：`orders.events` → `orders.retry.1m` → `orders.retry.10m` → `orders.dlq`。转发必须在业务提交/失败分类后可靠进行；若提交原 offset 成功但发 retry 失败会丢。

可用 Kafka transaction 原子发送 retry 记录和提交 consumed offset（链路都在 Kafka），或先持久化失败 intent/outbox。消息 headers/payload 保留 original topic/partition/offset 仅用于追踪，业务去重仍用 event ID。

重试 topic 破坏同 key 顺序：后续 v11 正常成功，v10 之后重试到达。consumer 必须使用 aggregate version/幂等，而不能依赖原 partition 顺序。

## 6. Poison Message

Poison message 在当前代码/数据下每次处理都确定失败，例如不支持 schema、损坏压缩、触发 bug。若一直原地重试，会阻塞 partition 或快速烧 CPU。

检测：相同 event ID/offset 多次同类失败；解码错误；非重试 error；超过 budget。隔离时保存足够诊断元数据和原始 payload 的安全引用，不在普通日志打印敏感内容。

一条消息导致进程崩溃时，外层 recovery/隔离必须能识别 offset；但数据反序列化应先限制大小/深度，防止 OOM 在 recovery 前杀进程。

## 7. Dead Letter Queue 契约

DLQ envelope：

```json
{
  "dlq_id":"01JZ...",
  "original":{"topic":"orders.v1.events","partition":7,"offset":421},
  "event_id":"evt_9",
  "schema_version":3,
  "attempts":8,
  "first_failed_at":"2026-07-17T08:00:00Z",
  "last_failed_at":"2026-07-17T09:30:00Z",
  "error":{"category":"unsupported_schema","code":"SCHEMA_99"},
  "payload_ref":"s3://restricted-dlq/sha256/..."
}
```

DLQ 有 owner、告警、保留、访问控制、加密、修复和 re-drive 流程。payload 可能含 PII；默认只保存加密对象引用/checksum，权限与审计更严。

## 8. DLQ 不是成功

进入 DLQ 表示主消费路径绕过了该业务结果。业务看板必须反映缺失投影/任务失败；不能只因 consumer lag 归零就显示健康。

指标：DLQ rate、oldest age、按 error category、未处置数、re-drive success、再次 DLQ。低频关键支付 DLQ 比高频日志 DLQ 更严重，告警按业务价值。

## 9. Re-drive

修复 schema/代码/数据后重放：

1. 选择明确问题类别和时间范围。
2. 在 staging/影子 consumer 验证修复。
3. 为 re-drive 建批次 ID、速率和并发上限。
4. 保持原 event ID，新增 re-drive ID/attempt。
5. 消费者幂等；已被人工修复的消息不会重复副作用。
6. 观察下游、失败率，支持暂停。
7. 成功后记录审计，不立即删除原 DLQ 证据。

整队一键重放会把旧 bug/容量峰值重新注入生产。

## 10. 顺序与阻塞策略

当失败消息 v10 阻塞同订单 v11：

- 严格事件序列：暂停该 key/partition，修复 v10，代价是共享 partition 其他 key 受阻。
- 每 aggregate parking lot：把该 aggregate 后续消息暂存，其他 key 继续；实现复杂。
- 状态快照投影：v11 完整状态可覆盖，v10 隔离并忽略旧版本。
- 增量 ledger：不能跳过，必须补齐或从事实重建。

策略来自事件语义，不是统一“跳过 poison”。

## 11. 优雅退避与消费背压

下游整体 503 时，继续 poll→retry topic 会把同样流量复制并增加 broker 写。consumer 应降低/暂停 partitions、减少并发、打开 circuit breaker，并保持 group 心跳的正确客户端方式。

恢复时缓慢增加并发，避免 retry backlog + 新流量同时冲击。为 retry 和新消息分配公平/优先容量。

## 12. 未知结果

支付请求 timeout 可能已扣款，标为 `unknown`：

1. 以稳定 idempotency key 查询供应商。
2. 使用 webhook/对账文件确认。
3. 在确认未执行后才重试。
4. 超过业务时间进入人工调查，不把 unknown 当 failed。

消息重试框架必须允许 unknown 状态，不能只有 success/fail。

## 13. Retry Storm

故障恢复时百万消息同时到期产生重试风暴。防护：full jitter、全局 token bucket、每 tenant 配额、下游 circuit breaker、retry schedule 分桶、恢复斜坡、最大积压容量。

多个微服务逐层各重试 3 次会放大为 3^N。选择一个拥有业务上下文的层负责重试，下游客户端只做极少连接级重试。

## 14. 应用案例一：搜索索引暂时失败

### 输入

OpenSearch 维护期间 30 分钟 503；订单事件 5000/s；搜索允许 2 小时 lag，事件是完整状态快照。

### 处理

1. 3 次秒级 jitter 后判依赖不可用，打开 circuit breaker。
2. pause 消费或转持久 retry topic，不让应用内存堆积。
3. retry 记录保留 event ID/version；v11 后到可覆盖 v10。
4. 恢复后按 20%→50%→100% 速率追赶，监控集群 reject/latency。
5. 超过 2 小时进入业务降级和告警，不悄悄清 offset。

### 验证

恢复期间新+旧写入不超 OpenSearch bulk capacity；最终文档最高 version；consumer lag、retry age 和搜索 freshness 同时可见。

### 失败注入

503 持续 3 小时，系统不无限 retry CPU、不 OOM；搜索 API显示数据更新时间/降级。若直接 DLQ 所有暂时错误，恢复要人工重放且失去自动追赶。

## 15. 应用案例二：不支持的 schema

### 输入

consumer 只支持 v1/v2，却收到 v99；同 partition 后有其他订单。

### 处理

1. 解码 envelope 后发现版本不支持，确定性错误不重试。
2. 把 payload 加密引用和 schema 信息写 DLQ，可靠提交后推进 offset。
3. 告警 producer/consumer compatibility；按业务决定 partition 是否允许继续。
4. 发布 v99 支持后，小批 re-drive，保持 event ID。

### 验证与失败分支

v99 不形成热循环；DLQ 写失败时不能提交原 offset，否则永久丢。re-drive 重复不会重复业务结果。

## 16. 应用案例三：邮件供应商限流

### 输入

供应商返回 429 + Retry-After 120 秒；每 tenant 每小时配额不同；邮件有发送截止时间。

### 处理

1. 解析可信 Retry-After，记录 next_attempt_at，不 sleep 持有 consumer lease。
2. 调度表按 tenant token bucket 和优先级领取。
3. deadline 前指数+jitter重试；过期后标 failed-expired，进入业务补偿而非 DLQ 无尽等待。
4. 供应商 timeout 标 unknown 先对账。

### 输出

限流不会让其他 tenant 饿死；高优先事务邮件保留容量；过期状态可查询。

## 17. 方案比较

| 机制 | 适合 | 风险 |
|---|---|---|
| Consumer 内重试 | 极短瞬态 | 阻塞 poll/lease |
| Retry topic | 大吞吐分级延迟 | 乱序、topic 数、转发原子性 |
| Queue visibility | 工作任务 | 旧 worker 并发、需续租 |
| Job 调度表 | 取消/优先级/任意时间 | DB 扫描与锁设计 |
| DLQ | 确定失败隔离 | 变成无人处理墓地 |

## 18. 调试和生产指标

按 error category 记录 attempt、next delay、first/last failure、retry age；观察 retry throughput、scheduled count、DLQ oldest、re-drive batch、circuit state、下游 429/5xx、unknown outcome。

日志用 event ID/dlq ID 关联，不打印 payload/secret。trace 对每次 attempt 建 link，不把几小时重试做一个不断增长 span。

### 容量与保留

重试容量按峰值失败率而非平均估算。若入口 2 万条/s、依赖最长故障 2 小时，理论积压 1.44 亿条；还要计 payload、索引、副本和 re-drive 双流量。容量不足时应在入口背压/暂停非关键生产者，而不是让 broker 磁盘耗尽。

DLQ 保留必须覆盖发现、修复、审批和重放周期。原 topic 先到期时，DLQ envelope 仍需足够的 schema 和安全 payload 引用才能重建。加密对象引用的生命周期不能短于 DLQ；密钥轮换要保证旧消息可解。

告警分两级：单条关键资金 poison 立即业务告警；大量同 error code 表示发布/依赖事故。用错误类别聚合，未知类别必须显式出现，不能落入无监控的 `other`。

### Re-drive 回滚

re-drive 发布错误时要能按 batch ID 停止，并识别已成功处理的 event。修复代码若再次失败，消息进入新的隔离状态但保留原始失败链；不能覆盖第一次证据。对产生外部副作用的重放先在 dry-run/影子模式计算将要改变的对象集合，再审批执行。

## 19. 综合练习与验收

实现含短重试、延迟调度、DLQ 和 re-drive 的 consumer。注入 503、429、schema v99、poison panic、DB commit 后崩溃和供应商 unknown。

验收：确定错误不原样重试；暂时错误有 jitter/总预算；重试不阻塞整个消费；DLQ 写与原 offset 推进不丢；re-drive 限速且幂等；顺序/gap 策略按事件类型区分；retry storm 不超过下游预算；所有 DLQ 有 owner 和 oldest-age 告警。

## 来源

- [Apache Kafka consumer configuration](https://kafka.apache.org/documentation/#consumerconfigs)（访问日期：2026-07-17）
- [Apache Kafka transactions](https://kafka.apache.org/documentation/#transactions)（访问日期：2026-07-17）
- [AWS Builders Library: Timeouts, retries and backoff with jitter](https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/)（访问日期：2026-07-17）
- [Redis Streams](https://redis.io/docs/latest/develop/data-types/streams/)（访问日期：2026-07-17）
- [RFC 9110: Retry-After](https://www.rfc-editor.org/rfc/rfc9110.html#name-retry-after)（访问日期：2026-07-17）
