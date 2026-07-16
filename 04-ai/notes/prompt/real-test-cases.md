---
type: ai-note
stage: junior
topic: real-test-cases
verified: 2026-07-16
tags: [ai, prompt, evaluation, test-cases]
---

# 真实测试样例

## 是什么

真实测试样例来自目标产品实际任务分布，包含输入、上下文、期望行为、评分规则、来源和风险标签。它可以由经过授权和脱敏的真实请求、领域专家编写的代表性任务或根据真实失败模式构造的合成数据组成。

## 为什么需要

通用 Benchmark 和开发者自拟演示通常缺少产品特有术语、脏数据、权限、无答案和失败损失。真实样例使 Prompt 优化面向产品任务而不是演示效果。

## 关键特性

- 代表目标语言、长度、用户类型、设备/渠道和时间分布。
- 覆盖正常、边界、无答案、冲突、权限和对抗输入。
- 保存来源、许可、脱敏方式和数据保留要求。
- 成功标准在运行前定义，不能看完输出再改变答案。
- 高风险样例单独设门槛，不被大量简单样例稀释。

## 实际怎么使用

```json
{
  "id": "support-no-answer-017",
  "input": "...",
  "context_refs": ["kb-v12:doc-8"],
  "expected_behavior": "abstain_and_request_information",
  "rubric": ["no_unsupported_claim", "mentions_missing_order_id"],
  "tags": ["no-answer", "zh-CN", "high-risk"],
  "source": "authorized-production-failure",
  "privacy": "de-identified"
}
```

线上失败先完成授权、脱敏和去重，再进入回归集。开发集用于迭代，保留测试集用于发布判断。

## 常见错误与边界

- 把真实用户原文直接提交仓库或发送到新服务。
- 只收集投诉和失败，数据分布与真实流量严重偏离。
- 合成样例全部由待评模型生成，语言模式过于一致。
- 期望答案过度具体，把其他正确表达判错。
- 样例缺版本，知识库或政策变化后参考答案已过时。

## 补充知识

数据集应支持切片分析和版本差异。删除请求必须能追踪到包含该数据的原始记录、派生集和日志。

## 来源

- [OpenAI：Evaluation Best Practices](https://developers.openai.com/api/docs/guides/evaluation-best-practices)（访问日期：2026-07-16）
- [Anthropic：Demystifying Evals for AI Agents](https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents)（访问日期：2026-07-16）
- [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework)（访问日期：2026-07-16）

