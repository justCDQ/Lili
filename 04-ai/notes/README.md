# AI 入门与初级模块笔记

覆盖 [AI Roadmap](../roadmap.md) 的阶段零、阶段一和阶段二。每个路线图知识点对应一篇独立笔记；实验性、时效性结论必须记录模型、版本和日期。

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

## 写作与维护要求

- 资料优先级：规范与官方文档 → 论文/模型卡 → 可复现实验 → 高质量工程资料。
- 每篇包含：是什么、为什么、关键特性、实际使用、常见错误与边界、补充知识、来源。
- 来源写直接链接和访问日期；模型行为不能只凭一次运行下结论。
