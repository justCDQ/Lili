# Runtime Schema Validation：类型不能替代运行时校验

## 是什么

TypeScript 类型编译后消失，网络、存储、表单和 JSON 输入仍是 unknown。运行时 schema 验证真实值并返回已验证数据或错误。

## 为什么需要

TypeScript 不检查运行时输入。API、环境变量、存储、消息和用户输入在进入可信代码前需要解析与校验，否则静态类型只是未经证明的断言。

## 关键特性与规则

外部边界从 unknown 开始；验证结构、范围、格式和业务约束；错误需可定位字段；schema 与静态类型保持单一来源或自动推导。

## 实际使用

```tsx
type User={id:string};
function isUser(v:unknown):v is User{return typeof v==='object'&&v!==null&&typeof (v as {id?:unknown}).id==='string';}
const raw:unknown=await response.json(); if(!isUser(raw)) throw new TypeError('invalid user');
```

## 常见错误与边界

as User 不做验证；JSON.parse 泛型不验证；客户端验证不能代替服务端校验。

## 相关补充知识

验证器应返回结构化错误路径并区分缺失、类型、范围和业务规则。Schema 与 TypeScript 类型应从同一来源派生或做一致性测试，避免二者长期漂移。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/narrowing.html#using-type-predicates)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/functions.html#unknown)
- [json-schema.org](https://json-schema.org/learn/getting-started-step-by-step)

访问日期：2026-07-16。
