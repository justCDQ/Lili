# 浏览器 DevTools：Elements、Console、Network 与 Sources

## 是什么与为什么需要

DevTools 是浏览器内置诊断工具。Elements 检查实时 DOM 和计算样式；Console 查看消息并在页面上下文执行 JavaScript；Network 记录请求、响应、时序和缓存；Sources 查看资源、设置断点和单步执行。它提供运行时证据，不能只靠阅读源码替代。

## 实际使用

1. Elements 选择节点，在 Styles 临时切换声明，在 Computed 看最终值和来源。
2. Console 运行 `document.querySelector('main')`，检查错误堆栈；不要粘贴不可信代码。
3. 打开 Network 后刷新，查看主文档的 Status、Headers、Timing、Response、Initiator。
4. Sources 在事件处理函数设断点，使用 Step over/into/out、Scope、Watch 定位状态变化。

## 关键规则

- Elements 中修改默认只存在于当前页面会话，刷新会丢失。
- Network 只记录面板打开后的活动；复现加载问题先打开再刷新。
- “Disable cache”通常仅在 DevTools 打开时生效。
- Console 当前执行上下文可能是 iframe 或扩展，不一定是顶层页面。
- 生产代码经过转译或压缩时，source map 决定能否映射到源文件。

## 常见错误与边界

404 是资源不存在，控制台的 CORS 报错还需看 Network 响应。缓存命中会使请求时序与首次访问不同。临时修改 DOM 不能修复源文件。DevTools 展现的是特定浏览器实现，跨浏览器问题需在目标浏览器复测。

## 补充知识

保留日志可跨导航观察；Network 节流是模拟，不等同真实网络。可从 Initiator 追踪是谁触发请求，从 Accessibility 树核对语义。

## 来源

- [Chrome DevTools](https://developer.chrome.com/docs/devtools/)
- [Chrome：Inspect network activity](https://developer.chrome.com/docs/devtools/network/)
- [Chrome：Console reference](https://developer.chrome.com/docs/devtools/console/reference/)
- [Chrome：Sources overview](https://developer.chrome.com/docs/devtools/sources/)

访问日期：2026-07-16。
