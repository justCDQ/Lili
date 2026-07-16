# 避免长事务与事务中调用慢外部接口

## 是什么

事务从 BEGIN 到 COMMIT 持有快照、锁、连接和未完成 WAL 语义。外部 HTTP/消息调用延迟不可控且通常不能与数据库原子提交。

## 为什么需要

长事务增加锁等待、表膨胀、死锁和连接池耗尽；提交失败时外部副作用无法自动撤销。

## 关键特性或规则

事务内仅做必要数据库操作；先准备外部输入；设置事务/语句超时；用 transactional outbox 把待发布事件与业务写同事务。

## 实际怎么使用

```go
err := withTx(ctx,db,func(tx *sql.Tx) error {
 if _,err:=tx.ExecContext(ctx,`UPDATE orders SET status='paid' WHERE id=$1`,id);err!=nil{return err}
 _,err:=tx.ExecContext(ctx,`INSERT INTO outbox(topic,payload) VALUES($1,$2)`,`order.paid`,payload)
 return err
})
// 提交后独立 worker 发布 outbox
```

## 常见错误与边界

提交后直接发消息会丢事件；先发消息再提交会产生幽灵事件；Saga/补偿不是强原子性，必须幂等和可观测。

## 补充知识

监控 oldest transaction、锁等待和连接占用；批处理拆成可恢复小批次。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/tutorial-transactions.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/routine-vacuuming.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://pkg.go.dev/database/sql#Tx)（访问日期：2026-07-16）
