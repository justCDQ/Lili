# 值、变量、类型、条件、循环、函数与错误

## 是什么

值是程序处理的数据；变量是绑定值的名称；类型规定值的表示、可用操作和约束。条件按布尔结果选择分支，循环重复执行，函数封装输入到输出的行为，错误表示操作未按约定完成。

JavaScript 是动态类型语言，变量声明不固定运行时值类型；优先 `const`，需要重新赋值才用 `let`，避免 `var` 的函数作用域和提升行为。条件会执行真值转换，业务判断应显式比较。函数应返回结果或抛出异常，不要同时依赖隐式全局状态。

```js
function average(values) {
  if (!Array.isArray(values) || values.length === 0) {
    throw new TypeError("values must be a non-empty array");
  }
  let total = 0;
  for (const value of values) {
    if (!Number.isFinite(value)) throw new TypeError("invalid number");
    total += value;
  }
  return total / values.length;
}
```

## 常见错误与边界

- 区分 `null`（明确无值）与 `undefined`（未提供/未初始化）。
- `Number` 使用 IEEE 754 双精度，金额不可直接依赖浮点精确相等。
- 只捕获能处理的错误；否则补充上下文后继续抛出。
- 循环必须有终止条件，函数参数和返回值要有明确契约。
- `NaN !== NaN`，使用 `Number.isNaN`；严格比较优先 `===`。

## 补充知识

JavaScript 原始类型包括 Number、BigInt、String、Boolean、Symbol、Undefined 和 Null，其余均为对象。函数本身也是可调用对象。生产代码通常配合 TypeScript 或静态检查降低动态类型错误。

## 为什么需要

这些基础决定程序如何表示数据、组织控制流、处理输入输出并报告失败。掌握它们才能明确函数契约、资源边界和可测试行为，而不是只让示例在单一输入下运行。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 实际怎么使用

运行本文代码，并至少加入正常、空值、非法输入、边界规模和外部资源失败五类用例。先写预期输出或错误，再用测试固定；对文件和命令行示例同时检查 stdout、stderr、退出码、权限和大输入。

## 来源

- [MDN：JavaScript language overview](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Language_overview)（访问日期：2026-07-16）
- [MDN：Control flow and error handling](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Control_flow_and_error_handling)（访问日期：2026-07-16）
- [ECMA-262：ECMAScript Language Specification](https://tc39.es/ecma262/)（访问日期：2026-07-16）
