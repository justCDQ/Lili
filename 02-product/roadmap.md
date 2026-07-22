# 产品 Roadmap：从产品小白到独立定义产品

这是一条适合个人执行的产品学习路线。用户访谈不是必修前提，学习者可以先通过公开证据、产品使用、功能拆解、数据和真实问题日志建立判断，再在有条件时把访谈作为可选验证手段。

产品能力的核心不是写文档，而是持续回答：为谁解决什么问题、为什么值得解决、如何验证、做什么与不做什么。

```text
用户价值 ∩ 商业价值 ∩ 技术可行性 = 可持续的产品机会
```

## 能力阶梯

| 阶段 | 能力目标 | 代表产出 |
| --- | --- | --- |
| 入门 | 看懂产品、功能、用户、场景和价值 | 产品定位卡 |
| 初级 | 能拆解功能并发现有证据的问题 | 功能拆解与问题清单 |
| 中级 | 能定义需求、MVP、指标和方案 | PRD、原型与验证报告 |
| 高级 | 能推动上线、实验、商业与迭代 | Roadmap、数据复盘和商业分析 |
| 专项 | 能处理 B 端或 AI 产品复杂性 | 权限流程或 AI 产品方案 |

---

## 阶段零：产品基础语言

### 基本概念

- [ ] 产品、项目、功能、服务和平台的区别。
- [ ] 用户、客户、使用者、购买者、决策者和管理员。
- [ ] 用户目标、使用场景、痛点、需求、方案和功能的区别。
- [ ] 用户价值、商业价值、实现成本和机会成本。
- [ ] 产品定位、价值主张、核心能力、辅助能力和产品边界。
- [ ] B2C、B2B、SaaS、平台、工具和内容产品的基本模式。
- [ ] 产品生命周期：发现、定义、设计、开发、发布、度量和迭代。

定位模板：

```text
我们为【目标用户】，
在【特定场景】下，
解决【关键问题】，
通过【核心方案】，
帮助用户获得【最终价值】。
```

必做：为五个不同类型的产品制作产品卡，记录定位、角色、核心任务、价值、收费方式、主要成本和风险。

验收：能清楚回答一个产品为谁服务、解决什么、为什么被选择、如何持续存在以及什么不属于它。

---

## 阶段一：产品观察、思考与拆解

不要只截图和罗列功能。每次拆解从一个具体任务或功能开始，追踪“问题 → 方案 → 交互 → 技术 → 指标 → 商业”的完整链条。

- [ ] 选择一个具体任务，明确入口、前置条件和完成标准。
- [ ] 还原主流程、分支、失败、权限和恢复路径。
- [ ] 推断目标用户、JTBD、使用频率和替代方案。
- [ ] 判断功能提供的用户价值以及对产品目标的作用。
- [ ] 分析默认值、限制、渐进披露和范围取舍。
- [ ] 推断可能的核心指标、守护指标和负面影响。
- [ ] 分析实现成本、运营成本、模型成本或数据依赖。
- [ ] 与两个竞品或非软件替代方案比较。
- [ ] 写出保留、删除、改进各一项，并说明证据与风险。

验收：连续完成 30 个功能拆解；结论能区分“看见的事实”“基于证据的推断”和“个人观点”。

---

## 阶段二：问题发现与证据收集

个人学习可以使用以下路径，不要求先找到访谈对象：

### 公开证据

- [ ] 应用商店评论、插件市场评价和公开评分变化。
- [ ] GitHub Issue、Discussion、公开 Roadmap 和 Release Notes。
- [ ] 帮助中心、FAQ、客服公开案例、状态页和已知问题。
- [ ] 公开社区中重复出现的问题、变通方案和迁移原因。
- [ ] 行业报告、公开数据、搜索趋势和政策变化。
- [ ] 竞品定价页、文档、演示、更新日志和下线功能。

### 主动观察与实验

- [ ] Dogfooding：把自己作为真实用户，记录任务、耗时、错误和绕路。
- [ ] 工作流观察：还原一个任务目前如何用表格、消息、纸笔或多个工具完成。
- [ ] 客服/销售/实施等二手材料：只在合法、脱敏且有权限时使用。
- [ ] Landing Page、原型、假门、问卷或等待名单等低成本验证。
- [ ] 搜索广告、内容点击、开源采用、下载和留存等行为信号。
- [ ] 已有产品的埋点、漏斗、搜索词和失败日志；没有数据时明确标注假设。

### 问题表达

```text
当【某类用户】处于【具体场景】时，
为了【完成任务或获得进展】，
目前使用【现有做法】，
但受到【有证据的问题】影响，
导致【时间、质量、成本或风险结果】。
```

- [ ] 写出 JTBD、当前工作流和替代方案。
- [ ] 记录问题频率、严重度、影响范围和证据可信度。
- [ ] 区分症状、根因、需求表达和解决方案请求。
- [ ] 使用多种来源交叉验证，不把单条评论当作普遍事实。
- [ ] 建立“问题证据库”，为每条结论保留来源和日期。

验收：面对一个想法，先提交问题陈述、证据、当前做法、替代方案和待验证假设，而不是直接列功能。

---

## 阶段三：需求分析与优先级

学习：用户问题、业务目标、产品目标、成功指标、范围、非目标、依赖、风险和上线策略。

Todo：

- [ ] 明确目标用户和场景。
- [ ] 将功能描述改写为用户问题。
- [ ] 定义业务目标和产品目标。
- [ ] 定义成功指标与守护指标。
- [ ] 明确功能范围和非目标。
- [ ] 识别技术、数据、合规和运营依赖。
- [ ] 使用 RICE 或 Value/Effort 排序。
- [ ] 记录为什么做、为什么不做。
- [ ] 识别最大产品风险。

需求模板：

```markdown
# 需求分析
## 背景
## 目标用户
## 用户场景
## 用户问题
## 业务目标
## 产品目标
## 成功指标
## 功能范围
## 非目标
## 依赖
## 风险
## 上线策略
```

验收：能有依据地拒绝或延后问题不清、范围过大、价值过低或风险未验证的需求。

---

## 阶段四：产品方案与 MVP

学习：方案发散、方案比较、核心假设、低成本验证、MVP 范围、成功与失败标准。

Todo：

- [ ] 同一问题提出至少三个方案。
- [ ] 比较用户价值、使用成本、开发成本和风险。
- [ ] 找到最关键、最不确定的假设。
- [ ] 设计一周内可以完成的验证。
- [ ] 明确 MVP 核心用户和核心场景。
- [ ] 删除与核心价值无关的功能。
- [ ] 定义继续投入、调整和停止条件。

验收：能把预计三个月的完整需求，拆成一周验证、两周原型、四周 MVP，再根据数据决定是否扩大。

---

## 阶段五：PRD 与产品表达

学习：目标、非目标、用户场景、流程、业务规则、异常、验收标准、依赖和上线策略。

Todo：

- [ ] 写清楚为什么做和不做什么。
- [ ] 使用图表示流程和状态。
- [ ] 写出可验证的用户故事。
- [ ] 使用 Given/When/Then 描述验收条件。
- [ ] 覆盖权限、冲突、失败和边界。
- [ ] 根据产品、设计、开发和测试调整文档内容。
- [ ] 使用 ADR 或决策记录保存重要变化。

验收：团队阅读后，对目标、范围、规则、异常和完成标准理解一致。

---

## 阶段六：产品数据与实验

学习：北极星指标、过程指标、质量指标、守护指标、漏斗、留存、Cohort、A/B Test、SQL。

Todo：

- [ ] 定义产品北极星指标。
- [ ] 为核心流程定义漏斗。
- [ ] 设计埋点事件和属性。
- [ ] 区分结果指标和过程指标。
- [ ] 按用户群、版本、来源分析。
- [ ] 识别虚荣指标。
- [ ] 建立产品指标看板。
- [ ] 根据数据提出新的产品假设。
- [ ] 理解相关性与因果关系的区别。

验收：上线后能回答有多少人使用、哪类用户使用、在哪流失、是否完成目标、是否持续使用，以及是否产生负面影响。

---

## 阶段七：市场、竞品和商业

学习：市场规模、用户类型、购买决策、直接竞品、间接替代、收费、成本、获客、留存、差异化和护城河。

Todo：

- [x] [分析三个直接竞品](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [分析两个替代方案](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [识别使用者、购买者和审批者](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [分析产品定价与收费方式](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [分析获客、交付、模型和基础设施成本](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [理解 MRR、ARR、ARPU、CAC、LTV、Churn、毛利](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [判断竞品哪些功能不应复制](notes/07-market-business/01-market-competition-business-system.md)。
- [x] [分析长期差异化和迁移成本](notes/07-market-business/01-market-competition-business-system.md)。

验收：能解释产品谁使用、谁付费、为什么付费、如何获客、主要成本、竞争优势和可能被谁替代。

---

## 阶段八：产品推进与持续迭代

学习：目标导向 Roadmap、产品机会、实验、项目、跨团队依赖、灰度、回滚、复盘。

Todo：

- [x] [Roadmap 围绕目标而不是功能列表](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。
- [x] [明确负责人、依赖和风险](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。
- [x] [组织需求、设计和上线评审](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。
- [x] [制定灰度、回滚和通知策略](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。
- [x] [上线后观察指标和用户反馈](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。
- [x] [完成无责复盘和行动项](notes/08-delivery-iteration/01-outcome-roadmap-release-postmortem.md)。

验收：能够推动一个中等复杂度功能完成发现、定义、设计、开发、灰度、数据观察和复盘的闭环。

---

## 阶段九：B 端产品专项

- [x] [用户角色与权限矩阵](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [核心业务对象与关系](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [状态机和审批流程](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [数据责任、操作日志和合规](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [标准化、配置与定制的边界](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [系统迁移、数据导入和实施流程](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。
- [x] [使用者、管理员、购买者不同的场景](notes/09-b2b-products/01-enterprise-domain-permission-workflow-migration.md)。

---

## 阶段十：AI 产品专项

- [x] [判断 AI 是否真的必要](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [定义 AI 输出质量和不可接受错误](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [设计引用、证据和人工确认](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [评估模型成本、延迟和隐私](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [设计失败降级和人工接管](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [建立真实评估数据集](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。
- [x] [平衡质量、体验、成本和商业模式](notes/10-ai-products/01-ai-product-quality-evaluation-operations.md)。

---

## 学习资源

书籍：启示录、用户体验要素、精益创业、用户故事地图、精益数据分析、商业模式新生代、好战略坏战略。

网站：SVPG、Product Talk、Mind the Product、Lenny’s Newsletter、Reforge、Y Combinator Library、Stratechery、Amplitude Blog、Mixpanel Blog。


---
