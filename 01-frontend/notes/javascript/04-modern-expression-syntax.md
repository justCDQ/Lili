# 解构、展开、模板字符串与可选链

## 是什么

解构从数组/对象读取绑定；展开在支持位置复制可枚举项；模板字符串支持插值和多行；可选链在 null/undefined 前短路。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

展开是浅复制；?? 只对 null/undefined 回退；可选链只保护链中明确写 ?. 的位置。

## 实际使用

```js
const {user:{name}={},items=[]}=payload;
const next={...state, items:[...state.items,newItem]};
const label=`${name ?? '匿名'}：${items.length}`;
const city=payload.user?.address?.city;
```

## 常见错误与边界

用 || 会把 0/空串当缺省；展开类实例会丢原型；嵌套解构缺默认值会抛错。

## 相关补充知识

展开语法只做浅复制，不会复制嵌套对象；可选链只在链左侧为空时短路；`??` 只处理 `null/undefined`，与会把 `0`、空串和 `false` 视为假值的 `||` 不同。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Destructuring)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Spread_syntax)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Optional_chaining)

访问日期：2026-07-16。
