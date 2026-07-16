# RBAC、ABAC、资源权限、多租户与审计

## 是什么

RBAC 按角色授予权限；ABAC 根据主体、资源、动作、环境属性决策；资源权限校验具体对象所有权/关系。多租户隔离租户数据与配额；审计记录谁在何时对什么执行什么结果。

## 为什么需要

认证只证明主体，所有业务读写仍必须服务端授权并隔离租户。

## 关键特性或规则

默认拒绝；每个 endpoint 和对象做校验；租户上下文来自可信凭据而非请求 body；审计日志追加写、时间统一、包含决策和结果并脱敏。

## 实际怎么使用

```go
func CanUpdate(actor Actor,doc Document) bool { return actor.TenantID==doc.TenantID && (actor.Role=="admin"||actor.ID==doc.OwnerID) }
// 查询必须同时带 tenant_id，且数据库策略作纵深防御
```

## 常见错误与边界

只隐藏按钮不是授权；先按全局 ID 查询再检查可能形成时序/日志泄漏；共享缓存缺 tenant key 会串数据；审计日志不能保存密码/token。

## 补充知识

PostgreSQL Row-Level Security 可做纵深隔离，但策略和管理员绕过行为必须测试。

## 来源

- [一手资料 1](https://csrc.nist.gov/projects/role-based-access-control)（访问日期：2026-07-16）
- [一手资料 2](https://csrc.nist.gov/projects/attribute-based-access-control)（访问日期：2026-07-16）
- [一手资料 3](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)（访问日期：2026-07-16）
- [一手资料 4](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)（访问日期：2026-07-16）
