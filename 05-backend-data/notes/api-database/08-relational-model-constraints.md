# 实体、关系、主键、外键与唯一约束

## 是什么

实体是需要持久标识的业务对象；关系表达实体间关联。主键唯一且非空标识行；外键要求引用值存在；唯一约束防止候选键重复。约束由数据库在并发写入下统一执行。

## 为什么需要

应用检查存在竞态，数据不变量必须尽可能落在数据库约束。

## 关键特性或规则

主键稳定且不复用；关系表用外键；多列唯一表达业务范围，如 UNIQUE(tenant_id,slug)；明确删除更新动作；NOT NULL 同样是重要约束。

## 实际怎么使用

```sql
CREATE TABLE users(id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,email text NOT NULL UNIQUE);
CREATE TABLE orders(id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,user_id bigint NOT NULL REFERENCES users(id) ON DELETE RESTRICT);
```

## 常见错误与边界

外键引用列缺索引会使删除/更新检查昂贵；自然键变化会传播；用全局唯一 ID 仍不能代替租户范围约束。

## 补充知识

约束命名便于 migration 和错误映射；复杂不变量可用 CHECK，但跨行规则通常需事务设计。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/ddl-default.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/datatype-numeric.html)（访问日期：2026-07-16）
