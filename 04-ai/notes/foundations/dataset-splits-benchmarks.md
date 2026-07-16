---
type: ai-note
stage: beginner
topic: dataset-splits-benchmarks
verified: 2026-07-16
tags: [ai, dataset, validation, test, benchmark]
---

# 训练集、验证集、测试集与 Benchmark

## 是什么

训练集用于更新模型参数；验证集用于选择模型、Prompt、超参数和停止条件；测试集只用于对最终方案做相对独立的估计。Benchmark 是具有任务、数据、运行规则和评分方法的标准化评测，用于在约定条件下比较系统。

## 为什么需要

如果开发者根据测试结果不断修改方案，测试集信息会进入决策，最终分数会高估真实泛化能力。公开 Benchmark 能提供外部参照，但通常不能覆盖具体产品的用户、数据、权限、成本和失败损失。

## 关键特性

- 数据划分应避免同一文档、同一用户、近重复内容或时间上相关样例跨集合泄漏。
- 时间敏感任务常按时间切分，确保测试数据发生在训练/开发数据之后。
- 三个集合应接近目标线上分布；若分布不同，指标解释必须注明。
- Benchmark 分数只有在相同数据版本、Prompt、工具、采样、运行次数和评分器下才可直接比较。
- 污染是模型训练数据包含测试题或答案，导致分数不能代表新任务能力。

## 实际怎么使用

AI 应用可以采用：

```text
development.jsonl  # 日常 Prompt、RAG 和流程调试
regression.jsonl   # 已修复线上/实验失败，持续扩充
test.jsonl         # 版本发布前评估，限制查看与调参
```

每条样例记录 ID、来源、日期、许可、任务标签、输入、期望、评分器版本和风险等级。报告同时给出样例数量、分组分布、置信区间或多次 Trial 结果，不只给一个总分。

## 常见错误与边界

- 随机按行切分包含同一文档多个 Chunk 的数据，造成内容泄漏。
- 用公开 Benchmark 排名直接决定产品模型，不测真实输入和成本。
- 在测试失败后逐条修改 Prompt，再继续报告同一测试集分数。
- 测试集只有正常样例，没有无答案、权限、边界和对抗样例。
- 更换 Judge 模型或 Rubric 后直接与旧分数比较。

## 补充知识

离线评测无法替代线上监控。线上数据受真实流量、界面、延迟、用户行为和系统集成影响；发布后仍需观察任务完成、人工接管、投诉、安全和费用指标。

## 来源

- [Google ML：Datasets, Generalization and Overfitting](https://developers.google.com/machine-learning/crash-course/overfitting)（访问日期：2026-07-16）
- [Google ML Glossary](https://developers.google.com/machine-learning/glossary/)（访问日期：2026-07-16）
- [OpenAI API：Evals](https://platform.openai.com/docs/api-reference/evals)（访问日期：2026-07-16）
- [Anthropic：Demystifying evals for AI agents](https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents)（访问日期：2026-07-16）

