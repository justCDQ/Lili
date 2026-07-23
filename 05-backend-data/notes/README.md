# 后端与数据学习笔记

本索引覆盖 [后端与数据 Roadmap](../roadmap.md) 的阶段零至阶段十三。路线图每个复选项对应一篇独立笔记；强耦合的中级主题在一篇深度文章内分节覆盖。算法日练仍放在 `daily/`；可迁移的原理和模式归档在此处。

## 覆盖表

| Roadmap 阶段 | 路线图知识点 | 笔记数 | 状态 |
| --- | ---: | ---: | --- |
| 阶段零：编程入门 | 5 | 5 | 已覆盖 |
| 阶段零：工具与环境 | 5 | 5 | 已覆盖 |
| 阶段零：服务与数据入门 | 5 | 5 | 已覆盖 |
| 阶段一：数据结构与算法 | 6 | 6 | 已覆盖 |
| 阶段二：Go 语言 | 2 | 2 | 已覆盖 |
| 阶段二：Go 并发 | 4 | 4 | 已覆盖 |
| 阶段二：Go 工具 | 2 | 2 | 已覆盖 |
| 阶段三：Linux | 3 | 3 | 已覆盖 |
| 阶段三：操作系统 | 4 | 4 | 已覆盖 |
| 阶段三：网络 | 4 | 4 | 已覆盖 |
| 阶段四：API 与服务 | 7 | 7 | 已覆盖 |
| 阶段五：数据库建模 | 3 | 3 | 已覆盖 |
| 阶段五：SQL | 2 | 2 | 已覆盖 |
| 阶段五：索引 | 3 | 3 | 已覆盖 |
| 阶段五：事务 | 3 | 3 | 已覆盖 |
| 阶段五：Migration | 3 | 3 | 已覆盖 |
| 阶段六：Redis 与缓存 | 6 | 6 | 已覆盖 |
| 阶段七：消息队列与异步任务 | 5 | 5 | 已覆盖 |
| 阶段八：搜索 | 4 | 4 | 已覆盖 |
| 阶段八：对象存储 | 3 | 3 | 已覆盖 |
| 阶段八：数据工程 | 5 | 5 | 已覆盖 |
| 阶段九：分布式系统 | 11 | 3 | 已覆盖 |
| 阶段十：系统设计 | 1 | 1 | 已覆盖 |
| 阶段十一：云原生交付 | 9 | 1 | 已覆盖 |
| 阶段十二：可观测性、SRE 与安全 | 9 | 1 | 已覆盖 |
| 阶段十三：领域建模与架构 | 5 | 1 | 已覆盖 |
| **合计** | **119** | **91** | **已覆盖阶段零至十三** |

## 阶段零：编程与计算机基础

### 编程入门

1. [值、变量、类型、条件、循环、函数与错误](00-programming-basics/01-values-variables-control-functions-errors.md)
2. [Array、Slice、Map、Struct、Object、String 与文件](00-programming-basics/02-collections-strings-files.md)
3. [模块、包、依赖、输入输出与命令行参数](00-programming-basics/03-modules-dependencies-io-cli.md)
4. [伪代码、日志、断点与测试调试](00-programming-basics/04-pseudocode-debugging-testing.md)
5. [从 JavaScript 算法表达迁移到 Go](00-programming-basics/05-javascript-to-go-comparison.md)

### 工具与环境

1. [文件、目录、路径、权限、进程与环境变量](01-computer-systems/01-files-paths-permissions-processes-env.md)
2. [终端、Shell 与 Git 基础](01-computer-systems/02-shell-git.md)
3. [编辑器、格式化、静态检查、测试与调试器](01-computer-systems/03-editor-format-lint-test-debug.md)
4. [二进制、十六进制、字符编码、时间与时区](01-computer-systems/04-binary-hex-encoding-time.md)
5. [CPU、内存、磁盘与网络](01-computer-systems/05-cpu-memory-disk-network.md)

### 服务与数据入门

1. [客户端、服务器、进程、端口、域名与 DNS](02-service-data-basics/01-client-server-process-port-domain-dns.md)
2. [HTTP、Header、JSON 与 Cookie](02-service-data-basics/02-http-json-cookie.md)
3. [关系数据库与事务直觉](02-service-data-basics/03-database-relational-transaction-intuition.md)
4. [同步、异步、并发、并行与队列](02-service-data-basics/04-sync-async-concurrency-parallelism-queue.md)
5. [本地服务、curl 与日志](02-service-data-basics/05-local-http-curl-logs.md)

## 阶段一：数据结构与算法

1. [复杂度与线性结构](03-algorithms/01-complexity-linear-structures.md)
2. [哈希表](03-algorithms/02-hash-table.md)
3. [树、堆、Trie 与 B+Tree](03-algorithms/03-trees-heap-trie-bplus.md)
4. [图算法](03-algorithms/04-graphs.md)
5. [二分、排序与外部排序](03-algorithms/05-search-sort-external-sort.md)
6. [双指针、窗口、分治、贪心、回溯与动态规划](03-algorithms/06-problem-solving-patterns.md)

## 阶段二：Go

### 语言

1. [Go 语言核心类型与组织](04-go/01-language-core.md)
2. [Go 错误模型](04-go/02-error-model.md)

### 并发

1. [Goroutine、Channel、Select 与 Context](04-go/03-goroutine-channel-select-context.md)
2. [Mutex、RWMutex、WaitGroup 与 Atomic](04-go/04-mutex-waitgroup-atomic.md)
3. [Worker Pool、并发限制与优雅停止](04-go/05-worker-pool-limits-graceful-stop.md)
4. [Race Condition、Deadlock 与 Goroutine Leak](04-go/06-race-deadlock-goroutine-leak.md)

### 工具

1. [Test、Table-driven Test 与 Benchmark](04-go/07-tests-table-benchmark.md)
2. [Race Detector、pprof、trace、vet、Staticcheck 与 Delve](04-go/08-diagnostic-tools.md)

## 阶段三：Linux、操作系统和网络

### Linux

1. [文件、权限、用户、进程、Signal、Pipe、Socket 与 systemd](05-linux-network/01-linux-files-users-processes-ipc-systemd.md)
2. [ps、top、htop、lsof、ss、curl、grep、awk 与 sed](05-linux-network/02-linux-diagnostic-commands.md)
3. [CPU、内存、磁盘、网络与日志排查](05-linux-network/03-resource-log-troubleshooting.md)

### 操作系统

1. [进程、线程、用户态、内核态与系统调用](05-linux-network/04-process-thread-syscall.md)
2. [虚拟内存、Page、Stack、Heap 与上下文切换](05-linux-network/05-virtual-memory-pages-stack-heap-switch.md)
3. [File Descriptor、阻塞/非阻塞、I/O Multiplexing 与 epoll](05-linux-network/06-file-descriptor-nonblocking-epoll.md)
4. [Lock、Deadlock、OOM、CPU 密集与 I/O 密集](05-linux-network/07-lock-deadlock-oom-workload.md)

### 网络

1. [DNS、TCP 握手、重传、滑动窗口与拥塞控制](05-linux-network/08-dns-tcp-reliability-congestion.md)
2. [TLS、证书与 HTTPS](05-linux-network/09-tls-certificates-https.md)
3. [HTTP/1.1、HTTP/2、HTTP/3、Keepalive 与 Connection Pool](05-linux-network/10-http-versions-connections.md)
4. [TIME_WAIT、反向代理与 L4/L7 负载均衡](05-linux-network/11-time-wait-proxy-load-balancing.md)

## 阶段四：API 与服务

1. [REST、RPC、gRPC、GraphQL、WebSocket、SSE 与 Webhook](06-api-database/01-api-styles.md)
2. [Resource、Method、Status Code、分页、筛选、排序、批量与版本](06-api-database/02-resource-method-pagination-version.md)
3. [统一错误模型](06-api-database/03-unified-errors.md)
4. [Request ID、Logging、Recovery、Timeout、CORS、Rate Limit、Metrics 与 Trace](06-api-database/04-service-middleware-observability.md)
5. [Session、Cookie、JWT、OAuth 2、OIDC 与 API Key](06-api-database/05-authentication-mechanisms.md)
6. [RBAC、ABAC、资源权限、多租户与审计](06-api-database/06-authorization-multitenancy-audit.md)
7. [幂等 Key、异步任务 API 与 OpenAPI](06-api-database/07-idempotency-async-openapi.md)

## 阶段五：关系型数据库

### 建模

1. [实体、关系、主键、外键与唯一约束](06-api-database/08-relational-model-constraints.md)
2. [规范化、反规范化、审计字段、状态与历史](06-api-database/09-normalization-history.md)
3. [软删除风险与多租户数据库设计](06-api-database/10-soft-delete-multitenancy.md)

### SQL

1. [Join、Group By、Having、Subquery、CTE 与 Window Function](06-api-database/11-advanced-sql.md)
2. [分页、游标分页、复杂筛选与 N+1](06-api-database/12-pagination-filter-nplusone.md)

### 索引

1. [B+Tree、复合索引、最左前缀与覆盖索引](06-api-database/13-btree-composite-covering.md)
2. [Selectivity、Index Scan、Full Table Scan 与 Write Amplification](06-api-database/14-selectivity-scans-write-amplification.md)
3. [EXPLAIN、EXPLAIN ANALYZE 与慢查询](06-api-database/15-explain-slow-queries.md)

### 事务

1. [ACID、隔离级别、MVCC、锁与死锁](06-api-database/16-acid-isolation-mvcc-locks.md)
2. [乐观锁、悲观锁、唯一约束冲突与库存并发](06-api-database/17-concurrency-control-inventory.md)
3. [避免长事务与事务中调用慢外部接口](06-api-database/18-transaction-boundaries.md)

### Migration

1. [Schema Migration、Data Migration 与 Expand and Contract](06-api-database/19-schema-data-expand-contract.md)
2. [向后兼容、零停机、大表回填与回滚](06-api-database/20-zero-downtime-backfill-rollback.md)
3. [备份、恢复、RPO 与 RTO](06-api-database/21-backup-recovery-rpo-rto.md)

## 阶段六：Redis 与缓存

1. [Redis 数据类型、原子命令与内存模型](07-cache-messaging-storage/01-redis-data-types.md)
2. [Cache Aside、Read Through、Write Through 与多级缓存](07-cache-messaging-storage/02-cache-patterns-multilevel.md)
3. [缓存 Key、TTL、命中率、随机过期与容量治理](07-cache-messaging-storage/03-cache-key-ttl-hit-rate.md)
4. [缓存穿透、击穿、雪崩与热点 Key 治理](07-cache-messaging-storage/04-cache-penetration-breakdown-avalanche-hot-key.md)
5. [缓存更新、删除、最终一致性与事件驱动失效](07-cache-messaging-storage/05-cache-invalidation-consistency.md)
6. [Redis 分布式锁、租约续期与 Fencing Token](07-cache-messaging-storage/06-redis-distributed-lock-fencing.md)

## 阶段七：消息队列与异步任务

1. [消息系统模型与 Consumer Group](07-cache-messaging-storage/07-messaging-model.md)
2. [投递语义、重复、乱序与幂等 Consumer](07-cache-messaging-storage/08-delivery-semantics-idempotent-consumer.md)
3. [Retry、指数退避、DLQ 与 Poison Message](07-cache-messaging-storage/09-retry-dlq-poison-message.md)
4. [消息 Schema、积压、优先级与任务取消](07-cache-messaging-storage/10-message-schema-backlog-priority-cancellation.md)
5. [Transactional Outbox、最终一致性与补偿](07-cache-messaging-storage/11-outbox-consistency-compensation.md)

## 阶段八：搜索、对象存储与数据工程

### 搜索

1. [倒排索引、Analyzer、Tokenizer 与 Mapping](07-cache-messaging-storage/12-search-index-analyzer-mapping.md)
2. [全文查询、过滤、高亮、聚合、自动补全与 Search After](07-cache-messaging-storage/13-search-query-aggregation-autocomplete-pagination.md)
3. [搜索索引的数据库同步、更新、删除与重建](07-cache-messaging-storage/14-search-sync-update-delete.md)
4. [搜索引擎作为派生系统](07-cache-messaging-storage/15-search-derived-system.md)

### 对象存储

1. [S3、Bucket、Object、Key 与 Presigned URL](07-cache-messaging-storage/16-s3-object-presigned-url.md)
2. [Multipart Upload、版本、生命周期、CDN 与 ETag](07-cache-messaging-storage/17-multipart-version-lifecycle-cdn-etag.md)
3. [上传类型与内容校验、权限、病毒扫描与审计](07-cache-messaging-storage/18-upload-security-scan-audit.md)

### 数据工程

1. [OLTP、OLAP、ETL、ELT 与 CDC](07-cache-messaging-storage/19-oltp-olap-etl-elt-cdc.md)
2. [Warehouse、Lake、Batch、Streaming、Lineage 与 Data Quality](07-cache-messaging-storage/20-warehouse-lake-batch-streaming-lineage-quality.md)
3. [事实表、维度表与星型模型](07-cache-messaging-storage/21-dimensional-modeling.md)
4. [ClickHouse、dbt、Airflow、Kafka、Flink 与 Spark 的定位](07-cache-messaging-storage/22-data-tools-positioning.md)
5. [漏斗、留存、错误率、性能与 AI 成本分析](07-cache-messaging-storage/23-product-analytics-metrics.md)

新增笔记时复制 [后端与数据笔记模板](../notes-template.md) 并更新本索引。

## 阶段九：分布式系统

1. [复制、故障切换与分片路由](08-distributed-systems/01-replication-sharding.md)：Leader-Follower、Multi-Leader、Leaderless、复制确认、Read-after-write、Failover、Range/Hash/Directory 分片、热点、迁移、跨分片与全局 ID。
2. [一致性、仲裁与 Raft：在网络分区中定义正确写入](08-distributed-systems/02-consistency-consensus-cap-raft.md)：CAP、线性一致性、最终一致性、读己之写、Quorum、Raft 选举、日志复制、多数派与 split brain 防护。
3. [分布式事务与韧性：把不确定结果变成可恢复流程](08-distributed-systems/03-distributed-transactions-failure-resilience.md)：2PC、Saga、TCC、Outbox、补偿、幂等、部分失败、重试、背压、熔断、隔离与负载削减。

## 阶段十：系统设计

1. [需求、规模估算与架构取舍](09-system-design/01-requirements-estimation-architecture-tradeoffs.md)

## 阶段十一：Docker、Kubernetes、IaC 与 CI/CD

1. [容器、Kubernetes、Terraform 与 CI/CD](10-cloud-delivery/01-containers-kubernetes-iac-cicd.md)

## 阶段十二：可观测性、SRE 与安全

1. [可观测性、SRE 与应用安全](11-observability-sre-security/01-observability-sre-security.md)

## 阶段十三：领域建模与架构

1. [领域建模与分层架构](12-domain-architecture/01-domain-modeling-layered-architecture.md)
