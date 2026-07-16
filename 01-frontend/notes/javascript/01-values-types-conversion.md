# 值、变量、类型、运算符、表达式与类型转换

## 是什么

JavaScript 值分为原始值和对象；变量绑定值；运算符组合操作数形成表达式；转换可显式调用 Number/String/Boolean，也会在部分运算中隐式发生。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

let/const 有块级作用域；=== 不做类型强制转换；null 与 undefined 语义不同；对象比较身份。

## 实际使用

```js
const raw='42'; const count=Number(raw);
if (Number.isNaN(count)) throw new TypeError('invalid number');
const total=count*2;
```

## 常见错误与边界

依赖 ==、+ 的隐式转换会产生边界错误；typeof null 为 object 是历史行为；浮点数不适合直接表示货币。

## 相关补充知识

`typeof null`、`NaN`、`-0` 和大整数边界容易产生误判。对外部字符串应显式解析并检查结果；金额和精确十进制不能直接假设二进制浮点运算无误差。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Grammar_and_types)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Expressions_and_operators)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Data_structures)

访问日期：2026-07-16。
