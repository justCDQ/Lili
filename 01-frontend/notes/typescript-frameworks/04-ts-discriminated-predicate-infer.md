# Discriminated Union、Type Predicate 与 infer

## 是什么

可辨识 union 共享字面量字段以支持穷尽分支；type predicate 声明自定义守卫结果；infer 在条件类型模式中提取部分类型。

## 为什么需要

可辨识联合把状态与其合法数据绑定，类型谓词封装运行时判断，`infer` 从已有类型结构提取部分。三者用于建立可穷尽状态机和可复用类型工具，避免无效状态组合。

## 关键特性与规则

判别字段稳定且互斥；predicate 实现必须与声明一致；infer 只能在条件类型 extends 子句使用。

## 实际使用

```ts
type State={kind:'loading'}|{kind:'done';data:string}|{kind:'error';message:string};
function render(s:State){switch(s.kind){case 'loading':return '...';case 'done':return s.data;case 'error':return s.message;default:{const neverS:never=s;return neverS;}}}
```

## 常见错误与边界

可选判别字段削弱 narrowing；错误 predicate 会制造不安全；用断言代替穷尽检查会隐藏新增状态。

## 相关补充知识

自定义类型谓词的返回声明不会被编译器证明，谓词实现错误会造成不安全收窄。跨信任边界优先使用能够返回错误路径的 Schema Validator，并用 `never` 检查联合分支是否穷尽。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/narrowing.html#discriminated-unions)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html#inferring-within-conditional-types)

访问日期：2026-07-16。
