# 前端模块学习笔记

按 [前端 Roadmap](../roadmap.md) 的知识点组织。以下索引中，一条路线图勾选项对应一篇可独立阅读的笔记。

## 阶段零：计算机与 Web 入门

### 开发环境

- [文件、目录、路径、扩展名、文本编码与压缩包](00-web-basics/01-files-paths-encoding-archives.md)
- [VS Code：文件、搜索、终端、插件与调试](00-web-basics/02-vscode-workspace.md)
- [命令行文件与目录操作](00-web-basics/03-command-line-file-operations.md)
- [Git 核心工作流](00-web-basics/04-git-core-workflow.md)
- [Node.js 与包管理](00-web-basics/05-node-package-management.md)
- [浏览器 DevTools 面板](00-web-basics/06-browser-devtools-panels.md)

### Web 如何工作

- [客户端、服务器、浏览器与搜索引擎](00-web-basics/07-clients-servers-browsers-search-engines.md)
- [URL、域名、DNS、IP、端口与 HTTP](00-web-basics/08-url-dns-ip-port-http.md)
- [HTML、CSS 与 JavaScript 的职责](00-web-basics/09-html-css-javascript-responsibilities.md)
- [静态网站、动态网站、前端、后端与 API](00-web-basics/10-static-dynamic-frontend-backend-api.md)
- [本地服务器、构建、部署与 HTTPS](00-web-basics/11-local-server-build-deploy-https.md)

## 阶段一：HTML 与内容结构

### 基础语法

- [HTML 基础语法](01-html/01-html-syntax.md)
- [常见内容元素](01-html/02-common-content-elements.md)
- [文档 head 与元数据](01-html/03-document-head-metadata.md)
- [块级与行内内容](01-html/04-block-inline-content.md)

### 语义、表单与可访问性

- [语义页面区域](01-html/05-semantic-page-regions.md)
- [表单控件、校验、自动填充与提交](01-html/06-forms-controls-validation-autofill-submit.md)
- [表单标签、分组、错误与帮助](01-html/07-form-labels-groups-errors-help.md)
- [标题、替代文本、焦点与键盘](01-html/08-headings-alt-focus-keyboard.md)
- [Dialog、Popover 与 Details](01-html/09-native-interactive-elements.md)
- [响应式图片与懒加载](01-html/10-responsive-images.md)
- [SEO、Open Graph 与结构化内容](01-html/11-seo-open-graph-structured-data.md)
- [ARIA 的使用边界](01-html/12-aria-boundaries.md)

## 阶段二：CSS 与视觉实现

### 基础与布局

- [选择器、声明、单位、颜色、背景、边框与字体](02-css/01-css-syntax-values.md)
- [Cascade、Specificity、Inheritance 与默认样式](02-css/02-cascade-specificity-inheritance.md)
- [Box Model、Normal Flow、display、溢出与滚动](02-css/03-box-model-flow-display.md)
- [Position、Containing Block、BFC 与 Stacking Context](02-css/04-position-containing-bfc-stacking.md)
- [Flexbox、Grid、多列布局与居中](02-css/05-flex-grid-multicol-centering.md)
- [响应式设计、媒体查询、移动优先与断点](02-css/06-responsive-media-mobile-first.md)

### 系统化与现代 CSS

- [Custom Properties、Design Token 与主题切换](02-css/07-custom-properties-tokens-themes.md)
- [Cascade Layers、Container Queries 与 Logical Properties](02-css/08-layers-container-logical.md)
- [深色模式、RTL、高对比度与打印样式](02-css/09-color-schemes-rtl-print.md)
- [Transition、Animation、Transform 与关键帧](02-css/10-transition-animation-transform.md)
- [prefers-reduced-motion 与动效性能](02-css/11-reduced-motion-performance.md)
- [BEM、CSS Modules、CSS-in-JS 与 Utility CSS](02-css/12-css-methodologies.md)

## 阶段三：JavaScript 编程基础到语言深入

### 编程基础

- [值、变量、类型、运算符、表达式与类型转换](03-javascript/01-values-types-conversion.md)
- [条件、循环、函数、参数、返回值与递归](03-javascript/02-control-flow-functions-recursion.md)
- [Array、Object、String、Map、Set 与 Date](03-javascript/03-built-in-collections-date.md)
- [解构、展开、模板字符串与可选链](03-javascript/04-modern-expression-syntax.md)
- [模块、导入导出、错误与异常处理](03-javascript/05-modules-errors.md)
- [伪代码、断点与日志调试](03-javascript/06-problem-solving-debugging.md)

### DOM、事件与网络

- [DOM 树与操作](03-javascript/07-dom-manipulation.md)
- [Event、传播、委托、默认行为与自定义事件](03-javascript/08-events.md)
- [表单状态、校验、URL 与 Web Storage](03-javascript/09-form-url-storage.md)
- [JSON、Fetch、HTTP 方法、状态码与请求错误](03-javascript/10-json-fetch-http-errors.md)
- [同步、异步、Callback、Promise 与 async/await](03-javascript/11-async-promises.md)

### 语言深入

- [Scope、Closure、Execution Context 与 Hoisting](03-javascript/12-scope-closure-context-hoisting.md)
- [Prototype、原型链、Class、this 与对象模型](03-javascript/13-object-model-this.md)
- [Iterator、Generator、Proxy 与 Reflect](03-javascript/14-iterators-generators-proxy-reflect.md)
- [Event Loop、Call Stack、Task 与 Microtask](03-javascript/15-event-loop-tasks.md)
- [Promise 组合、Async Iterator 与 AbortController](03-javascript/16-promise-composition-async-iterator-abort.md)
- [Streams、Structured Clone 与 Transferable Object](03-javascript/17-streams-structured-clone-transfer.md)
- [GC、WeakMap、WeakSet 与常见内存泄漏](03-javascript/18-gc-weak-collections-leaks.md)

## 阶段四：TypeScript、框架与完整应用

### TypeScript

- [基础类型、函数、对象、接口、类型别名与枚举](04-typescript-frameworks/01-ts-basic-types.md)
- [Union、Intersection、Narrowing 与 Generic](04-typescript-frameworks/02-ts-unions-generics.md)
- [Conditional、Mapped 与 Template Literal Types](04-typescript-frameworks/03-ts-type-transformations.md)
- [Discriminated Union、Type Predicate 与 infer](04-typescript-frameworks/04-ts-discriminated-predicate-infer.md)
- [Declaration File、Module Resolution、Compiler Options 与 Project References](04-typescript-frameworks/05-ts-declarations-resolution-projects.md)
- [Runtime Schema Validation](04-typescript-frameworks/06-runtime-validation.md)

### 框架基础

- [组件、Props、State、事件、条件、列表与生命周期](04-typescript-frameworks/07-framework-components.md)
- [响应式更新与渲染模型](04-typescript-frameworks/08-framework-rendering-model.md)
- [路由、Layout、表单、请求、错误边界与懒加载](04-typescript-frameworks/09-routing-layout-forms-errors-lazy.md)
- [客户端、服务端、URL 与持久化状态](04-typescript-frameworks/10-state-categories.md)
- [SSR、SSG、CSR 与 Hydration](04-typescript-frameworks/11-rendering-strategies.md)
- [框架掌握与 React、Vue、Svelte 比较](04-typescript-frameworks/12-framework-selection.md)

### 应用工程基础

- [Vite 等开发与构建工具](04-typescript-frameworks/13-vite-build-tools.md)
- [ESLint、Formatter、类型检查与 Git Hooks](04-typescript-frameworks/14-quality-gates.md)
- [Unit、Component 与 E2E 测试](04-typescript-frameworks/15-testing-levels.md)
- [环境变量、配置、构建、部署与错误监控](04-typescript-frameworks/16-config-deploy-monitoring.md)
- [依赖选择、升级、锁文件与供应链](04-typescript-frameworks/17-dependencies-lockfiles-supply-chain.md)

## 阶段五：浏览器与运行时

- [从 DNS 到 HTTP：浏览器请求的连接路径与诊断](05-browser-runtime/01-dns-tcp-tls-http.md)
- [HTML Parser 与 Speculative Parser](05-browser-runtime/02-html-parser-preload-scanner.md)
- [CSS 与 Script 的阻塞行为](05-browser-runtime/03-css-script-blocking.md)
- [preload、modulepreload、prefetch 与 preconnect](05-browser-runtime/04-resource-hints.md)
- [HTTP Cache 与 Service Worker Cache](05-browser-runtime/05-http-service-worker-cache.md)
- [Network 瀑布图](05-browser-runtime/06-network-waterfall.md)
- [Style、Layout、Paint 与 Composite](05-browser-runtime/07-rendering-pipeline.md)
- [主线程、合成线程与光栅线程](05-browser-runtime/08-rendering-threads-raster.md)
- [Long Task 与 Layout Thrashing](05-browser-runtime/09-long-task-layout-thrashing.md)
- [requestAnimationFrame 与任务切片](05-browser-runtime/10-raf-task-yielding.md)
- [GPU 合成层](05-browser-runtime/11-gpu-layers.md)
- [JavaScript Heap 与分配分析](05-browser-runtime/12-heap-allocation-profiling.md)
- [Detached DOM 与 Retainer Path](05-browser-runtime/13-detached-dom-retainer.md)
- [事件、定时器、Worker 与 Blob URL 泄漏](05-browser-runtime/14-resource-leaks.md)
- [无界缓存](05-browser-runtime/15-unbounded-cache.md)
- [内存泄漏实验](05-browser-runtime/16-memory-leak-lab.md)
- [卡顿实验](05-browser-runtime/17-jank-lab.md)
- [性能优化前后报告](05-browser-runtime/18-performance-before-after-report.md)

## 阶段六：应用与组件架构

- [单一职责与组合](06-application-architecture/01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled 组件](06-application-architecture/02-controlled-uncontrolled.md)
- [Headless 与 Compound Component](06-application-architecture/03-headless-compound-components.md)
- [State Machine 与 Context Boundary](06-application-architecture/04-state-machine-context-boundary.md)
- [依赖倒置与 API 稳定性](06-application-architecture/05-dependency-inversion-api-stability.md)
- [基础、通用与业务组件边界](06-application-architecture/06-component-layer-boundaries.md)
- [URL State](06-application-architecture/07-url-state.md)
- [Server State](06-application-architecture/08-server-state.md)
- [Form State](06-application-architecture/09-form-state.md)
- [Interaction State](06-application-architecture/10-interaction-state.md)
- [Global State](06-application-architecture/11-global-state.md)
- [Persistent State](06-application-architecture/12-persistent-state.md)
- [状态所有权与单一数据源](06-application-architecture/13-state-ownership-single-source.md)
- [按业务领域组织模块](06-application-architecture/14-organize-by-domain.md)
- [单向依赖规则](06-application-architecture/15-one-way-dependency-rules.md)
- [封装第三方基础设施](06-application-architecture/16-wrap-third-party-infrastructure.md)
- [统一错误、请求、缓存与重试](06-application-architecture/17-unified-error-request-cache-retry.md)
- [Architecture Decision Record](06-application-architecture/18-architecture-decision-records.md)
- [循环依赖检查](06-application-architecture/19-circular-dependency-checks.md)

## 覆盖表

| 路线阶段 | 路线图知识点 | 笔记数 | 状态 |
| --- | ---: | ---: | --- |
| 阶段零：计算机与 Web 入门 | 11 | 11 | 完成 |
| 阶段一：HTML 与内容结构 | 12 | 12 | 完成 |
| 阶段二：CSS 与视觉实现 | 12 | 12 | 完成 |
| 阶段三：JavaScript | 18 | 18 | 完成 |
| 阶段四：TypeScript、框架与应用工程 | 17 | 17 | 完成 |
| 阶段五：浏览器与运行时 | 18 | 18 | 完成 |
| 阶段六：应用与组件架构 | 19 | 19 | 完成 |
| 合计 | 107 | 107 | 完成 |

## 维护约定

- 文件名使用稳定主题，不按日期命名。
- 每篇包含“是什么、为什么需要、关键规则、实际使用、常见错误与边界、补充知识、来源”。
- 修改技术结论时同步核对官方规范或官方文档，并更新访问日期。
- 每完成后续路线图知识点，继续按“一条勾选项一篇”补充本索引。
