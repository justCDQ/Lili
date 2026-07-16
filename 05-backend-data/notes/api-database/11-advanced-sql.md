# Join、Group By、Having、Subquery、CTE 与 Window Function

## 是什么

JOIN 组合关系；GROUP BY 聚合组；HAVING 过滤聚合后结果；subquery 嵌入查询；CTE 命名中间结果；窗口函数在不折叠行的情况下跨行计算。

## 为什么需要

复杂报表和业务查询需要明确行数、基数与执行阶段，避免应用层重复聚合。

## 关键特性或规则

JOIN 条件写完整；LEFT JOIN 右表过滤放 ON 或明确接受退化；WHERE 在聚合前，HAVING 在聚合后；窗口 ORDER BY 与最终输出排序独立。

## 实际怎么使用

```sql
WITH paid AS (SELECT * FROM orders WHERE status='paid')
SELECT u.id,count(p.id) AS orders,sum(p.total) AS revenue,rank() OVER(ORDER BY sum(p.total) DESC)
FROM users u LEFT JOIN paid p ON p.user_id=u.id
GROUP BY u.id HAVING count(p.id)>0;
```

## 常见错误与边界

多对多 JOIN 会乘行导致重复聚合；NOT IN 遇 NULL 语义易错，常用 NOT EXISTS；CTE 是否物化取决于版本和写法。

## 补充知识

先写预期行粒度，再逐步 JOIN 并核对 count；参数化值，不拼接 SQL。

## 来源

- [PostgreSQL/标准资料 1](https://www.postgresql.org/docs/current/queries-table-expressions.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 2](https://www.postgresql.org/docs/current/queries-with.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 3](https://www.postgresql.org/docs/current/tutorial-window.html)（访问日期：2026-07-16）
- [PostgreSQL/标准资料 4](https://www.postgresql.org/docs/current/functions-subquery.html)（访问日期：2026-07-16）
