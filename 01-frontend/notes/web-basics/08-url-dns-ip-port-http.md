# URL、域名、DNS、IP、端口与 HTTP

## 是什么

URL 标识资源，典型结构为 `scheme://host:port/path?query#fragment`。域名是便于管理的层级名称；DNS 把名称解析为 IP 地址。IP 标识网络接口，端口标识主机上的服务。HTTP 定义客户端请求与服务器响应的语义。

## 为什么需要

请求失败、跨域、部署与 API 调试都需要拆解 URL 和通信阶段。域名可解析不代表端口可达；连接成功不代表 HTTP 响应成功。

## 关键规则

- `https` 默认端口 443，`http` 默认 80；显式非默认端口属于 origin 的组成部分。
- URL fragment 不随 HTTP 请求发送，浏览器在本地定位文档片段。
- DNS 可能返回多个 IPv4/IPv6 地址并被缓存。
- 请求含方法、目标、头和可选正文；响应含状态码、头和可选正文。
- 常见状态：200 成功、301 永久重定向、400 请求有误、403 拒绝、404 不存在、500 服务端错误。

## 实际使用

```sh
nslookup example.com
curl -i https://example.com/
curl -i -X POST -H 'Content-Type: application/json' -d '{"name":"Li"}' https://api.example.test/users
```

在 Network 查看 Request URL、Remote Address、Method、Status、Request/Response Headers 和 Timing。

## 常见错误与边界

DNS 不保存网页；IP 不一定唯一对应一个站点。查询参数顺序和重复键的业务含义由服务端决定。HTTP 是无状态协议，但 Cookie 等机制可建立会话。状态 200 的响应体仍可能是业务错误。

## 补充知识

Origin 由 scheme、host、port 组成，路径不同仍是同源。URL 中非 ASCII 字符会经过编码；不要手工重复编码。

## 来源

- [MDN：How the web works](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Web_standards/How_the_web_works)
- [WHATWG：URL Standard](https://url.spec.whatwg.org/)
- [MDN：Overview of HTTP](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Overview)

访问日期：2026-07-16。
