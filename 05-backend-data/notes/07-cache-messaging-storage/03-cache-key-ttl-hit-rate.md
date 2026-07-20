---
stage: intermediate
direction: backend-data
topic: cache-key-ttl-hit-rate
---

# 缓存 Key、TTL、命中率、随机过期与容量治理

缓存可维护性取决于 key 是否稳定表达数据边界、TTL 是否匹配陈旧预算、命中率是否按真实请求权重解释，以及内存/淘汰是否在故障前可观测。只设置一个统一 TTL 无法形成可靠策略。

## 1. Key 是缓存契约

key 同时决定身份、隔离、路由、失效和观测。建议结构：

```text
<domain>:<schema-version>:<tenant>:<entity>:<id>:<variant>
catalog:v3:t_7:product:p_9:locale:zh-CN
```

逐段含义：

- domain 避免不同业务碰撞，并用于 ACL/监控分组。
- schema version 隔离不兼容序列化格式。
- tenant 从可信认证上下文取得，防跨租户复用。
- entity/id 指向事实对象，不把可变标题当身份。
- variant 包含 locale、货币、实验组等真正改变表示的输入。

key 不应含密码、token、邮箱等敏感数据；SLOWLOG、MONITOR、备份和运维工具可能看到它。必要时对标识做稳定不可逆映射，但仍要支持按事件定位失效。

## 2. Redis Cluster hash slot

Cluster 对 key 计算 CRC16 后映射到 16384 slots。多 key 原子命令通常要求同 slot。hash tag 使用第一个非空 `{...}` 内容：

```text
cart:{t7:u9}:items
cart:{t7:u9}:summary
```

只给必须一起操作的同一购物车加 tag。把 `{tenant}` 放所有 key 会让大租户集中一个 slot；把常量 `{cache}` 放全部 key 会彻底破坏分片。迁移 key schema 要考虑新旧 slot 和双读成本。

## 3. TTL 的三种语义

### 新鲜期

在该时间内允许直接返回，不主动验证事实源。

### 最大陈旧期

依赖失败时允许返回旧值的 hard deadline。不是所有数据都允许 stale。

### 保留期

对象在 Redis 中最多保留多久，用于内存和隐私治理；它可能长于新鲜期。

单个 Redis TTL 只能表达“key 何时消失”，不能同时表达 fresh/hard 两个边界。需要在 value 保存 `refreshed_at`/version，或用两个 key/元数据。

## 4. 绝对过期与滑动过期

绝对 TTL 在写入时确定，到期后消失，适合验证码、价格快照。滑动 TTL 每次访问续期，适合闲置会话；若每次 GET 都续期，会让高频 key 永不淘汰并增加写/复制压力。

会话同时需要 idle timeout 和 absolute lifetime：活跃不能让被盗 session 永久有效。Redis TTL 可管理 idle，值中保存 absolute expiry，由服务端每次验证。

## 5. Redis 过期行为

`EXPIRE`/`PEXPIRE` 设置相对时间，`EXPIREAT`/`PEXPIREAT` 设置绝对 Unix 时间，`TTL`/`PTTL` 返回剩余时间。返回值要区分 key 不存在和无过期。

```text
SET catalog:v3:t7:p9 '{...}' EX 300
EXPIRE catalog:v3:t7:p9 600 GT
PTTL catalog:v3:t7:p9
PERSIST catalog:v3:t7:p9
```

NX/XX/GT/LT 条件能限制何时修改 TTL。具体组合按命令文档和客户端版本验证。

过期由两条路径完成：访问到已过期 key 时惰性删除，以及后台主动采样过期。到达 TTL 不保证内存物理上在同一毫秒释放；逻辑读取会把已过期 key 当不存在。大量同刻过期仍会增加 CPU 和回源。

## 6. 随机过期 jitter

若基础 TTL 为 `T`，可以在安全陈旧预算内加入随机量：

```text
ttl = T × random(0.8, 1.2)
```

jitter 目标是打散同批 key 的过期时间，不是掩盖没有回源保护。范围必须受业务最大陈旧期约束；验证码、签名等严格到期数据不能随意延长。

随机源无需密码学强度，但同一批实例不能固定相同 seed 产生相同序列。可使用稳定 hash(key) 映射抖动，保证重建时分布稳定。

## 7. TTL 选择方法

输入：允许陈旧时间、更新频率、读取频率、回源成本、内存预算、失效事件可靠性。

```text
候选 TTL ≤ 业务最大陈旧期
候选 TTL ≈ 能把回源 QPS 控制在事实源余量内的值
候选 TTL × 活跃 key 写入速率 ≈ 同时驻留 key 数量
```

低频大对象可能直到 TTL 都未再次访问，缓存收益小；高频极热对象适合 refresh-ahead。更新频率高于读取频率的数据不适合缓存，维护成本超过命中收益。

## 8. 命中率的精确定义

请求命中率：

```text
request_hit_rate = hits / (hits + misses)
```

字节命中率：命中的响应字节 / 总请求响应字节。成本加权命中率：避免的事实源成本 / 无缓存时总成本。不同指标回答不同问题。

Redis `keyspace_hits`/`keyspace_misses` 是实例级命令统计，不一定等于业务缓存请求：EXISTS、脚本、不同 DB/业务会混入。应用按 cache name、结果状态和规范化路由记录指标，但 label 不用原始 key。

### 为什么高命中率仍可能失败

- 99% hit 的 1% miss 集中在一个高成本查询。
- L1 hit 掩盖 L2 大量 miss。
- negative cache hit 被当成功，却掩盖过度查询不存在 ID。
- stale hit 被计为 fresh hit。
- 热 key 命中高，但单节点网络/CPU 饱和。

至少拆分 fresh_hit、stale_hit、negative_hit、cold_miss、expired_miss、evicted_miss、error。

## 9. Memory 与 eviction

`maxmemory` 控制用于数据集的内存阈值，复制/AOF buffers 的部分内存不计入 eviction 比较；为系统、fork copy-on-write、buffers 和 allocator 留余量。

常见策略：

| 策略 | 候选 | 适用 | 风险 |
|---|---|---|---|
| `noeviction` | 不淘汰 | 不能丢的 Redis 数据 | 新写报错，需严格容量 |
| `allkeys-lru` | 全 key 近似 LRU | 普通缓存 | 不等于精确 LRU |
| `allkeys-lfu` | 全 key 近似频率 | 稳定热点 | 新热点升温速度需测 |
| `allkeys-random` | 全 key 随机 | 访问近均匀 | 可能删热 key |
| `volatile-ttl` | 有 TTL 且最短 | 应用能给价值型 TTL | 无 TTL key 不淘汰 |
| `volatile-*` | 仅带 TTL | 混合数据 | 无候选时表现近 noeviction |

缓存和不可淘汰 Stream/session 混在同实例会让策略难以满足两类目标，优先拆部署。

## 10. 过期与淘汰不是同一原因

- expired：业务 TTL 到期，通常预期行为。
- evicted：内存压力下按 policy 删除，可能提前于 TTL。

若 miss 增加，分别看 `expired_keys` 与 `evicted_keys`。eviction 持续说明 working set 超容量、策略不匹配或大 key；盲目加长 TTL 会更糟。expired 高可能是正常更新，也可能 TTL 过短导致 thrash。

## 11. 容量估算

抽样 `MEMORY USAGE key SAMPLES n`，按 key 类型/大小分桶。粗略：

```text
dataset ≈ Σ(active_keys_per_class × p95_bytes_per_key)
required_RAM ≈ dataset / target_utilization + buffers + COW_headroom
```

用 p95 而非平均避免长尾低估。RDB/AOF rewrite 和主从同步可能出现额外内存/网络峰值。容量测试要包含真实 key 长度、序列化、TTL 元数据和 allocator。

## 12. 大 Key 与热 Key

大 key 指值/成员多导致单命令、复制、迁移或删除昂贵；热 key 指访问集中导致单分片 CPU/网络热点。它们可能是不同 key。

治理：

- 大 String 拆按访问边界，不按固定字节随意切。
- 大 Hash/Set/ZSet 按实体或时间桶拆分，维护跨桶查询上限。
- 热只读可 L1、客户端缓存或副本读（接受一致性限制）。
- 热写计数可分片计数再聚合，但精确实时值成本上升。
- 删除大 key 用 UNLINK，并在低峰限速；仍监控后台释放。

使用 `SCAN`/`--bigkeys` 生产巡检要限速；`KEYS *` 会阻塞扫描整个 keyspace。

## 13. Keyspace notification 边界

Redis 可发布 key 事件通知，但需要配置，使用 Pub/Sub，消费者断线会漏事件，且通知本身增加开销。它适合触发辅助动作，不适合作为唯一可靠业务事件源。可靠失效使用数据库 outbox/持久 Stream/Kafka，再配 TTL 收敛。

## 14. 应用案例一：日报表缓存

### 输入

每个租户的日报表查询 3 秒；报表在当天持续变化，历史日冻结；工作日 09:00 大量用户同时打开。允许当天数据陈旧 60 秒，历史数据 24 小时。

### 处理

1. key：`report:v2:{tenant}:daily:2026-07-17:tz:Asia-Shanghai`。
2. 当天 fresh TTL 45–60 秒，hard stale 120 秒；历史 TTL 24 小时 ±20%。
3. 08:55 后台只预热最近活跃租户，并设全局并发/数据库 QPS 上限。
4. miss 用跨实例租约 + 数据库容量闸门；未持租约者短暂返回 stale 或等待。
5. 报表修正事件按 tenant/date 删除；TTL 作为漏事件后的上界。

### 输出与验证

按 tenant tier 记录 hit/miss/fill；09:00 数据库 QPS 不超过预算；历史 key 的 byte hit 高。随机暂停失效 consumer，120 秒内当天数据仍收敛。

### 失败注入

清空 Redis 后受控预热，不让所有用户并发生成报表。若租约持有者崩溃，短 TTL 释放；数据库闸门继续保护事实源。

## 15. 应用案例二：一次性下载 token

### 输入

token 必须 10 分钟到期且只能消费一次；Redis failover 可能丢最近写；下载权限敏感。

### 方案与取舍

仅 Redis `GETDEL` 提供单节点命令原子，但异步复制 failover 可能让已消费 token 在旧副本状态中重新出现，或刚创建 token 丢失。若业务不能接受重放，应把消费记录/nonce 唯一约束放 PostgreSQL，Redis 只做加速拒绝。

1. 数据库保存 token hash、expires_at、consumed_at。
2. Redis 缓存“可能有效”10 分钟，key 不含原 token。
3. 消费事务执行 `UPDATE ... WHERE consumed_at IS NULL AND expires_at > now()`。
4. 成功后删除缓存；失败返回统一错误防枚举。

### 验证与失败注入

并发消费 100 次只有一个数据库 UPDATE 成功。Redis 在消费前后 failover，不会让第二次下载通过。严格时间由数据库时钟判断，随机 TTL 不得延长安全到期。

## 16. 应用案例三：Schema 版本切换

### 输入

商品缓存 v2 使用 `price` 浮点元，v3 使用 `price_cents` 整数；滚动发布 30 分钟。

### 处理

1. 新代码读写 `catalog:v3...`，miss 时查数据库；不覆盖 v2。
2. 旧代码继续使用 v2；数据库是共同事实源。
3. 控制 v3 预热速率，观察 memory 双份峰值。
4. 旧实例归零后让 v2 按 TTL 淘汰，不执行 KEYS 全删。
5. 回滚只切回仍被维护的 v2；若新写语义不能无损转换，先禁用新功能。

### 验证

两个版本并存不发生反序列化错误；内存峰值在容量内；回滚窗口结束前 v2 更新链路仍工作。

## 17. 调试与告警

指标至少包括：

- fresh/stale/negative hit rate；
- miss 原因和 fill p95/p99；
- key 数、dataset bytes、RSS、fragmentation；
- expired/s、evicted/s、rejected writes；
- hot key command QPS、网络字节；
- DB fallback QPS 与限流拒绝；
- schema 版本 key 比例；
- TTL 分布而非单个平均。

告警要基于持续时间和业务影响。短时 eviction 在批量加载期间可能正常，持续 eviction + DB QPS/p99 上升才是高风险链路。

## 18. 生产检查

1. 每类 key 有 owner、schema、TTL、事实源和删除策略。
2. 过期 jitter 不突破安全/业务截止时间。
3. eviction policy 与实例用途一致。
4. key 不含 secret/PII，ACL 限制业务前缀。
5. 冷启动/清库/节点故障做过容量演练。
6. SCAN、大 key 删除和批量写有速率限制。
7. 缓存 miss 不能绕过授权或业务校验。
8. 删除请求覆盖所有版本 key 与 CDN。

## 19. 综合练习与验收

为商品、日报表、会话和一次性 token 分别设计 key/TTL/eviction。输出容量模型、TTL 分布、Cluster slot 设计和故障矩阵。

验收标准：报表缓存清空时数据库不超预算；会话同时有 idle/absolute expiry；一次性 token 的最终单次消费由数据库保证；v2→v3 滚动发布可回退；指标能把 expired、evicted、error miss 分开；大 key 和热 key 各有独立治理方案。

## 来源

- [Redis EXPIRE](https://redis.io/docs/latest/commands/expire/)（访问日期：2026-07-17）
- [Redis key eviction](https://redis.io/docs/latest/develop/reference/eviction/)（访问日期：2026-07-17）
- [Redis Cluster specification](https://redis.io/docs/latest/operate/oss_and_stack/reference/cluster-spec/)（访问日期：2026-07-17）
- [Redis memory optimization](https://redis.io/docs/latest/operate/oss_and_stack/management/optimization/memory-optimization/)（访问日期：2026-07-17）
- [Redis keyspace notifications](https://redis.io/docs/latest/develop/pubsub/keyspace-notifications/)（访问日期：2026-07-17）
