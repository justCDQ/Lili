---
type: ai-note
stage: beginner
topic: tokenization-context-cost
verified: 2026-07-16
tags: [ai, tokenization, context-window, cost]
---

# Tokenization、Context Window 与输入输出成本

## 是什么

Tokenization 把文本或其他输入转换为模型处理的离散 Token ID。Token 可能是字、子词、标点、空格片段或字节组合，不等于字符或单词。Context Window 是一次推理能处理的 Token 总范围，通常包括系统指令、消息历史、工具定义、检索内容、当前输入和已生成内容。输入输出成本是供应商根据各类 Token 或模态用量计算的费用。

## 为什么需要

Token 数影响请求能否进入上下文、推理延迟和费用。字符数相同的不同语言、空白、代码和编码可能产生不同 Token 数。上下文过长还可能降低关键信息被正确使用的概率。

## 关键特性

- Tokenizer 属于具体模型或模型族；不能用一个模型的 Token 估算器精确计算另一个模型。
- 上下文上限不是全部可用于输入，输出预算也占用容量，具体规则以 API 为准。
- 工具 Schema、图片、音频、缓存 Token 和推理 Token可能有独立统计或计费方式。
- Prompt caching 可降低重复前缀的成本或延迟，但有命中条件、保留策略和隐私边界。
- 截断策略可能直接报错，也可能删除早期内容；自动截断必须显式理解。

## 实际怎么使用

建立 Token Budget：

```text
系统/安全规则：固定上限
工具定义：按任务选择，不全量注入
历史：保留关键事实，旧内容摘要
检索：按相关性、权限和新鲜度筛选
当前输入：完整保留或先明确拒绝过长输入
输出：为任务设置最大值
安全余量：处理估算误差和协议开销
```

在请求前使用官方 Tokenizer 或供应商计数接口估算；响应后以 API Usage 为准。日志按组成部分记录 Token，才能定位增长来源。

## 常见错误与边界

- 使用“一个 Token 约等于几个字”的经验值做硬限制。
- 每轮把全部历史、全部文档和全部工具重复发送。
- 只限制输入，不限制最大输出和 Agent 总步骤。
- 达到上限时直接从最早消息删除，丢失安全约束、用户决策或任务状态。
- 以为更长上下文必然提高质量，没有用评测集比较相关性和干扰。

## 补充知识

上下文管理应同时考虑权限、来源可信度、冲突和时间有效性。压缩与摘要会损失信息，重要事实应以结构化状态或数据库事实保存，而不是仅依赖对话摘要。

## 来源

- [Hugging Face LLM Course：Tokenizers](https://huggingface.co/learn/llm-course/en/chapter2/4)（访问日期：2026-07-16）
- [Google ML：Introduction to Large Language Models](https://developers.google.com/machine-learning/resources/intro-llms)（访问日期：2026-07-16）
- [OpenAI API：Responses Usage](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [Anthropic Docs：Token Counting](https://docs.anthropic.com/en/docs/build-with-claude/token-counting)（访问日期：2026-07-16）

