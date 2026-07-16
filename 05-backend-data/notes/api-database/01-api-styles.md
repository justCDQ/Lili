# REST、RPC、gRPC、GraphQL、WebSocket、SSE 与 Webhook

## 是什么

REST 以资源和统一 HTTP 语义建接口；RPC 暴露操作；gRPC 用 Protobuf 契约和 HTTP/2 调用；GraphQL 由客户端声明字段；WebSocket 提供双向消息通道；SSE 是服务端到浏览器的文本事件流；Webhook 是服务端向订阅方回调 HTTP。

## 为什么需要

不同交互方向、契约、浏览器支持、缓存与运维要求不同，选错会增加协议和故障处理成本。

## 关键特性或规则

普通 CRUD 优先 HTTP 资源；内部强契约低延迟可评估 gRPC；聚合读取可评估 GraphQL；双向实时用 WebSocket；单向更新用 SSE；Webhook 必须签名、重试、幂等。

## 实际怎么使用

```go
func events(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    if _, err := fmt.Fprint(w, "event: ping\ndata: {}\n\n"); err != nil {
        return
    }
    flusher.Flush()
}
```

## 常见错误与边界

WebSocket 不自带消息确认/重连语义；SSE 仅文本且浏览器连接有约束；GraphQL 单端点仍需授权与成本限制；Webhook 不能假定同步成功。

流式端点还必须监听 `r.Context().Done()`、限制每个主体的连接数，并设置代理空闲超时与心跳；示例只发送单个事件，不是完整生产循环。

## 补充知识

先写通信矩阵：调用方、方向、频率、消息大小、延迟、失败与演进，再选择协议。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）
- [一手资料 2](https://grpc.io/docs/what-is-grpc/introduction/)（访问日期：2026-07-16）
- [一手资料 3](https://spec.graphql.org/)（访问日期：2026-07-16）
- [一手资料 4](https://www.rfc-editor.org/rfc/rfc6455.html)（访问日期：2026-07-16）
- [一手资料 5](https://html.spec.whatwg.org/multipage/server-sent-events.html)（访问日期：2026-07-16）
