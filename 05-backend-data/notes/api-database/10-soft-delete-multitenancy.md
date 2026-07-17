# 软删除与多租户数据库设计

软删除保留原行并记录删除状态；多租户设计让一套产品安全承载多个客户的数据。两者都会改变查询默认值、唯一性、外键、索引、授权、备份和清除策略，不能只增加 `deleted_at` 或 `tenant_id` 一列就结束。

部分唯一索引、行级安全策略和 `FOR UPDATE SKIP LOCKED` 示例使用 PostgreSQL 18.4 语法。

## 软删除究竟保存什么

硬删除执行 `DELETE`，行在事务语义上消失，底层空间随后由 vacuum 回收。软删除通常执行 `UPDATE`：

```sql
UPDATE app_users
SET deleted_at = now(), deleted_by = $2
WHERE tenant_id = $1 AND id = $3 AND deleted_at IS NULL;
```

软删除适合需要短期撤销、业务归档或保留引用上下文的对象。它不等于：

- 法律意义上的个人数据删除；敏感字段、搜索索引、缓存、导出和备份都有独立处理要求。
- 审计日志；一行的 `deleted_at` 不能说明此前所有变更。
- 防误操作的唯一手段；重要删除仍需授权、确认、幂等和审计。
- 免费存储；行和索引仍占空间，每次软删是一次更新并产生新行版本和 WAL。

### 建议字段

| 字段 | 作用 | 边界 |
|---|---|---|
| `deleted_at timestamptz` | 删除发生的绝对时刻；`NULL` 表示活动 | 不记录删除原因和主体 |
| `deleted_by bigint` | 发起删除的可信主体 | 主体本身可能后续停用 |
| `delete_reason text` | 受控原因或工单号 | 避免存入不必要敏感信息 |
| `purge_after timestamptz` | 允许硬清除的最早时刻 | 法务冻结时必须覆盖清除任务 |

布尔 `is_deleted` 只能表示两个状态，不能回答何时删除或计算保留期。若业务有归档、冻结、封禁、取消等不同状态，应分别建模，不能都挤进“已删除”。

## 软删除对唯一性的影响

普通约束：

```text
UNIQUE (tenant_id, email)
```

会让已删除用户继续占用邮箱。若产品允许删除后重新注册，应对活动行建立唯一部分索引：

```sql
CREATE UNIQUE INDEX app_users_active_email_uq
ON app_users (tenant_id, lower(email))
WHERE deleted_at IS NULL;
```

部分索引只包含满足谓词的行，因此软删后旧索引项退出，新活动用户可以使用同一邮箱。需要注意：

- `lower(email)` 定义的是大小写不敏感业务规则；Unicode、排序规则和邮箱规范必须与产品要求一致。
- 查询若要利用该索引，其条件必须能让规划器在计划阶段证明 `deleted_at IS NULL`。任意参数化谓词不一定能推出部分索引谓词。
- 部分唯一索引不是普通唯一约束，外键不能把它当作所有行都存在的通用候选键。
- 恢复旧用户可能撞上新用户的活动邮箱；恢复命令必须有明确冲突处理。

恢复前锁定目标并检测冲突：

```sql
BEGIN;

SELECT id, email
FROM app_users
WHERE tenant_id = $1 AND id = $2
FOR UPDATE;

SELECT id
FROM app_users
WHERE tenant_id = $1
  AND lower(email) = lower($3)
  AND deleted_at IS NULL
FOR KEY SHARE;

UPDATE app_users
SET deleted_at = NULL, deleted_by = NULL, delete_reason = NULL
WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NOT NULL;

COMMIT;
```

即使先查询无冲突，也必须由唯一索引裁决并发恢复或注册；应用捕获 `unique_violation` 返回冲突结果。

## 查询默认值与关联语义

活动查询必须显式表达：

```sql
SELECT id, display_name
FROM app_users
WHERE tenant_id = $1 AND deleted_at IS NULL;
```

把过滤隐藏在 ORM 默认 scope 中容易产生以下错误：

- 管理员恢复页面需要已删除行，却被默认 scope 隐藏。
- 原生 SQL、后台任务或新代码路径忘记 scope。
- 关联查询只过滤父表，没有过滤已删除子表。
- `LEFT JOIN` 的过滤条件放错位置，改变结果行数。

若要保留所有项目，同时只连接活动成员，应把子表过滤写在 `ON`：

```sql
SELECT p.id, count(m.id) AS active_members
FROM projects AS p
LEFT JOIN project_members AS m
  ON m.project_id = p.id
 AND m.deleted_at IS NULL
WHERE p.tenant_id = $1
  AND p.deleted_at IS NULL
GROUP BY p.id;
```

把 `m.deleted_at IS NULL` 放在 `WHERE` 时，SQL 的三值逻辑仍可能保留无匹配行，因为连接产生的 `NULL IS NULL` 为真；但更复杂的右表条件常把左连接意外变成内连接。把“哪些子行有资格匹配”放在 `ON` 最清晰。

## 多租户隔离方式

| 方式 | 隔离单元 | 优点 | 成本与风险 |
|---|---|---|---|
| 共享表 | 每行 `tenant_id` | 资源利用率高，迁移统一 | 任一漏过滤可能跨租户，噪声邻居明显 |
| 每租户 schema | schema | 名称空间分离，可独立导出部分对象 | 大量 schema 的迁移、连接池和目录管理复杂 |
| 每租户数据库 | database/cluster | 权限、备份和资源隔离更强 | 连接、升级、监控、容量和跨租户统计成本高 |
| 混合 | 按规模或合规分层 | 大租户可独立隔离 | 路由、迁移和运维模型必须同时支持多种形态 |

选择依据包括合规、爆炸半径、单租户规模、定制需求、恢复粒度和运维能力。共享表并非天然不安全，但要求每层都把租户边界作为不变量。

## 共享表中的复合租户键

只给表增加 `tenant_id` 仍可能建立跨租户引用：

```text
-- 不充分：只能证明 project_id 存在
project_id bigint REFERENCES projects(id)
```

应该让父表暴露租户范围内的候选键，再使用复合外键：

```sql
CREATE TABLE projects (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  tenant_id bigint NOT NULL REFERENCES tenants(id),
  name text NOT NULL,
  deleted_at timestamptz,
  UNIQUE (tenant_id, id)
);

CREATE TABLE tasks (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  tenant_id bigint NOT NULL,
  project_id bigint NOT NULL,
  title text NOT NULL,
  deleted_at timestamptz,
  UNIQUE (tenant_id, id),
  FOREIGN KEY (tenant_id, project_id)
    REFERENCES projects (tenant_id, id)
    ON DELETE RESTRICT
);
```

这样任何写入路径都不能把租户 1 的任务指向租户 2 的项目。复合外键仍不等于访问授权；它只保证数据所属关系一致。

索引通常以租户列开头以匹配常见查询：

```sql
CREATE INDEX tasks_active_project_idx
ON tasks (tenant_id, project_id, id)
WHERE deleted_at IS NULL;
```

是否把 `tenant_id` 放首位取决于查询矩阵、数据分布和索引策略，但共享表的大多数在线查询都应限定租户。

## PostgreSQL 行级安全策略

RLS（Row-Level Security）在表层限制普通查询和修改命令能看到或写入哪些行。它是数据库纵深防御，不替代应用授权、表权限和可信身份建立。

### 角色和会话上下文

不要让应用使用表所有者、超级用户或具有 `BYPASSRLS` 的角色，因为这些角色可绕过策略。典型分工：

```sql
CREATE ROLE app_owner NOLOGIN;
CREATE ROLE app_runtime LOGIN NOSUPERUSER NOBYPASSRLS;

ALTER TABLE tasks OWNER TO app_owner;
GRANT SELECT, INSERT, UPDATE, DELETE ON tasks TO app_runtime;
```

连接取得可信租户身份后，在事务内设置本地参数：

```sql
BEGIN;
SELECT set_config('app.tenant_id', '42', true);
-- true 表示 SET LOCAL 语义，只持续到当前事务结束
SELECT id, title FROM tasks WHERE deleted_at IS NULL;
COMMIT;
```

租户值必须由已经验证的会话、令牌或服务身份产生，不能直接信任请求体中的 `tenantId`。连接池下使用事务局部设置，避免租户上下文泄漏到下一个请求。

### 策略函数和默认拒绝

先定义缺少上下文时返回 `NULL` 的稳定读取函数：

```sql
CREATE FUNCTION app_current_tenant_id()
RETURNS bigint
LANGUAGE sql
STABLE
AS $$
  SELECT NULLIF(current_setting('app.tenant_id', true), '')::bigint
$$;
```

`current_setting(..., true)` 在参数不存在时返回 `NULL`，而不是抛错。策略中的比较遇到 `NULL` 不为真，因此默认看不到任何行。

```sql
ALTER TABLE tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE tasks FORCE ROW LEVEL SECURITY;

CREATE POLICY tasks_tenant_select
ON tasks
FOR SELECT
TO app_runtime
USING (tenant_id = app_current_tenant_id());

CREATE POLICY tasks_tenant_insert
ON tasks
FOR INSERT
TO app_runtime
WITH CHECK (tenant_id = app_current_tenant_id());

CREATE POLICY tasks_tenant_update
ON tasks
FOR UPDATE
TO app_runtime
USING (tenant_id = app_current_tenant_id())
WITH CHECK (tenant_id = app_current_tenant_id());

CREATE POLICY tasks_tenant_delete
ON tasks
FOR DELETE
TO app_runtime
USING (tenant_id = app_current_tenant_id());
```

- `USING` 限制命令可读取、更新或删除的现有行。
- `WITH CHECK` 限制插入或更新后的新行，防止把行改到别的租户。
- 启用 RLS 后若没有适用策略，普通角色默认拒绝访问。
- 表所有者通常绕过 RLS；`FORCE ROW LEVEL SECURITY` 让所有者在普通访问中也受策略约束，但超级用户和 `BYPASSRLS` 仍可绕过。
- 多个 permissive 策略默认用 `OR` 合并；restrictive 策略用 `AND` 与其他适用策略组合。增加策略可能意外扩大访问范围，必须测试组合结果。

RLS 策略表达式作为查询的一部分执行。策略使用的函数要控制 `search_path`、权限和副作用；复杂连接会影响性能，还可能引入安全边界。敏感表应使用独立运行角色执行集成测试，而不是以 owner 测试后误判策略有效。

## 完整案例：删除与恢复租户项目

### 输入

租户 `42` 的管理员删除项目 `9001`。规则为：

1. 只能操作自己租户的项目。
2. 删除项目时，其活动任务同时进入已删除状态。
3. 项目 slug 在活动项目中租户内唯一。
4. 30 天内可恢复；若 slug 已被新项目占用，恢复失败并要求改名。
5. 超过保留期且没有法务冻结才允许硬清除。

### 建表和索引

```sql
ALTER TABLE projects
  ADD COLUMN slug text,
  ADD COLUMN deleted_by bigint,
  ADD COLUMN purge_after timestamptz,
  ADD COLUMN legal_hold boolean NOT NULL DEFAULT false;

UPDATE projects
SET slug = 'project-' || id::text
WHERE slug IS NULL;

ALTER TABLE projects ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX projects_active_slug_uq
ON projects (tenant_id, lower(slug))
WHERE deleted_at IS NULL;

GRANT SELECT, INSERT, UPDATE, DELETE ON projects TO app_runtime;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects FORCE ROW LEVEL SECURITY;

CREATE POLICY projects_tenant_all
ON projects
FOR ALL
TO app_runtime
USING (tenant_id = app_current_tenant_id())
WITH CHECK (tenant_id = app_current_tenant_id());
```

上述更新展示迁移顺序：先允许空值、回填、再收紧为 `NOT NULL`。大表应采用可恢复的小批次回填，并在每步观察锁、WAL、复制延迟和剩余空值，不能直接照搬为一次无界生产更新。

### 删除步骤

```sql
BEGIN;
SELECT set_config('app.tenant_id', '42', true);

WITH deleted_project AS (
  UPDATE projects
  SET deleted_at = now(),
      deleted_by = 77,
      purge_after = now() + interval '30 days'
  WHERE tenant_id = 42
    AND id = 9001
    AND deleted_at IS NULL
  RETURNING id, deleted_at
)
UPDATE tasks AS t
SET deleted_at = p.deleted_at
FROM deleted_project AS p
WHERE t.tenant_id = 42
  AND t.project_id = p.id
  AND t.deleted_at IS NULL;

COMMIT;
```

生产命令还应记录审计事件，并检查第一条更新是否返回项目。零行可能表示不存在、已删除或无权限；对外错误不应泄露其他租户对象是否存在。

### 输出与验证

以运行角色执行：

```sql
BEGIN;
SELECT set_config('app.tenant_id', '42', true);

SELECT id, slug, deleted_at
FROM projects
WHERE id = 9001;

SELECT count(*) AS active_tasks
FROM tasks
WHERE project_id = 9001 AND deleted_at IS NULL;
ROLLBACK;
```

管理恢复路径可显式查询已删除项目；普通活动列表必须加 `deleted_at IS NULL`。验证 `active_tasks = 0`，并检查审计记录、RLS 测试和唯一索引定义。

### 恢复步骤与冲突分支

```sql
BEGIN;
SELECT set_config('app.tenant_id', '42', true);

UPDATE projects
SET deleted_at = NULL,
    deleted_by = NULL,
    purge_after = NULL
WHERE tenant_id = 42
  AND id = 9001
  AND deleted_at IS NOT NULL
  AND purge_after > now()
RETURNING id, slug;

UPDATE tasks
SET deleted_at = NULL
WHERE tenant_id = 42
  AND project_id = 9001
  AND deleted_at IS NOT NULL;

COMMIT;
```

若另一个活动项目已经使用同一 slug，第一条更新触发 `unique_violation`，整个事务回滚，任务不会被部分恢复。产品可以要求管理员先选择新 slug，但不能静默覆盖现有项目。

该示例把所有随项目删除的任务都恢复；真实系统若任务可在项目删除前单独删除，必须记录删除原因或删除批次，只恢复由本次级联软删的任务。

### 跨租户失败分支

把事务上下文设为租户 `41` 后请求项目 `9001`：RLS 的 `USING` 使该行不可见，更新返回零行；复合外键则继续阻止任何把租户 41 的任务关联到租户 42 项目的写入。二者分别保护访问路径与数据关系。

## 硬清除、外键和备份

清除任务应采用小批次、可恢复检查点和独立审计，不能无界执行大删除：

```sql
WITH purge_batch AS (
  SELECT id
  FROM projects
  WHERE deleted_at IS NOT NULL
    AND purge_after <= now()
    AND legal_hold = false
  ORDER BY purge_after, id
  LIMIT 500
  FOR UPDATE SKIP LOCKED
)
DELETE FROM projects AS p
USING purge_batch AS b
WHERE p.id = b.id
RETURNING p.id;
```

实际运行前必须处理子表外键、审计保留、搜索和对象存储数据。数据库硬删不能让历史备份中的数据立即消失；备份加密、访问控制、保留到期和密钥销毁属于完整删除策略。

## 调试与验证矩阵

至少以三种数据库角色测试：应用运行角色、后台维护角色、表所有者。覆盖：

| 场景 | 预期 |
|---|---|
| 无租户上下文查询 | 返回零行或受控失败 |
| 租户 42 查询自己的活动行 | 只返回租户 42 且未删除行 |
| 租户 42 按 ID 查询租户 41 行 | 不可见 |
| 插入其他 `tenant_id` | `WITH CHECK` 拒绝 |
| 更新时改变 `tenant_id` | `WITH CHECK` 拒绝 |
| owner 与 `BYPASSRLS` 角色测试 | 明确认知会绕过策略 |
| 删除后注册相同邮箱 | 唯一部分索引允许 |
| 恢复时邮箱或 slug 被占用 | 唯一冲突，事务回滚 |

使用 `EXPLAIN (ANALYZE, BUFFERS)` 时只能在安全环境运行实际查询；`ANALYZE` 会执行语句。检查 RLS 后计划是否仍能使用 `(tenant_id, ...)` 索引，并监控软删造成的表膨胀、vacuum 压力和索引增长。

## 练习：多租户知识库

设计 `spaces`、`documents` 和 `document_revisions`：文档支持 14 天恢复，租户内活动文档路径唯一，修订历史不能软删除，运行角色必须受 RLS 限制。

完成标准：

- 每张共享表都有租户键和必要的复合外键。
- 用唯一部分索引保证活动路径唯一。
- 分别定义 `SELECT`、`INSERT`、`UPDATE`、`DELETE` 策略，并测试无上下文行为。
- 说明恢复时路径冲突、单独删除子对象、法务冻结三个分支。
- 给出硬清除批次 SQL 及备份保留边界。
- 使用非 owner、非 `BYPASSRLS` 角色执行集成测试。

## 来源

- [PostgreSQL 18：Row Security Policies](https://www.postgresql.org/docs/18/ddl-rowsecurity.html)（访问日期：2026-07-17）
- [PostgreSQL 18：CREATE POLICY](https://www.postgresql.org/docs/18/sql-createpolicy.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Partial Indexes](https://www.postgresql.org/docs/18/indexes-partial.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Constraints](https://www.postgresql.org/docs/18/ddl-constraints.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Routine Vacuuming](https://www.postgresql.org/docs/18/routine-vacuuming.html)（访问日期：2026-07-17）
