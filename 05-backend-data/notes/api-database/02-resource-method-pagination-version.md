# HTTP API：资源、方法、状态码、分页与版本

HTTP API 的契约不只是一组 URL。资源标识、方法语义、状态码、条件请求、分页游标和兼容策略共同决定客户端能否正确重试、缓存和演进。

## 1. 资源与 URI

资源是可被标识的业务对象或集合，例如订单、用户的地址集合、导出任务。URI 标识资源，响应正文是资源在某时刻的一种表示。

```text
/users/u_123                 单个用户
/users/u_123/addresses       用户的地址集合
/orders?status=paid          订单集合的筛选表示
/exports/ex_456              异步导出任务
```

路径应表达稳定身份和层级，不把数据库表名、内部主键结构或部署服务名暴露成永久契约。动作无法自然表示为 CRUD 时可建立命令资源，例如 `POST /orders/o_1/cancellations`，它还能拥有创建时间、原因和处理状态。

## 2. HTTP 方法逐项说明

| 方法 | 语义 | 安全 | 幂等 | 常见成功状态 | 正文用途 |
|---|---|---:|---:|---|---|
| `GET` | 获取资源表示 | 是 | 是 | 200、206、304 | 响应表示 |
| `HEAD` | 与 GET 相同但无响应正文 | 是 | 是 | 200、304 | 无正文 |
| `POST` | 让目标资源按自身语义处理正文 | 否 | 否 | 200、201、202、204 | 命令或新成员 |
| `PUT` | 创建或完整替换目标资源状态 | 否 | 是 | 200、201、204 | 完整表示 |
| `PATCH` | 对资源应用部分修改 | 否 | 不保证 | 200、204 | 补丁文档 |
| `DELETE` | 删除目标资源的关联 | 否 | 是 | 202、204 | 通常无正文 |
| `OPTIONS` | 查询通信选项 | 是 | 是 | 200、204 | 可选能力描述 |

“幂等”比较服务端预期效果，而不是每次响应：第一次 `DELETE` 可返回 204，第二次返回 404，删除后的状态仍相同。`POST` 可通过业务幂等键获得重复请求安全，但方法本身仍不是规范定义的幂等方法。

`PUT` 与 `PATCH` 的边界：`PUT` 正文通常代表目标资源的完整期望状态，遗漏字段可能表示清除；`PATCH` 正文是变更指令。JSON Merge Patch 的 `null` 通常表示删除成员，JSON Patch 则显式使用 `add`、`remove`、`replace`、`move`、`copy`、`test` 操作，二者不能混用。

## 3. 状态码按结果分类

### 3.1 成功与异步

- `200 OK`：有成功响应正文，读取或同步操作常用。
- `201 Created`：创建了资源；应通过 `Location` 或正文给出新资源标识。
- `202 Accepted`：请求被接受但尚未完成；应返回任务资源，不能让客户端误认为业务已成功。
- `204 No Content`：成功且无正文；响应不能携带消息正文。
- `206 Partial Content`：响应满足 Range 请求的一部分，需配合 `Content-Range`。

### 3.2 客户端请求问题

- `400 Bad Request`：语法、解析或通用请求错误。
- `401 Unauthorized`：缺少或无效认证凭据，通常配合 `WWW-Authenticate`；它不是“已登录但没权限”。
- `403 Forbidden`：已理解请求但拒绝执行，可能是权限或策略原因。
- `404 Not Found`：目标不存在；为避免资源枚举，也可对无权发现的对象返回 404。
- `405 Method Not Allowed`：资源存在但方法不允许，必须返回 `Allow`。
- `409 Conflict`：与资源当前状态冲突，如版本冲突或状态迁移不允许。
- `412 Precondition Failed`：`If-Match` 等前置条件失败。
- `415 Unsupported Media Type`：请求 `Content-Type` 不支持。
- `422 Unprocessable Content`：语法正确但内容无法按指令处理，如字段约束失败。
- `429 Too Many Requests`：限流拒绝，可返回 `Retry-After`。

### 3.3 服务端与上游问题

- `500 Internal Server Error`：未被更具体语义覆盖的服务端故障。
- `502 Bad Gateway`：作为网关收到无效上游响应。
- `503 Service Unavailable`：暂时不可用或过载，可返回 `Retry-After`。
- `504 Gateway Timeout`：网关等待上游超时。

状态码表达 HTTP 层分类，领域细节放结构化错误正文。不要所有失败都返回 200，也不要把堆栈或 SQL 错误暴露给客户端。

## 4. 条件请求与并发写

服务返回实体标签：

```http
HTTP/1.1 200 OK
ETag: "order-v17"
Content-Type: application/json

{"id":"o_123","status":"draft","version":17}
```

客户端修改时携带 `If-Match`：

```http
PATCH /orders/o_123 HTTP/1.1
If-Match: "order-v17"
Content-Type: application/merge-patch+json

{"status":"submitted"}
```

只有当前 ETag 仍匹配才更新，否则返回 412。服务端必须在同一个受控事务中比较版本并写入，不能只在前端检查。`If-None-Match: *` 可用于“仅当资源不存在时创建”。

## 5. 筛选、排序与字段

集合接口应规定允许字段、操作符、类型、默认值和上限：

```text
GET /orders?status=paid&created_after=2026-07-01T00:00:00Z
    &sort=-created_at,id&limit=50
```

排序必须有唯一且稳定的最终键，例如 `created_at DESC, id DESC`。只按非唯一时间排序会在同一时间戳处丢项或重复。筛选字段必须白名单映射到 SQL 参数，不能把查询字符串直接拼进列名或表达式。

## 6. 偏移分页与游标分页

### 6.1 Offset/limit

`?offset=100&limit=50` 容易跳页和显示页码，适合小型、相对静态数据。数据库往往仍需跳过前 100 行；并发插入或删除还会让下一页重复或遗漏。

### 6.2 Keyset/cursor

游标分页保存上一页最后一条的排序键：

```sql
SELECT id, created_at, total_cents
FROM orders
WHERE tenant_id = $1
  AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT $4;
```

响应使用不透明游标：

```json
{
  "items": [{"id":"o_123","created_at":"2026-07-17T02:00:00Z"}],
  "page": {"next_cursor":"eyJjcmVhdGVkX2F0IjoiLi4uIiwiaWQiOiJvXzEyMyJ9","has_more":true}
}
```

游标应包含排序值、方向和必要的筛选绑定，并做签名或使用服务端状态，避免客户端篡改。它不是数据库 offset 的 Base64 包装；否则仍有同样的扫描和一致性问题。改变筛选或排序后旧游标应被拒绝。

`has_more` 可通过多取一条计算：查询 `limit + 1`，返回前 `limit` 条。如果需要精确总数，应明确其额外代价；大表实时 `COUNT(*)` 可能比取一页更昂贵。

## 7. 版本与兼容演进

版本策略包括路径 `/v1/orders`、媒体类型参数、请求头或域名。路径版本易观察；头版本保持 URI，但调试和缓存键更复杂。没有策略能替代兼容设计。

通常兼容的变化：新增可选响应字段、新增可选请求字段、新增端点。可能破坏客户端的变化：删除或改名、改变类型/单位/时区、收紧枚举而客户端穷举、改变默认排序、把可空改为必填、改变错误语义。

服务端可以先“扩展”再“收缩”：

1. 新增字段并同时支持旧字段。
2. 发布迁移说明和弃用时间。
3. 观测旧字段真实流量。
4. 客户端迁移完成后停止旧写入。
5. 在新的大版本中删除旧字段。

弃用通知可使用文档、响应头和开发者控制台，但停止日期必须与客户端迁移周期匹配。服务器对未知请求字段的策略需要明确：静默忽略利于兼容，却可能掩盖拼写错误；严格拒绝更易发现错误，却使新增字段经旧网关时失败。

## 8. 完整案例：订单列表和提交

### 输入

- 租户 `t_7` 有百万级订单。
- 列表按 `created_at DESC, id DESC`，每页最多 100 条。
- 草稿提交可能被两台设备同时操作。
- 客户端需要区分验证失败、版本冲突和暂时过载。

### 步骤

1. `GET /v1/orders?status=draft&limit=50` 返回 items、next_cursor，并为单个订单返回 ETag。
2. 后端把租户、状态和游标值作为参数执行 keyset SQL；查询 51 条判断 `has_more`。
3. 客户端提交 `POST /v1/orders/o_9/submissions`，正文含提交参数，`If-Match` 带当前 ETag。
4. 事务中验证租户权限、订单状态和版本，再创建 submission 资源并更新订单。
5. 成功返回 201 与 `Location: /v1/orders/o_9/submissions/s_2`。

### 输出

```http
HTTP/1.1 412 Precondition Failed
Content-Type: application/problem+json

{"type":"https://api.example.com/problems/stale-version","title":"Resource version is stale","status":412,"current_etag":"\"order-v18\""}
```

### 验证

- 连续取十页，ID 无重复；在第一页后插入新订单，后续旧快照范围不漏项。
- 使用相同 ETag 并发提交两次，只有一次状态迁移成功。
- `limit=1000` 被明确拒绝或截断，并在契约中一致。
- 修改 cursor 中的筛选值后签名失败，返回 400。

### 失败分支

如果只按 `created_at` 构造游标，时间相同的订单会跨页丢失。修正为加入唯一 `id`，让 WHERE、ORDER BY 和索引列顺序一致。如果更新前只在应用内先查版本、后写入，两个请求仍可能同时通过；必须用条件 UPDATE 或事务锁保证原子性。

## 9. 调试与常见错误

- 名词路径却用 `POST /getOrders`：方法和路径重复表达，缓存与工具语义变弱。
- 创建返回 200 且没有资源标识：客户端无法稳定定位新对象。
- 认证失败返回 403：客户端不知道是否需要刷新凭据。
- `PATCH` 把缺失字段当清空：部分更新误删数据。
- 游标只做 Base64、不签名、不绑定筛选：可被篡改并跨查询误用。
- 每页同步计算精确总数：在大表上产生不必要延迟。
- 前端隐藏按钮代替授权：攻击者仍可直接调用接口。

验证时记录方法、规范化路由、状态码、耗时和响应大小；不要把原始 token、Cookie 或敏感查询值写入日志。对每个写接口测试重复、并发、超时后重试和非法租户。

## 10. 练习

设计项目任务 API：创建任务、按负责人和状态筛选、游标分页、条件更新、归档。输出接口表、请求/响应样例、错误表和兼容变更方案。

完成标准：路径表达资源；方法和状态码符合语义；排序有唯一终结键；并发更新只成功一次；游标不能跨筛选复用；所有字段校验和权限在服务端执行。

## 来源

- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-17）
- [RFC 9111: HTTP Caching](https://www.rfc-editor.org/rfc/rfc9111.html)（访问日期：2026-07-17）
- [RFC 5789: PATCH Method for HTTP](https://www.rfc-editor.org/rfc/rfc5789.html)（访问日期：2026-07-17）
- [RFC 7396: JSON Merge Patch](https://www.rfc-editor.org/rfc/rfc7396.html)（访问日期：2026-07-17）
- [RFC 6902: JSON Patch](https://www.rfc-editor.org/rfc/rfc6902.html)（访问日期：2026-07-17）
