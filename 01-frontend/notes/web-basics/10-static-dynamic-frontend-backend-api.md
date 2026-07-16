# 静态网站、动态网站、前端、后端与 API

## 是什么

静态服务器把已有文件按请求直接返回；动态服务器运行应用逻辑，可能读取数据库并生成响应。前端是用户设备中呈现和交互的部分，后端处理服务端业务、数据与权限。API 是组件间约定的接口；Web API 常通过 HTTP 暴露数据或操作。

## 为什么需要

架构选择影响部署、缓存、安全和开发边界。静态不等于无交互：静态页面仍可运行 JavaScript 并调用 API；动态不等于每次返回 HTML，也可返回 JSON。

## 实际使用

静态响应：`GET /about.html` 直接读取文件。动态响应：`GET /users/42` 校验权限、查询数据库、返回 JSON。前端调用：

```js
const response = await fetch('/api/products');
if (!response.ok) throw new Error(`HTTP ${response.status}`);
const products = await response.json();
```

## 关键规则与边界

- 客户端代码和请求可被用户查看、修改，安全校验必须在后端执行。
- API 契约至少明确方法、URL、输入、成功响应、错误、鉴权和版本策略。
- 前后端是职责划分，不必是不同仓库或团队。
- 静态文件更易被 CDN 缓存；动态响应能否缓存取决于语义和响应头。

## 常见错误与边界

不要把 API key 写入浏览器包。`fetch` 只在网络失败时拒绝，HTTP 404/500 需检查 `response.ok`。页面内容来自数据库不代表浏览器直接连接数据库。

## 补充知识

服务端渲染在服务器生成 HTML，客户端渲染在浏览器用数据建立 UI；混合方案可同时使用。API 不限 REST，还包括 GraphQL、RPC、浏览器 Web API 等。

## 来源

- [MDN：What is a web server](https://developer.mozilla.org/en-US/docs/Learn_web_development/Howto/Web_mechanics/What_is_a_web_server)
- [MDN：Client-server overview](https://developer.mozilla.org/en-US/docs/Learn_web_development/Extensions/Server-side/First_steps/Client-Server_overview)
- [MDN：Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)

访问日期：2026-07-16。
