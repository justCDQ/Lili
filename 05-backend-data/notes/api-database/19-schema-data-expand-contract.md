# Schema 与数据迁移：Expand–Migrate–Contract

数据库迁移是数据、约束、应用版本和运维状态的协调变更。滚动部署期间旧代码和新代码会同时访问数据库，因此不能把“执行一条 ALTER”当成完整迁移。

## 1. 三阶段模型

```mermaid
flowchart LR
    A["Expand：新增兼容结构"] --> B["Migrate：双写/回填/验证"]
    B --> C["Switch：读路径切换"]
    C --> D["Contract：停止旧写后删除旧结构"]
```

- Expand：只增加旧应用能忽略的新列、表、索引或约束的非强制形态。
- Migrate：新代码开始写新旧结构，历史数据分批回填并验证。
- Switch：读取切到新结构，保留快速回退能力。
- Contract：确认旧代码、任务和脚本不再使用后，删除旧字段/约束。

Contract 通常与 Expand 分属不同发布，间隔由部署、数据规模和回退窗口决定。

## 2. 迁移风险不是只看 SQL 长度

评估维度：锁模式与等待、表重写、WAL/复制延迟、I/O、事务时长、索引构建、磁盘临时空间、应用兼容、失败后状态、回滚是否丢新数据。

任何 DDL 都可能在锁队列中造成影响。先设置短 `lock_timeout`，避免迁移无限等待并阻塞后续请求：

```sql
SET lock_timeout = '2s';
SET statement_timeout = '30s';
ALTER TABLE orders ADD COLUMN customer_note text;
```

值只是示例，应按变更和窗口制定。超时后检查数据库状态，而不是无脑循环重试。

## 3. 迁移文件与状态管理

每个迁移有全局唯一、单调顺序标识，内容进入版本控制。迁移工具维护已应用版本与校验和；已在共享环境执行的 migration 不应原地修改，应追加修正迁移。

迁移分事务型和非事务型。PostgreSQL 多数 DDL 可事务回滚，但 `CREATE INDEX CONCURRENTLY` 不能在事务块中运行。工具必须支持显式 no-transaction 文件，失败后检查 invalid index 并清理/重试。

生产应用角色不应自动拥有任意 DDL 权限；迁移使用受控角色、审计和流水线。应用启动时自动迁移在多副本并发、长迁移和失败恢复上风险大。

## 4. 新增必填列

目标：`users.display_name` 最终 NOT NULL。安全路径：

1. Expand：新增可空列。
2. 部署新代码：所有新写入同时填充 display_name。
3. 分批回填历史 NULL。
4. 添加 `NOT VALID` CHECK 约束，避免初始扫描阻塞组合：

```sql
ALTER TABLE users
ADD CONSTRAINT users_display_name_nn
CHECK (display_name IS NOT NULL) NOT VALID;
```

5. `VALIDATE CONSTRAINT`，允许并发读写但仍需评估资源/锁。
6. 设置列 NOT NULL；在已验证等价约束时 PostgreSQL 可避免完整表扫描的部分成本。
7. 删除临时 CHECK（若无需保留）。

添加带常量默认值的列在现代 PostgreSQL 某些条件下无需重写全表，但易变默认表达式、版本与具体 DDL 不同。不能把该优化当所有默认值的通则；用测试库和官方版本行为验证。

## 5. 重命名列不是兼容变更

直接 `RENAME old TO new` 会让旧应用立即失败。兼容做法：

1. 新增 `new_name`。
2. 应用双写 old/new；读取优先 new，缺失回退 old。
3. 回填并比较两列。
4. 所有读切 new，监控 old 读取/写入归零。
5. 停止 old 写入。
6. 独立 contract 发布删除 old。

双写放在应用会有遗漏后台脚本的风险；数据库 trigger 能覆盖所有写入口，但增加隐藏行为和迁移复杂度。选择后列出全部写入者并测试。双写顺序必须在同一事务或由数据库生成列等机制保证。

## 6. 改变类型与单位

`amount` 从浮点元改为整数分不是简单 cast：要定义舍入、负数、最大范围和历史异常。新增 `amount_cents bigint`，新写入从同一业务输入计算两列；回填使用明确舍入规则；比较 `amount_cents` 与旧值换算；切读后再删旧列。

`ALTER COLUMN TYPE ... USING` 可能重写大表并持强锁。离线小表可直接做，大表采用新列/新表分阶段迁移。

枚举演进也需兼容。数据库 enum 添加值与应用序列化/客户端穷举有发布顺序；删除/重命名更复杂。状态字段可用 text + CHECK，但约束变更同样需要阶段设计。

## 7. 添加外键和检查约束

大表可先 `NOT VALID` 添加外键/CHECK：新写入仍被约束，历史行暂不扫描；随后 `VALIDATE CONSTRAINT`。`NOT VALID` 不是“约束完全不生效”。

```sql
ALTER TABLE order_items
ADD CONSTRAINT order_items_order_fk
FOREIGN KEY (tenant_id, order_id)
REFERENCES orders (tenant_id, order_id)
NOT VALID;

ALTER TABLE order_items
VALIDATE CONSTRAINT order_items_order_fk;
```

外键引用列必须有唯一/主键支持。引用方列不会自动总是获得适合删除/更新检查的索引，需按访问和外键操作设计。多租户外键携带 tenant，防跨租户引用。

## 8. 索引迁移

普通 `CREATE INDEX` 会阻塞写；`CREATE INDEX CONCURRENTLY` 允许并发写但执行多个阶段、耗时更长、消耗资源，并且不能在事务块中。唯一索引并发创建后，可将其附加为约束（满足条件时），减少锁窗口。

删除索引也要先确认使用和依赖。`DROP INDEX CONCURRENTLY` 有限制。索引上线后等待统计更新并观察真实计划，不能仅因存在就假定使用。

## 9. 回填设计

回填必须有稳定游标、批量上限、速率控制、可恢复 checkpoint、幂等写和指标。避免大 OFFSET；按主键 keyset：

```sql
UPDATE users
SET display_name = derive_display_name(first_name, last_name)
WHERE user_id > $1
  AND user_id <= $2
  AND display_name IS NULL;
```

每批短事务提交，记录 last ID、rows、duration、WAL、replica lag 和错误。`WHERE target IS NULL` 让重跑幂等；但若转换算法变化，要版本化，不能静默覆盖已确认数据。

回填期间新写入必须由双写覆盖，否则扫描过的范围又产生空值。完成判断是全量验证查询为零且持续观察，而不是 worker 到达最大 ID。

## 10. 读切换与回退

采用 feature flag 分批切读：内部流量→少量租户→全部。比较旧/新读结果，记录差异但控制高基数和敏感数据。回退时新代码可能已写只有新结构能表达的数据；如果旧列仍同步写并能无损表示，才能安全回退。

因此 rollback 分两类：

- 应用回退：在 contract 前切回旧读路径。
- 数据/DDL 逆迁移：可能丢数据或不可快速执行，不能假定每个 up migration 都有安全 down。

生产回滚计划应明确“停止发布、禁用新写、切旧读、修复数据”的实际步骤。

## 11. 完整案例：姓名字段拆分

### 输入

旧列 `users.full_name NOT NULL`；目标新增 `given_name`、`family_name`，但姓名无法对所有文化可靠自动拆分。三版应用会滚动共存，用户可自行修正。

### 步骤

1. Expand 新增两个可空列和 `name_format_version`。
2. 新应用编辑表单同时写 full_name 与结构化字段；旧应用仍只写 full_name。
3. 回填只对有可靠结构来源的用户自动填充；不可靠的保持 NULL，不能按空格武断拆分。
4. 新读路径优先结构化字段，缺失时显示 full_name；记录缺失率。
5. 提供用户确认/修正流程，把确认来源和版本写入。
6. 所有写入口迁移后停止旧应用写；若业务仍需展示完整姓名，可保留 full_name 作为独立展示字段，而非强制删除。

### 输出

迁移不会伪造姓名结构。结构化字段的完整度以“已确认/可靠来源”衡量，而不是强行达到 100%。

### 验证

- 旧版、新版和后台任务在 expand 阶段都能写。
- 双写不一致率按来源统计并可定位。
- 回填重跑不覆盖用户已确认值。
- lock timeout 下迁移失败不会阻塞业务长队列。
- contract 前连续一个发布周期无旧列唯一写入。

### 失败分支

如果直接 rename full_name，旧版立刻报列不存在。若把姓名按第一个空格拆分，会制造错误数据且难以回滚。修正为兼容新增、明确未知值、用户确认和分阶段切换。

## 12. 迁移上线清单

上线前：估算行数/大小/WAL；确认 PostgreSQL 18.4 行为；在生产规模副本演练；列出锁模式；设置 timeout；确认磁盘和复制余量；备份/恢复能力可用；定义观测和停止阈值。

上线中：监控锁等待、DB CPU/I/O、WAL、replica lag、长事务、错误率、池等待和回填速度；每批可暂停；DDL 失败检查残留状态。

上线后：运行数据不变量验证；等待旧写归零；记录何时允许 contract；不要立即删除旧数据以保留回退窗口。

## 13. 常见错误与练习

错误包括：同一发布新增 NOT NULL 又只让新代码填；直接重命名；单事务回填全表；在应用启动自动并发迁移；修改已执行 migration；忽略 `CREATE INDEX CONCURRENTLY` 失败后的 invalid 索引；把 down migration 当无损保证。

练习：把订单 `status` 拆成状态与原因，并为 5000 万行设计 expand、双写、回填、验证、切读和 contract。

完成标准：旧/新应用可共存；DDL 有锁与超时说明；回填 keyset、短事务、幂等、可暂停；验证不变量和差异；应用回退不丢新数据；contract 有流量证据和独立发布。

## 来源

- [PostgreSQL 18: ALTER TABLE](https://www.postgresql.org/docs/18/sql-altertable.html)（访问日期：2026-07-17）
- [PostgreSQL 18: CREATE INDEX](https://www.postgresql.org/docs/18/sql-createindex.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Constraints](https://www.postgresql.org/docs/18/ddl-constraints.html)（访问日期：2026-07-17）
- [PostgreSQL 18: Explicit Locking](https://www.postgresql.org/docs/18/explicit-locking.html)（访问日期：2026-07-17）
- [PostgreSQL 18.4 Release Notes](https://www.postgresql.org/docs/release/18.4/)（访问日期：2026-07-17）
