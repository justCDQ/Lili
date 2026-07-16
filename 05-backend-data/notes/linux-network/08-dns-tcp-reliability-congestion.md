# DNS、TCP 握手、重传、滑动窗口与拥塞控制

## 是什么

DNS 解析名称与记录；TCP 用三次握手建立双向字节流，通过序号、确认和重传提供可靠有序传输。接收窗口限制接收端容量，拥塞窗口限制网络注入量，实际发送受二者较小值约束。

## 为什么需要

连接延迟、吞吐、超时和重传问题必须区分名称解析、路径损失、接收端与网络拥塞。

## 关键特性或规则

DNS 有 TTL 与缓存；TCP 是字节流无消息边界；重传可能由超时或重复 ACK 触发；拥塞控制依据反馈调整窗口。

## 实际怎么使用

```sh
dig +trace example.com
ss -tin dst 203.0.113.10
tcpdump -ni any 'host 203.0.113.10 and tcp port 443'
```

## 常见错误与边界

DNS 命中不同 IP 不一定错误；一次 read 不对应一次 write；随意调低重传或 TIME_WAIT 参数可能破坏可靠性。

## 补充知识

Happy Eyeballs、IPv4/IPv6、CDN 和负载均衡会使不同客户端路径不同。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc1034.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc9293.html)（访问日期：2026-07-16）
- [一手资料 3](https://www.rfc-editor.org/rfc/rfc5681.html)（访问日期：2026-07-16）
- [一手资料 4](https://man7.org/linux/man-pages/man8/ss.8.html)（访问日期：2026-07-16）
