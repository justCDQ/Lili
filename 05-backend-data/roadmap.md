# 后端与数据系统 Roadmap

## 一、数据结构与算法

- [ ] 复杂度、数组、链表、栈、队列、双端队列。
- [ ] 哈希表、冲突、扩容、负载因子。
- [ ] 二叉树、搜索树、堆、Trie、B+Tree 基础。
- [ ] 图、DFS、BFS、拓扑排序、最短路径、并查集。
- [ ] 二分、快速、归并、稳定排序和外部排序基础。
- [ ] 双指针、滑动窗口、分治、贪心、回溯和动态规划。

必做：LRU Cache、任务依赖调度器、Token Bucket 限流器、Top-K、简化 HashMap。

验收：能根据数据规模、操作频率和内存约束选择数据结构，不需要追求算法竞赛水平。

---

## 二、Go

### 语言

- [ ] Slice、Map、Struct、Interface、Pointer、Error、Package、Module、Generic。
- [ ] Error Wrapping、业务错误、系统错误、可重试错误。

### 并发

- [ ] Goroutine、Channel、Select、Context。
- [ ] Mutex、RWMutex、WaitGroup、Atomic。
- [ ] Worker Pool、并发限制、优雅停止。
- [ ] Race Condition、Deadlock、Goroutine Leak。

### 工具

- [ ] Test、Table-driven Test、Benchmark。
- [ ] Race Detector、pprof、trace、vet、Staticcheck、Delve。

必做：并发下载器、Worker Pool、优雅退出 HTTP 服务、CPU/Memory Profile 报告。

验收：能够编写可靠 Go 服务，正确处理取消、超时、竞争、错误和性能问题。

---

## 三、Linux、操作系统和网络

### Linux

- [ ] 文件、权限、用户、进程、Signal、Pipe、Socket、Systemd。
- [ ] ps、top、htop、lsof、ss、curl、grep、awk、sed。
- [ ] CPU、内存、磁盘、网络和日志排查。

### 操作系统

- [ ] 进程、线程、用户态、内核态、系统调用。
- [ ] 虚拟内存、Page、Stack、Heap、上下文切换。
- [ ] File Descriptor、阻塞/非阻塞、I/O Multiplexing、epoll。
- [ ] Lock、Deadlock、OOM、CPU 密集和 I/O 密集。

### 网络

- [ ] DNS、TCP 握手、重传、滑动窗口、拥塞控制。
- [ ] TLS、证书和 HTTPS。
- [ ] HTTP/1.1、HTTP/2、HTTP/3、Keepalive、Connection Pool。
- [ ] TIME_WAIT、反向代理、L4/L7 负载均衡。

验收：服务不可访问时，能沿 DNS、网络、端口、进程、日志、依赖逐层排查。

---

## 四、API 与服务

- [ ] REST、RPC、gRPC、GraphQL、WebSocket、SSE、Webhook 的适用场景。
- [ ] Resource、Method、Status Code、分页、筛选、排序、批量和版本。
- [ ] 统一错误：输入、认证、权限、不存在、冲突、限流、系统和外部依赖。
- [ ] Request ID、Logging、Recovery、Timeout、CORS、Rate Limit、Metrics、Trace。
- [ ] Session、Cookie、JWT、OAuth2、OIDC、API Key。
- [ ] RBAC、ABAC、资源权限、多租户和审计。
- [ ] 幂等 Key、异步任务 API 和 OpenAPI。

验收：API 不直接暴露数据库结构，错误可定位，权限由服务端控制，并具备测试、文档和监控。

---

## 五、关系型数据库

### 建模

- [ ] 实体、关系、主键、外键、唯一约束。
- [ ] 规范化、反规范化、审计字段、状态和历史。
- [ ] 软删除的风险和多租户设计。

### SQL

- [ ] Join、Group By、Having、Subquery、CTE、Window Function。
- [ ] 分页、游标分页、复杂筛选和 N+1。

### 索引

- [ ] B+Tree、复合索引、最左前缀、覆盖索引。
- [ ] Selectivity、Index Scan、Full Table Scan、Write Amplification。
- [ ] Explain、Explain Analyze 和慢查询。

### 事务

- [ ] ACID、隔离级别、MVCC、锁、死锁。
- [ ] 乐观锁、悲观锁、唯一约束冲突和库存并发。
- [ ] 避免长事务和事务中调用慢外部接口。

### Migration

- [ ] Schema、Data Migration、Expand and Contract。
- [ ] 向后兼容、零停机、大表回填和回滚。
- [ ] 备份、恢复、RPO 和 RTO。

验收：能设计中等复杂度数据库、选择索引、分析慢查询、处理并发并完成安全迁移。

---

## 六、Redis 与缓存

- [ ] String、Hash、List、Set、Sorted Set、Bitmap、HyperLogLog、Stream、Geo。
- [ ] Cache Aside、Read Through、Write Through、多级缓存。
- [ ] Key、TTL、命中率、随机过期。
- [ ] 穿透、击穿、雪崩和热点 Key。
- [ ] 缓存更新、删除、最终一致性和事件驱动失效。
- [ ] 分布式锁：唯一值、Lua 释放、过期、续期、Fencing Token。

原则：数据库是事实来源，缓存是派生数据。

验收：能说明什么时候不应该使用缓存，Redis 故障时系统可以降级。

---

## 七、消息队列与异步任务

- [ ] Producer、Consumer、Topic、Queue、Partition、Offset、Ack、Consumer Group。
- [ ] At-most-once、At-least-once、重复、乱序、幂等 Consumer。
- [ ] Retry、指数退避、Dead Letter Queue、Poison Message。
- [ ] Schema 版本、消息积压、优先级和取消。
- [ ] Transactional Outbox、最终一致性和补偿。

验收：能处理重复消息、幂等、失败、死信、消息积压以及数据库与事件的一致性。

---

## 八、搜索、对象存储与数据分析

### 搜索

- [ ] 倒排索引、Analyzer、Tokenizer、Mapping。
- [ ] 全文、过滤、高亮、聚合、自动补全、Search After。
- [ ] 数据库同步、更新和删除。
- [ ] 搜索引擎不作为主数据库。

### 对象存储

- [ ] S3、Bucket、Object、Key、Presigned URL。
- [ ] Multipart Upload、断点续传、版本、生命周期、CDN、ETag。
- [ ] 类型和内容校验、权限、病毒扫描和审计。

### 数据工程

- [ ] OLTP、OLAP、ETL、ELT、CDC。
- [ ] Warehouse、Lake、Batch、Streaming、Lineage、Data Quality。
- [ ] 事实表、维度表、星型模型。
- [ ] ClickHouse、dbt、Airflow、Kafka、Flink/Spark 的定位。
- [ ] 漏斗、留存、错误率、性能和 AI 成本分析。

---

## 九、分布式系统

演进顺序：单体、单机优化、垂直扩容、水平扩容、负载均衡、缓存、复制、消息、服务拆分、多区域。

### 复制与分片

- [ ] Leader-Follower、Multi-Leader、Leaderless。
- [ ] 同步/异步复制、复制延迟、Read-after-write、Failover。
- [ ] Range、Hash、Directory Sharding、热点、扩容、跨分片、全局 ID。

### 一致性与共识

- [ ] CAP 的正确理解。
- [ ] Linearizability、Eventual Consistency、Read-your-writes、Quorum。
- [ ] Raft 的 Leader Election、Log Replication、Term、Majority。

### 分布式事务

- [ ] 2PC、Saga、TCC、Outbox、Compensation、Idempotency。

### 故障模式

- [ ] Network Partition、Partial Failure、Duplicate Request。
- [ ] Retry Storm、Thundering Herd、Split Brain、Hot Key。
- [ ] Backpressure、Cascading Failure。
- [ ] Timeout、Backoff、Jitter、Circuit Breaker、Bulkhead、Load Shedding。

验收：能处理部分失败和最终一致性，并知道什么时候不应该拆微服务或分库分表。

---

## 十、系统设计

固定步骤：需求、非功能、规模估算、API、数据模型、高层架构、读写流程、扩展性、一致性、高可用、安全、可观测性、成本和 Trade-off。

每周设计：URL Shortener、File Storage、Notification、Task Queue、Order、Payment、Search、IM、Google Docs、Figma、ChatGPT、RAG Platform。

验收：先问需求和规模，不直接堆中间件；能说明方案适用范围和为什么不用更复杂方案。

---

## 十一、Docker、Kubernetes、IaC 与 CI/CD

### Docker

- [ ] Dockerfile、Layer、Multi-stage、非 Root、Volume、Network、Health Check、Signal、Compose、Registry、扫描。

### Kubernetes

- [ ] Pod、Deployment、StatefulSet、Service、Ingress。
- [ ] ConfigMap、Secret、Job、CronJob、PV/PVC。
- [ ] Startup、Readiness、Liveness。
- [ ] Rolling、Rollback、Canary、HPA、Helm。

### Terraform

- [ ] Provider、Resource、State、Plan、Apply、Module、Remote State、Drift。

### CI/CD

- [ ] Lint、测试、构建、依赖和镜像扫描。
- [ ] Registry、测试环境、Smoke Test、生产审批、灰度和回滚。
- [ ] Commit、镜像、发布版本可追溯。

验收：能将完整系统自动部署到集群，处理健康检查、扩容、灰度、回滚和基础设施变更。

---

## 十二、可观测性、SRE 与安全

### 可观测性

- [ ] 结构化日志、Request ID、Trace ID、版本和脱敏。
- [ ] RED、USE、业务指标、分位数和告警。
- [ ] OpenTelemetry、Prometheus、Grafana、Loki、Tempo/Jaeger。

### SRE

- [ ] SLI、SLO、SLA、Error Budget。
- [ ] Timeout、Retry、Backoff、Jitter、熔断、限流、隔舱、降级。
- [ ] Backup、Restore、RPO、RTO、故障演练和无责复盘。

### 安全

- [ ] 密码哈希、Session、MFA、RBAC、ABAC、租户隔离。
- [ ] SQL/Command Injection、Path Traversal、SSRF、IDOR、Replay、Mass Assignment。
- [ ] 传输和静态加密、Secret 管理、日志脱敏、审计和数据删除。

验收：能从告警定位到请求、服务、数据库或依赖；系统有明确 SLO、恢复方案和安全边界。

---

## 十三、领域建模与架构

- [ ] Interface、Application、Domain、Infrastructure 分层。
- [ ] Entity、Value Object、Aggregate、Repository、Domain Event。
- [ ] Bounded Context 和统一业务语言。
- [ ] 模块化单体优先于微服务。
- [ ] 业务规则不堆在 Handler，数据库模型不等于领域模型。

---

## 十四、推荐资源

书籍：Learning Go、The Go Programming Language、100 Go Mistakes、How Linux Works、OSTEP、Computer Networking: A Top-Down Approach、高性能 MySQL、Database Internals、DDIA、Release It!、System Design Interview、Site Reliability Engineering、Kubernetes in Action、Terraform: Up & Running。

网站：Go Blog、Go by Example、PostgreSQL Docs、Use The Index Luke、Percona、Cockroach Labs、ClickHouse、Jepsen、High Scalability、ByteByteGo、AWS Architecture、Cloudflare、Netflix Tech、Uber Engineering、Kubernetes Blog、CNCF、OpenTelemetry、Google SRE。


---
