# Schema Migration、Data Migration 与 Expand and Contract

## 是什么

schema migration 改结构/约束；data migration 转换或回填行。Expand and Contract 先添加兼容结构，双读/双写或回填，再切换读取，最后删除旧结构。

## 为什么需要

应用和数据库无法总在同一瞬间升级，分阶段迁移保证滚动部署期间新旧版本共存。

## 关键特性或规则

每个步骤可观测、可暂停、可重入；先部署兼容代码再改变约束；schema 与 data 分开；记录 migration 版本且单写执行。

## 实际怎么使用

```sql
-- expand
ALTER TABLE users ADD COLUMN display_name text;
-- 分批回填
UPDATE users SET display_name=name WHERE id>$1 AND id<=$2 AND display_name IS NULL;
-- 验证后约束/contract 在后续发布执行
```

## 常见错误与边界

同一 migration 大表 UPDATE 会长锁/WAL；立即 rename/drop 破坏旧实例；回滚 DDL 不等于能恢复已转换数据。

## 补充知识

PostgreSQL 某些 ALTER 只改元数据，另一些重写表；按当前版本官方文档确认锁级别。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/ddl-alter.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/sql-altertable.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/explicit-locking.html)（访问日期：2026-07-16）
