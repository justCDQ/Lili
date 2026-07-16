# Union、Intersection、Narrowing 与 Generic

## 是什么

union 表示多选一，intersection 合并约束；narrowing 根据控制流缩小类型；generic 用类型参数保留输入输出关系。

## 为什么需要

联合类型表达有限候选，交叉类型组合约束，收窄让分支内类型更精确，泛型保留输入与输出之间的关系。正确使用可减少断言和重复重载，同时保持 API 可推导。

## 关键特性与规则

对 union 先判别再使用成员；intersection 不是对象运行时合并；generic 约束写 T extends ...；不要用泛型掩盖固定类型。

## 实际使用

```ts
function first<T>(xs:readonly T[]):T|undefined{return xs[0];}
function print(x:string|number){if(typeof x==='string') return x.toUpperCase(); return x.toFixed(2);}
```

## 常见错误与边界

错误断言可绕过 narrowing；过宽 union 使调用者难处理；泛型只出现一次通常没有建立关系的价值。

## 相关补充知识

泛型参数应表示真实关系，而不是把所有类型变成可配置项。无法从参数推导、只出现一次或始终需要调用方手写的泛型，通常应改为具体类型或联合类型。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/generics.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/unions-and-intersections.html)

访问日期：2026-07-16。
