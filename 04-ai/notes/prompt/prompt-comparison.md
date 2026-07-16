---
type: ai-note
stage: junior
topic: prompt-comparison
verified: 2026-07-16
tags: [ai, prompt, experiment, comparison]
---

# 同一任务的 Prompt 版本对比

## 是什么

Prompt 对比是在相同模型、参数、数据集、Schema、工具和评分规则下运行多个 Prompt 版本，只把 Prompt 作为主要变量，比较任务指标、失败类型、延迟和成本。

## 为什么需要

凭少量手工输入判断容易受到选例和主观偏差影响。一个 Prompt 可能改善正常样例但破坏无答案或边界场景，必须用固定数据集逐项比较。

## 关键特性与规则

- 除 Prompt 外应固定模型、参数、工具、检索、Schema、样例和评分规则。
- 非确定输出需要多次 Trial，并报告分布而非只选择最好一次。
- 指标按正常、边界、无答案、高风险和语言等切片检查，避免总平均掩盖退化。
- Prompt 修改不得针对测试集逐题拟合；测试集与开发集应分离。
- 质量提升要同时报告 Token、延迟、成本和安全护栏变化。

## 实际步骤

1. 定义基线 Prompt 和目标，例如提高事实抽取准确率且不增加无依据字段。
2. 每个新版本只改变一个可说明的策略：增加示例、调整上下文、增加失败规则或改 Schema。
3. 固定完整模型 ID、参数、检索版本、工具和评测集。
4. 对非确定输出运行多次 Trial，保存全部结果。
5. 使用确定性检查、人工 Rubric 或已校准 Grader。
6. 按正常、边界、无答案、长输入、对抗和语言分组比较。
7. 记录回退条件；新版本未达到门槛不替换生产版本。

```text
v1：基线
v2：只增加失败行为
v3：在 v2 上只增加两个边界示例
指标：字段准确率、无答案准确率、Schema 通过率、P95 延迟、平均成本
```

## 常见错误与边界

- 同时换模型、参数和 Prompt，无法归因。
- 只展示成功案例，不保存失败和全部 Trial。
- 对测试集逐题修改 Prompt，形成测试泄漏。
- 只比较总平均分，忽略高风险分组退化。
- 新 Prompt 更长导致 Token 和延迟增加，却只报告质量。

## 补充知识

Prompt 对比结果只对指定模型和数据分布成立。供应商更新别名或业务输入变化后应重新评估。

## 来源

- [OpenAI：Evaluation Best Practices](https://developers.openai.com/api/docs/guides/evaluation-best-practices)（访问日期：2026-07-16）
- [Anthropic：Define Success Criteria and Build Evaluations](https://platform.claude.com/docs/en/test-and-evaluate/define-success)（访问日期：2026-07-16）
- [Anthropic：Prompt Engineering Overview](https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/overview)（访问日期：2026-07-16）
