---
type: ai-note
stage: junior
topic: core-api-capabilities
verified: 2026-07-16
tags: [ai, multi-turn, streaming, structured-output, tool-calling, multimodal]
---

# 多轮、Streaming、Structured Output、Tool Calling 与多模态

## 是什么

- 多轮：后续请求包含或引用先前交互状态。
- Streaming：生成过程中以事件增量返回结果。
- Structured Output：按 Schema 返回机器可处理的数据。
- Tool Calling：模型选择工具并生成参数，应用验证并执行，再把结果返回模型。
- 多模态：请求或响应包含文本以外的图像、音频、视频或文件内容。

## 为什么需要

这些能力分别解决状态延续、响应等待、接口稳定、外部行动和非文本输入。它们不是同一层能力，组合后需要明确状态机和失败恢复。

## 关键特性

### 多轮

对话历史只是上下文，不是可靠数据库。长期事实、权限和任务状态应持久化为结构化数据。历史要处理 Token 预算、过期信息和删除要求。

### Streaming

事件可能包含文本、工具参数、引用、完成和错误。只有完成事件后才能确认整体状态，流中断可能缺少最终 Usage。

### Structured Output

Schema 约束结构，不保证语义正确。供应商支持的 JSON Schema 子集可能不同，仍需运行时与业务校验。

### Tool Calling

模型提出调用，不拥有执行权限。服务端检查工具名、参数、身份、资源权限、幂等和风险；写操作按影响请求用户确认。

### 多模态

输入需要记录 MIME 类型、尺寸、分辨率/采样率、顺序和数据保留策略。模型可能遗漏细节；OCR、音频转写和视觉理解都要建立独立评测。

## 实际怎么使用

组合状态机：

```text
Preparing → Streaming
Streaming → ToolRequested → Validating → WaitingApproval
WaitingApproval → ToolRunning → Streaming
任意状态 → Failed / Cancelled
Streaming → Completed / Incomplete
```

每个事件保存 `response_id`、`item_id`、序号和类型。Tool 结果与原请求关联；页面刷新后从持久化状态恢复，而不是重新执行副作用。

## 常见错误与边界

- 把多轮历史当永久记忆，未处理删除和冲突。
- 在 Tool 参数仍流式生成时提前执行。
- 结构化输出通过 Schema 后直接写生产数据。
- 多模态文件只检查扩展名，不检查实际类型、大小、恶意内容和权限。
- 对不同能力使用同一错误提示，无法判断是解析、模型、工具还是网络失败。

## 补充知识

上线前应分别测试能力，再测试组合：单独验证 Streaming、Schema、Tool 和多模态，最后覆盖中断、重复事件、部分成功和人工接管。

## 来源

- [OpenAI API：Responses](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [OpenAI API：Streaming Events](https://platform.openai.com/docs/api-reference/responses-streaming)（访问日期：2026-07-16）
- [Anthropic Docs：Tool Use](https://docs.anthropic.com/en/docs/agents-and-tools/tool-use/overview)（访问日期：2026-07-16）
- [MCP Specification：Tools](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)（访问日期：2026-07-16）

