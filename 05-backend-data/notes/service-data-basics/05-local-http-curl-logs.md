# 本地运行服务、使用 curl 请求与查看日志

## 是什么

```go
func main() {
    h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("method=%s path=%q remote=%q", r.Method, r.URL.Path, r.RemoteAddr)
        w.Header().Set("Content-Type", "application/json")
        io.WriteString(w, `{"status":"ok"}`)
    })
    log.Fatal(http.ListenAndServe("127.0.0.1:8080", h))
}
```

```sh
go run .
curl --fail-with-body --include http://127.0.0.1:8080/health
curl --verbose http://127.0.0.1:8080/health
```

`--include` 显示响应 Header，`--verbose` 显示连接与请求细节，`--fail-with-body` 在 HTTP 4xx/5xx 返回失败码同时保留响应体。

## 实际怎么使用

确认进程存在、监听地址正确、端口未占用；再检查 curl 的 DNS/连接/TLS/HTTP 输出；最后关联服务日志。日志至少含时间、级别、方法、路径、状态、耗时和请求 ID，避免记录认证 Header、Cookie 与正文秘密。

## 常见错误与边界

本地成功不证明容器、反向代理或公网可用；`localhost` 可能解析 IPv4/IPv6 不同地址；curl 默认不会替你验证业务结果。生产服务还必须设置超时、限制正文、恢复 panic 和优雅退出。

## 为什么需要

这些概念组成一次服务请求从寻址、传输、处理到持久化的最小闭环。缺少其中任一层的明确契约，都会让连接失败、协议错误、并发问题或数据不一致难以定位。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 补充知识

本地成功只验证单进程与本机网络条件。进入容器、代理或远端数据库后，还要显式处理超时、连接池、取消、重试、幂等、事务边界和敏感日志。

## 来源

- [curl Manual](https://curl.se/docs/manpage.html)（访问日期：2026-07-16）
- [Go：net/http](https://pkg.go.dev/net/http)（访问日期：2026-07-16）
- [OpenTelemetry：Logs](https://opentelemetry.io/docs/specs/otel/logs/)（访问日期：2026-07-16）
