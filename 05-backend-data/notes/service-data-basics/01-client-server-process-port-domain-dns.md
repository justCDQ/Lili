# 客户端、服务器、进程、端口、域名与 DNS

## 是什么

客户端发起请求；服务器程序在进程中运行并监听网络地址。IP 标识网络接口，端口在主机上区分服务端点，域名是便于管理的名称，DNS 将域名查询为 IP 等记录。一次访问通常经历 DNS 解析、建立传输连接、发送应用协议请求和接收响应。

```go
srv := &http.Server{Addr: "127.0.0.1:8080", Handler: http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") },
)}
log.Fatal(srv.ListenAndServe())
```

## 关键特性或规则

- `127.0.0.1` 仅本机可达；`0.0.0.0` 监听所有 IPv4 接口，不是客户端目标地址。
- 端口属于传输协议和网络命名空间；同一地址/端口通常只能由一个监听者占用。
- DNS 结果可缓存且有 TTL，一个域名可返回多个地址。
- 域名解析成功不代表端口开放或应用健康。
- 服务端必须设置读写、空闲和处理超时。

## 常见错误与边界

DNS 不保证立即一致；本机 hosts、缓存、代理和容器网络会改变解析路径。端口可达不证明协议正确。不要将开发服务无认证暴露到公网。

## 为什么需要

这些概念组成一次服务请求从寻址、传输、处理到持久化的最小闭环。缺少其中任一层的明确契约，都会让连接失败、协议错误、并发问题或数据不一致难以定位。

## 实际怎么使用

运行本文 Go 服务或数据库示例，使用 curl 发出正常、非法方法、错误 JSON、超大正文和并发请求。逐层记录 DNS/地址、端口、请求 Header、状态码、日志、数据变化和错误恢复，并为核心处理函数添加测试。

## 补充知识

本地成功只验证单进程与本机网络条件。进入容器、代理或远端数据库后，还要显式处理超时、连接池、取消、重试、幂等、事务边界和敏感日志。

## 来源

- [RFC 1034：Domain Names](https://www.rfc-editor.org/rfc/rfc1034)（访问日期：2026-07-16）
- [IANA：Service Name and Port Number Registry](https://www.iana.org/assignments/service-names-port-numbers/)（访问日期：2026-07-16）
- [Go：net/http Server](https://pkg.go.dev/net/http#Server)（访问日期：2026-07-16）
