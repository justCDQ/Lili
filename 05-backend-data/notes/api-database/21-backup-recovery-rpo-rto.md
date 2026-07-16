# 备份、恢复、RPO 与 RTO

## 是什么

备份保存可恢复的数据副本；恢复验证副本能重建服务。RPO 是可接受的数据丢失时间窗口，RTO 是可接受的恢复耗时。PostgreSQL 可做逻辑备份、物理 base backup 和 WAL 持续归档/PITR。

## 为什么需要

没有经过恢复演练的备份不能证明可用，RPO/RTO 决定备份频率、复制和恢复架构。

## 关键特性或规则

备份加密、访问隔离、异地和不可变保留；记录版本与校验；定期恢复到隔离环境并做行数/关键业务校验；PITR 同时保存 base backup 与完整 WAL 链。

## 实际怎么使用

```sh
pg_dump --format=custom --file=app.dump appdb
createdb restore_test
pg_restore --dbname=restore_test --clean --if-exists app.dump
pg_basebackup -D /backup/base -Fp -Xs -P
```

## 常见错误与边界

复制不是备份，误删会复制；只测试 pg_dump 成功不足；备份可能缺角色、扩展、配置和外部对象；恢复时间随数据量变化。

`pg_restore --clean` 会删除目标库中待恢复对象，只能对确认无生产流量的隔离恢复库执行。物理备份目录必须为空且权限受控；命令成功后仍需检查日志、对象、关键行数和应用读写验证。

## 补充知识

恢复演练测量真实 RTO，并文档化 DNS、密钥、应用依赖和流量切换步骤。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/backup.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/app-pgdump.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/continuous-archiving.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/app-pgbasebackup.html)（访问日期：2026-07-16）
