# Conditional、Mapped 与 Template Literal Types

## 是什么

conditional type 按可赋值关系选分支；mapped type 遍历属性键生成新对象类型；template literal type 组合字符串字面量 union。

## 为什么需要

条件类型、映射类型和模板字面量类型可从单一源类型派生新契约，减少 DTO、事件名和配置类型重复。它们适合稳定、可解释的变换，不适合隐藏复杂业务规则。

## 关键特性与规则

注意条件类型对裸类型参数 union 的分布；mapped modifier 可增删 readonly/?；字符串组合可能指数膨胀。

## 实际使用

```ts
type Api<T>=T extends Error?{ok:false;error:T}:{ok:true;data:T};
type Readonlyish<T>={[K in keyof T]:Readonly<T[K]>};
type EventName<K extends string>=`${K}Changed`;
```

## 常见错误与边界

复杂类型会降低错误可读性和编译性能；不能替代运行时转换；大规模字符串 union 更适合代码生成。

## 相关补充知识

分配条件类型会对裸类型参数的联合成员分别计算；用元组包裹可关闭分配。递归类型和大联合会提高编译成本，应通过命名中间类型和限制递归深度保持错误信息可读。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/mapped-types.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/template-literal-types.html)

访问日期：2026-07-16。
