# 软删除风险与多租户数据库设计

## 是什么

软删除保留行并用 deleted_at/status 排除；多租户可共享表加 tenant_id、分 schema 或分数据库。两者都改变唯一约束、查询、权限和运维。

## 为什么需要

遗漏过滤会暴露已删除或其他租户数据，恢复/合规删除也必须有明确语义。

## 关键特性或规则

所有唯一和索引考虑 deleted_at/tenant_id；服务端从可信身份设置租户；后台任务同样隔离；定义恢复、硬删除、级联、保留和备份清除策略。

## 实际怎么使用

```sql
CREATE UNIQUE INDEX users_active_email_uq ON users(tenant_id,lower(email)) WHERE deleted_at IS NULL;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_orders ON orders USING (tenant_id=current_setting('app.tenant_id')::bigint);
```

## 常见错误与边界

全局查询默认漏 tenant 条件风险高；RLS 表 owner 默认可绕过；软删除使外键仍存在且表/索引膨胀；隐私删除不能只设置标记。

## 补充知识

按合规和规模选择隔离级别；共享表成本低但爆炸半径大，独立库隔离强但运维成本高。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/indexes-partial.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-16）
