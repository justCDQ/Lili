# 规范化、反规范化、审计字段、状态与历史

## 是什么

规范化通过拆表减少更新异常；反规范化复制/预计算数据换取读取性能。审计字段记录创建/修改主体和时间；状态表示生命周期；历史表或事件保存过去版本。

## 为什么需要

模型决定一致性边界和变更可追溯性，必须在读写模式和约束基础上设计。

## 关键特性或规则

先规范化并以 profile 证明反规范化收益；复制字段必须定义事实来源和更新机制；时间用 timestamptz；状态转换在事务中校验并记录。

## 实际怎么使用

```sql
CREATE TYPE order_status AS ENUM('pending','paid','canceled');
CREATE TABLE order_status_history(order_id bigint REFERENCES orders,id bigint GENERATED ALWAYS AS IDENTITY,from_status order_status,to_status order_status NOT NULL,changed_at timestamptz NOT NULL DEFAULT now(),PRIMARY KEY(id));
```

## 常见错误与边界

updated_at 不能回答谁改了什么；枚举新增/删除有迁移成本；历史表无限增长需保留和索引策略；反规范化双写会漂移。

## 补充知识

事件溯源与普通审计历史不同，前者以事件重建状态，复杂度更高。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/ddl-basics.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/datatype-enum.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/datatype-datetime.html)（访问日期：2026-07-16）
