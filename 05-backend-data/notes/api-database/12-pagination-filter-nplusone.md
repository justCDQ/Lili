# 分页、游标分页、复杂筛选与 N+1

## 是什么

offset 分页跳过前 N 行；游标分页从稳定排序键继续。复杂筛选组合可选条件。N+1 是先取 N 个对象再逐个查询关联，导致请求数线性增长。

## 为什么需要

数据量和并发写入会使分页变慢或漂移；N+1 放大数据库往返与连接占用。

## 关键特性或规则

游标包含全部排序键和唯一 tie-breaker；查询条件与索引顺序匹配；限制 page size；用查询计数测试发现 N+1。

## 实际怎么使用

```sql
SELECT id,created_at,total FROM orders
WHERE tenant_id=$1 AND status=ANY($2)
AND (created_at,id)<($3,$4)
ORDER BY created_at DESC,id DESC LIMIT $5;
-- 关联数据一次 JOIN 或 WHERE user_id=ANY($1) 批量取
```

## 常见错误与边界

offset 越大扫描/丢弃越多；只用 created_at 游标会遗漏同时间行；动态 ORDER BY 标识符不能用普通参数替代，需白名单。

## 补充知识

游标应不透明、签名或校验，避免客户端注入任意查询状态。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/queries-limit.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/functions-comparisons.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/using-explain.html)（访问日期：2026-07-16）
