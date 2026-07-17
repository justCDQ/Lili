# Join、Group By、Having、Subquery、CTE 与窗口函数

复杂 SQL 的核心不是关键字数量，而是每一步的行粒度。`JOIN` 组合关系，`GROUP BY` 把多行折叠成组，`HAVING` 过滤组，子查询和 CTE 表达中间关系，窗口函数则在保留明细行的同时计算跨行结果。

查询示例以 PostgreSQL 18.4 为主；基础 JOIN、聚合和窗口案例也使用 SQLite 3.51 执行兼容子集验证。

## 查询的逻辑处理路径

理解下列逻辑顺序有助于解释别名可见性、聚合时机和结果行数；它不是执行计划的物理执行顺序：

```mermaid
flowchart LR
    F["FROM / JOIN：形成输入行"] --> W["WHERE：过滤输入行"]
    W --> G["GROUP BY：划分组"]
    G --> H["HAVING：过滤组"]
    H --> S["SELECT：计算输出"]
    S --> D["DISTINCT：去重"]
    D --> O["ORDER BY：排序"]
    O --> L["LIMIT / OFFSET：截取"]
```

规划器可以在保证语义相同的前提下改写和重排物理操作，因此性能判断必须看 `EXPLAIN`，不能把逻辑顺序当执行顺序。

## JOIN：先声明输出行粒度

假设有客户、订单、明细三层：

```sql
CREATE TABLE customers (
  id integer PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE orders (
  id integer PRIMARY KEY,
  customer_id integer NOT NULL REFERENCES customers(id),
  status text NOT NULL,
  created_at timestamp NOT NULL
);

CREATE TABLE order_items (
  order_id integer NOT NULL REFERENCES orders(id),
  product_name text NOT NULL,
  quantity integer NOT NULL CHECK (quantity > 0),
  unit_price numeric(12, 2) NOT NULL,
  PRIMARY KEY (order_id, product_name)
);
```

`customers JOIN orders` 的自然粒度是一行一订单；再连接 `order_items` 后是一行一订单明细。若此时直接 `count(orders.id)`，一张含三条明细的订单会被计数三次。

### JOIN 类型

| 类型 | 返回规则 | 常见用途 | 关键风险 |
|---|---|---|---|
| `INNER JOIN` | 只保留两边匹配行 | 查询确实存在关联的数据 | 无匹配父行会消失 |
| `LEFT JOIN` | 保留左侧全部行，右侧无匹配时补 `NULL` | 包含零订单客户 | `WHERE` 右表条件可能排除补出的空值 |
| `RIGHT JOIN` | 保留右侧全部行 | 可改写为调换左右的左连接 | 阅读方向不统一 |
| `FULL JOIN` | 保留两侧未匹配行 | 对账两个来源 | 需处理两边键为空 |
| `CROSS JOIN` | 返回笛卡尔积 | 生成所有组合、日期骨架 | 行数为两边乘积 |

### `ON` 与 `WHERE` 的差别

查询所有客户及其已支付订单：

```sql
SELECT c.id, c.name, o.id AS order_id
FROM customers AS c
LEFT JOIN orders AS o
  ON o.customer_id = c.id
 AND o.status = 'paid';
```

`status` 放在 `ON` 中表示只有已支付订单可以匹配，但没有已支付订单的客户仍保留。改成 `WHERE o.status = 'paid'` 会排除右侧为 `NULL` 的行，结果等同于这一目的下的内连接。

### 半连接与反连接意图

只关心“是否存在”时用 `EXISTS`，避免连接后重复父行：

```sql
SELECT c.id, c.name
FROM customers AS c
WHERE EXISTS (
  SELECT 1
  FROM orders AS o
  WHERE o.customer_id = c.id
    AND o.status = 'paid'
);
```

查询没有已支付订单的客户用 `NOT EXISTS`。`NOT IN (subquery)` 若子查询含 `NULL`，比较结果可能变成未知，从而一行也不返回；`NOT EXISTS` 通常更直接表达反连接意图。

## GROUP BY：把明细折叠成组

每个分组键组合产生一行。选择列表中未聚合的列通常必须属于分组键，或能被 PostgreSQL 证明函数依赖于分组键。

```sql
SELECT
  o.customer_id,
  count(*) AS paid_order_count,
  sum(oi.quantity * oi.unit_price) AS revenue
FROM orders AS o
JOIN order_items AS oi ON oi.order_id = o.id
WHERE o.status = 'paid'
GROUP BY o.customer_id;
```

这里 `count(*)` 数的是连接后的明细行，不是订单。正确统计订单数可先在订单级汇总金额，再按客户汇总：

```sql
WITH order_totals AS (
  SELECT
    o.id,
    o.customer_id,
    sum(oi.quantity * oi.unit_price) AS total
  FROM orders AS o
  JOIN order_items AS oi ON oi.order_id = o.id
  WHERE o.status = 'paid'
  GROUP BY o.id, o.customer_id
)
SELECT
  customer_id,
  count(*) AS paid_order_count,
  sum(total) AS revenue
FROM order_totals
GROUP BY customer_id;
```

### 聚合对 `NULL` 的处理

- `count(*)` 统计输入行，包括列值为空的行。
- `count(column)` 只统计该列非空的行。
- `sum`、`avg`、`min`、`max` 忽略空值；没有输入行时 `sum` 返回 `NULL`，不是 `0`。
- 需要零值时用 `coalesce(sum(amount), 0)`，并确认零与未知在业务上确实等价。

条件聚合可以用标准 `CASE`，PostgreSQL 也支持更清晰的 `FILTER`：

```sql
SELECT
  customer_id,
  count(*) FILTER (WHERE status = 'paid') AS paid_count,
  count(*) FILTER (WHERE status = 'canceled') AS canceled_count
FROM orders
GROUP BY customer_id;
```

## WHERE 与 HAVING

`WHERE` 在分组前过滤行；`HAVING` 在聚合后过滤组：

```sql
SELECT customer_id, count(*) AS paid_count
FROM orders
WHERE status = 'paid'
GROUP BY customer_id
HAVING count(*) >= 2;
```

能在分组前确定的条件应写 `WHERE`，既清楚又可减少聚合输入。`HAVING customer_id = 7` 可能被规划器下推，但语义上仍不如 `WHERE customer_id = 7` 明确。

## 子查询的四种形态

### 标量子查询

标量子查询必须返回一列且最多一行；零行产生 `NULL`，多于一行报错：

```sql
SELECT
  o.id,
  (SELECT c.name FROM customers AS c WHERE c.id = o.customer_id) AS customer_name
FROM orders AS o;
```

能用普通连接清晰表达时通常优先连接。标量相关子查询不必然逐行执行，规划器可能改写；是否高效以执行计划为准。

### `EXISTS` 子查询

`EXISTS` 只关心是否至少一行，选择列表内容无关，惯例写 `SELECT 1`。数据库找到满足条件的行即可停止该存在性判断。

### 集合比较子查询

`IN`、`ANY`、`ALL` 把左值与子查询集合比较。必须明确空集合和 `NULL` 的三值逻辑。例如 `x = ANY(empty_array)` 为假，而 `x <> ALL(empty_array)` 为真。

### `FROM` 子查询与 `LATERAL`

普通 `FROM` 子查询不能引用同级前面的表；加 `LATERAL` 后可为每个左侧行产生相关关系。查询每位客户最近一张订单：

```sql
SELECT c.id, c.name, latest.id AS latest_order_id
FROM customers AS c
LEFT JOIN LATERAL (
  SELECT o.id
  FROM orders AS o
  WHERE o.customer_id = c.id
  ORDER BY o.created_at DESC, o.id DESC
  LIMIT 1
) AS latest ON true;
```

合适的 `(customer_id, created_at DESC, id DESC)` 索引能支持每个客户的有界查找。大量客户下仍要用 `EXPLAIN (ANALYZE, BUFFERS)` 验证循环次数与 I/O。

## CTE：给中间关系命名

CTE 使用 `WITH` 定义只在当前语句可见的辅助语句：

```sql
WITH paid_orders AS (
  SELECT id, customer_id
  FROM orders
  WHERE status = 'paid'
),
order_totals AS (
  SELECT p.customer_id, p.id, sum(i.quantity * i.unit_price) AS total
  FROM paid_orders AS p
  JOIN order_items AS i ON i.order_id = p.id
  GROUP BY p.customer_id, p.id
)
SELECT customer_id, count(*) AS order_count, sum(total) AS revenue
FROM order_totals
GROUP BY customer_id;
```

CTE 不是永久表，也不是天然性能优化。PostgreSQL 18 对无副作用、非递归 CTE：

- 父查询只引用一次时，默认可折叠进父查询共同优化。
- 被多次引用时，默认通常物化并只计算一次。
- `MATERIALIZED` 强制形成优化边界，可能避免昂贵函数重复计算，也可能阻止条件下推。
- `NOT MATERIALIZED` 允许合并，可能减少无用行，也可能重复计算。

选择必须依据实际计划和函数副作用，不能沿用“CTE 一定物化”的旧规则。

递归 CTE 由非递归项、`UNION [ALL]` 和递归项组成，适合层级遍历；必须有终止条件和环检测策略。组织结构、依赖图若可能成环，应使用路径数组或 PostgreSQL 的 `CYCLE` 子句标记环，避免无界递归。

## 窗口函数：计算跨行信息但保留明细

窗口函数必须带 `OVER`。它不会像 `GROUP BY` 那样把一组折叠为一行：

```sql
SELECT
  o.id,
  o.customer_id,
  o.created_at,
  row_number() OVER (
    PARTITION BY o.customer_id
    ORDER BY o.created_at DESC, o.id DESC
  ) AS recency_rank
FROM orders AS o;
```

窗口定义包含：

| 部分 | 作用 | 缺失时 |
|---|---|---|
| `PARTITION BY` | 把输入分成互不影响的分区 | 全部输入属于一个分区 |
| 窗口 `ORDER BY` | 定义分区内顺序和 peer 组 | 排名和前后值可能无稳定含义 |
| frame | 限定当前行参与计算的窗口范围 | 使用函数对应的默认 frame |

常用函数：

- `row_number()`：每行唯一序号；要稳定必须提供唯一排序终结键。
- `rank()`：并列值同名次，后续名次跳号。
- `dense_rank()`：并列值同名次，后续不跳号。
- `lag()` / `lead()`：访问前一行或后一行，用于变化量和时间间隔。
- `first_value()` / `last_value()`：在当前 frame 内取首尾值；默认 frame 下 `last_value()` 常只到当前 peer 组末尾，若要整个分区末值需显式 frame。
- 聚合函数加 `OVER`：计算运行总计、移动平均或分区总计。

运行累计收入：

```sql
SELECT
  id,
  created_at,
  total,
  sum(total) OVER (
    ORDER BY created_at, id
    ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
  ) AS running_total
FROM order_totals;
```

使用 `ROWS` 明确按物理排序行累计；默认 `RANGE` 会把排序键相同的 peer 行一起纳入，可能让累计值一次跳过多行。

窗口函数的输出不能直接在同层 `WHERE` 中过滤，因为 `WHERE` 逻辑上更早。用子查询或 CTE：

```sql
WITH ranked AS (
  SELECT
    o.*,
    row_number() OVER (
      PARTITION BY customer_id
      ORDER BY created_at DESC, id DESC
    ) AS rn
  FROM orders AS o
)
SELECT id, customer_id, created_at
FROM ranked
WHERE rn <= 3;
```

## 完整案例：客户收入排行榜

### 输入数据

```sql
INSERT INTO customers (id, name) VALUES
  (1, 'Alpha'), (2, 'Beta'), (3, 'Gamma');

INSERT INTO orders (id, customer_id, status, created_at) VALUES
  (101, 1, 'paid',     '2026-07-01 09:00:00'),
  (102, 1, 'paid',     '2026-07-03 09:00:00'),
  (103, 2, 'paid',     '2026-07-02 09:00:00'),
  (104, 2, 'canceled', '2026-07-04 09:00:00');

INSERT INTO order_items (order_id, product_name, quantity, unit_price) VALUES
  (101, 'Keyboard', 2, 100.00),
  (101, 'Mouse',    1,  50.00),
  (102, 'Monitor',  1, 300.00),
  (103, 'Keyboard', 1, 100.00),
  (104, 'Mouse',    9,  50.00);
```

目标：输出所有客户，包括没有已支付订单的客户；展示已支付订单数、收入、收入排名和最近支付订单时间；只保留收入至少 `100` 的客户，但仍显示无收入客户用于运营跟进。

### 步骤一：先汇总到一行一订单

```sql
WITH order_totals AS (
  SELECT
    o.id AS order_id,
    o.customer_id,
    o.created_at,
    sum(i.quantity * i.unit_price) AS total
  FROM orders AS o
  JOIN order_items AS i ON i.order_id = o.id
  WHERE o.status = 'paid'
  GROUP BY o.id, o.customer_id, o.created_at
)
SELECT * FROM order_totals ORDER BY order_id;
```

可验证的中间结果为订单 `101 = 250`、`102 = 300`、`103 = 100`；已取消订单 `104` 不进入结果。

### 步骤二：汇总到一行一客户

```sql
WITH order_totals AS (
  SELECT o.id, o.customer_id, o.created_at,
         sum(i.quantity * i.unit_price) AS total
  FROM orders AS o
  JOIN order_items AS i ON i.order_id = o.id
  WHERE o.status = 'paid'
  GROUP BY o.id, o.customer_id, o.created_at
),
customer_totals AS (
  SELECT
    c.id,
    c.name,
    count(ot.id) AS paid_order_count,
    coalesce(sum(ot.total), 0) AS revenue,
    max(ot.created_at) AS latest_paid_at
  FROM customers AS c
  LEFT JOIN order_totals AS ot ON ot.customer_id = c.id
  GROUP BY c.id, c.name
)
SELECT
  id,
  name,
  paid_order_count,
  revenue,
  latest_paid_at,
  dense_rank() OVER (ORDER BY revenue DESC) AS revenue_rank
FROM customer_totals
WHERE revenue >= 100 OR paid_order_count = 0
ORDER BY revenue DESC, id;
```

### 输出与验证

结果顺序和数值：

| id | name | paid_order_count | revenue | revenue_rank |
|---:|---|---:|---:|---:|
| 1 | Alpha | 2 | 550.00 | 1 |
| 2 | Beta | 1 | 100.00 | 2 |
| 3 | Gamma | 0 | 0 | 3 |

验证步骤：

1. 手工按明细计算每个订单金额。
2. 确认订单计数为 `2` 而不是 Alpha 的明细数 `3`。
3. 确认取消订单金额 `450` 未进入 Beta 收入。
4. 确认 `LEFT JOIN` 保留 Gamma。
5. 用 `EXPLAIN (ANALYZE, BUFFERS)` 在测试库检查估算行数与实际行数、排序和聚合节点；该命令会真正执行查询。

### 失败分支：直接连接后聚合

错误查询直接从客户连接订单和明细，并使用 `count(o.id)`。Alpha 的订单 `101` 有两条明细，因此被重复计算，订单数会得到 `3`。修正方法是在 `order_totals` 先恢复“一行一订单”的粒度，再向客户级聚合，或只在确实合理时使用 `count(DISTINCT o.id)`；后者不会自动修复其他被重复相乘的金额。

## 参数、安全与调试

- 值必须使用驱动参数，不拼接用户输入。
- 表名、列名、排序方向属于 SQL 结构，普通参数不能代替；使用服务端白名单映射。
- 每加入一个 JOIN，先运行 `count(*)` 和按主键重复计数，观察基数是否符合预期。
- 对中间 CTE 单独查询少量样本，确认行粒度。
- 排名、分页和 `LIMIT` 的 `ORDER BY` 应包含唯一终结键。
- 对计划比较估算 `rows` 与 `EXPLAIN ANALYZE` 的实际 `rows`；巨大偏差通常指向统计信息、相关列或表达式问题。

## 练习：商品品类月度榜单

基于订单、明细、商品和品类，输出每个品类收入前三的商品，同时保留没有销量的品类，并展示商品较上月收入变化。

完成标准：

- 逐层写出明细、商品月汇总、品类排名三个行粒度。
- 使用 `LEFT JOIN` 正确保留零销量品类。
- 使用 `lag()` 计算环比，使用 `dense_rank()` 或 `row_number()` 并解释并列规则。
- 提供输入数据、确定输出、验证计算和一个多对多 JOIN 导致重复聚合的失败分支。
- 用索引和执行计划验证常用时间范围，而不是凭 SQL 外观判断性能。

## 来源

- [PostgreSQL 18：Table Expressions](https://www.postgresql.org/docs/18/queries-table-expressions.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Subquery Expressions](https://www.postgresql.org/docs/18/functions-subquery.html)（访问日期：2026-07-17）
- [PostgreSQL 18：WITH Queries](https://www.postgresql.org/docs/18/queries-with.html)（访问日期：2026-07-17）
- [PostgreSQL 18：Window Functions](https://www.postgresql.org/docs/18/functions-window.html)（访问日期：2026-07-17）
- [PostgreSQL 18：SELECT](https://www.postgresql.org/docs/18/sql-select.html)（访问日期：2026-07-17）
