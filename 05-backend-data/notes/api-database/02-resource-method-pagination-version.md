# Resource、Method、Status Code、分页、筛选、排序、批量与版本

## 是什么

资源 API 用稳定标识符表示业务对象；HTTP 方法和状态码表达请求语义。分页限制结果集，筛选/排序明确查询，批量减少往返，版本管理不兼容演进。

## 为什么需要

清晰契约可使客户端正确缓存、重试、导航并处理兼容性。

## 关键特性或规则

GET 安全且幂等；PUT/DELETE 应幂等，POST 通常非幂等；创建成功常用 201+Location；游标包含稳定排序键和唯一 tie-breaker；限制最大 page size。

## 实际怎么使用

```http
GET /v1/orders?status=paid&sort=-created_at&limit=50&after=eyJpZCI6MTIzfQ
POST /v1/orders:batchCancel
Idempotency-Key: 0190...
```

## 常见错误与边界

offset 在大偏移时慢且并发写入会重复/遗漏；批量请求需逐项结果或原子语义；用 200 包装所有错误会破坏协议工具。

## 补充知识

不兼容变更才升主版本；新增可选字段通常向后兼容，但严格客户端仍需测试。

## 来源

- [一手资料 1](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）
- [一手资料 2](https://www.rfc-editor.org/rfc/rfc9457.html)（访问日期：2026-07-16）
- [一手资料 3](https://spec.openapis.org/oas/latest.html)（访问日期：2026-07-16）
