# 交互设计入门、初级与中级学习笔记

本目录覆盖 [交互设计 Roadmap](../roadmap.md) 的阶段零至阶段八。每篇对应一个独立知识点；日常产品任务分析仍放在 [`daily/`](../daily/README.md)。

## 阅读顺序

1. `00-foundations/`：建立任务、界面结构、流程、状态和交付物的共同语言。
2. `01-principles/`：掌握分析设计决策所需的基础原则及适用边界。
3. `02-deconstruction/`：对真实任务进行可复现拆解、比较和重新设计。
4. `03-validation/`：使用公开证据、走查、审计与原型验证方案。
5. `04-information-architecture/`：建立、审计、重组并验证复杂产品结构。
6. `05-flows-states/`：为流程中的正常、异常、中断和恢复状态建立契约。
7. `06-interaction-patterns/`：按导航、输入、数据、反馈、操作和协作建立模式库。
8. `07-enterprise-interactions/`：处理字段与数据范围权限、临时访问和可恢复审批工作流。
9. `08-visual-motion-accessibility/`：建立视觉系统、状态动效、键盘与焦点规范。
10. `09-ai-product-interaction/`：设计可控制、可恢复、可核验的 AI 对话、产物、工具、Agent、引用与记忆交互。

## 阶段零：交互设计基础

### 基本语言

- [用户目标、任务、场景、入口、路径与完成标准](00-foundations/01-user-goal-task-context-path.md)
- [页面、弹窗、抽屉与浮层](00-foundations/02-page-dialog-drawer-popover.md)
- [导航、表单、列表与详情](00-foundations/03-navigation-form-list-detail.md)
- [信息架构、用户流程、任务流程、状态与反馈](00-foundations/04-ia-user-flow-task-flow-state-feedback.md)
- [线框图、原型保真度与设计规范](00-foundations/05-wireframe-prototype-fidelity-spec.md)
- [区分产品、交互、视觉与工程问题](00-foundations/06-product-interaction-visual-engineering-problems.md)

### 基础原则

- [心智模型、概念模型与现实世界映射](01-principles/01-mental-conceptual-real-world-mapping.md)
- [可见性、可供性、反馈、一致性与约束](01-principles/02-visibility-affordance-feedback-consistency-constraints.md)
- [识别优于回忆、认知负荷与信息分组](01-principles/03-recognition-cognitive-load-grouping.md)
- [用户控制、撤销、容错、错误预防与恢复](01-principles/04-control-undo-tolerance-error-prevention-recovery.md)
- [渐进披露、默认值、熟悉性与学习成本](01-principles/05-progressive-disclosure-defaults-familiarity-learning.md)
- [Fitts 定律与 Hick 定律的适用边界](01-principles/06-fitts-hick-laws-boundaries.md)

## 阶段一：思考、产品拆解与重新设计

- [产品拆解的上下文记录](02-deconstruction/01-record-context.md)
- [还原入口、前置条件、主流程、分支、异常与完成状态](02-deconstruction/02-reconstruct-flow.md)
- [初始、空、加载、成功、失败、权限、离线与冲突状态清单](02-deconstruction/03-state-inventory.md)
- [信息层级、操作层级、默认值、文案与渐进披露分析](02-deconstruction/04-hierarchy-defaults-copy-disclosure.md)
- [识别交互模式及适用条件](02-deconstruction/05-identify-patterns.md)
- [用启发式原则分析并区分事实与推断](02-deconstruction/06-heuristics-facts-inferences.md)
- [与同类产品或现实工作流比较](02-deconstruction/07-compare-product-workflow.md)
- [重新设计流程与线框图](02-deconstruction/08-redesign-flow-wireframe.md)
- [评估改动收益、代价、风险与工程注意事项](02-deconstruction/09-benefit-cost-risk-engineering.md)
- [实现关键交互并验证键盘、响应式与状态切换](02-deconstruction/10-implement-keyboard-responsive-states.md)

## 阶段二：个人可执行的研究与验证

### 无需访谈的证据来源

- [用产品评论、公开 Issue、帮助文档、更新日志和社区问题取证](03-validation/01-public-reviews-issues-docs-changelogs.md)
- [竞品、历史版本与跨平台任务对比](03-validation/02-competitive-version-platform-comparison.md)
- [Cognitive Walkthrough：认知走查](03-validation/03-cognitive-walkthrough.md)
- [Heuristic Evaluation：启发式评估](03-validation/04-heuristic-evaluation.md)
- [Accessibility Audit：无障碍审计](03-validation/05-accessibility-audit.md)
- [状态审计：空、错、慢、断、冲突与权限](03-validation/06-state-audit.md)
- [Dogfooding 与任务日志](03-validation/07-dogfooding-task-log.md)
- [使用埋点、热图、搜索词、错误日志与客服材料](03-validation/08-analytics-heatmaps-search-errors-support.md)

### 原型验证

- [为原型定义真实任务、成功条件和观察项](03-validation/09-prototype-task-success-observations.md)
- [原型自测中的认知走查](03-validation/10-self-test-cognitive-walkthrough.md)
- [用键盘、窄屏、慢网、空数据和错误响应制造边界](03-validation/11-manufacture-boundary-conditions.md)
- [比较改版前后的步骤、完成率、时间、错误与认知成本](03-validation/12-compare-before-after.md)
- [可选的可用性测试：增强证据而非学习准入](03-validation/13-optional-usability-testing.md)
- [根据证据迭代并维护假设清单](03-validation/14-evidence-iteration-assumption-log.md)

## 阶段三：信息架构

- [绘制复杂产品的现有站点地图](04-information-architecture/01-current-sitemap.md)
- [审计重复、混乱与层级过深的入口](04-information-architecture/02-entry-problem-audit.md)
- [按任务、角色或业务对象重新分类](04-information-architecture/03-reclassify-task-role-object.md)
- [设计顶部、侧边、Tabs、面包屑与搜索的边界](04-information-architecture/04-navigation-boundaries.md)
- [验证目标功能是否容易找到](04-information-architecture/05-findability-validation.md)

## 阶段四：用户流程与状态

- [初始状态](05-flows-states/01-initial.md)
- [空状态](05-flows-states/02-empty.md)
- [加载状态](05-flows-states/03-loading.md)
- [成功状态](05-flows-states/04-success.md)
- [失败状态](05-flows-states/05-failure.md)
- [部分成功状态](05-flows-states/06-partial-success.md)
- [无权限状态](05-flows-states/07-no-permission.md)
- [离线状态](05-flows-states/08-offline.md)
- [数据过期状态](05-flows-states/09-stale-data.md)
- [取消状态](05-flows-states/10-cancel.md)
- [重试状态](05-flows-states/11-retry.md)
- [并发冲突状态](05-flows-states/12-concurrent-conflict.md)

## 阶段五：交互模式库

### 导航模式

- [Sidebar 侧边导航](06-interaction-patterns/00-navigation/sidebar.md)
- [Top Navigation 顶部导航](06-interaction-patterns/00-navigation/top-navigation.md)
- [Tabs 标签页](06-interaction-patterns/00-navigation/tabs.md)
- [Breadcrumb 面包屑](06-interaction-patterns/00-navigation/breadcrumb.md)
- [Command Palette 命令面板](06-interaction-patterns/00-navigation/command-palette.md)
- [Stepper 步骤器](06-interaction-patterns/00-navigation/stepper.md)

### 输入模式

- [Form 表单](06-interaction-patterns/01-input/form.md)
- [Inline Edit 行内编辑](06-interaction-patterns/01-input/inline-edit.md)
- [Auto Complete 自动完成](06-interaction-patterns/01-input/autocomplete.md)
- [Tag Input 标签输入](06-interaction-patterns/01-input/tag-input.md)
- [Date Picker 日期选择](06-interaction-patterns/01-input/date-picker.md)
- [Upload 文件上传](06-interaction-patterns/01-input/upload.md)
- [Rich Text Editor 富文本编辑器](06-interaction-patterns/01-input/rich-text-editor.md)

### 数据模式

- [Table 数据表格](06-interaction-patterns/02-data/table.md)
- [List 列表](06-interaction-patterns/02-data/list.md)
- [Card 卡片](06-interaction-patterns/02-data/card.md)
- [Tree 树形数据](06-interaction-patterns/02-data/tree.md)
- [Timeline 时间线](06-interaction-patterns/02-data/timeline.md)
- [Dashboard 仪表盘](06-interaction-patterns/02-data/dashboard.md)
- [Chart 图表](06-interaction-patterns/02-data/chart.md)

### 反馈模式

- [Toast 轻提示](06-interaction-patterns/03-feedback/toast.md)
- [Alert 警示](06-interaction-patterns/03-feedback/alert.md)
- [Inline Error 行内错误](06-interaction-patterns/03-feedback/inline-error.md)
- [Progress 进度](06-interaction-patterns/03-feedback/progress.md)
- [Skeleton 骨架屏](06-interaction-patterns/03-feedback/skeleton.md)
- [Empty State 空状态](06-interaction-patterns/03-feedback/empty-state.md)
- [Result Page 结果页](06-interaction-patterns/03-feedback/result-page.md)

### 操作模式

- [Create 创建](06-interaction-patterns/04-operations/create.md)
- [Edit 编辑](06-interaction-patterns/04-operations/edit.md)
- [Delete 删除](06-interaction-patterns/04-operations/delete.md)
- [Batch 批量操作](06-interaction-patterns/04-operations/batch.md)
- [Drag 拖拽](06-interaction-patterns/04-operations/drag.md)
- [Undo 撤销](06-interaction-patterns/04-operations/undo.md)
- [Retry 重试](06-interaction-patterns/04-operations/retry.md)
- [Save Draft 保存草稿](06-interaction-patterns/04-operations/save-draft.md)

### 协作模式

- [Comment 评论](06-interaction-patterns/05-collaboration/comment.md)
- [Mention 提及](06-interaction-patterns/05-collaboration/mention.md)
- [Share 分享](06-interaction-patterns/05-collaboration/share.md)
- [Notification 通知](06-interaction-patterns/05-collaboration/notification.md)
- [Presence 在线状态](06-interaction-patterns/05-collaboration/presence.md)
- [Version History 版本历史](06-interaction-patterns/05-collaboration/version-history.md)
- [Conflict Resolution 冲突解决](06-interaction-patterns/05-collaboration/conflict-resolution.md)

## 阶段六：复杂 B 端交互

### 表格与表单

- [复杂数据表格：列、排序、筛选、分页、跨页选择、编辑、导出与性能](06-interaction-patterns/02-data/table.md)
- [复杂表单：分组、依赖、校验、草稿、错误保留与冲突](06-interaction-patterns/01-input/form.md)

### 权限与审批

- [B 端权限状态、字段权限与数据范围](07-enterprise-interactions/01-permission-states-data-scope.md)
- [临时权限、权限申请与到期恢复](07-enterprise-interactions/02-temporary-access-request-expiry.md)
- [审批流程、并行决策与异常恢复](07-enterprise-interactions/03-approval-workflow-state-exceptions.md)

## 阶段七：视觉基础、动效与无障碍

- [视觉层级、间距、字体、色彩与图标系统](08-visual-motion-accessibility/01-visual-hierarchy-spacing-typography-color-icons.md)
- [界面动效、状态转移与减少动画](08-visual-motion-accessibility/02-motion-state-transitions-reduced-motion.md)
- [键盘、焦点与可访问交互](08-visual-motion-accessibility/03-keyboard-focus-accessible-interaction.md)

## 阶段八：AI 产品交互

- [AI 对话与流式响应：澄清、多轮、停止、重试、继续与断线恢复](09-ai-product-interaction/01-chat-streaming.md)
- [AI Artifact 编辑与版本：创建、局部修改、Diff、恢复、导出与修改来源](09-ai-product-interaction/02-artifact-editing-versioning.md)
- [AI 工具调用与风险确认：工具、参数、影响、结果与高风险写入批准](09-ai-product-interaction/03-tool-call-approval.md)
- [Agent 任务与人工接管：目标、计划、进度、暂停、取消、恢复与失败步骤](09-ai-product-interaction/04-agent-task-control.md)
- [AI 引用与记忆交互：原文定位、日期、推断和记忆作用范围控制](09-ai-product-interaction/05-citation-memory.md)

## 覆盖表

| Roadmap 范围 | Roadmap 条目 | 笔记数 | 覆盖状态 |
| --- | ---: | ---: | --- |
| 阶段零 · 基本语言 | 5 | 6 | 完成；界面容器与内容/任务结构拆成两篇 |
| 阶段零 · 基础原则 | 6 | 6 | 完成 |
| 阶段一 | 10 | 10 | 完成 |
| 阶段二 · 无需访谈的证据来源 | 8 | 8 | 完成 |
| 阶段二 · 原型验证 | 6 | 6 | 完成 |
| 阶段三 · 信息架构 | 5 | 5 | 完成 |
| 阶段四 · 用户流程与状态 | 12 | 12 | 完成 |
| 阶段五 · 交互模式库 | 42 | 42 | 完成 |
| 阶段六 · 复杂 B 端交互 | 16 | 5 | 完成；表格与表单沿用模式库深度文章，权限与审批新增 3 篇 |
| 阶段七 · 视觉、动效与无障碍 | 6 | 3 | 完成；按视觉系统、动效状态机与可访问交互聚合 |
| 阶段八 · AI 产品交互 | 10 | 5 | 完成；按五组紧密耦合的状态与控制契约聚合 |
| **合计** | **126** | **106** | **阶段零至阶段八全部覆盖** |

## 维护规则

- 新笔记必须说明是什么、为什么、实际方法、检查或步骤、常见错误与边界、补充知识及来源。
- 来源优先使用 W3C/WAI、平台官方 HIG、公共服务设计规范和方法原始资料；每篇保留 2–5 个直接链接与访问日期。
- 观察事实、推断和待验证假设必须分开记录；版本变化后复核旧结论。
- 用户访谈不是本阶段必修内容，不在本索引创建访谈路线。
