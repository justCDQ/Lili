# Array、Object、String、Map、Set 与 Date

## 是什么

Array 是有序可变集合；Object 是属性表；String 是不可变 UTF-16 序列；Map 支持任意键且保序；Set 保存唯一值；Date 表示自纪元起的毫秒时间点。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

数组方法区分是否变异；对象键除 symbol 外转字符串；Map/Set 用 SameValueZero 比较；Date 解析非标准字符串不可靠。

## 实际使用

```js
const byId=new Map(users.map(u=>[u.id,u]));
const tags=[...new Set(users.flatMap(u=>u.tags))];
const iso=new Date().toISOString();
```

## 常见错误与边界

数组空槽不同于 undefined；字符串索引按 UTF-16 code unit；Date 同时受时区和历法显示影响。

## 相关补充知识

Map 的键可为任意值，Object 的属性键主要是字符串或 Symbol；Set 适合成员关系和去重。Date 表示时间点但不自带业务时区，解析非 ISO 字符串和跨时区显示需明确规则。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Indexed_collections)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Keyed_collections)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date)

访问日期：2026-07-16。
