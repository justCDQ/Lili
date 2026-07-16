# B+Tree、复合索引、最左前缀与覆盖索引

## 是什么

PostgreSQL B-tree 是平衡树索引，适合等值与范围排序。复合索引按列顺序组织；前导列约束最能缩小扫描范围。INCLUDE 保存非键列，使满足可见性条件时可能 index-only scan。

## 为什么需要

索引设计必须与实际 WHERE、JOIN、ORDER BY 和返回列匹配。

## 关键特性或规则

等值列通常置前，范围/排序列随后；一条复合索引不等于所有单列组合；INCLUDE 列不参与搜索和唯一性；索引列序依据查询矩阵和数据分布。

## 实际怎么使用

```sql
CREATE INDEX orders_tenant_status_created_idx ON orders(tenant_id,status,created_at DESC,id DESC) INCLUDE(total);
```

## 常见错误与边界

覆盖索引增大存储和写成本；index-only scan 仍依赖 visibility map；最左前缀是简化规则，优化器也可能 skip scan。

## 补充知识

用 pg_stat_user_indexes 检查使用情况，删除索引前核对所有查询与约束依赖。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/indexes-types.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/indexes-multicolumn.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/indexes-index-only-scans.html)（访问日期：2026-07-16）
