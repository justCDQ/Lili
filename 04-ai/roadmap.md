# AI 工程 Roadmap：从零基础到生产级 AI 系统

AI 工具、模型与框架会快速变化，这条路线以相对稳定的原理和工程闭环为主干：任务定义、模型接口、上下文、检索、工具、工作流、评估、可靠性与安全。具体厂商和框架作为可替换实现持续更新。

## 能力阶梯

| 阶段 | 能力目标 | 代表产出 |
| --- | --- | --- |
| 入门 | 能调用模型并理解输入、输出、Token 和限制 | 命令行 AI 小程序 |
| 初级 | 能构建稳定的结构化、流式 AI 功能 | AI Chat 与结构化抽取 |
| 中级 | 能构建 RAG、Tool 和可评估 Workflow | 文档助手与工具助手 |
| 高级 | 能构建可恢复、可观测、安全的 Agent | 长任务 Agent Workflow |
| 专项 | 能治理多模型、成本、质量和生产风险 | AI 平台或企业工作台 |

---

## 阶段零：编程、数据与 API 入门

### 必要编程基础

- [ ] 选择 JavaScript/TypeScript 或 Python 作为第一门 AI 应用语言。
- [ ] 变量、函数、对象/字典、数组/列表、模块、异常和异步。
- [ ] 命令行、Git、包管理、虚拟环境和环境变量。
- [ ] HTTP、JSON、REST API、状态码、认证和流式响应。
- [ ] 文件读写、文本编码、CSV/JSONL 和基础数据清洗。
- [ ] Secret 不进入代码和 Git；理解最小权限和费用上限。

### 第一次模型实验

- [ ] 使用 SDK 和原始 HTTP 各完成一次模型调用。
- [ ] 记录模型标识、日期、输入、输出、Token、延迟和费用。
- [ ] 实现 Streaming、取消、超时和错误展示。
- [ ] 使用 Schema 获取结构化结果并进行运行时校验。
- [ ] 为同一任务准备五条固定样例，比较两种提示或模型。

验收：能独立完成一个命令行或网页 AI 小程序；错误、密钥、超时和输出校验都有明确处理。

---

## 阶段一：基础认知

- [ ] 训练、推理、参数、损失、泛化、过拟合。
- [ ] 训练集、验证集、测试集和 Benchmark。
- [ ] Tokenization、Context Window、输入输出成本。
- [ ] Embedding、向量、点积、余弦相似度。
- [ ] Transformer、Attention、Q/K/V、位置、自回归生成。
- [ ] Pretraining、SFT、偏好优化、Prompt、RAG、Fine-tuning 的区别。

必做：Token 实验室、内存语义搜索、模型参数对比实验。

验收：能解释 LLM 为什么能生成、为什么会幻觉，以及 Prompt、RAG、Tool、Fine-tuning 分别解决什么问题。

---

## 阶段二：模型 API 与 Prompt

### 模型调用

- [ ] System、User、Assistant。
- [ ] 多轮、Streaming、Structured Output、Tool Calling、多模态。
- [ ] 请求取消、超时、有限重试、限流和 Usage。
- [ ] 统一模型 Client，不让业务直接依赖厂商结构。

### Prompt

```text
Role + Task + Context + Constraints + Examples + Output Format + Failure Behavior
```

- [ ] 为同一任务编写三个版本。
- [ ] 建立真实测试样例。
- [ ] 使用 JSON Schema 和运行时校验。
- [ ] Prompt 版本化。
- [ ] 记录模型、参数、Token、延迟和成本。
- [ ] 区分格式错误与内容错误。

验收：不凭感觉判断 Prompt；能判断问题应由 Prompt、代码、RAG、Tool 或 Workflow 解决。

---

## 阶段三：Context Engineering

上下文可能来自系统规则、当前输入、历史、文件、知识库、用户资料、项目状态、工具结果、搜索和记忆。

- [ ] 区分指令与数据。
- [ ] 区分可信和不可信内容。
- [ ] 建立 Token Budget。
- [ ] 对长对话摘要，同时保留关键事实。
- [ ] 去重、过滤过期和冲突信息。
- [ ] 按用户和租户权限获取上下文。
- [ ] 记录最终发给模型的内容。
- [ ] 建立上下文调试页面。
- [ ] 记忆可查看、修改、删除和过期。
- [ ] 区分会话记忆、工作记忆、用户记忆、项目记忆和数据库事实。

验收：能解释一次请求为什么包含这些上下文，以及为什么排除其他内容。

---

## 阶段四：AI 应用与 AI UX

状态：Idle、Preparing、Queued、Streaming、Tool Calling、Waiting Approval、Completed、Failed、Cancelled、Partial Success、Needs Input。

- [ ] 状态机。
- [ ] Streaming、Stop、Retry、Continue。
- [ ] 断线、重复事件、半截 JSON、部分结果保存。
- [ ] 流式 Markdown 和未闭合代码块。
- [ ] Artifact、Diff、History、Restore。
- [ ] Citation、证据和无答案状态。
- [ ] 不确定性、人工确认和人工接管。
- [ ] 页面刷新后的长任务恢复。

验收：用户始终知道系统在做什么，可以停止、修改、重试、确认和恢复。

---

## 阶段五：RAG

```text
文档 → 解析 → 清洗 → Chunk → Embedding → 索引 → 查询处理 → 检索 → Rerank → Context → 生成 → 引用 → 评估
```

### 文档解析

- [ ] PDF、Word、Markdown、HTML、表格、扫描件。
- [ ] 保留标题、页码、来源和原文定位。
- [ ] 去除重复页眉页脚。
- [ ] 分别记录不同文件的解析质量。

### Chunk

- [ ] 固定长度、段落、标题、语义、滑动窗口、父子分块。
- [ ] 比较长度、重叠和结构保留。
- [ ] 处理列表、表格和超长段落。
- [ ] 支持版本、更新、删除和重新索引。

### Retrieval

- [ ] Dense、Keyword、Hybrid。
- [ ] Metadata Filter、权限过滤。
- [ ] Query Rewrite、Multi-query、实体与时间过滤。
- [ ] Top-K、Threshold、Rerank。
- [ ] 无相关结果时不强行回答。

### Evaluation

- [ ] 至少 50 条真实问题。
- [ ] 每条定义相关文档和参考答案。
- [ ] Recall@K、Context Relevance、Groundedness、Citation Accuracy。
- [ ] 无答案、过期资料、冲突来源和权限测试。
- [ ] 修改 Chunk、检索或模型后运行回归。

验收：能定位问题发生在解析、检索、重排、上下文、生成还是引用阶段。

---

## 阶段六：Tool Calling 与 MCP

### Tool 设计

- [ ] 名称和描述清晰。
- [ ] 单一职责、Schema 简洁、输出稳定。
- [ ] 输入验证、超时、有限重试、幂等。
- [ ] 只读与写入工具分开。
- [ ] 写操作显示影响范围并请求确认。
- [ ] 权限、审计、脱敏和错误隔离。

### MCP

- [ ] Host、Client、Server。
- [ ] Transport、Tools、Resources、Prompts、Capabilities。
- [ ] 使用 Inspector 调试。
- [ ] 开发本地 Server。
- [ ] 增加认证、授权和日志。
- [ ] 防止任意文件、URL、命令访问和 SSRF。

验收：模型不拥有最终权限；所有参数和权限由服务端重新校验。

---

## 阶段七：Workflow 与 Agent

优先使用可预测 Workflow，固定流程可解决时不优先使用自由 Agent。

### Workflow 模式

- [ ] Prompt Chaining。
- [ ] Routing。
- [ ] Parallelization。
- [ ] Orchestrator-Workers。
- [ ] Evaluator-Optimizer。

### Agent 组件

- [ ] Goal、State、Planner、Tools、Executor、Evaluator、Memory、Stop Condition、Human Approval。
- [ ] 最大步骤、Token、成本和总超时。
- [ ] 暂停、取消、恢复、失败步骤和部分完成。
- [ ] 防止循环、重复操作和重试风暴。

### 长任务架构

```text
创建 Task → 数据库 → Queue → Worker → 持久化步骤 → 推送进度 → 完成
```

验收：任务可观测、可恢复、高风险动作需确认，能量化最终任务完成率。

---

## 阶段八：Evaluation

评估维度：Correctness、Relevance、Completeness、Format、Safety、Latency、Cost、Groundedness、Tool Selection、Task Completion。

- [ ] 从真实请求收集正常、边界、无答案、权限和对抗样例。
- [ ] 为 Bug 增加回归样例。
- [ ] 使用确定性检查、人工评估和 LLM Judge。
- [ ] 测试 Judge 的稳定性。
- [ ] 同时记录质量、延迟和成本。
- [ ] 设置上线门槛和 CI 回归。
- [ ] 线上失败回流到离线评估集。

验收：每次改 Prompt、模型、检索、Tool 或 Agent，都能量化质量变化。

---

## 阶段九：生产工程

### Model Gateway

- [x] [Generate、Stream、Embed、Structured Output、Tool、Usage 和 Error 的统一接口](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#2-模型网关统一能力而非抹平差异)。
- [x] [多模型路由、Fallback、超时、重试](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#路由决策)。
- [x] [按任务、复杂度、上下文、延迟、成本和合规选择模型](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#路由决策)。

### Cache 与 Queue

- [x] [Exact、Semantic、Retrieval、Tool Result 和 Prefix Cache](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#3-缓存复用结果但不跨越权限和新鲜度)。
- [x] [租户隔离、过期、失效和命中率](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#3-缓存复用结果但不跨越权限和新鲜度)。
- [x] [文档解析、Embedding、批量生成和 Agent 长任务进入队列](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#4-queue把可等待工作从同步路径移走)。
- [x] [重试、死信、幂等、取消、进度和配额](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#队列合同与幂等)。

### Observability

- [x] [Request ID、User、Tenant、Model、Prompt Version](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#5-可观测性从一次答案追到每个决定)。
- [x] [Token、Cost、首 Token 延迟和总时长](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#5-可观测性从一次答案追到每个决定)。
- [x] [Retrieval Query、过滤、结果与分数](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#5-可观测性从一次答案追到每个决定)。
- [x] [Tool 参数、结果、审批和错误](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#5-可观测性从一次答案追到每个决定)。
- [x] [Agent Step、状态、停止原因](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#5-可观测性从一次答案追到每个决定)。

### Reliability

- [x] [Timeout、Backoff、Jitter、Circuit Breaker、Rate Limit、Fallback、Partial Success](notes/14-production-engineering/01-model-gateway-cache-queue-observability-reliability.md#6-可靠性机制的适用边界)。

验收：线上问题能定位到具体模型、Prompt、上下文、检索、Tool 或 Agent 步骤。

---

## 阶段十：安全与治理

- [x] [Prompt Injection 和间接注入](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#1-威胁模型谁能控制什么)。
- [x] [外部网页、邮件、PDF、数据库和 Tool 结果视为不可信数据](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#2-把指令数据和能力拆开)。
- [x] [服务端最小权限和租户隔离](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#3-最小权限与租户隔离)。
- [x] [写操作确认、二次验证、幂等和审计](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#4-写操作确认不是一个按钮而是受约束的交易)。
- [x] [日志脱敏、数据保存周期和删除机制](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#5-日志保留和删除)。
- [x] [限制文件路径、网络目标、命令和执行沙箱](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#6-文件网络命令与沙箱)。
- [x] [对第三方 MCP Server 做安全审查](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#7-第三方-mcp-server-审查)。
- [x] [Red Team：泄露 Prompt、跨租户、删除数据、内网访问、重复支付、资源消耗](notes/15-security-governance/01-prompt-injection-permissions-audit-and-red-team.md#8-红队把攻击假设变成可重复测试)。

验收：模型无法绕过服务端权限；所有高风险动作可审计且需要确认。

---

## 阶段十一：进阶选修

### 开源模型与本地推理

[Hugging Face、量化、显存、Batch、Throughput、Latency、KV Cache、Serving、许可证和总持有成本](notes/16-advanced-options/01-local-inference-and-fine-tuning-decisions.md#2-开源模型与本地推理的组成)。

### Fine-tuning

[Fine-tuning 的适用条件、PEFT/LoRA、数据与评估、训练恢复和发布](notes/16-advanced-options/01-local-inference-and-fine-tuning-decisions.md#7-fine-tuning-的输入输出与风险)：仅在 Prompt、上下文、RAG、Tool 和模型选择都不能解决，并拥有高质量数据和评估集时使用。

---

## 持续更新机制

### 每日：一个可复现的 AI 实验

- [ ] 从 Prompt、Context、Structured Output、RAG、Tool、Workflow、Agent、Evaluation 或 Safety 中选择一个小问题。
- [ ] 固定任务与样例，只改变一个变量。
- [ ] 记录模型完整标识、接口/框架版本、日期和配置。
- [ ] 同时记录质量、延迟、Token、成本、失败类型和安全影响。
- [ ] 保存输入、输出、Trace 或截图等证据，不只记录“效果不错”。
- [ ] 得出可证伪结论，并说明适用范围和下一次实验。

### 每周：能力与案例更新

- [ ] 精读一个官方文档、技术报告、模型卡、协议变更或高质量工程案例。
- [ ] 用现有评估集复测一个新模型、新能力或框架升级。
- [ ] 把线上失败或实验失败加入回归样例。
- [ ] 更新一张能力卡：它解决什么、不能解决什么、成本和替代方案是什么。

### 每季度：路线图维护

- [ ] 检查官方 API、模型生命周期、SDK 和协议版本变化。
- [ ] 检查 RAG、Agent、Evaluation、Safety 与多模态是否出现稳定的新实践。
- [ ] 删除只绑定过时框架的内容，保留稳定原理和迁移说明。
- [ ] 在仓库变更记录中写明新增、修改、废弃及证据来源。

更新原则：官方文档、规范、论文和可复现实验优先于二手资讯；新能力必须经过固定评估集，而不是仅凭演示判断。

---

## 学习资源

书籍：动手学深度学习、Hands-On Large Language Models、Build a Large Language Model From Scratch、AI Engineering、Designing Machine Learning Systems、Machine Learning Design Patterns、DDIA、Release It!。

网站与博客：[Hugging Face Learn](https://huggingface.co/learn)、[Model Context Protocol](https://modelcontextprotocol.io/)、[Anthropic Engineering](https://www.anthropic.com/engineering)、3Blue1Brown、Andrej Karpathy、fast.ai、Simon Willison、Lilian Weng、Chip Huyen、Eugene Yan、Hamel Husain、Latent Space。优先阅读所用模型和 SDK 的官方文档与模型卡。


---
