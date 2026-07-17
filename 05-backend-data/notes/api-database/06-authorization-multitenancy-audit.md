# 授权模型、资源权限、多租户隔离与审计日志

授权是服务端根据主体、操作、资源、环境和策略作出的允许/拒绝决策。UI 隐藏按钮、前端路由守卫、JWT 中的角色名都不能替代 API 和数据库边界内的强制检查。

## 1. 授权决策的输入

```text
decision = policy(subject, action, resource, context)
```

- `subject`：用户、服务账号或客户端，包含可信租户成员关系与认证强度。
- `action`：`invoice.read`、`invoice.approve` 等稳定业务动作，不只是 HTTP 方法。
- `resource`：具体对象及其 owner、tenant、状态、敏感级别。
- `context`：时间、网络、设备、最近认证、请求来源等环境条件。
- `policy`：由受控配置/代码定义并版本化的规则。

认证中间件只建立 principal。handler 或领域服务必须在加载目标资源并知道动作后授权。对“列出集合”还要把权限转成查询条件，不能先取全量再由前端过滤。

## 2. RBAC：基于角色

RBAC 把权限赋给角色，再把角色赋给主体。典型关系：

```text
role: expense_viewer  -> expense.read
role: expense_editor  -> expense.read, expense.update_draft
role: expense_approver -> expense.read, expense.approve
```

RBAC 易解释、易审计，适合岗位稳定的组织权限。角色不应直接写成大量 UI 页面名；权限使用领域动作，UI 再查询能力。角色爆炸通常来自把地区、部门、所有权和资源状态全部编码进角色名，例如 `apac_finance_manager_draft_only`。

角色分配要限定作用域：`tenant_id + role`，有时还包括项目/部门。全局管理员和租户管理员必须是不同权限，避免跨租户提升。

## 3. ABAC：基于属性

ABAC 按属性表达条件：

```text
allow expense.read when
  subject.tenant_id == resource.tenant_id AND
  (subject.department_id == resource.department_id OR subject.has("expense.audit"))
```

ABAC 能表达所有权、部门、数据级别、时间和金额，但策略更难测试与解释。属性必须来自可信来源；客户端提交的 `tenant_id`、`role`、`owner_id` 只是请求数据，不能直接成为授权事实。策略应有明确默认拒绝、类型、缺失值行为和冲突合并规则。

RBAC 与 ABAC 常组合：角色给出粗粒度动作，属性限制具体对象。例如“approver 可批准，但不能批准自己的报销，且金额超过 5 万需二级审批”。

## 4. 资源级与字段级权限

仅检查 `expense.read` 会产生 IDOR/BOLA：攻击者把 `/expenses/e_1` 改为 `/expenses/e_2` 读取别人的对象。正确流程：

1. 从认证 principal 取得可信 tenant/subject。
2. 按 `tenant_id + resource_id` 查询，或在数据库策略内强制租户条件。
3. 对具体动作评估 owner、状态和角色。
4. 对返回字段做最小披露，例如普通成员看不到银行账号。

写入也要字段级授权。客户端提交完整对象时，不能把 `role`、`tenant_id`、`approved_by`、`balance` 等受保护字段直接绑定并持久化。为每个操作定义专用输入 DTO，只接受允许字段，受控字段由服务端生成。

## 5. 拒绝、隐藏存在性与缓存

未认证返回 401；已认证但策略拒绝通常返回 403；需要隐藏对象存在性时可返回 404。团队应为同类资源保持一致，避免通过状态、时间或响应大小枚举。

授权结果可以短时缓存，但缓存键至少包含主体、租户、动作、资源/策略版本；角色撤销和资源所有权变更应使缓存失效。不要把“用户是管理员”的结果复用于另一租户，也不要在 CDN 缓存包含私人数据的响应而缺失正确 cache key/`private` 设置。

## 6. 多租户隔离模型

### 6.1 共享表、共享 schema

每个租户表行含 `tenant_id`，主键/唯一约束/索引通常都要带租户：

```sql
CREATE TABLE expenses (
    tenant_id uuid NOT NULL,
    expense_id uuid NOT NULL,
    owner_id uuid NOT NULL,
    status text NOT NULL,
    amount_cents bigint NOT NULL CHECK (amount_cents >= 0),
    PRIMARY KEY (tenant_id, expense_id)
);
CREATE UNIQUE INDEX expenses_tenant_external_ref_uq
ON expenses (tenant_id, external_ref)
WHERE external_ref IS NOT NULL;
```

优点是运营和资源利用简单；风险是每条查询都必须隔离。单列全局 `expense_id` 即使随机，也不能代替 tenant 条件。外键也应携带 tenant，避免一行引用另一租户对象。

### 6.2 每租户 schema 或数据库

独立 schema 提升逻辑隔离但迁移数量增大；独立数据库提供更强故障/备份/区域边界，但连接池、升级、观测和成本更复杂。可按规模和合规采用混合模型，大租户独库、小租户共享。

选择维度包括隔离强度、恢复粒度、数据驻留、邻居噪声、租户数量、成本和迁移能力，不存在对所有产品通用的最佳模型。

## 7. PostgreSQL Row-Level Security（RLS）

RLS 能在数据库层按行强制策略，是纵深防御，不替代应用动作授权。启用后若无适用策略通常默认拒绝，但表所有者和具有 `BYPASSRLS` 的角色通常绕过；生产应用角色不能拥有这些能力。

```sql
ALTER TABLE expenses ENABLE ROW LEVEL SECURITY;
ALTER TABLE expenses FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON expenses
USING (tenant_id = current_setting('app.tenant_id')::uuid)
WITH CHECK (tenant_id = current_setting('app.tenant_id')::uuid);
```

每个事务用 `SET LOCAL app.tenant_id = ...` 绑定当前租户，并保证连接归还池前状态不会泄漏。`USING` 控制可见/可更新旧行，`WITH CHECK` 控制 INSERT/UPDATE 后的新行。策略函数的权限、稳定性和性能都需审计。

RLS 常见失败是用表 owner 连接导致策略未实际生效，或只写 `USING` 忽略写入后 tenant 被改变。测试要使用与生产相同权限角色。

## 8. 策略执行点与默认拒绝

可以把策略写在应用代码、策略引擎或数据库 RLS。无论位置，都要保证：

- 所有入口经过同一强制点，包括后台任务、批量脚本和 GraphQL resolver；
- 策略不可用时对敏感操作默认拒绝，而不是默认允许；
- 决策带策略版本和理由码，便于审计；
- 策略变更先做测试与影子评估；
- 超级管理员使用单独流程、强认证、短时提权和完整审计。

授权检查和写入之间可能发生 TOCTOU。把依赖的状态条件放入同一事务或条件 UPDATE，例如 `WHERE tenant_id=$1 AND id=$2 AND status='submitted'`，根据影响行数确认。

## 9. 审计日志与普通日志的区别

审计日志记录“谁在何时以什么权限对哪个对象做了什么、结果如何”，用于追责和合规。普通应用日志用于调试，可采样和轮转；审计记录通常要求更严格保留、访问控制、完整性和导出流程。

一条审计事件建议包含：

```json
{
  "event_id":"ae_01JZ",
  "occurred_at":"2026-07-17T03:14:15.123Z",
  "actor":{"type":"user","id":"u_9","tenant_id":"t_7"},
  "action":"expense.approve",
  "resource":{"type":"expense","id":"e_42","tenant_id":"t_7"},
  "decision":"allow",
  "reason_code":"role_and_amount_policy",
  "policy_version":"expense-authz-2026-07-01",
  "request_id":"rq_8",
  "changes":{"status":{"from":"submitted","to":"approved"}}
}
```

不记录密码、token、Cookie、完整银行卡或无关正文。敏感变更可记录字段名、掩码值或加密后的受限快照。时间统一 UTC 且高精度，事件 ID 全局唯一。

审计写入要与业务结果一致。关键操作可在同一数据库事务写审计/outbox，再异步投递不可变存储；直接向外部日志服务同步写入会把其故障变成业务故障，异步又必须防丢。删除和更正审计数据应通过追加补正事件，访问审计日志本身也要审计。

## 10. 完整案例：多租户报销审批

### 输入

- 员工只能看本租户自己的报销；财务可看本部门；审计员可看本租户全部。
- 审批人不能批准自己的报销；金额超过 50,000 元需 `senior_approver`。
- `tenant_id` 会出现在 URL/正文，但不可信。

### 步骤

1. 从验证后的 session 得到 `subject=u_9` 和当前租户成员关系 `t_7`。
2. 查询用 `tenant_id=t_7 AND expense_id=e_42`，数据库 RLS 再执行同一租户隔离。
3. 策略先检查 `expense.approve` 角色权限，再检查 `owner_id != u_9`、状态为 submitted、金额阈值。
4. 条件 UPDATE 在同一事务中改变状态，并插入审计/outbox 事件。
5. 返回批准后的表示；前端按钮只是展示策略能力，不是强制点。

### 输出

允许事件记录主体、动作、资源、策略版本和状态变化。拒绝时返回 403 或为不可发现对象返回 404，审计记录 `deny` 与低敏感 reason code。

### 验证

- 把 URL 中 tenant 改为 `t_8`，结果不可发现且数据库无跨租户行。
- 普通 approver 尝试批准自己的报销被拒绝。
- 金额边界 50,000 与 50,000.01 按以分存储的整数规则正确判断。
- 使用生产应用角色测试 RLS，不用表 owner。
- 每个成功状态迁移恰有一条可关联审计事件。

### 失败分支

若服务先按 `expense_id` 查询，再在内存检查 tenant，查询日志、缓存或错误时间可能已暴露跨租户对象。修正为租户条件进入最初查询、主/外键、唯一约束和索引，并用 RLS 纵深保护。若只依赖 JWT 中旧角色，撤权后 token 过期前仍可审批；高风险动作需查询当前授权状态或使用短期提权。

## 11. 常见错误

- 前端隐藏按钮作为权限控制。
- 只校验角色，不校验对象 owner/tenant/状态。
- 接受正文的 `tenant_id` 或 `approved_by` 作为可信值。
- 共享表的唯一约束遗漏 tenant，导致一个租户占用另一租户名称。
- RLS 使用表 owner 连接，测试时误以为策略有效。
- 审计日志含 token 或完整敏感正文。
- 策略服务超时后默认允许。
- 管理员永久全局权限，无提权期限和审计。

## 12. 练习

为文档系统定义 viewer/editor/owner/auditor 四类角色，并加入文档所有权、敏感级别和租户属性规则。实现集合查询、单对象读取和转移所有权的授权测试。

完成标准：默认拒绝；所有查询带 tenant；跨租户主外键不可建立；字段级敏感信息按权限裁剪；策略不可用时敏感写入拒绝；成功和拒绝均有不含 secret 的审计事件；用非 owner 数据库角色证明 RLS 生效。

## 来源

- [NIST SP 800-162: Attribute Based Access Control](https://csrc.nist.gov/pubs/sp/800/162/upd2/final)（访问日期：2026-07-17）
- [NIST RBAC Model](https://csrc.nist.gov/projects/role-based-access-control)（访问日期：2026-07-17）
- [PostgreSQL 18: Row Security Policies](https://www.postgresql.org/docs/18/ddl-rowsecurity.html)（访问日期：2026-07-17）
- [OWASP API Security Top 10](https://owasp.org/API-Security/)（访问日期：2026-07-17）
- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)（访问日期：2026-07-17）
