---
type: ai-note
stage: junior
topic: format-vs-content-errors
verified: 2026-07-16
tags: [ai, prompt, errors, validation]
---

# 格式错误与内容错误

## 是什么

格式错误表示输出不符合协议或 Schema，例如非法 JSON、缺字段、类型错误。内容错误表示结构正确但事实、推理、分类、引用或业务判断错误。两类错误需要不同检测和修复方式。

## 为什么需要

只统计 Schema 通过率会高估系统质量；把内容错误当格式问题反复改 JSON 指令也不会改善事实。明确分类才能定位 Prompt、模型、数据、检索或代码责任。

## 关键特性与规则

- 格式校验是确定性的协议检查；内容校验依赖事实来源、业务规则或评分标准。
- 拒绝、取消、截断和安全阻断是完成状态，不应强行归为 JSON 格式错误。
- Schema 通过只证明结构满足约束，不证明事实、引用和业务判断正确。
- 错误分类必须互斥程度足够、版本稳定，并保留原始输出与检测证据。

## 实际检测

```text
1. 传输/完成状态
2. JSON 解析
3. Schema 校验
4. 业务规则校验
5. 事实与引用校验
6. 任务 Rubric / 最终状态
```

格式错误使用解析器和 Validator 确定性检测。内容错误使用参考答案、数据库事实、引用对齐、单元测试、人工 Rubric或经校准的模型 Grader。

## 修复路径

- 格式：Structured Output、简化 Schema、运行时校验、处理不完整状态。
- 内容缺资料：补 Context、RAG 或 Tool。
- 确定性规则错误：移到代码。
- 模型能力不足：更换模型或在评测支持下微调。
- 需求不清：重写成功标准和样例。

## 常见错误与边界

- 合法 JSON 就标记成功。
- Grader 只读最终文本，不检查 Tool 最终状态。
- 自动修复无效字段后不记录原始失败。
- 把模型拒绝或安全阻断算作格式错误。
- 内容错误没有细分为事实、遗漏、相关性、引用和安全。

## 补充知识

失败分类本身需要版本化。新增类别后历史数据不应被无说明地重新解释。

## 来源

- [JSON Schema：Getting Started](https://json-schema.org/learn/getting-started-step-by-step)（访问日期：2026-07-16）
- [OpenAI：Evaluation Best Practices](https://developers.openai.com/api/docs/guides/evaluation-best-practices)（访问日期：2026-07-16）
- [Anthropic：Increase Output Consistency](https://platform.claude.com/docs/en/test-and-evaluate/strengthen-guardrails/increase-consistency)（访问日期：2026-07-16）
