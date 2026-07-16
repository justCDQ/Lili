---
type: ai-note
stage: beginner
topic: training-inference-generalization
verified: 2026-07-16
tags: [ai, training, inference, loss, generalization, overfitting]
---

# 训练、推理、参数、损失、泛化与过拟合

## 是什么

训练是使用数据和优化算法调整模型参数，使目标损失下降的过程。参数是模型从数据中学到的数值。损失函数把预测与训练目标的差异转换为可优化数值。推理是在参数固定后，给模型输入并计算输出。泛化是模型在未见过但来自目标分布的数据上仍能完成任务。过拟合是训练数据表现改善而新数据表现不再改善或变差。

## 为什么需要

AI 应用开发者通常不训练基础模型，但这些概念决定如何解释模型版本、评测结果和微调风险。训练损失低不等于产品任务成功；推理输出也不是从数据库检索的固定答案。

## 关键特性

- 参数由训练更新，Prompt 不会永久修改基础模型参数。
- 损失是训练优化信号，不一定等于业务指标。不同损失会强调不同错误。
- 推理受输入、参数、解码策略和运行实现影响，可能具有随机性。
- 泛化只能相对于目标数据分布和任务定义讨论。
- 过拟合常见信号是训练损失继续下降，验证损失开始上升。
- 数据不代表真实场景、模型复杂度过高、重复调参和数据泄漏都会破坏泛化。

## 实际怎么使用

应用开发中应把供应商模型视为固定版本的外部依赖：

1. 使用独立评测集验证目标任务，不用公开 Benchmark 分数替代产品评测。
2. Prompt 或微调迭代只查看训练/开发集，保留最终测试集。
3. 同时记录总体指标和失败分组，例如语言、长度、权限、无答案与对抗输入。
4. 发现训练/开发分数提高而测试或线上变差时，检查过拟合和数据分布变化。
5. 基础模型更新后重新运行回归集，不假设新版本所有任务都更好。

## 常见错误与边界

- 把参数数量当作质量的充分条件；数据、架构、训练方法和任务适配同样重要。
- 用训练损失直接比较采用不同目标函数的模型。
- 把温度调低称为消除随机性；推理服务和模型仍可能存在非确定因素。
- 在同一测试集上反复改 Prompt，测试集实际上已经变成开发集。
- 认为过拟合只发生在模型训练；Prompt 也可能被过度针对少量样例优化。

## 补充知识

欠拟合表示模型连训练数据中的规律也没有充分学到。正则化、更多代表性数据、降低复杂度和早停可缓解过拟合，但选择应由验证曲线和真实任务决定。

## 来源

- [Google ML Crash Course](https://developers.google.com/machine-learning/crash-course/)（访问日期：2026-07-16）
- [Google ML：Loss](https://developers.google.com/machine-learning/crash-course/linear-regression/loss)（访问日期：2026-07-16）
- [Google ML：Overfitting](https://developers.google.com/machine-learning/crash-course/overfitting/overfitting)（访问日期：2026-07-16）
- [Google ML Glossary](https://developers.google.com/machine-learning/glossary/)（访问日期：2026-07-16）

