# Dogfooding 与任务日志

## 是什么

Dogfooding 是在真实工作中持续使用自己的产品；任务日志按时间记录目标、条件、动作、等待、误操作、绕路、失败和恢复。它能发现长期、跨会话和真实数据下的问题，但内部人员的权限、知识和动机与目标用户不同。

## 为什么需要

验证的价值在于减少未经检查的设计假设。明确这一方法的适用问题、输入条件和输出，可以避免把单次观察、评估者判断或工具结果扩张为普遍结论，并为改版前后复测提供一致基线。

## 关键特性或原则

### 日志字段

- 日期、版本、平台、角色和任务目标。
- 起止时间、等待时间、步骤和中断。
- 预期、实际结果、误操作、回退和求助。
- 错误文案、截图、恢复方式与数据影响。
- 事实、推断、严重度和后续验证。

## 实际怎么使用

### 执行方法

1. 选高频和高风险真实任务，连续记录固定周期。
2. 不为记录而制造虚假任务；异常测试另行标识。
3. 每周聚类重复问题，区分学习期错误与稳定问题。
4. 与公开反馈、Issue、日志或任务指标交叉验证。
5. 将修复前后使用同一任务复测。

## 常见错误与边界

- 团队因熟悉产品而忽略发现性问题。
- 内部账户权限过高，无法覆盖普通用户路径。
- 只记“感受”，没有时间、步骤和复现条件。
- 把内部高频工作流直接等同于外部需求。
- 日志中的业务数据和同事信息需最小化、脱敏并限制访问。

## 补充知识

单一验证方法通常只能覆盖部分问题。原则性检查适合发现风险，可运行测试适合确认行为，日志与埋点适合描述已发生事件；需要判断原因或推广范围时，应组合不同证据并记录反证、样本限制和版本时效。

## 来源

- [NN/g：Diary Studies](https://www.nngroup.com/articles/diary-studies/)（访问日期：2026-07-16）
- [GOV.UK Service Manual：Quality assurance](https://www.gov.uk/service-manual/technology/quality-assurance-testing-your-service-regularly)（访问日期：2026-07-16）
- [ICO：Data minimisation](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/data-protection-principles/a-guide-to-the-data-protection-principles/data-minimisation/)（访问日期：2026-07-16）
