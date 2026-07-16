# Request ID、Logging、Recovery、Timeout、CORS、Rate Limit、Metrics 与 Trace

## 是什么

这些横切能力建立请求身份、结构化事件、panic 隔离、截止时间、跨源策略、流量保护和可观测信号。middleware 按顺序包裹 handler。

## 为什么需要

可靠服务必须能限制资源、定位失败并量化延迟/错误，而不是只完成业务处理。

## 关键特性或规则

request ID 接受外部值前校验长度/字符或重生成；日志结构化且脱敏；recovery 返回 500 并记录栈；超时向下传 context；指标低基数；trace 传播标准上下文。

## 实际怎么使用

```go
handler:=requestID(logging(recoverer(timeout(rateLimit(cors(mux))))))
srv:=&http.Server{Addr:":8080",Handler:handler,ReadHeaderTimeout:5*time.Second,IdleTimeout:60*time.Second}
```

## 常见错误与边界

recover 不能修复数据不变量；仅 handler timeout 不限制慢读；CORS 不是认证；按 user/IP 限流需考虑代理和 NAT；把 request ID 当 metric label 会爆基数。

## 补充知识

区分日志事件、指标聚合和 trace 因果链，三者通过 request/trace ID 关联。

## 来源

- [一手资料 1](https://pkg.go.dev/net/http)（访问日期：2026-07-16）
- [一手资料 2](https://opentelemetry.io/docs/specs/otel/)（访问日期：2026-07-16）
- [一手资料 3](https://www.w3.org/TR/trace-context/)（访问日期：2026-07-16）
- [一手资料 4](https://fetch.spec.whatwg.org/#http-cors-protocol)（访问日期：2026-07-16）
- [一手资料 5](https://www.rfc-editor.org/rfc/rfc6585.html)（访问日期：2026-07-16）
