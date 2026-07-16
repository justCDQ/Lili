# 基础类型、函数、对象、接口、类型别名与枚举

## 是什么

TypeScript 在 JavaScript 上做静态检查。基础类型描述值，函数类型描述参数/返回，对象类型描述属性；interface 可声明合并，type 可表达任意类型别名；enum 生成运行时代码。

## 为什么需要

基础类型把函数、对象和模块之间的数据契约交给编译器检查，能在运行前发现属性缺失、参数错误和不安全返回值。它也是后续联合类型、泛型和类型收窄的前提。

## 关键特性与规则

优先推断局部值，公共边界显式标注；interface 常用于可扩展对象契约，type 适合 union/映射；多数场景可用 as const 对象替代 enum。

## 实际使用

```ts
type ID=string; interface User{id:ID;name:string}
function label(user:User):string{return user.name;}
const role={Admin:'admin',User:'user'} as const;
```

## 常见错误与边界

类型会在运行时擦除；any 关闭检查；对象结构类型允许额外能力但对象字面量有 excess property check。

## 相关补充知识

接口和类型别名都只存在于编译阶段，不能验证网络或存储中的未知数据。边界输入先用运行时 Schema 校验，再把结果交给静态类型；对外库还要考虑声明文件兼容性。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/objects.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/enums.html)

访问日期：2026-07-16。
