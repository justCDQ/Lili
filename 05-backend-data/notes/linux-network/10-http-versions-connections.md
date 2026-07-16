# HTTP/1.1、HTTP/2、HTTP/3、Keepalive 与 Connection Pool

## 是什么

HTTP/1.1 使用文本消息并默认持久连接；HTTP/2 在单个 TCP 连接多路复用二进制流并压缩字段；HTTP/3 把 HTTP 语义映射到 QUIC。keepalive 复用连接，连接池管理空闲/活跃连接上限和生命周期。

## 为什么需要

连接建立、TLS 和拥塞窗口有成本，复用与多路复用影响延迟、资源和故障范围。

## 关键特性或规则

HTTP 版本不改变方法和状态语义；响应体必须读完并关闭才便于 Go 复用；池设置 idle/max per host/timeouts；HTTP/2 仍有 TCP 层队头阻塞。

## 实际怎么使用

```sh
curl --http1.1 -v https://example.com/
curl --http2 -v https://example.com/
curl --http3 -v https://example.com/
# Go: 为 http.Client 复用一个 Transport，不要每请求新建
```

## 常见错误与边界

keep-alive 不是无限存活；池过大耗 fd/服务端资源，过小排队；连接复用下 DNS 变化不会立刻生效。

## 补充知识

HTTP keepalive 与 TCP keepalive 目的不同；前者复用应用连接，后者探测长时间空闲 TCP 对端。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9112.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc9113.html)（访问日期：2026-07-16）
- [一手资料 3](https://www.rfc-editor.org/rfc/rfc9114.html)（访问日期：2026-07-16）
- [一手资料 4](https://pkg.go.dev/net/http#Transport)（访问日期：2026-07-16）
