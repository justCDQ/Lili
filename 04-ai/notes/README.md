# AI 学习笔记

覆盖 [AI Roadmap](../roadmap.md) 的学习笔记。每个路线图知识点对应一篇独立笔记；实验性、时效性结论必须记录模型、版本和日期。

## 阶段零：编程、数据与 API 入门

- [x] [JavaScript、TypeScript 与 Python 的选择](foundations/language-choice.md)
- [x] [AI 应用需要的编程基础](foundations/programming-basics.md)
- [x] [命令行、Git、包管理、虚拟环境与环境变量](foundations/development-environment.md)
- [x] [HTTP、JSON、REST API、状态码、认证与流式响应](foundations/http-json-api.md)
- [x] [文件、文本编码、CSV、JSONL 与数据清洗](foundations/files-encoding-data-formats.md)
- [x] [Secret、最小权限与费用上限](foundations/secrets-permissions-cost.md)
- [x] [SDK 与原始 HTTP 模型调用](model-api/sdk-vs-http.md)
- [x] [模型标识、输入输出、Token、延迟与费用记录](model-api/request-metadata-usage.md)
- [x] [Streaming、取消、超时与错误展示](model-api/streaming-cancellation-timeout-errors.md)
- [x] [Structured Output、Schema 与运行时校验](model-api/structured-output-validation.md)
- [x] [固定样例与模型/Prompt 对比](evaluation/fixed-cases-comparison.md)

## 阶段一：基础认知

- [x] [训练、推理、参数、损失、泛化与过拟合](foundations/training-inference-generalization.md)
- [x] [训练集、验证集、测试集与 Benchmark](foundations/dataset-splits-benchmarks.md)
- [x] [Tokenization、Context Window 与输入输出成本](foundations/tokenization-context-cost.md)
- [x] [Embedding、向量、点积与余弦相似度](foundations/embeddings-vector-similarity.md)
- [x] [Transformer、Attention、Q/K/V、位置与自回归生成](foundations/transformer-attention-generation.md)
- [x] [Pretraining、SFT、偏好优化、Prompt、RAG 与 Fine-tuning](foundations/model-adaptation-methods.md)

## 阶段二：模型 API 与 Prompt

- [x] [System、User 与 Assistant 消息](model-api/message-roles.md)
- [x] [多轮、Streaming、Structured Output、Tool Calling 与多模态](model-api/core-api-capabilities.md)
- [x] [取消、超时、有限重试、限流与 Usage](model-api/reliability-rate-limits-usage.md)
- [x] [统一模型 Client 与厂商隔离](model-api/unified-model-client.md)
- [x] [Prompt 的 Role、Task、Context、Constraints、Examples、Output Format 与 Failure Behavior](prompt/prompt-anatomy.md)
- [x] [同一任务的 Prompt 版本对比](prompt/prompt-comparison.md)
- [x] [真实测试样例](prompt/real-test-cases.md)
- [x] [JSON Schema 与运行时校验](prompt/schema-in-prompt-workflow.md)
- [x] [Prompt 版本化](prompt/prompt-versioning.md)
- [x] [模型、参数、Token、延迟与成本记录](prompt/prompt-observability.md)
- [x] [格式错误与内容错误](prompt/format-vs-content-errors.md)
- [x] [Prompt、代码、RAG、Tool 与 Workflow 的选择](prompt/solution-selection.md)

## 阶段三：Context Engineering

- [x] [指令与数据的边界](context-engineering/01-instruction-data-boundary.md)
- [x] [可信边界与不可信上下文](context-engineering/02-trust-boundaries-untrusted-context.md)
- [x] [Token Budget 分配](context-engineering/03-token-budget-allocation.md)
- [x] [长对话摘要与关键事实](context-engineering/04-conversation-summary-key-facts.md)
- [x] [上下文去重、过期与冲突](context-engineering/05-dedup-staleness-conflicts.md)
- [x] [上下文权限与租户隔离](context-engineering/06-context-permission-tenant-isolation.md)
- [x] [最终模型输入的记录与可重放性](context-engineering/07-final-model-input-recording.md)
- [x] [上下文调试页面](context-engineering/08-context-debugger.md)
- [x] [记忆生命周期与用户控制](context-engineering/09-memory-lifecycle-user-control.md)
- [x] [记忆类型与数据库事实](context-engineering/10-memory-types-and-database-facts.md)

## 阶段四：AI 应用与 AI UX

- [x] [AI 任务状态机](ai-ux/01-ai-task-state-machine.md)
- [x] [Streaming、Stop、Retry 与 Continue](ai-ux/02-stream-stop-retry-continue.md)
- [x] [断线、重复事件与半截 JSON](ai-ux/03-disconnect-duplicate-events-partial-json.md)
- [x] [流式 Markdown 与未闭合代码块](ai-ux/04-streaming-markdown-incomplete-code-blocks.md)
- [x] [Artifact、Diff、History 与 Restore](ai-ux/05-artifact-diff-history-restore.md)
- [x] [Citation、证据与无法回答](ai-ux/06-citations-evidence-no-answer.md)
- [x] [不确定性、人工确认与人工接管](ai-ux/07-uncertainty-confirmation-human-handoff.md)
- [x] [页面刷新后的长任务恢复](ai-ux/08-long-task-refresh-recovery.md)

## 阶段五：RAG

### 文档解析

- [x] [PDF、Word、Markdown、HTML、表格与扫描件解析](rag-parsing/01-document-formats-and-parsing.md)
- [x] [标题、页码、来源与原文定位](rag-parsing/02-structure-page-source-locators.md)
- [x] [重复页眉页脚的识别与清洗](rag-parsing/03-remove-repeated-headers-footers.md)
- [x] [按文件记录解析质量](rag-parsing/04-parsing-quality-by-file.md)

### Chunk

- [x] [固定、段落、标题、语义、滑动窗口与父子分块](rag-chunking/01-chunking-strategies.md)
- [x] [Chunk 长度、重叠与结构保留的比较](rag-chunking/02-length-overlap-structure.md)
- [x] [列表、表格与超长段落的分块](rag-chunking/03-lists-tables-long-paragraphs.md)
- [x] [文档版本、更新、删除与重新索引](rag-chunking/04-version-update-delete-reindex.md)

## 写作与维护要求

- 资料优先级：规范与官方文档 → 论文/模型卡 → 可复现实验 → 高质量工程资料。
- 每篇包含：是什么、为什么、关键特性、实际使用、常见错误与边界、补充知识、来源。
- 来源写直接链接和访问日期；模型行为不能只凭一次运行下结论。
