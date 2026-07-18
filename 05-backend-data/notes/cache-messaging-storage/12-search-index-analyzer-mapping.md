---
stage: intermediate
direction: backend-data
topic: search-index-analysis-mapping
---

# 倒排索引、Analyzer、Tokenizer 与 Mapping

全文搜索把文本转换为可检索 term，并用倒排表从 term 定位文档。Mapping 决定字段类型、索引方式和可查询能力；字段一旦以错误类型建立，通常需要新索引重建，而不是原地修正历史 posting。

## 1. 能力边界与机制图

```mermaid
flowchart LR
    A["事实输入"] --> B["索引/处理机制"]
    B --> C["查询或派生输出"]
    C --> D["验证与观测"]
    D --> E["失败恢复/重建"]
```

这条链路的关键不变量是“索引时与查询时使用兼容的分析规则”。排查时保存原文、Analyzer 名称与版本、Analyze API token 序列、mapping 版本和目标索引别名，才能判断问题发生在切词、字段类型还是发布切换。

## 2. 倒排索引

机制：term dictionary 保存词项，posting list 保存包含该词项的文档及可选词频、位置和 offset。

实际使用：适合从词找文档、布尔组合和相关性评分。

主要失败：把它当主键数据库会失去事务、约束和权威写入。

验证方法：用 _termvectors、profile、段统计观察 term 与 posting。

取舍：写入需要分析、段合并和额外存储，换取读检索。

生产实现 倒排索引 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 3. Segment

机制：Lucene 将写入形成不可变 segment，refresh 后可搜索，后台 merge 合并小段。

实际使用：理解 near-real-time、删除标记和合并 I/O。

主要失败：把 refresh 等同持久提交或立即 merge。

验证方法：观察 refresh time、segment count、merge throttle 与磁盘。

取舍：过短 refresh 提升可见性但增加小段和资源。

生产实现 Segment 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 4. Character filter

机制：在 tokenizer 前转换字符并维护 offset 映射。

实际使用：清理 HTML、统一特定字符或映射规则。

主要失败：无边界删除字符会让高亮 offset 与原文错位。

验证方法：用 Analyze API 查看每阶段 token 与 offset。

取舍：预处理越强，原文语义损失越大。

生产实现 Character filter 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 5. Tokenizer

机制：把字符流切为 token，并产生 position 与 offsets。

实际使用：standard、keyword、whitespace 等按语言和字段选择。

主要失败：中文直接 whitespace 会把无空格句子当整体。

验证方法：对真实多语言语料运行 Analyze API。

取舍：更细粒度召回高但索引和噪声增加。

生产实现 Tokenizer 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 6. Token filter

机制：对 token 小写化、停用、同义、词干或 n-gram。

实际使用：统一大小写、形态或业务同义词。

主要失败：索引时同义扩展会放大索引且更新需重建。

验证方法：比较索引 analyzer 与 search analyzer 的 token graph。

取舍：召回提升常伴随精度下降与运维成本。

生产实现 Token filter 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 7. Analyzer

机制：由 char filter、tokenizer、token filters 构成，应用于 text 索引与查询。

实际使用：为标题、正文、代码、中文分别定义。

主要失败：索引与查询 analyzer 不兼容会导致查不到。

验证方法：保存 golden corpus 的输入、tokens 与预期命中。

取舍：一个万能 analyzer 无法同时优化全部字段。

生产实现 Analyzer 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 8. text 与 keyword

机制：text 经分析用于全文；keyword 保留整体用于精确过滤、排序与聚合。

实际使用：用 multi-fields 同时提供全文和精确能力。

主要失败：在 text 上做聚合或把长正文设 keyword。

验证方法：检查 mapping 与 field capabilities。

取舍：多字段增加索引和 doc values 空间。

生产实现 text 与 keyword 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 9. Mapping

机制：定义字段类型、analyzer、doc_values、index、dynamic 等。

实际使用：固定公开索引契约并控制动态字段。

主要失败：无限 dynamic keys 造成 mapping explosion。

验证方法：设置 mapping limit、检查 rejected 与 field count。

取舍：严格 mapping 增加发布协调但保护集群。

生产实现 Mapping 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 10. Doc values

机制：面向列式访问的磁盘结构，用于排序、聚合和脚本。

实际使用：keyword、numeric、date 的聚合排序。

主要失败：把大型不需要聚合字段全部保留 doc_values。

验证方法：用 field usage/profile 与磁盘评估。

取舍：关闭节省空间但失去高效聚合排序。

生产实现 Doc values 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 11. Nested 与 object

机制：object 数组会扁平化字段关系，nested 将每项作为隐藏文档保持配对。

实际使用：订单 items 需保持 sku 与 quantity 同项关系时用 nested。

主要失败：用普通 object 查询导致跨数组元素错误匹配。

验证方法：构造交叉数组反例验证 nested query。

取舍：nested 写入、查询和文档数成本更高。

生产实现 Nested 与 object 时必须把输入类型、处理上限、输出状态和失败类别写入接口契约；调用方不能根据一次成功响应推导未承诺的持久性、一致性或权限语义。

## 12. 方案比较

| 方案 | 主要收益 | 关键边界 |
|---|---|---|
| text | 全文相关性 | 不能直接精确聚合 |
| keyword | 精确过滤排序 | 不分词 |
| object | 简单对象 | 数组字段关系会扁平 |
| nested | 保持数组项关系 | 隐藏文档与查询成本 |
| flattened | 动态键集合 | 类型/查询能力受限 |

## 13. 完整案例：多语言商品搜索

### 输入与约束

中文、英文商品标题与品牌；必须支持品牌精确过滤和标题召回；原始数据库仍是事实源。

### 处理步骤

1. 为 title 定义中文/英文策略并以语言字段路由，不假定同一 tokenizer 适合所有语言。
2. title 使用 text，增加 title.keyword；brand 使用 keyword 并配置规范化。
3. 用 500 条真实查询建立 golden corpus，记录 tokens、期望结果和不应命中。
4. 创建版本化索引 products-v3，批量回填后别名切换。
5. 观测零结果率、nDCG/人工判定、segment 与 merge。

### 输出

标题按语言检索，品牌精确过滤，mapping 可版本化回滚。

### 验证

Analyze API token 与预期一致；同义词更新不破坏旧查询；别名切换前后文档版本对账。

### 失败分支

将中文标题用 whitespace 后大量零结果；修正 tokenizer 并新建索引重建，不能只改 mapping 期待旧 posting 改变。

### 完成判断

多语言商品索引完成的证据包括：中英混合 golden query 的 token 与命中集稳定、`brand.keyword` 聚合不受文本切词影响、旧索引到新索引的源文档版本无缺口，并且别名回退能够恢复旧查询行为。

## 14. 完整案例：订单 items 查询

### 输入与约束

一个订单含 [{sku:A,qty:1},{sku:B,qty:10}]，查询 sku=A 且 qty>=10 不应命中。

### 处理步骤

1. 先用 object mapping 建反例，确认数组字段扁平后会交叉匹配。
2. 新索引将 items 映射 nested。
3. nested query 在同一 nested path 组合 sku 与 quantity。
4. 数据库主键/version 作为外部版本重建。
5. 比较 nested 文档数量与查询延迟。

### 输出

只有同一 item 同时满足条件才命中。

### 验证

反例订单不命中；正确 A:10 订单命中；profile 显示可接受成本。

### 失败分支

只在应用结果后过滤会让分页总数和聚合错误；必须在正确 nested 模型中查询。

### 完成判断

订单 items 案例只有在 object 反例确实误命中、nested 版本消除交叉匹配、正确的 `A:10` 样例仍命中，并完成 nested 文档膨胀与查询 profile 基线后才验收。

## 15. 失败注入矩阵

| 注入点 | 应观察的系统行为 | 不允许的结果 |
|---|---|---|
| 输入中加入未知字段/极端长度 | indexing_rate 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 暂停一个处理节点 | refresh_time 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 让下游返回429或503 | segment_count 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 制造重复与乱序版本 | merge_time 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 使磁盘达到高水位 | mapping_field_count 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 让请求在提交后断线 | rejected_writes 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 撤销一个租户权限 | disk_watermark 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |
| 使后台任务积压超过SLO | analyze_golden_failures 出现可解释变化并触发受控降级 | 静默丢数据、跨权限返回或无限重试 |

## 16. 调试路径与观测信号

1. `indexing_rate`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
2. `refresh_time`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
3. `segment_count`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
4. `merge_time`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
5. `mapping_field_count`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
6. `rejected_writes`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
7. `disk_watermark`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
8. `analyze_golden_failures`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
9. `zero_result_rate`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。
10. `reindex_lag`：定义采样点、单位、维度和告警窗口；只使用低基数标签，具体资源 ID 进入受控日志或 trace。

零结果或召回突变时，先对同一原文运行 Analyze API，比对 character filter、tokenizer、token filter 的逐阶段输出；再核对字段实际 mapping、目标 alias 和 segment 中的文档版本。只有 token 一致后才进入评分与查询排障，避免用 BM25 参数掩盖切词错误。

## 17. 性能、安全与运维边界

1. 所有用户输入有长度、结构、复杂度和超时上限；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
2. 租户与权限条件来自认证主体，不信任请求正文；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
3. 重试有总预算、退避、抖动和幂等保护；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
4. 派生数据能从事实源按明确水位重建；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
5. 敏感正文不进入普通日志、指标标签或公开快照；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
6. 容量按峰值写入、查询、恢复和副本同时估算；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
7. 发布采用版本化结构、canary、对账和回退窗口；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
8. 删除覆盖索引、缓存、快照与下游派生系统；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
9. 依赖不可用时核心正确性不由降级路径破坏；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。
10. 运行手册包含暂停、恢复、验证和人工升级条件；对 倒排索引、Analyzer、Tokenizer 与 Mapping 的实现要给出责任人和可执行验证。

## 18. 综合练习与验收

为“多语言商品搜索”实现最小生产链路，并把“订单 items 查询”作为不同约束下的对照。提交配置、请求/响应、数据样例、故障脚本和观测面板。

- [ ] 验收 1：倒排索引 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 2：Segment 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 3：Character filter 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 4：Tokenizer 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 5：Token filter 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 6：Analyzer 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 7：text 与 keyword 的机制、使用、失败与验证均能用真实输入复现。
- [ ] 验收 8：Mapping 的机制、使用、失败与验证均能用真实输入复现。
- [ ] Analyzer golden set 覆盖中文、英文、标点、大小写、同义词和空 token，并保存预期 token 序列。
- [ ] Mapping 变更通过新索引重建、源版本对账、alias canary 和可执行回退完成发布。
- [ ] Nested 方案记录内部文档放大倍数、磁盘水位与 p95 profile，避免只验证查询正确性。

## 来源

- [OpenSearch 3.7 Mapping](https://docs.opensearch.org/latest/mappings/)（访问日期：2026-07-17）
- [OpenSearch Analyze API](https://docs.opensearch.org/latest/api-reference/analyze-apis/)（访问日期：2026-07-17）
- [Elasticsearch Mapping](https://www.elastic.co/docs/manage-data/data-store/mapping)（访问日期：2026-07-17）
- [Apache Lucene Index File Formats](https://lucene.apache.org/core/)（访问日期：2026-07-17）
- [OpenSearch 3.7 version history](https://docs.opensearch.org/latest/version-history/)（访问日期：2026-07-17）
