# ACID、隔离级别、MVCC、锁与死锁

## 是什么

原子性保证事务全成或全败；一致性依赖约束和正确事务；隔离控制并发可见性；持久性保证提交后恢复。MVCC 保存行版本，让读写减少互阻；锁保护冲突资源；死锁是循环等待。

## 为什么需要

并发正确性不能靠单请求测试，必须知道隔离级别允许的现象并处理重试。

## 关键特性或规则

PostgreSQL 默认 Read Committed 每语句新快照；Repeatable Read/Serializable 更强但可能 abort；所有事务保持固定锁顺序；捕获 serialization/deadlock 错误并重试整个事务。

## 实际怎么使用

```sql
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
SELECT balance FROM accounts WHERE id=$1 FOR UPDATE;
UPDATE accounts SET balance=balance-$2 WHERE id=$1 AND balance >= $2;
COMMIT;
```

## 常见错误与边界

MVCC 不等于无锁；长事务阻止 vacuum 清理旧版本；重试单条语句会破坏事务语义；锁等待无超时可拖垮连接池。

应用必须检查条件 `UPDATE` 的受影响行数；0 行可能表示余额不足或记录不存在。真实转账还必须在同一事务增加目标账户余额、记录账务流水，并对整个可重试事务使用稳定幂等键。

## 补充知识

设置 statement_timeout、lock_timeout、idle_in_transaction_session_timeout 并监控 pg_stat_activity。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/mvcc-intro.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/transaction-iso.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/explicit-locking.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/runtime-config-client.html)（访问日期：2026-07-16）
