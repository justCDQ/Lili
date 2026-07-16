# 乐观锁、悲观锁、唯一约束冲突与库存并发

## 是什么

乐观锁用版本/条件更新检测冲突；悲观锁先 SELECT FOR UPDATE 锁行；唯一约束在数据库原子判重。库存扣减可用条件 UPDATE，保证不为负。

## 为什么需要

读后写存在 lost update 和超卖竞态，应用内 mutex 无法覆盖多实例。

## 关键特性或规则

冲突少、工作短用乐观；必须串行且锁集合小可悲观；库存不变量写入单条条件更新或锁定事务；约束冲突映射 409。

## 实际怎么使用

```sql
UPDATE products SET stock=stock-$2,version=version+1
WHERE id=$1 AND stock >= $2 AND version=$3;
-- RowsAffected=0 表示库存不足或版本冲突，需区分/重读
INSERT INTO users(email) VALUES($1) ON CONFLICT(email) DO NOTHING;
```

## 常见错误与边界

先 SELECT 检查再 UPDATE 会竞态；持行锁等待用户输入会长事务；无限重试在热点下形成风暴；ON CONFLICT 不应吞掉未知冲突。

## 补充知识

热点库存可分段、排队或预留，但必须定义超时释放和最终事实来源。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/explicit-locking.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/sql-update.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/sql-insert.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-16）
