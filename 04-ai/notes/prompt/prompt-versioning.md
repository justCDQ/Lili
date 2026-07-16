---
type: ai-note
stage: junior
topic: prompt-versioning
verified: 2026-07-16
tags: [ai, prompt, versioning, release]
---

# Prompt 版本化

## 是什么

Prompt 版本化为 Prompt 模板、示例、Schema、变量和发布状态分配不可变版本，使一次请求能够追溯到确切配置。版本可以存在 Git，也可以存在带发布历史的 Prompt Registry。

## 为什么需要

Prompt 是生产行为的一部分。直接覆盖字符串会导致问题无法复现、灰度无法比较、回滚没有目标。仅保存当前 Prompt 也无法解释历史输出。

## 关键特性

- 版本不可变；修改产生新版本。
- 记录变更原因、作者、评测结果、兼容模型、发布日期和回退版本。
- Prompt、Schema、工具和数据集分别版本化，并在发布记录中组合。
- 环境别名如 `production` 可移动，但请求日志保存解析后的具体版本。

## 实际怎么使用

```yaml
id: contract-extract
version: 4
schema: contract-v2
compatible_models: [model-id-a]
change: add explicit missing-field behavior
evaluation: eval-run-2026-07-16-04
rollback: 3
```

发布流程：代码评审 → 固定评测 → 安全检查 → 小流量灰度 → 指标比较 → 扩大或回滚。Prompt 中包含的示例也属于版本内容。

## 常见错误与边界

- 用日期或“final-v2”文件名但没有不可变内容哈希和发布记录。
- 日志只写 `production`，别名移动后无法还原。
- Prompt 更新不运行回归集。
- Schema 已变化但 Prompt 版本未关联，输出失败原因不清。
- 将真实敏感样例直接存入 Prompt Registry。

## 补充知识

版本号表达的是配置变化，不承诺质量单调提升。是否升级由评测和发布策略决定。

## 来源

- [OpenAI：Prompt Engineering](https://developers.openai.com/api/docs/guides/prompt-engineering)（访问日期：2026-07-16）
- [OpenAI API：Responses Prompt Version](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [Semantic Versioning](https://semver.org/)（访问日期：2026-07-16）

