---
stage: intermediate
direction: backend-data
topic: search-query-aggregation
---

# 全文查询、过滤、高亮、聚合、自动补全与 Search After

搜索请求由召回、过滤、评分、排序和分页组成。相关性查询与结构化过滤承担不同职责；深分页必须使用稳定排序、PIT 和 search_after，而不是无界 from/size。

## 1. 能力边界与机制图

```mermaid
flowchart LR
    A["事实输入"] --> B["索引/处理机制"]
    B --> C["查询或派生输出"]
    C --> D["验证与观测"]
    D --> E["失败恢复/重建"]
```

查询链路要同时固定 query DSL、权限 filter、索引/alias、PIT 与 sort 值。相关性分数只在同一索引快照和相同请求条件下可比较；聚合、总数和高亮也必须受同一个权限过滤约束。

## 2. Query context

机制：匹配文档并计算 score。

实际使用：match、multi_match、bool should。

主要失败：把权限条件放评分分支导致泄露。

验证方法：用 explain/profile 与命中集合验证。

取舍：相关性灵活但计算成本高。

生产实现 Query context 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 3. Filter context

机制：只判断匹配，不计算相关性并可能缓存 bitset。

实际使用：tenant、状态、时间、权限边界。

主要失败：把用户输入 query_string 当权限过滤。

验证方法：对每个主体做资源级测试。

取舍：精确快速但不负责排序质量。

生产实现 Filter context 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 4. BM25

机制：按词频、逆文档频率和长度归一评分。

实际使用：文本相关性基线。

主要失败：把 score 当跨查询可比绝对值。

验证方法：固定语料比较 ranking evaluation。

取舍：需要字段 boost 与业务信号校准。

生产实现 BM25 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 5. Bool query

机制：must/filter/should/must_not 组合条件。

实际使用：表达结构化召回。

主要失败：minimum_should_match 默认误解。

验证方法：用边界用例列真值表。

取舍：复杂嵌套降低可解释性。

生产实现 Bool query 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 6. Highlight

机制：基于匹配 term 和 offset 生成片段。

实际使用：展示命中上下文。

主要失败：直接把高亮 HTML 注入 DOM 造成 XSS。

验证方法：转义原文并只允许受控标签。

取舍：增加查询/存储 offset 成本。

生产实现 Highlight 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 7. Aggregation

机制：在匹配文档上做 bucket/metric/pipeline。

实际使用：分类面板、统计分布。

主要失败：高基数 terms size 无界。

验证方法：观察 bucket 数、内存、误差/其他桶。

取舍：实时聚合便利但可能压集群。

生产实现 Aggregation 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 8. Autocomplete

机制：edge n-gram、completion、search-as-you-type 等候选策略。

实际使用：前缀建议与拼写体验。

主要失败：任意 n-gram 导致索引爆炸。

验证方法：测索引大小、候选精度和输入延迟。

取舍：预计算快但更新和个性化复杂。

生产实现 Autocomplete 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 9. from/size

机制：跳过前面 hits 后取页面。

实际使用：浅分页和后台小数据。

主要失败：深页每 shard 维护大量候选。

验证方法：限制 max window 并压测。

取舍：支持页码但成本随深度上升。

生产实现 from/size 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 10. search_after

机制：以上一页完整 sort values 作为下一页边界。

实际使用：稳定向后遍历。

主要失败：sort 不唯一导致重复/遗漏。

验证方法：加入唯一 tie-breaker 并检测重复。

取舍：不能随意跳页，状态需客户端保存。

生产实现 search_after 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 11. PIT

机制：固定索引读取视图避免 refresh 改变分页集合。

实际使用：长列表与导出的一致分页。

主要失败：忘记关闭/keep_alive 过长占资源。

验证方法：监控 PIT context、过期恢复。

取舍：提高一致性但不是数据库事务快照。

生产实现 PIT 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 12. 方案比较

| 方案 | 主要收益 | 关键边界 |
|---|---|---|
| match | 分析后全文 | 计算score |
| term | 精确term | 不分析输入 |
| filter | 结构约束 | 无score |
| from/size | 页码简单 | 深分页昂贵 |
| PIT+search_after | 稳定游标 | 不能随机跳页 |

## 13. 完整案例：商品搜索页

### 输入与约束

关键词、品牌、价格区间、tenant权限；显示高亮和品牌聚合；p95 150ms。

### 处理步骤

1. tenant/权限放 filter，关键词用 multi_match 并控制字段boost。
2. 聚合只在允许字段，terms size 有上限。
3. 高亮转义原文并限制片段数/长度。
4. 查询模板白名单，不开放任意脚本/query_string。
5. 用真实查询集做ranking evaluation和负载测试。

### 输出

结果、facet和高亮在权限范围内，延迟有界。

### 验证

无权商品不进入hits或聚合计数；高亮不执行HTML；慢查询profile可定位。

### 失败分支

先搜索全局再应用层过滤会泄漏聚合和总数；修正为权限filter进入搜索请求。

### 完成判断

商品搜索页的验收以整份响应为单位：golden query 的排序达到约定质量，权限外商品不进入 hits、total 或 facet，高亮经安全渲染，并且慢查询能由 profile 定位到具体 clause 与 shard。

## 14. 完整案例：千万结果导出

### 输入与约束

按 created_at,id 排序导出，索引持续写入，不能重复遗漏。

### 处理步骤

1. 创建PIT并设置有界keep_alive。
2. sort使用created_at asc与唯一_id asc。
3. 每页把最后hit sort数组作为search_after。
4. 每页续PIT keep_alive并记录导出checkpoint。
5. 结束/失败关闭PIT，超时则从业务水位重新开始。

### 输出

同一PIT内稳定遍历，导出可检查重复ID。

### 验证

并发写入时导出集合稳定；一百万条ID无重复；PIT资源回收。

### 失败分支

只按created_at会在同毫秒记录处重复/遗漏；加入唯一tie-breaker。

### 完成判断

导出验收要求在固定 PIT 内遍历全部预期 ID，`created_at + _id` sort 值严格前进，重复与遗漏均为零；中断恢复要么继续有效 checkpoint，要么从新的业务水位完整重启，并确认 PIT 已释放。

## 15. 失败注入矩阵

| 注入点 | 应观察的系统行为 | 不允许的结果 |
|---|---|---|
| 输入中加入未知字段/极端长度 | search_p95 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 暂停一个处理节点 | query_cache_hit 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 让下游返回429或503 | fetch_p95 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 制造重复与乱序版本 | aggregation_buckets 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 使磁盘达到高水位 | highlight_bytes 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 让请求在提交后断线 | zero_result_rate 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 撤销一个租户权限 | autocomplete_precision 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 使后台任务积压超过SLO | deep_page_reject 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |

## 16. 调试路径与观测信号

1. `search_p95`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
2. `query_cache_hit`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
3. `fetch_p95`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
4. `aggregation_buckets`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
5. `highlight_bytes`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
6. `zero_result_rate`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
7. `autocomplete_precision`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
8. `deep_page_reject`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
9. `pit_contexts`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
10. `partial_search_results`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。

搜索变慢时先保存完整 DSL、索引、routing、PIT ID 和 shard failure，再比较 query/fetch 两阶段耗时与 profile tree。排序异常检查 BM25 输入和 sort tie-breaker；facet 泄漏检查 filter 是否在搜索执行层，而不是只看应用返回的 hits。

## 17. 性能、安全与运维边界

1. 所有用户输入有长度、结构、复杂度和超时上限；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
2. 租户与权限条件来自认证主体，不信任请求正文；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
3. 重试有总预算、退避、抖动和幂等保护；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
4. 派生数据能从事实源按明确水位重建；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
5. 敏感正文不进入普通日志、指标标签或公开快照；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
6. 容量按峰值写入、查询、恢复和副本同时估算；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
7. 发布采用版本化结构、canary、对账和回退窗口；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
8. 删除覆盖索引、缓存、快照与下游派生系统；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
9. 依赖不可用时核心正确性不由降级路径破坏；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。
10. 运行手册包含暂停、恢复、验证和人工升级条件；对 全文查询、过滤、高亮、聚合、自动补全与 Search After 的实现要给出责任人和可执行验证。

## 18. 综合练习与验收

为“商品搜索页”实现最小生产链路，并把“千万结果导出”作为不同约束下的对照。提交配置、请求/响应、数据样例、故障脚本和观测面板。

- [ ] 验收 1：Query context 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 2：Filter context 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 3：BM25 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 4：Bool query 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 5：Highlight 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 6：Aggregation 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 7：Autocomplete 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 8：from/size 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 商品查询集同时验证 hits、total、aggregation、highlight 与权限过滤，不能只检查首条结果。
- [ ] Autocomplete 记录 prefix 长度分布、precision、索引大小和请求频率，并对空串和超长前缀限流。
- [ ] Search After 导出用稳定唯一排序键，验证 PIT 过期、checkpoint 失效和资源回收路径。

## 来源

- [OpenSearch Search API](https://docs.opensearch.org/latest/api-reference/search-apis/search/)（访问日期：2026-07-17）
- [OpenSearch Paginate results](https://docs.opensearch.org/latest/search-plugins/searching-data/paginate/)（访问日期：2026-07-17）
- [OpenSearch Aggregations](https://docs.opensearch.org/latest/aggregations/)（访问日期：2026-07-17）
- [Elasticsearch Search APIs](https://www.elastic.co/docs/reference/elasticsearch/rest-apis/search-apis)（访问日期：2026-07-17）
- [OpenSearch 3.7 version history](https://docs.opensearch.org/latest/version-history/)（访问日期：2026-07-17）
