# 交互设计入门与初级学习笔记

本目录覆盖 [交互设计 Roadmap](../roadmap.md) 的阶段零、阶段一和阶段二。每篇对应一个独立路线图条目；日常产品任务分析仍放在 [`daily/`](../daily/README.md)。

## 阅读顺序

1. `foundations/`：建立任务、界面结构、流程、状态和交付物的共同语言。
2. `principles/`：掌握分析设计决策所需的基础原则及适用边界。
3. `deconstruction/`：对真实任务进行可复现拆解、比较和重新设计。
4. `validation/`：使用公开证据、走查、审计与原型验证方案。

## 阶段零：交互设计基础

### 基本语言

- [用户目标、任务、场景、入口、路径与完成标准](foundations/01-user-goal-task-context-path.md)
- [页面、弹窗、抽屉与浮层](foundations/02-page-dialog-drawer-popover.md)
- [导航、表单、列表与详情](foundations/03-navigation-form-list-detail.md)
- [信息架构、用户流程、任务流程、状态与反馈](foundations/04-ia-user-flow-task-flow-state-feedback.md)
- [线框图、原型保真度与设计规范](foundations/05-wireframe-prototype-fidelity-spec.md)
- [区分产品、交互、视觉与工程问题](foundations/06-product-interaction-visual-engineering-problems.md)

### 基础原则

- [心智模型、概念模型与现实世界映射](principles/01-mental-conceptual-real-world-mapping.md)
- [可见性、可供性、反馈、一致性与约束](principles/02-visibility-affordance-feedback-consistency-constraints.md)
- [识别优于回忆、认知负荷与信息分组](principles/03-recognition-cognitive-load-grouping.md)
- [用户控制、撤销、容错、错误预防与恢复](principles/04-control-undo-tolerance-error-prevention-recovery.md)
- [渐进披露、默认值、熟悉性与学习成本](principles/05-progressive-disclosure-defaults-familiarity-learning.md)
- [Fitts 定律与 Hick 定律的适用边界](principles/06-fitts-hick-laws-boundaries.md)

## 阶段一：思考、产品拆解与重新设计

- [产品拆解的上下文记录](deconstruction/01-record-context.md)
- [还原入口、前置条件、主流程、分支、异常与完成状态](deconstruction/02-reconstruct-flow.md)
- [初始、空、加载、成功、失败、权限、离线与冲突状态清单](deconstruction/03-state-inventory.md)
- [信息层级、操作层级、默认值、文案与渐进披露分析](deconstruction/04-hierarchy-defaults-copy-disclosure.md)
- [识别交互模式及适用条件](deconstruction/05-identify-patterns.md)
- [用启发式原则分析并区分事实与推断](deconstruction/06-heuristics-facts-inferences.md)
- [与同类产品或现实工作流比较](deconstruction/07-compare-product-workflow.md)
- [重新设计流程与线框图](deconstruction/08-redesign-flow-wireframe.md)
- [评估改动收益、代价、风险与工程注意事项](deconstruction/09-benefit-cost-risk-engineering.md)
- [实现关键交互并验证键盘、响应式与状态切换](deconstruction/10-implement-keyboard-responsive-states.md)

## 阶段二：个人可执行的研究与验证

### 无需访谈的证据来源

- [用产品评论、公开 Issue、帮助文档、更新日志和社区问题取证](validation/01-public-reviews-issues-docs-changelogs.md)
- [竞品、历史版本与跨平台任务对比](validation/02-competitive-version-platform-comparison.md)
- [Cognitive Walkthrough：认知走查](validation/03-cognitive-walkthrough.md)
- [Heuristic Evaluation：启发式评估](validation/04-heuristic-evaluation.md)
- [Accessibility Audit：无障碍审计](validation/05-accessibility-audit.md)
- [状态审计：空、错、慢、断、冲突与权限](validation/06-state-audit.md)
- [Dogfooding 与任务日志](validation/07-dogfooding-task-log.md)
- [使用埋点、热图、搜索词、错误日志与客服材料](validation/08-analytics-heatmaps-search-errors-support.md)

### 原型验证

- [为原型定义真实任务、成功条件和观察项](validation/09-prototype-task-success-observations.md)
- [原型自测中的认知走查](validation/10-self-test-cognitive-walkthrough.md)
- [用键盘、窄屏、慢网、空数据和错误响应制造边界](validation/11-manufacture-boundary-conditions.md)
- [比较改版前后的步骤、完成率、时间、错误与认知成本](validation/12-compare-before-after.md)
- [可选的可用性测试：增强证据而非学习准入](validation/13-optional-usability-testing.md)
- [根据证据迭代并维护假设清单](validation/14-evidence-iteration-assumption-log.md)

## 覆盖表

| Roadmap 范围 | Roadmap 条目 | 笔记数 | 覆盖状态 |
| --- | ---: | ---: | --- |
| 阶段零 · 基本语言 | 5 | 6 | 完成；界面容器与内容/任务结构拆成两篇 |
| 阶段零 · 基础原则 | 6 | 6 | 完成 |
| 阶段一 | 10 | 10 | 完成 |
| 阶段二 · 无需访谈的证据来源 | 8 | 8 | 完成 |
| 阶段二 · 原型验证 | 6 | 6 | 完成 |
| **合计** | **35** | **36** | **全部覆盖** |

## 维护规则

- 新笔记必须说明是什么、为什么、实际方法、检查或步骤、常见错误与边界、补充知识及来源。
- 来源优先使用 W3C/WAI、平台官方 HIG、公共服务设计规范和方法原始资料；每篇保留 2–5 个直接链接与访问日期。
- 观察事实、推断和待验证假设必须分开记录；版本变化后复核旧结论。
- 用户访谈不是本阶段必修内容，不在本索引创建访谈路线。
