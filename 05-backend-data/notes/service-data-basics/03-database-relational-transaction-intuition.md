# 数据库、表、行、列、主键、查询与事务

## 是什么

数据库管理持久数据及并发访问。关系表具有列定义和多行记录；主键唯一标识一行；查询读取或修改满足条件的数据；事务把一组操作作为一个一致性单元提交或回滚。

```sql
CREATE TABLE todo (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title TEXT NOT NULL,
  done BOOLEAN NOT NULL DEFAULT FALSE
);
SELECT id, title FROM todo WHERE done = FALSE ORDER BY id LIMIT 20;
```

## 关键特性或规则

- 类型、`NOT NULL`、唯一和外键约束应在数据库中保护不变量。
- 查询必须明确列和条件；无条件 UPDATE/DELETE 风险高。
- 主键负责身份，不应包含频繁变化的业务含义。
- 事务边界围绕业务不变量，保持短小；错误时回滚。
- 索引加速特定读但占空间并增加写成本。

## 常见错误与边界

事务不表示所有外部系统一起原子提交；隔离级别不同会允许不同并发现象。查询返回顺序只有 `ORDER BY` 才有保证。应用检查不能替代数据库唯一约束处理并发。

## 为什么需要

这些概念组成一次服务请求从寻址、传输、处理到持久化的最小闭环。缺少其中任一层的明确契约，都会让连接失败、协议错误、并发问题或数据不一致难以定位。

## 实际怎么使用

运行本文 Go 服务或数据库示例，使用 curl 发出正常、非法方法、错误 JSON、超大正文和并发请求。逐层记录 DNS/地址、端口、请求 Header、状态码、日志、数据变化和错误恢复，并为核心处理函数添加测试。

## 补充知识

本地成功只验证单进程与本机网络条件。进入容器、代理或远端数据库后，还要显式处理超时、连接池、取消、重试、幂等、事务边界和敏感日志。

## 来源

- [PostgreSQL：Table Basics](https://www.postgresql.org/docs/current/ddl-basics.html)（访问日期：2026-07-16）
- [PostgreSQL：Transactions](https://www.postgresql.org/docs/current/tutorial-transactions.html)（访问日期：2026-07-16）
- [PostgreSQL：Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)（访问日期：2026-07-16）
