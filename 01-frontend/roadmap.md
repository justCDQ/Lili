# 前端深化 Roadmap

## 一、最终能力

完成本模块后，应能够：

- 独立设计中大型前端应用架构。
- 使用 DevTools 定位网络、运行时、渲染和内存问题。
- 建立 Monorepo、Design System、测试和发布体系。
- 开发 CLI、Codemod、构建插件和自动化工具。
- 处理复杂状态、权限、国际化和可观测性。
- 在 AI Native UI、编辑器、数据可视化等方向形成专项能力。

---

## 二、Web 基础

### HTML

- [ ] 语义化元素与文档结构。
- [ ] 表单、原生校验、自动填充。
- [ ] 焦点顺序和键盘操作。
- [ ] `label`、`fieldset`、`legend`。
- [ ] Dialog、Popover、Details。
- [ ] 响应式图片、`picture`、`srcset`。
- [ ] SEO 与可访问性基础。
- [ ] ARIA 的使用边界。

验收：不使用 UI 库，实现一个完全支持键盘操作、错误关联和浏览器自动填充的注册流程。

### CSS

- [ ] Normal Flow、Box Model。
- [ ] Flexbox、Grid。
- [ ] Containing Block、BFC、Stacking Context。
- [ ] Cascade、Specificity、Inheritance。
- [ ] Custom Properties、Cascade Layers。
- [ ] Container Queries、Logical Properties。
- [ ] 深色模式、RTL、高对比度模式。
- [ ] Transition、Animation、Transform。
- [ ] `prefers-reduced-motion`。

验收：能解释尺寸、滚动、定位和 `z-index` 的来源；能设计一套 Token，而不是到处写任意值。

### JavaScript

- [ ] Scope、Closure、Prototype、`this`。
- [ ] Module、Iterator、Generator、Proxy。
- [ ] Event Loop、Task、Microtask。
- [ ] Promise、Async Iterator、AbortController。
- [ ] Streams、Structured Clone、Transferable Object。
- [ ] GC、WeakMap、WeakSet。
- [ ] 事件监听器、计时器、闭包造成的泄漏。

必做：

- [ ] 并发任务调度器。
- [ ] 可取消、可重试请求器。
- [ ] SSE 流解析器。
- [ ] EventEmitter。
- [ ] LRU Cache。

验收：遇到异步问题时能画出调用栈、任务队列和生命周期，而不是通过增加 `setTimeout` 试错。

### TypeScript

- [ ] Narrowing、Generics。
- [ ] Conditional、Mapped、Template Literal Types。
- [ ] Discriminated Union、Type Predicate、`infer`。
- [ ] Declaration File、Module Resolution。
- [ ] Project References、Compiler Options。
- [ ] Runtime Schema Validation。

必做：

- [ ] 类型安全 API Client。
- [ ] 类型安全事件系统。
- [ ] 使用联合类型表达状态机。
- [ ] 为 JavaScript 库编写声明文件。

验收：能设计稳定公共 API；不滥用 `any` 和 `as`；知道类型系统不能替代运行时验证。

---

## 三、浏览器与运行时

### 页面加载

- [ ] DNS、TCP、TLS、HTTP。
- [ ] HTML Parser、Preload Scanner。
- [ ] CSS 和 Script 的阻塞行为。
- [ ] preload、prefetch、preconnect。
- [ ] HTTP Cache、Service Worker Cache。
- [ ] Network 瀑布图分析。

### 渲染流水线

- [ ] Style、Layout、Paint、Composite。
- [ ] 主线程、合成线程、光栅化。
- [ ] Long Task、Layout Thrashing。
- [ ] `requestAnimationFrame` 和任务切片。
- [ ] GPU 合成层的收益与成本。

### 内存

- [ ] Heap Snapshot、Allocation Timeline。
- [ ] Detached DOM、Retainer Path。
- [ ] Listener、Timer、Worker、Blob URL 泄漏。
- [ ] 无界缓存治理。

必做：

- [ ] 制造并修复六类内存泄漏。
- [ ] 制作一个故意很卡的页面。
- [ ] 记录优化前后 CPU、FPS、Heap、Long Task、LCP、INP 和 CLS。

验收：能区分网络慢、JS 慢、渲染慢和后端慢；所有结论都有 Profile 或指标证据。

---

## 四、应用与组件架构

### 组件设计

- [ ] 单一职责、Composition。
- [ ] Controlled / Uncontrolled。
- [ ] Headless、Compound Component。
- [ ] State Machine、Context Boundary。
- [ ] Dependency Inversion、API Stability。
- [ ] 业务组件、通用组件、基础组件的边界。

### 状态设计

- [ ] URL 状态。
- [ ] Server State。
- [ ] Form State。
- [ ] Interaction State。
- [ ] Global State。
- [ ] Persistent State。
- [ ] 状态所有权与单一数据源。

### 应用分层

- [ ] 按业务领域组织模块。
- [ ] 建立单向依赖规则。
- [ ] 封装第三方基础设施。
- [ ] 统一错误、请求、缓存和重试模型。
- [ ] 使用 ADR 记录决策。
- [ ] 建立循环依赖检查。

必做：设计一个企业后台，输出架构图、模块依赖图、状态分类、权限模型、错误模型和测试策略。

验收：新功能能明确归属；替换请求库、状态库或 UI 库时，不需要重写全部业务。

---

## 五、性能工程

### 指标

- [ ] TTFB、FCP、LCP、INP、CLS。
- [ ] Long Task、FPS、JS Heap。
- [ ] Bundle Size、请求数量、缓存命中率。
- [ ] 实验室数据与真实用户数据。
- [ ] 分位数与性能预算。

### 加载性能

- [ ] 关键请求链。
- [ ] 首屏 JavaScript。
- [ ] 路由和组件拆包。
- [ ] 图片、字体和资源优先级。
- [ ] CDN 与缓存。
- [ ] SSR、SSG、CSR、Streaming SSR。
- [ ] Hydration 和第三方脚本。

### 运行时性能

- [ ] 长任务切片。
- [ ] Worker。
- [ ] 虚拟列表。
- [ ] 无效渲染。
- [ ] 高频输入调度。
- [ ] Canvas 和图表更新。
- [ ] 内存边界。

必做：为一个真实项目建立基线、性能预算、优化报告、线上采集和 CI 回归检查。

验收：不使用“感觉快了”作为结论；能说明每项优化的收益、成本和适用范围。

---

## 六、工程化与工具链

### 构建与模块

- [ ] ESM、CJS、Package Exports。
- [ ] Tree Shaking、Code Splitting、Chunk。
- [ ] Source Map、HMR。
- [ ] Bundler Plugin、AST、Compiler Pipeline。
- [ ] ESM/CJS/类型声明发布。

### Monorepo

- [ ] Workspace、Package Boundary。
- [ ] Task Graph、Remote Cache。
- [ ] Incremental Build。
- [ ] Versioning、Changelog、Release。
- [ ] 跨包依赖约束。

### CLI 与自动化

- [ ] 交互式和非交互式命令。
- [ ] 模板、Dry Run、回滚。
- [ ] Migration、Codemod。
- [ ] 日志等级、单元和集成测试。

必做：

- [ ] 一个构建插件。
- [ ] 一个 AST Codemod。
- [ ] 一个团队 CLI。
- [ ] 一套 Monorepo 依赖规则。
- [ ] 一套自动发布流程。

验收：从工具使用者变成工具建设者，并能量化工具对团队效率的提升。

---

## 七、测试、安全与可观测性

### 测试

- [ ] Unit、Integration、Component、E2E。
- [ ] Contract、Visual Regression。
- [ ] Accessibility、Performance Test。
- [ ] 测试加载、失败、重试、权限和边界。
- [ ] Bug 修复附带回归测试。

### 安全

- [ ] XSS、CSRF、CSP、CORS。
- [ ] Cookie、SameSite、OAuth。
- [ ] Token Storage、Open Redirect、Clickjacking。
- [ ] 依赖风险、Source Map 暴露。
- [ ] 前端权限不等于后端鉴权。

### 可观测性

- [ ] Runtime Error、Promise Error、API Error。
- [ ] Release、Source Map、Web Vitals。
- [ ] User Action、Trace、Business Event。
- [ ] 错误聚合、采样、告警和版本关联。

验收：线上问题可以回答谁、在哪、何时、哪个版本、执行了什么、影响多少用户以及是否由本次发布引入。

---

## 八、专项方向

选择一到两个长期深入：

### 企业级前端平台

Design System、CLI、权限、国际化、埋点、监控、Feature Flag、发布平台。

### 复杂交互工程

编辑器、Canvas、白板、PDF 标注、可视化、大表格、拖拽、Undo/Redo、协作。

### AI Native Frontend

Streaming、Citation、Artifact、Tool Call、Agent Task、Approval、Memory、长任务。

### Web Platform

WebAssembly、WebGPU、Worker、PWA、WebRTC、浏览器扩展。

---

## 九、推荐资源

### 书籍

- You Don’t Know JS Yet
- JavaScript 高级程序设计
- CSS in Depth
- High Performance Browser Networking
- Web Performance in Action
- A Philosophy of Software Design
- Clean Architecture

### 网站

- MDN
- web.dev
- Chrome Developers
- V8 Blog
- TypeScript Handbook
- React 官方文档
- Playwright
- OWASP


---
