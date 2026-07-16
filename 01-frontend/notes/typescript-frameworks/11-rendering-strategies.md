# SSR、SSG、CSR 与 Hydration

## 是什么

CSR 在浏览器构建 UI；SSR 每请求在服务器生成 HTML；SSG 构建期生成 HTML；hydration 在已有服务端 HTML 上附加客户端行为。

## 为什么需要

CSR、SSR、SSG 与 Hydration 决定 HTML 在何处何时产生、数据何时取得以及客户端如何接管交互。选择会影响首屏、缓存、服务器成本、个性化和故障模式。

## 关键特性与规则

选择取决于更新频率、首屏、SEO、个性化和基础设施；服务端与客户端初始输出必须一致；只给需交互区域发送 JS。

## 实际使用

```tsx
// 概念流程
const html=renderToString(<App/>); // server
hydrateRoot(document.getElementById('root'),<App/>); // client
```

## 常见错误与边界

hydration mismatch 会重建或报错；SSR 不自动减少客户端 JS；SSG 内容更新需要重建或再验证。

## 相关补充知识

Hydration 要求服务端首屏与客户端初次渲染一致，时间、随机数、浏览器专有状态和地区格式都可能造成不匹配。静态生成还需定义内容更新后的失效与重新生成策略。

## 来源

- [React Documentation](https://react.dev/reference/react-dom/server)
- [React Documentation](https://react.dev/reference/react-dom/client/hydrateRoot)
- [web.dev](https://web.dev/articles/rendering-on-the-web)

访问日期：2026-07-16。
