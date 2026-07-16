# EXPLAIN、EXPLAIN ANALYZE 与慢查询

## 是什么

EXPLAIN 展示优化器计划和估算；ANALYZE 实际执行并显示真实行数/时间；慢查询治理从采样、归一化、计划、等待与数据分布查根因。

## 为什么需要

只看 SQL 文本或总时间无法区分估算错误、I/O、锁等待、CPU、网络或返回数据过多。

## 关键特性或规则

先在安全数据副本或只读事务评估；比较 estimated rows 与 actual rows；看 loops、buffers、sort spill、scan；同时观察锁和 I/O。

## 实际怎么使用

```sql
EXPLAIN (ANALYZE, BUFFERS, TIMING OFF, FORMAT TEXT)
SELECT id, created_at, total
FROM orders
WHERE tenant_id = 42 AND status = 'paid'
ORDER BY created_at DESC, id DESC
LIMIT 50;
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
SELECT query,calls,total_exec_time,mean_exec_time,rows FROM pg_stat_statements ORDER BY total_exec_time DESC LIMIT 20;
```

## 常见错误与边界

EXPLAIN ANALYZE 会真实执行 UPDATE/DELETE；生产执行重查询有风险；单次快不代表 p99；计划会因参数、统计和缓存改变。

示例中的 `CREATE EXTENSION` 需要相应权限，应由数据库变更流程执行，不应由应用请求动态执行。生产运行 `EXPLAIN ANALYZE` 前先在只读副本或等量测试环境验证成本，并设置 `statement_timeout`。

## 补充知识

使用 auto_explain 或 pg_stat_statements 需评估开销和敏感 SQL 文本；优化后用业务负载验证。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/sql-explain.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/using-explain.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/pgstatstatements.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/auto-explain.html)（访问日期：2026-07-16）
