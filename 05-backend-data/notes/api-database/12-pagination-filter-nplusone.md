# 分页、动态筛选与 N+1 查询

分页限制一次返回的数据量；筛选把请求条件转换为受控查询；N+1 指先执行一次父查询，再为 N 个结果逐个查询关联数据，使数据库往返次数随结果数增长。三者必须一起设计，因为排序、筛选和加载关联共同决定结果正确性与性能。

复合行比较、数组参数、显式类型转换、JSON 聚合和覆盖索引示例采用 PostgreSQL 18.4 语法。

## 分页必须先有全序

没有 `ORDER BY` 时，SQL 结果没有承诺顺序。只按非唯一列排序也不构成稳定全序：多行拥有同一 `created_at` 时，数据库可以用任意相对顺序返回。

```text
ORDER BY created_at DESC, id DESC
```

唯一 `id` 是 tie-breaker，使任意两行都能比较先后。排序方向、`NULLS FIRST/LAST`、表达式和排序规则都属于分页协议，必须在下一页保持一致。

## OFFSET/LIMIT 分页

```sql
SELECT id, created_at, total
FROM orders
WHERE tenant_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3;
```

### 优点

- 接口容易表达页码，可直接跳到第 N 页。
- 在小结果集、后台管理和变化不频繁的数据上足够实用。
- 可配合单独 `count(*)` 展示总页数。

### 成本与漂移

数据库仍需定位并丢弃 offset 之前的行。即使索引支持顺序，深页也要经过大量索引项；若还需排序，成本更高。

并发写入会让页边界漂移。第一次请求后若头部插入一行，第二页的 `OFFSET 20` 可能再次包含第一页末行；若头部删除一行，原第二页首行可能被跳过。单条语句内部有一致快照，但多次 HTTP 请求通常不共享同一数据库快照。

`LIMIT/OFFSET` 每页都必须使用同一确定性排序。不同 `LIMIT` 值可能产生不同计划和无序输出，不能依赖当前观察到的物理行顺序。

## Keyset/Cursor 分页

Keyset 分页不计算“跳过多少行”，而是表达“从上一页最后一个排序键之后继续”。对于降序 `(created_at, id)`：

```sql
SELECT id, created_at, total
FROM orders
WHERE tenant_id = $1
  AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT $4;
```

PostgreSQL 行构造器按从左到右的词典序比较：先比较 `created_at`；相等时比较 `id`。游标必须包含两个值，漏掉 `id` 会在相同时间戳处漏行或重复。

升序则使用 `>`；上一页需要反转比较和数据库查询排序，再在服务端恢复显示方向。不要只改一个方向。

### 游标内容

游标通常编码为不透明字符串，内部至少包含：

```json
{
  "v": 1,
  "createdAt": "2026-07-17T08:30:00.123456Z",
  "id": 9021,
  "sort": "created_at_desc",
  "filterHash": "sha256:..."
}
```

- `v` 支持协议演进。
- 时间保留数据库所需精度，不能经毫秒格式化后丢失微秒。
- `sort` 防止把一种排序的游标用于另一种排序。
- `filterHash` 绑定筛选条件，防止中途更换范围造成含义不一致。
- 编码不是防篡改；需要信任游标内容时用服务器密钥签名或使用服务端游标存储，并限制长度和版本。

### Keyset 的边界

- 不能高效任意跳到第 1000 页。
- 若排序键会修改，一行可能移动到游标前后，仍可能漏读或重复。适合使用不可变 `created_at` 和唯一 ID。
- 新插入到第一页之前的数据不会出现在沿旧游标向后的遍历中，这通常符合时间线语义。
- 要获得固定快照，需要增加不可变上界，例如首次请求记录 `snapshot_max_id` 或业务时间截止点，并让后续页都带上界。
- 可空排序键需要显式 `NULLS` 规则和游标编码；初学实现宜优先选择非空键。

## 索引要匹配筛选与排序

常用查询：

```text
WHERE tenant_id = $1
  AND status = $2
  AND (created_at, id) < ($3, $4)
ORDER BY created_at DESC, id DESC
LIMIT $5
```

对应索引候选：

```sql
CREATE INDEX orders_tenant_status_created_idx
ON orders (tenant_id, status, created_at DESC, id DESC)
INCLUDE (total);
```

前两列为等值范围，后两列支持游标边界和排序。`INCLUDE` 是否值得取决于返回列、表更新频率和 visibility map；它不保证一定出现 index-only scan。

若 `status` 是可选筛选，单个复合索引未必同时最佳支持“有 status”和“无 status”。应从真实查询矩阵、频率和 `EXPLAIN` 决定是增加 `(tenant_id, created_at, id)`、依赖 skip scan/bitmap 组合，还是限制接口组合。

## 动态筛选的安全构建

值使用参数，查询结构使用白名单。请求：

```json
{
  "statuses": ["paid", "refunded"],
  "createdFrom": "2026-07-01T00:00:00Z",
  "createdTo": "2026-08-01T00:00:00Z",
  "minTotal": "100.00",
  "sort": "newest",
  "pageSize": 50
}
```

服务端步骤：

1. 验证枚举、时间、金额精度和页大小上限。
2. 从可信身份得到 `tenant_id`，不取请求体租户值。
3. 只为存在的筛选添加谓词和参数。
4. 把 `sort` 映射到固定 SQL 片段。
5. 解码并验证游标版本、签名和筛选绑定。

生成的 SQL 可以是：

```sql
SELECT id, status, created_at, total
FROM orders
WHERE tenant_id = $1
  AND status = ANY($2::text[])
  AND created_at >= $3::timestamptz
  AND created_at < $4::timestamptz
  AND total >= $5::numeric
  AND (created_at, id) < ($6::timestamptz, $7::bigint)
ORDER BY created_at DESC, id DESC
LIMIT $8;
```

时间范围常用半开区间 `[from, to)`，相邻区间不重叠且无需制造“当天最后一微秒”。金额不要用浮点解析后再拼字符串。

### 可选条件的两种写法

固定 SQL 有时写成：

```text
WHERE ($2::text IS NULL OR status = $2)
```

它减少 SQL 形态，但 `OR`、通用计划和参数分布可能影响索引选择。动态加入 `status = $2` 会产生多个受控 SQL 形态，却更容易让规划器看到直接谓词。没有一种写法总是更快，使用代表性参数比较计划。

### 不可直接参数化的结构

下面是错误意图：

```text
ORDER BY $1 $2
```

普通参数代表值，不能代表列标识符或 `ASC/DESC` 关键字。应做白名单映射：

```text
newest      -> ORDER BY created_at DESC, id DESC
oldest      -> ORDER BY created_at ASC, id ASC
total_high  -> ORDER BY total DESC, id DESC
```

每种排序都需要对应游标字段和索引评估。不要把客户端字符串直接插入 SQL。

## N+1 查询怎样发生

父查询返回 50 张订单：

```sql
SELECT id, customer_id, created_at
FROM orders
WHERE tenant_id = $1
ORDER BY created_at DESC, id DESC
LIMIT 50;
```

随后代码循环 50 次：

```sql
SELECT product_name, quantity, unit_price
FROM order_items
WHERE order_id = $1;
```

总查询数为 `1 + N = 51`。即使每条查询只耗 2ms，串行网络往返、连接占用、解析/执行开销和不同快照都被放大。开发环境数据少、数据库在本机时往往不明显。

## 解决 N+1 的三种方式

### JOIN 一次读取

```sql
SELECT
  o.id,
  o.created_at,
  i.product_name,
  i.quantity,
  i.unit_price
FROM orders AS o
LEFT JOIN order_items AS i ON i.order_id = o.id
WHERE o.tenant_id = $1
  AND o.id = ANY($2::bigint[])
ORDER BY o.created_at DESC, o.id DESC, i.product_name;
```

优点是一次往返；缺点是父列在每条明细重复，多个一对多关联同时连接会产生笛卡尔乘积。必须在分页父 ID 之后再连接，否则 `LIMIT 50` 限制的是明细行而不是订单。

### 两次批量查询

先取一页订单，再用 `ANY` 批量取所有明细：

```sql
SELECT order_id, product_name, quantity, unit_price
FROM order_items
WHERE order_id = ANY($1::bigint[])
ORDER BY order_id, product_name;
```

应用以 `order_id` 分组，固定为两次查询。它避免父行重复，适合 ORM DataLoader 和多个一对多集合。需要处理空 ID 列表，并限制批量大小。

### 数据库聚合子对象

PostgreSQL 可用 `jsonb_agg` 把子对象聚合为每个父行一个 JSON 数组：

```sql
SELECT
  o.id,
  o.created_at,
  coalesce(
    jsonb_agg(
      jsonb_build_object(
        'productName', i.product_name,
        'quantity', i.quantity,
        'unitPrice', i.unit_price
      ) ORDER BY i.product_name
    ) FILTER (WHERE i.order_id IS NOT NULL),
    '[]'::jsonb
  ) AS items
FROM orders AS o
LEFT JOIN order_items AS i ON i.order_id = o.id
WHERE o.tenant_id = $1
  AND o.id = ANY($2::bigint[])
GROUP BY o.id, o.created_at
ORDER BY o.created_at DESC, o.id DESC;
```

这减少应用组装，但把 JSON 构建 CPU 和响应体内存放到数据库。大型子集合仍需独立分页，不能无限聚合。

## 完整案例：订单检索接口

### 输入与数据

请求租户 42 的已支付订单，每页 2 条，按新到旧：

```text
orders:
id=105, created_at=10:05, status=paid
id=104, created_at=10:04, status=paid
id=103, created_at=10:03, status=draft
id=102, created_at=10:02, status=paid
id=101, created_at=10:01, status=paid

order_items:
order_id=105, product_name=Keyboard, quantity=1
order_id=104, no rows
```

### 第一步：请求第一页时多取一条

```sql
SELECT id, created_at, total
FROM orders
WHERE tenant_id = 42
  AND status = 'paid'
ORDER BY created_at DESC, id DESC
LIMIT 3;
```

返回 `105, 104, 102`。接口只输出前两条，把第三条作为 `hasNextPage = true` 的证据；下一游标编码输出页最后一条 `104, 10:04`，不是多取的 `102`。

### 第二步：批量加载两张订单的明细

```sql
SELECT order_id, product_name, quantity, unit_price
FROM order_items
WHERE order_id = ANY(ARRAY[105, 104]::bigint[])
ORDER BY order_id, product_name;
```

应用按订单 ID 组装，整个页面固定两次数据库查询。

### 第三步：请求下一页

```sql
SELECT id, created_at, total
FROM orders
WHERE tenant_id = 42
  AND status = 'paid'
  AND (created_at, id) < ('2026-07-17 10:04:00+00', 104)
ORDER BY created_at DESC, id DESC
LIMIT 3;
```

返回 `102, 101`，数量不超过页大小，因此 `hasNextPage = false`。

### 输出

```json
{
  "items": [
    {"id": 105, "items": [{"productName": "Keyboard", "quantity": 1}]},
    {"id": 104, "items": []}
  ],
  "pageInfo": {
    "hasNextPage": true,
    "endCursor": "eyJ2IjoxLCJjcmVhdGVkQXQiOiIyMDI2LTA3LTE3VDEwOjA0OjAwWiIsImlkIjoxMDQsInNvcnQiOiJjcmVhdGVkX2F0X2Rlc2MiLCJmaWx0ZXJIYXNoIjoic2hhMjU2OmQ1MTQxMjEzYjVjOGZkMGZkY2ZmMmQyOTUxZDM2MGVlMTFkNzRlZTgwMzlhYTQ0NDFjYzJjMTkxNzgxMWEyMjMifQ.kLt5W6NJnPzQ6flXmW7AUa_HWVAcITeukSVdbG-BZk8"
  }
}
```

游标是无空格 JSON 的 base64url 编码与 HMAC-SHA-256 签名，示例密钥为 `demo-only-key`；生产密钥必须由密钥管理系统生成和轮换，不能写入代码仓库。游标解码后最后位置是 `createdAt=2026-07-17T10:04:00Z, id=104`，筛选摘要绑定租户 42 与 `paid` 状态。实际金额和其他字段由查询结果填充；示例只展示分页和关联结构。

### 验证

1. 连续遍历所有页，收集 ID，确认无重复且等于同一筛选下的基准 ID 集合。
2. 插入一条比 `105` 更新的数据后继续旧游标，确认旧遍历不会重复现有页。
3. 为一页 1、2、50 条分别断言查询数固定为 2，而不是 `N+1`。
4. 解码游标时改变筛选或排序，确认服务拒绝。
5. 用 `EXPLAIN (ANALYZE, BUFFERS)` 检查 keyset 条件读取行数不会随页深度线性增加。

### 失败分支

- 只把 `created_at` 放进游标：两条相同时间的订单可能被跳过。修复为复合键。
- 先 JOIN 明细再 `LIMIT 2`：一张有两条明细的订单占满结果，页面不是两张订单。修复为先分页订单 ID。
- 接受 `sort=created_at desc; drop table...` 并拼接：产生 SQL 注入。修复为固定白名单映射。
- ORM 序列化时懒加载 `items`：测试断言查询数立刻发现从 2 增至 51。

## 计数和缓存边界

精确 `count(*)` 可能扫描大量可见行和索引项，不能因为 UI 有“总页数”就无条件每页计算。可选方案：

- 只返回 `hasNextPage`，不返回总数。
- 用户明确请求时计算精确总数。
- 使用有更新时间的预计算计数，并标注它是近似或延迟值。
- 对后台小范围筛选继续使用精确计数。

不能把 PostgreSQL `pg_class.reltuples` 的估算值包装成精确业务总数。

## 调试与测试清单

- `ORDER BY` 是否构成全序。
- 游标是否包含全部排序键、方向、空值规则和筛选绑定。
- 首尾页、空页、相同排序值、删除和插入是否有测试。
- page size 是否有服务端最大值。
- 动态结构是否白名单，值是否参数化。
- 关联加载查询数是否由测试或 tracing 记录。
- 深 offset 与 keyset 是否用真实分布比较计划、延迟和缓冲区。
- 响应体和子集合是否有独立上限。

## 练习：工单列表

实现按 `priority DESC, created_at DESC, id DESC` 排序的租户工单列表，支持状态数组、负责人、半开时间范围，并批量加载标签。

完成标准：

- cursor 包含三个排序键、协议版本和筛选摘要。
- SQL 使用正确的复合行比较；若 priority 可空，定义空值顺序和游标处理。
- 给出支持主查询的候选索引并说明可选筛选的取舍。
- 标签加载总查询数固定，不随工单数增长。
- 用包含相同 priority 和 created_at 的数据验证无重复、无遗漏。
- 注入篡改游标、超大 page size 和非法排序时得到受控错误。

## 来源

- [PostgreSQL 18：LIMIT and OFFSET](https://www.postgresql.org/docs/18/queries-limit.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Row and Array Comparisons](https://www.postgresql.org/docs/18/functions-comparisons.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Indexes and ORDER BY](https://www.postgresql.org/docs/18/indexes-ordering.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Using EXPLAIN](https://www.postgresql.org/docs/18/using-explain.html)（访问日期：2026-07-17）
