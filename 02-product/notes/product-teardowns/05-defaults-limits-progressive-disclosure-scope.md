# 默认值、限制、渐进披露与范围取舍

## 是什么

默认值是在用户未选择时生效的设置；限制是容量、权限、速率、格式或业务规则边界；渐进披露先显示当前任务必需内容，再按需要展示高级能力；范围取舍决定本次覆盖与明确不覆盖的场景。

## 为什么需要

默认值影响多数未修改设置的用户，也可能改变隐私和费用。限制控制可靠性、成本及滥用风险。渐进披露降低初始复杂度，但隐藏关键状态会降低可发现性。范围取舍保证核心任务完整。

## 关键特性或判断规则

- 默认值必须说明适用人群、修改方式、追溯影响和安全后果。
- 限制必须在触发前可发现，并提供当前用量、重置时间和恢复路径。
- 渐进披露只隐藏低频或高级细节，不能隐藏风险、费用或任务状态。
- 范围取舍以核心任务闭环为单位，不能只交付无结果的前半段。

## 实际怎么使用

记录每个默认值、修改位置、受影响对象、是否追溯生效；记录每项限制的数值、计量窗口、达到限制时行为和升级路径；列出首屏、二级入口和隐藏条件；用“核心场景/边界场景/非目标”解释范围。

## 常见错误与边界

- 隐私、安全和付费默认值要保守且透明。
- 限制必须在操作前可发现，错误后给出恢复方式。
- 高频核心操作不应因“界面简洁”被深藏。
- 免费套餐限制是商业机制，也会产生用户和运营成本。
- 观察具体产品时标日期，不能把界面行为推断为内部技术原因。

## 补充知识

默认值具有路径依赖：用户往往不修改它。对隐私、安全和付费有影响的默认应更保守；套餐限制还要区分技术极限、成本控制和商业包装。

## 来源

- [UK ICO：Data protection by design and default](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/accountability-and-governance/guide-to-accountability-and-governance/accountability-and-governance/data-protection-by-design-and-default/)（访问日期：2026-07-16）
- [GitHub Docs：Rate limits](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api)（访问日期：2026-07-16）
- [GOV.UK Design System：Details component](https://design-system.service.gov.uk/components/details/)（访问日期：2026-07-16）
