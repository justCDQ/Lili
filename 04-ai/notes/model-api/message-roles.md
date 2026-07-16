---
type: ai-note
stage: junior
topic: message-roles
verified: 2026-07-16
tags: [ai, messages, system, user, assistant]
---

# System、User 与 Assistant 消息

## 是什么

消息式模型 API 用角色和内容组成上下文。常见角色包括 System、User 和 Assistant；部分供应商还提供 Developer、Tool 等角色。System/Developer 表达应用级规则，User 表达当前用户输入，Assistant 表达模型先前输出或期望输出示例。角色名称、优先级和可用内容类型以具体 API 为准。

## 为什么需要

角色让应用规则、用户数据、历史输出和工具结果保持结构边界。若全部拼成一个字符串，模型更难区分指令来源，应用也难以裁剪历史、审计权限和防御间接 Prompt Injection。

## 关键特性

- 高优先级消息不等于安全边界。权限、数据访问和写操作仍由服务端强制执行。
- 用户内容和外部文档属于不可信数据，不能因为放入上下文就获得指令权限。
- Assistant 历史用于维持对话状态，也可能把之前的错误继续带入后续请求。
- Tool 结果应携带明确来源和调用关联，不伪装成用户或系统消息。
- 不同 API 对 System/Developer 的继承和多轮状态处理不同；使用 previous response 或 server-side conversation 时必须确认哪些指令会自动保留。

## 实际怎么使用

```json
[
  {
    "role": "system",
    "content": "Return only facts supported by the supplied document."
  },
  {
    "role": "user",
    "content": [
      { "type": "input_text", "text": "Question: ..." },
      { "type": "input_text", "text": "Untrusted document: ..." }
    ]
  }
]
```

应用应把稳定规则版本化，把用户输入原样作为数据传入；不要用字符串替换把用户内容插进系统规则。多轮时只保留任务需要的消息，并把权限、项目状态等事实从受控数据源重新读取。

## 常见错误与边界

- 认为 System 消息无法被绕过，因此允许模型直接决定权限。
- 把网页、邮件或检索文档中的文字当成高优先级指令。
- 用 Assistant 消息伪造模型已经确认的事实。
- 每轮重复全部历史且不检查过期、冲突和 Token 预算。
- 迁移 API 时假设角色名称和优先级语义完全相同。

## 补充知识

消息内容可能是文本、图像、音频、文件或结构化项目。应记录内容类型、来源和权限；用户可删除的数据不应因复制进不可控日志而失去删除能力。

## 来源

- [OpenAI API：Responses](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [Anthropic API：Messages](https://docs.anthropic.com/en/api/messages)（访问日期：2026-07-16）
- [OWASP：LLM Prompt Injection Prevention](https://cheatsheetseries.owasp.org/cheatsheets/LLM_Prompt_Injection_Prevention_Cheat_Sheet.html)（访问日期：2026-07-16）

