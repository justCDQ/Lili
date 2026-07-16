# 统一错误：输入、认证、权限、不存在、冲突、限流、系统与外部依赖

## 是什么

统一错误模型把机器可读 code、HTTP status、用户可读 message、字段详情、request ID 和可选重试信息分开。错误类别对应不同恢复动作。

## 为什么需要

客户端才能决定修正输入、登录、申请权限、刷新、退避或联系支持，服务端才能聚合指标。

## 关键特性或规则

400/422 输入，401 缺少/无效认证，403 已认证但禁止，404 不存在，409 状态冲突，429 限流，5xx 服务故障；内部堆栈只写受控日志。

## 实际怎么使用

```go
type Problem struct { Type string `json:"type"`; Title string `json:"title"`; Status int `json:"status"`; Code string `json:"code"`; RequestID string `json:"request_id"` }
// Content-Type: application/problem+json
```

## 常见错误与边界

把数据库错误原样返回会泄漏结构；错误 message 作为稳定程序接口会破坏本地化；外部依赖超时不应都映射 500 且无限重试。

## 补充知识

RFC 9457 定义 Problem Details，可扩展业务 code 和字段错误但保持标准成员语义。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9457.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）
- [一手资料 3](https://pkg.go.dev/errors)（访问日期：2026-07-16）
