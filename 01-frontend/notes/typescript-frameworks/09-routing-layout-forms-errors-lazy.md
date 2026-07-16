# 路由、Layout、表单、请求、错误边界与懒加载

## 是什么

路由映射 URL 到页面，layout 复用页面骨架；表单管理输入与提交；请求状态区分 idle/loading/success/error；错误边界隔离渲染故障；懒加载按需获取代码。

## 为什么需要

完整应用需要把 URL、嵌套布局、数据加载、表单提交、错误恢复和代码分割组织成连续状态机。缺少统一路由边界会导致刷新丢失、错误无归属和加载状态重复。

## 关键特性与规则

URL 是导航状态源；layout 保留语义地标；提交防重复；请求错误可重试；路由级切包并提供加载/失败 UI。

## 实际使用

```tsx
const Settings=lazy(()=>import('./Settings.js'));
<Suspense fallback={<p>加载中</p>}><Settings/></Suspense>
```

## 常见错误与边界

只显示 spinner 会遗漏错误/空状态；错误边界通常不捕获事件异步错误；懒加载过细增加请求瀑布。

## 相关补充知识

路由级懒加载只解决代码下载，不自动解决数据瀑布。表单提交需防重复并保留输入；错误边界应区分路由未找到、权限、数据失败和渲染异常，同时提供可恢复动作。

## 来源

- [React Documentation](https://react.dev/reference/react/lazy)
- [React Documentation](https://react.dev/reference/react/Component#catching-rendering-errors-with-an-error-boundary)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/History_API)

访问日期：2026-07-16。
