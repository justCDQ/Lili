---
type: ai-note
stage: beginner
topic: http-json-api
verified: 2026-07-16
tags: [ai, http, json, rest, authentication, streaming]
---

# HTTP、JSON、REST API、状态码、认证与流式响应

## 是什么

HTTP 是无状态的应用层请求/响应协议。客户端向资源发送包含方法、目标、Header 和可选内容的请求，服务器返回状态码、Header 和可选内容。JSON 是常见的数据交换格式。REST 是一组网络系统架构约束，实际模型 API 通常使用基于 HTTP 和 JSON 的接口，但不应把所有 HTTP API 都等同于严格 REST。

## 为什么需要

官方 SDK 最终仍通过网络协议访问模型服务。理解 HTTP 能直接定位认证、限流、超时、代理、流中断和响应解析问题，也能在 SDK 缺少新能力时阅读原始 API 文档。

## 关键特性

### 请求与响应

- 方法表达意图；模型生成通常使用 `POST`，因为需要发送较大的请求体并产生计算。
- Header 携带媒体类型、认证、追踪和内容协商信息。
- `Content-Type: application/json` 表示内容按 JSON 解释。
- HTTP 无状态意味着每个请求语义可独立理解；会话历史通常由客户端在请求中显式传入或由应用服务保存。

### 状态码

- `2xx`：请求被成功处理，但业务结果仍需验证。
- `4xx`：请求、认证、权限、配额或速率限制问题；不是所有 `4xx` 都可重试。
- `5xx`：服务端失败，可能适合有限重试。
- 必须读取服务提供方的错误结构和 `Retry-After` 等 Header，不能只看状态码文本。

### 认证

API Key 或短期 Token 通常放在 `Authorization` 等 Header。认证回答调用者是谁，授权决定调用者能做什么。密钥不得进入 URL、前端代码和日志。

### Streaming

流式响应让客户端在完整结果生成前逐块处理数据。HTTP 连接成功不代表流一定完整；客户端必须处理分片边界、解码、取消、网络中断、重复事件和结束标记。

## 实际怎么使用

```bash
curl --fail-with-body \
  --request POST \
  --header "Authorization: Bearer $MODEL_API_KEY" \
  --header "Content-Type: application/json" \
  --data '{"model":"MODEL_ID","input":"Return one JSON object"}' \
  https://provider.example/v1/generate
```

程序处理顺序：

1. 设置连接和总超时。
2. 发送认证、内容类型、请求 ID 和请求体。
3. 检查 HTTP 状态并解析厂商错误结构。
4. 非流式响应解析 JSON，再做 Schema 校验。
5. 流式响应逐块解码，按协议组装事件；完成后检查结束原因和 Usage。
6. 日志记录请求 ID、模型、延迟和错误类别，不记录 Secret 和敏感原文。

## 常见错误与边界

- `fetch` 获得 `404` 或 `500` 时仍可能正常返回 Response，必须检查 `response.ok`。
- 把网络超时、服务端错误、模型拒绝和 Schema 错误归为同一种失败。
- 对 `401`、`403` 或无效参数自动重试，既无效又增加请求量。
- 假设一个网络 Chunk 就是一条完整 JSON 事件；传输分块和应用事件边界不同。
- 只保存最后文本，丢失结束原因、工具调用、引用、Usage 和服务端请求 ID。

## 补充知识

HTTP/1.1、HTTP/2 和 HTTP/3 使用不同传输方式，但共享 RFC 9110 定义的核心语义。代理、网关和负载均衡可能增加超时和缓冲，流式接口部署时需要单独验证链路。

## 来源

- [RFC 9110：HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）
- [MDN：Using Fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)（访问日期：2026-07-16）
- [MDN：JSON](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Scripting/JSON)（访问日期：2026-07-16）
- [MDN：Streams API](https://developer.mozilla.org/en-US/docs/Web/API/Streams_API)（访问日期：2026-07-16）

