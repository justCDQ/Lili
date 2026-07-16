# 客户端、服务器、浏览器与搜索引擎

## 是什么与为什么需要

客户端发起网络请求；服务器监听请求并返回资源或计算结果。浏览器是 Web 客户端，负责请求、解析、渲染和执行页面。搜索引擎通过爬取、索引、排序提供检索，搜索结果页不是浏览器也不是目标网站。区分角色有助于判断故障发生在页面代码、网络、服务端还是索引系统。

## 关键特性与规则

- 同一程序可在一次通信中是客户端、另一次是服务器；角色由交互决定。
- Web 服务器既可指机器，也可指处理 HTTP 的软件。
- 浏览器先请求 HTML，再根据文档继续请求 CSS、脚本、图片等资源。
- 搜索引擎能否发现页面取决于链接、抓取许可、可访问响应与索引策略；提交 URL 不保证收录。

## 实际使用

在 Network 刷新页面：第一行通常是主 HTML 文档，查看其状态；再检查后续资源。直接输入 URL 是导航；在搜索框输入关键词是向搜索服务查询，两者路径不同。

## 常见错误与边界

“服务器正常”不代表应用响应正确；需检查状态、响应体和控制台。浏览器缓存可能隐藏服务端变化。搜索结果摘要可能过时。浏览器也可展示 PDF、图片等非 HTML 资源。

## 补充知识

互联网是通信基础设施，Web 是其上的服务；邮件等服务不属于 Web。一个网站可由 CDN、反向代理、应用服务器和数据库共同提供。

## 来源

- [MDN：Browsing the web](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Environment_setup/Browsing_the_web)
- [MDN：How the web works](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Web_standards/How_the_web_works)
- [MDN：What is a web server](https://developer.mozilla.org/en-US/docs/Learn_web_development/Howto/Web_mechanics/What_is_a_web_server)

访问日期：2026-07-16。
