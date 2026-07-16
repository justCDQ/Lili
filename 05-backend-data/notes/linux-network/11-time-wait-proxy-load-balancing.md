# TIME_WAIT、反向代理与 L4/L7 负载均衡

## 是什么

主动关闭 TCP 的一端进入 TIME_WAIT，防止旧报文污染新连接并允许重传最终 ACK。反向代理代表服务接收请求。L4 基于传输层地址/端口转发，L7 理解 HTTP 等应用协议并可按路径、Header 路由。

## 为什么需要

连接生命周期和代理层决定客户端 IP、TLS 终止、健康检查、重试与容量位置。

## 关键特性或规则

TIME_WAIT 是正确 TCP 状态，先减少短连接再考虑内核调优；代理设置连接/请求超时；健康检查区分存活与就绪；L7 重试只针对安全/幂等操作。

## 实际怎么使用

```sh
ss -ant state time-wait | wc -l
curl -H 'Host: api.example.com' http://127.0.0.1/health
# Go 仅在可信代理链解析 Forwarded/X-Forwarded-For
```

## 常见错误与边界

盲信 X-Forwarded-For 可伪造来源；代理重试 POST 会重复写；只检查进程存活可能把无依赖能力实例加入流量。

## 补充知识

TLS 可在负载均衡器终止、透传或再次加密；选择影响可观测性和信任边界。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9293.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc9210.html)（访问日期：2026-07-16）
- [一手资料 3](https://www.rfc-editor.org/rfc/rfc7239.html)（访问日期：2026-07-16）
- [一手资料 4](https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/)（访问日期：2026-07-16）
