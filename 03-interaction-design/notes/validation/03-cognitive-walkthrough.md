# Cognitive Walkthrough：认知走查

## 是什么

认知走查由评估者按任务逐步模拟首次或低经验用户的判断，检查目标、操作发现、动作与目标的关联及反馈理解。它评估可学性风险，不证明真实用户行为。

## 为什么需要

验证的价值在于减少未经检查的设计假设。明确这一方法的适用问题、输入条件和输出，可以避免把单次观察、评估者判断或工具结果扩张为普遍结论，并为改版前后复测提供一致基线。

## 关键特性或原则

### 准备

定义目标用户已有知识、真实任务、入口、前置条件、正确动作序列和界面版本。每一步都必须能在原型或产品中实际执行。

### 四个核心问题

1. 用户会尝试实现当前步骤的正确效果吗？
2. 用户会注意到正确操作可用吗？
3. 用户会把该操作与想要的效果联系起来吗？
4. 执行后，用户能从反馈判断自己取得进展吗？

## 实际怎么使用

### 执行与输出

逐步记录答案、界面证据、失败原因、严重度和修改建议。失败原因应具体到术语、可见性、映射、先验知识或反馈，不写“用户不懂”。改版后按相同人物知识和任务复走。

## 常见错误与边界

- 评估者利用自己已知的入口和结果替用户作答。
- 只走主路径，遗漏错误恢复。
- 把任何不喜欢的设计都归入认知走查。
- 不适合单独评估长期效率、满意度或真实发生率；需结合任务日志、可用性测试或行为数据。

## 补充知识

单一验证方法通常只能覆盖部分问题。原则性检查适合发现风险，可运行测试适合确认行为，日志与埋点适合描述已发生事件；需要判断原因或推广范围时，应组合不同证据并记录反证、样本限制和版本时效。

## 来源

- [Interaction Design Foundation：Cognitive Walkthrough](https://www.interaction-design.org/literature/topics/cognitive-walkthrough)（访问日期：2026-07-16）
- [NN/g：Cognitive Walkthroughs](https://www.nngroup.com/articles/cognitive-walkthroughs/)（访问日期：2026-07-16）
- [GOV.UK Service Manual：Quality assurance](https://www.gov.uk/service-manual/technology/quality-assurance-testing-your-service-regularly)（访问日期：2026-07-16）
