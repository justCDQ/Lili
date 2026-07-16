---
type: ai-note
stage: beginner
topic: fixed-cases-comparison
verified: 2026-07-16
tags: [ai, evaluation, dataset, prompt, model]
---

# 固定样例与模型、Prompt 对比

## 是什么

固定样例是一组版本化的输入、期望行为和评分标准。对比实验在相同样例和配置下改变一个因素，例如 Prompt 或模型，重复运行并比较质量、延迟、Token、成本和失败类型。

## 为什么需要

模型输出具有非确定性，单次演示不能代表稳定能力。没有固定样例时，测试者会无意选择更适合新方案的问题，导致结果不可比较。先定义成功标准还能暴露产品要求是否足够具体。

## 样例结构

```json
{
  "id": "extract-001",
  "input": "原始输入",
  "expected": {
    "required_facts": ["..."],
    "forbidden_claims": ["..."]
  },
  "tags": ["normal", "zh-CN"],
  "source": "synthetic",
  "grader": "schema-and-facts-v1"
}
```

初始五条至少覆盖：正常、边界、输入不足/无答案、格式压力和对抗/不可信内容。真实上线后，把经过脱敏和授权的失败案例加入回归集。

## 实际怎么使用

1. 在修改 Prompt 或选择模型之前写样例和通过标准。
2. 固定 Schema、工具、检索数据和采样参数，只改变一个变量。
3. 每个样例必要时运行多次，保存每次 Trial，而不是只挑最好结果。
4. 优先使用确定性 Grader：Schema、正则、单元测试、数据库最终状态。
5. 语义质量使用有明确 Rubric 的人工或模型 Grader，并定期人工校准。
6. 同时输出逐样例结果和分组汇总，避免平均分掩盖关键场景失败。

```text
方案 A：Prompt v2 + Model X
方案 B：Prompt v3 + Model X
固定：数据集 v4、Schema v2、温度、工具、检索索引
比较：通过率、关键失败、P95 延迟、平均/总成本
```

## 常见错误与边界

- 修改 Prompt、模型和检索后把变化全部归因于模型。
- 只用“看起来更好”的主观判断，不写 Rubric。
- 让同一模型生成参考答案又评分，且没有人工校准。
- 反复针对测试集调 Prompt，最终分数提高但真实请求没有改善。
- 只汇总平均分，不单独检查安全、权限和高损失样例。

## 补充知识

Agent 评估除最终文本外还要检查 Tool 调用轨迹和环境最终状态。评估集需要能力集和回归集：能力集测试当前还做不到但希望跟踪的能力，回归集防止已解决问题再次出现。

## 来源

- [OpenAI API：Evals](https://platform.openai.com/docs/api-reference/evals)（访问日期：2026-07-16）
- [OpenAI API：Graders](https://platform.openai.com/docs/api-reference/graders)（访问日期：2026-07-16）
- [Anthropic：Demystifying evals for AI agents](https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents)（访问日期：2026-07-16）

