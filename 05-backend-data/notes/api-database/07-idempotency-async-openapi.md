# 幂等 Key、异步任务 API 与 OpenAPI

## 是什么

幂等 key 让服务识别同一逻辑写请求并返回首次结果；异步 API 先创建任务资源再轮询/推送状态；OpenAPI 用机器可读文档描述 HTTP 接口。

## 为什么需要

网络超时使客户端不知道写入是否成功，安全重试和长任务都需要显式协议。

## 关键特性或规则

key 有作用域、TTL 和请求摘要；相同 key 不同 body 返回冲突；首次处理与记录结果同一原子边界；任务状态含 queued/running/succeeded/failed/canceled；OpenAPI 与实现一起评审测试。

## 实际怎么使用

```http
POST /payments
Idempotency-Key: 0190...

HTTP/1.1 202 Accepted
Location: /jobs/abc
Retry-After: 3

# 服务端按 actor+route+key 原子保存请求摘要、状态和响应
```

## 常见错误与边界

只在内存去重会在重启失效；先执行业务再记录 key 有竞态；异步任务失败不能只写日志；文档生成不等于契约正确。

## 补充知识

支付等关键写入还需数据库唯一约束；Webhook/队列消费者也需独立幂等。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc7240.html)（访问日期：2026-07-16）
- [一手资料 3](https://spec.openapis.org/oas/latest.html)（访问日期：2026-07-16）
- [一手资料 4](https://datatracker.ietf.org/doc/draft-ietf-httpapi-idempotency-key-header/)（访问日期：2026-07-16）
