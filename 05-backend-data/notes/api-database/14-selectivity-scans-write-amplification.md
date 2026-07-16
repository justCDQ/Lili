# Selectivity、Index Scan、Full Table Scan 与 Write Amplification

## 是什么

选择性描述条件能排除多少行。index scan 通过索引定位 heap 行，sequential scan 顺序读表。每个索引让 INSERT/UPDATE/DELETE 额外维护页面、WAL 和缓存，形成写放大。

## 为什么需要

索引不是越多越快，优化器按估算成本选择扫描方式。

## 关键特性或规则

低选择性且返回大量行时顺序扫描可能更快；统计信息需新鲜；更新索引列产生更多维护；批量写入前评估索引/WAL/复制影响。

## 实际怎么使用

```sql
SELECT attname,n_distinct,most_common_vals FROM pg_stats WHERE tablename='orders';
SELECT idx_scan,idx_tup_read,idx_tup_fetch FROM pg_stat_user_indexes WHERE relname='orders';
```

## 常见错误与边界

看到 Seq Scan 不代表缺索引；小表顺序扫描合理；仅按 idx_scan=0 删除索引会误删约束或低频关键索引。

## 补充知识

扩展统计可改善相关列的基数估算；autovacuum/analyze 参数需按大表变化率调整。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/using-explain.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/planner-stats.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/monitoring-stats.html)（访问日期：2026-07-16）
