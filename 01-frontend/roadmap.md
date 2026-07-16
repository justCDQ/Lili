# 前端 Roadmap：从零基础到前端专家

这是一条可以独立学习的前端路线。顺序是“计算机与 Web 入门 → 页面 → 编程 → 工程应用 → 浏览器原理 → 架构与平台”，每个阶段都同时包含知识、练习、项目和验收。

## 能力阶梯

| 阶段 | 能力目标 | 代表产出 |
| --- | --- | --- |
| 入门 | 会使用开发环境，理解网页如何被访问 | 第一个可访问网页 |
| 初级 | 能独立还原页面并编写交互 | 多页面响应式网站 |
| 中级 | 能用框架和 TypeScript 构建完整应用 | 带 API、路由、状态和测试的应用 |
| 高级 | 能处理性能、架构、质量和复杂场景 | 中大型应用与性能报告 |
| 专家 | 能建设平台、工具和标准并影响团队 | Design System、CLI、基础设施或专项产品 |

---

## 阶段零：计算机与 Web 入门

### 开发环境

- [ ] 文件、目录、路径、扩展名、文本编码和压缩包。
- [ ] VS Code 的文件、搜索、终端、插件和调试面板。
- [ ] 命令行的 `cd`、`pwd`、`ls`、`mkdir`、`cp`、`mv` 和删除操作。
- [ ] Git 的仓库、工作区、暂存区、提交、分支、合并和远端。
- [ ] Node.js、包管理器、`package.json`、依赖和脚本。
- [ ] 使用浏览器 DevTools 查看 Elements、Console、Network 和 Sources。

### Web 如何工作

- [ ] 客户端、服务器、浏览器和搜索引擎的角色。
- [ ] URL、域名、DNS、IP、端口、HTTP 请求与响应。
- [ ] HTML、CSS、JavaScript 各自负责什么。
- [ ] 静态网站、动态网站、前端、后端和 API 的区别。
- [ ] 本地开发服务器、构建产物、部署和 HTTPS 的基本概念。

必做：创建一个 Git 仓库，用纯文本写出第一个 HTML 页面，在本地服务器打开并部署到公开地址。

验收：能从输入 URL 开始，口述浏览器获取并显示页面的大致过程；能独立创建、运行、调试和提交一个网页。

---

## 阶段一：HTML 与内容结构

### 基础语法

- [ ] 元素、属性、嵌套、注释、空元素和字符实体。
- [ ] 标题、段落、列表、链接、图片、音视频和表格。
- [ ] `head`、元数据、favicon、语言和页面标题。
- [ ] 块级与行内内容的基本区别。

### 语义、表单与可访问性

- [ ] `header`、`nav`、`main`、`article`、`section`、`aside`、`footer`。
- [ ] 表单控件、按钮、原生校验、自动填充和提交行为。
- [ ] `label`、`fieldset`、`legend`、错误提示与帮助文本。
- [ ] 标题层级、替代文本、焦点顺序和键盘操作。
- [ ] Dialog、Popover、Details 等原生交互元素。
- [ ] 响应式图片、`picture`、`srcset` 和懒加载。
- [ ] SEO、Open Graph 和结构化内容基础。
- [ ] ARIA 的使用边界：优先使用正确的原生元素。

必做：个人主页、文章详情页、数据表格和支持键盘操作及自动填充的注册流程。

验收：关闭 CSS 和 JavaScript 后内容仍然结构清楚；能使用键盘完成表单；HTML 校验无关键错误。

---

## 阶段二：CSS 与视觉实现

### 基础与布局

- [ ] 选择器、声明、单位、颜色、背景、边框和字体。
- [ ] Cascade、Specificity、Inheritance 和默认样式。
- [ ] Box Model、Normal Flow、`display`、溢出和滚动。
- [ ] Position、Containing Block、BFC 和 Stacking Context。
- [ ] Flexbox、Grid、多列布局和常见居中方法。
- [ ] 响应式设计、媒体查询、移动优先和断点选择。

### 系统化与现代 CSS

- [ ] Custom Properties、Design Token 和主题切换。
- [ ] Cascade Layers、Container Queries、Logical Properties。
- [ ] 深色模式、RTL、高对比度和打印样式。
- [ ] Transition、Animation、Transform 和关键帧。
- [ ] `prefers-reduced-motion` 与动效性能。
- [ ] BEM、CSS Modules、CSS-in-JS、Utility CSS 的边界与取舍。

必做：从设计稿还原三个不同布局；实现响应式营销页、Dashboard 和一组可复用基础组件样式。

验收：能解释元素尺寸、滚动、定位和 `z-index` 的来源；缩放和不同屏幕下不依赖大量临时覆盖；能建立一套 Token。

---

## 阶段三：JavaScript 编程基础到语言深入

### 编程基础

- [ ] 值、变量、类型、运算符、表达式和类型转换。
- [ ] 条件、循环、函数、参数、返回值和递归基础。
- [ ] Array、Object、String、Map、Set 和 Date 的常用操作。
- [ ] 解构、展开、模板字符串和可选链。
- [ ] 模块、导入导出、错误与异常处理。
- [ ] 使用伪代码拆解问题，使用断点和日志调试。

### DOM、事件与网络

- [ ] DOM 树、查询、创建、更新、删除和样式操作。
- [ ] Event、冒泡、捕获、委托、默认行为和自定义事件。
- [ ] 表单状态、校验、URL 与 Web Storage。
- [ ] JSON、Fetch、HTTP 方法、状态码和请求错误处理。
- [ ] 同步、异步、Callback、Promise 和 `async/await`。

### 语言深入

- [ ] Scope、Closure、Execution Context、Hoisting。
- [ ] Prototype、原型链、Class、`this` 和对象模型。
- [ ] Iterator、Generator、Proxy 和 Reflect。
- [ ] Event Loop、Call Stack、Task 和 Microtask。
- [ ] Promise 组合、Async Iterator、AbortController。
- [ ] Streams、Structured Clone、Transferable Object。
- [ ] GC、WeakMap、WeakSet 与常见内存泄漏。

必做：Todo、天气应用、分页搜索、图片懒加载、EventEmitter、并发调度器、可取消重试请求器和 SSE 流解析器。

验收：不依赖框架完成一个调用真实 API 的小型单页应用；能画出异步代码的调用栈、任务队列和生命周期。

---

## 阶段四：TypeScript、框架与完整应用

### TypeScript

- [ ] 基础类型、函数、对象、接口、类型别名和枚举的取舍。
- [ ] Union、Intersection、Narrowing、Generic。
- [ ] Conditional、Mapped、Template Literal Types。
- [ ] Discriminated Union、Type Predicate、`infer`。
- [ ] Declaration File、Module Resolution、Compiler Options 和 Project References。
- [ ] Runtime Schema Validation；类型系统不能替代运行时校验。

### 框架基础

- [ ] 组件、Props、State、事件、条件、列表和生命周期。
- [ ] 响应式更新或渲染模型，理解框架如何把状态映射为 UI。
- [ ] 路由、Layout、表单、请求、错误边界和懒加载。
- [ ] 客户端状态、服务端状态、URL 状态和持久化状态。
- [ ] SSR、SSG、CSR 和 Hydration 的入门概念。
- [ ] 先掌握一种主流框架，再比较 React、Vue、Svelte 等方案。

### 应用工程基础

- [ ] Vite 等开发与构建工具的基本配置。
- [ ] ESLint、Formatter、类型检查和 Git Hooks。
- [ ] Unit、Component 和 E2E 测试入门。
- [ ] 环境变量、配置、构建、部署和错误监控。
- [ ] 依赖选择、升级、锁文件和供应链基础。

必做：使用 TypeScript 和一种主流框架实现一个包含登录、路由、表单、CRUD、请求状态、错误处理和测试的完整应用。

验收：能从空目录搭建并部署完整应用；公共 API 类型稳定，不滥用 `any` 和类型断言；核心流程有自动化测试。

---

## 阶段五：浏览器与运行时

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

## 阶段六：应用与组件架构

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

## 阶段七：性能工程

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

## 阶段八：工程化与工具链

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

## 阶段九：测试、安全与可观测性

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

## 阶段十：专家专项

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

## 学习资源

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
