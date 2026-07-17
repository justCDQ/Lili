# 框架掌握与 React、Vue、Svelte 比较

框架选择不是语法偏好排名，而是把产品约束、团队能力、运行环境、生态依赖和长期维护映射到可验证方案。先用一个框架完成完整应用，再比较相同问题在其他框架中的机制和成本。

## 1. “掌握框架”的验收

能独立完成：

- 组件、props、state、事件、列表和表单；
- 路由、嵌套 layout、参数和错误页；
- 请求状态、运行时校验、缓存与 mutation；
- SSR/CSR/SSG 中至少一种部署；
- 类型检查、lint、单元、组件与 E2E；
- 性能记录、无障碍检查和错误监控；
- 依赖升级、安全公告与回滚。

只会脚手架和计数器不构成框架掌握。

## 2. 当前能力边界

### React 19.2

React 用组件函数与元素树描述 UI，state 更新触发 render/commit，生态通常由路由或全栈框架补齐数据、构建和部署。19.2 包含 Activity、useEffectEvent、Performance Tracks、部分预渲染与 SSR 改进。实际安装应使用已经包含安全修复的最新 19.2 patch；使用 Server Components 的框架还要跟随框架公告。

### Vue 3

Vue 提供 template/JSX、Composition API、Proxy/ref 响应式、computed/watch 和单文件组件。官方生态包括 Vue Router、Pinia、Vite 与 Nuxt。模板编译可获得静态提升和 patch flag 等优化。

### Svelte 5

Svelte 编译组件并生成针对性更新代码。Svelte 5 runes 模式用 `$state`、`$derived`、`$effect`、`$props` 等显式表达响应式；SvelteKit 提供路由、数据加载和服务端能力。旧版 legacy 语法仍可能出现在项目中，迁移成本要单独评估。

## 3. 比较同一功能

计数器不是足够基准，但可展示状态表达差异。

React：

```tsx
function Counter() {
  const [count, setCount] = useState(0);
  return <button onClick={() => setCount((value) => value + 1)}>{count}</button>;
}
```

Vue：

```vue
<script setup lang="ts">
import { ref } from "vue";
const count = ref(0);
</script>
<template><button @click="count++">{{ count }}</button></template>
```

Svelte：

```svelte
<script lang="ts">
let count = $state(0);
</script>
<button onclick={() => count++}>{count}</button>
```

这些代码只说明局部语法和默认响应式模型，不说明大型应用性能、可维护性或招聘成本。

## 4. 决策维度

| 维度 | 要收集的证据 |
|---|---|
| 产品形态 | SEO、首屏、离线、实时、编辑器、嵌入式 widget |
| 团队 | 现有经验、培训时间、值班能力、人员流动 |
| 生态 | 必需组件、图表、编辑器、原生桥接、企业库 |
| 平台 | Node/edge/serverless/static、浏览器范围、CSP |
| 性能 | 真实页面 bundle、交互、内存与服务端成本 |
| 类型 | 模板/插件/Compiler API 兼容，TS7/TS6 工具链 |
| 维护 | LTS/发布节奏、迁移文档、安全响应、治理 |
| 测试 | 测试工具成熟度、SSR/hydration 能力 |
| 招聘 | 本地候选池、内部轮岗、知识集中风险 |

### 4.1 把维度变成可验证问题

“生态成熟”不可直接评分，应拆成项目证据：目标日期选择器是否支持时区和键盘；数据表能否在 10,000 行、服务端分页和虚拟滚动下工作；监控 SDK 是否支持 SSR source map；身份库能否在目标 edge runtime 验证会话。

“团队熟悉”也要量化：完成垂直切片时间、review 缺陷、值班恢复时间、只有一人能维护的模块数。学习新框架不是零成本，但已有经验也不保证当前架构正确。

### 4.2 硬约束与偏好

硬约束不应被其他高分抵消。例如必须嵌入现有 React host、关键 SDK 只支持 Vue、部署平台禁止长运行 Node、license 不允许商业分发。偏好如模板风格、语法简短可加权但权重较低。

## 5. 框架与元框架

React/Vue/Svelte 主要解决 UI 与响应式；Next/React Router framework mode、Nuxt、SvelteKit 等进一步决定：文件路由、loader/action、服务端组件、缓存、构建目标和部署 adapter。

“选择 React”并没有决定 SSR cache 语义；必须比较实际应用栈。元框架越多地控制服务器和构建，升级时越要检查平台耦合与安全公告。

## 6. TypeScript 7 兼容性

TypeScript 7.0 CLI/LSP 已正式可用，但没有 Compiler API。React 的普通 `.ts/.tsx` + tsc 流程可直接使用 7；Vue、Svelte、MDX、Astro、Angular 模板和 typescript-eslint 等嵌入或 API 依赖工具可能仍要求 TypeScript 6。

评估时执行真实工具链，而不是只跑 `tsc`：模板 typecheck、lint、IDE 补全、测试转换、声明生成和构建都要通过。必要时使用 `@typescript/typescript6` 与 TS7 并存。

## 7. 性能比较方法

建立同一代表性切片：

- 1000 行可排序表格；
- 包含验证和异步提交的表单；
- SSR 页面和 hydration；
- 路由懒加载；
- 第三方编辑器或图表；
- 错误、空、加载和权限状态。

固定数据、浏览器、网络、CPU 和生产构建。记录 JS 传输/解压、解析执行、LCP、INP、内存、服务端 TTFB、cache hit、构建时间。重复多次并给分布，不用单次 Lighthouse 分数下结论。

SSR 还要记录服务端 CPU、内存、冷启动、首个流式 chunk 和 cache hit。无障碍比较不能只跑 axe：实际测试 Tab 顺序、焦点陷阱、错误播报、Dialog/Combobox 和 200% 缩放。框架本身不保证组件实现无障碍，主要证据来自所选设计系统和应用代码。

## 8. 生态验证

对每个关键依赖建立清单：

1. 是否支持目标框架当前主版本；
2. SSR 是否访问 window；
3. TypeScript 类型和模板检查是否可用；
4. 无障碍和键盘行为；
5. bundle 与 tree shaking；
6. 最近发布、安全和 issue 响应；
7. license 和商业限制；
8. 替换成本。

组件库数量多不等于所需组件质量高。用实际 PoC 验证最难的两个依赖。

## 9. 团队和迁移成本

迁移成本包括：

- 业务组件重写；
- 路由、数据和 cache 语义；
- 测试与选择器；
- 设计系统和无障碍回归；
- SSR/部署/监控；
- 招聘培训；
- 双栈期间重复基础设施。

若现有栈能满足目标，局部性能问题通常先测量和优化。全量重写需要可量化收益和渐进边界。

### 9.1 招聘与知识风险

统计岗位供给时区分“写过 demo”和“能维护 SSR、性能、测试与升级”。同时评估内部培养：文档、示例、代码 review 和轮岗能否让第二个人接管。热门框架可能候选多但薪酬竞争高；小框架可能招聘少但团队留存稳定。结论来自组织所在地区和岗位数据。

### 9.2 维护与安全

检查治理主体、security policy、受支持版本、patch 速度、迁移工具和生态 CI。React Server Components 的安全修复要求跟踪框架集成版本；Vue/Svelte 的编译器和服务端框架也属于攻击面。选择记录列出升级责任人、月度依赖窗口和紧急 patch SLA。

### 9.3 渐进迁移

可通过路由边界、Web Component、iframe 或独立部署逐步迁移。每种边界都有重复 runtime、样式隔离、身份和路由同步、监控拆分及焦点跨边界成本。先迁移可独立交付且高收益区域，并定义停止条件，不默认最终全量重写。

## 10. 决策实验

假设要做多租户数据后台，要求 SSR 登录页、10k 表格、复杂表单、图表、五年维护。

### 输入证据

- 团队 6 人，4 人熟悉 React、2 人 Vue、无人 Svelte；
- 企业设计系统已有 React 版本；
- 表格库只有 React 适配经过内部无障碍审计；
- 部署平台支持 Node；
- SEO 只用于登录前页面。

### PoC

分别实现登录、表格、编辑抽屉、SSR 和 E2E。测量 bundle、INP、开发时长、bug、模板/类型错误和部署。

### 输出

React 方案可能因已有设计系统和团队经验成本最低；这不是 React 普遍优于其他框架。若产品是大量静态内容加少量独立交互，SvelteKit 的 PoC 可能在 JS 体积和实现复杂度占优；若组织已有 Vue/Nuxt 平台，Vue 的维护成本可能最低。

### 验证与失败分支

六周后复查生产指标和开发周期。若关键表格库在 SSR 下崩溃或安全维护停止，触发替代方案；若性能差来自 50MB 图表数据，换框架不会解决，应改数据与虚拟化。

### 10.1 证据化评分矩阵

团队先过滤硬约束，再为剩余方案记录评分、证据和置信度：

| 维度/权重 | React 栈 | Vue 栈 | Svelte 栈 |
|---|---|---|---|
| 现有设计系统 20% | 5：30 个组件已审计 | 2：仅基础封装 | 1：需重写 |
| 表格与编辑器 15% | 5：PoC 全功能 | 4：编辑器待适配 | 3：wrapper 有 SSR 问题 |
| 团队交付 15% | 5：中位 3 天 | 3：中位 5 天 | 2：中位 7 天 |
| SSR/部署 10% | 4：Node 平台成熟 | 4：Nuxt PoC 通过 | 4：SvelteKit PoC 通过 |
| 无障碍 15% | 5：人工通过 | 3：两个组件需修 | 2：四个组件需自建 |
| 性能 10% | 3：INP 180ms | 4：150ms | 4：140ms |
| 维护/招聘 15% | 5：多人值班 | 3：两人值班 | 1：无人生产经验 |

性能差异没有大到抵消设计系统、无障碍和知识风险，因此该项目选 React。这些数值只对给定证据成立；设计系统跨框架或团队结构改变后应重新评分。

### 10.2 完整失败分支

PoC 后三个月，表格厂商停止维护并曝出高危漏洞。adapter 将 vendor API 限制在 `DataGrid` 边界，因此团队可以替换；若业务页面直接导入 vendor hooks 和类型，迁移成本会扩散。ADR 触发条件要求 7 天内评估修复或替代，期间用 feature flag 关闭有风险功能。

另一个失败分支是部署目标改为 edge runtime。Node-only session SDK 成为硬阻塞；此时重验元框架 adapter 和 SDK 替代，不因为之前 React 总分高就维持原方案。

## 11. 完整选型案例：内容与协作编辑产品

### 11.1 约束与证据

- 公共文档需要 SSG、结构化数据和多语言；
- 登录后编辑器需要协作、离线草稿和 100k 字文档；
- 目标是 Node + CDN，未来可能 edge；
- 团队 5 人：React 生产经验 3、Vue 1、Svelte 0；
- 现有编辑器插件只提供 React/Vue；
- WCAG 2.2 AA 是发布门；
- 首年必须在 12 周上线。

### 11.2 垂直切片

每个候选实现公共文档 SSG、登录、编辑器加载、协作断线恢复、评论表单、SSR 错误与 Playwright E2E。固定 50KB/1MB 文档、4× CPU slowdown 和相同 CDN。不同工程师交叉实现与 review，避免把个人熟练度全算成框架能力。

### 11.3 输出和决策

记录开发人日、TS/模板诊断、客户端 JS、编辑器可交互、内存、断线恢复、无障碍缺陷、SSR TTFB、构建 10k 文档时间。若 React 在交付和编辑器生态明显领先，采用 React 元框架；公共内容保持静态/服务端，编辑器是懒加载客户端边界。领域模型、协作协议、schema 和 API client 放框架无关包，vendor 编辑器通过 adapter 隔离。

### 11.4 验收和复查

上线 8 周复查编辑器 INP、内存崩溃、错误率、每功能周期和新人上手时间。bundle 预算超标时先分析编辑器和数据；关键生态不维护或 edge 成硬需求时重开 ADR。

## 12. 加权评分但不伪造精确性

```text
总分 = Σ(维度权重 × 有证据的评分)
```

权重来自项目约束；评分附证据链接和不确定度。不要用 8.3 对 8.2 宣称客观胜负。高风险硬约束可直接一票否决，例如目标平台无法部署、关键依赖 license 不可接受。

## 13. 常见错误

1. 用 stars、下载量或社交热度代替项目证据。
2. 比较框架 hello world，不比较完整栈。
3. 把语法短当作维护成本低。
4. 把 benchmark 微差推广到真实产品。
5. 忽略安全 patch 和元框架版本。
6. 只考虑开发，不考虑部署、监控和升级。
7. 在 TS7 下只测 tsc，忽略嵌入语言工具仍依赖 TS6。
8. 因个人偏好把建议写成普遍规律。

## 14. 练习

为一个公共内容 + 登录后台产品完成框架 ADR。验收：

1. 明确 10 个约束和权重；
2. React/Vue/Svelte 各做同一垂直切片；
3. 使用生产构建和同一测试数据；
4. 包含 SSR、表单、表格、错误和 E2E；
5. 验证 TS7/TS6 全工具链兼容；
6. 记录性能分布、实现时长和缺陷；
7. 给出选择、风险、回滚和复查日期；
8. 结论只由项目证据支持。

## 来源

- [React Versions](https://react.dev/versions)（访问日期：2026-07-17）
- [React 19.2](https://react.dev/blog/2025/10/01/react-19-2)（访问日期：2026-07-17）
- [Vue：Introduction](https://vuejs.org/guide/introduction.html)（访问日期：2026-07-17）
- [Svelte：What are runes?](https://svelte.dev/docs/svelte/what-are-runes)（访问日期：2026-07-17）
- [TypeScript Team：Announcing TypeScript 7.0](https://devblogs.microsoft.com/typescript/announcing-typescript-7-0/)（访问日期：2026-07-17）
