# 向后兼容、零停机、大表回填与回滚

## 是什么

零停机迁移让服务持续处理请求；大表回填按稳定键小批更新并节流；回滚包括应用回退、schema 前滚修复和数据恢复，不只是执行 down SQL。

## 为什么需要

大事务、强锁和不兼容字段会造成请求阻塞、复制延迟或无法回退。

## 关键特性或规则

发布前测量锁、WAL、磁盘和复制延迟；回填幂等并保存 cursor；设置 lock_timeout；新旧列过渡期定义写入真相；删除列延后多个版本。

## 实际怎么使用

```sql
-- 先 NOT VALID，降低建立约束时影响，再验证
ALTER TABLE orders ADD CONSTRAINT total_nonnegative CHECK(total>=0) NOT VALID;
ALTER TABLE orders VALIDATE CONSTRAINT total_nonnegative;
-- worker: WHERE id > cursor ORDER BY id LIMIT batch FOR UPDATE SKIP LOCKED
```

## 常见错误与边界

DDL 事务等待锁时也会阻塞后续锁队列；按 offset 回填越来越慢；双写顺序失败会漂移；down migration 可能不可逆。

## 补充知识

高风险迁移预演生产规模数据，并准备 kill switch、暂停点和数据一致性查询。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/sql-altertable.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/explicit-locking.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/queries-select-lists.html)（访问日期：2026-07-16）
