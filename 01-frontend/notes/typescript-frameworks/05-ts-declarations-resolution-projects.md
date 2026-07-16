# Declaration File、Module Resolution、Compiler Options 与 Project References

## 是什么

.d.ts 描述已有 JavaScript 的类型；模块解析模拟运行时/打包器定位说明符；tsconfig 控制检查与输出；project references 用多个复合项目表达依赖。

## 为什么需要

声明文件决定 JavaScript API 如何被类型系统理解，模块解析决定导入实际指向哪个文件，编译选项决定检查语义和输出，Project References 则用于拆分大型仓库的构建边界。

## 关键特性与规则

应用按实际宿主选择 resolution；开启 strict；声明必须真实描述运行时；不同 DOM/Node/Worker 环境分 tsconfig；大型项目用 references 与 tsc -b。

## 实际使用

```ts
// tsconfig.json
{"compilerOptions":{"strict":true,"module":"ESNext","moduleResolution":"Bundler","noEmit":true},"include":["src"]}
```

## 常见错误与边界

paths 不修改运行时导入；错误 resolution 可能类型通过但运行失败；手写宽泛 declare module 会退化为 any。

## 相关补充知识

`module`、`moduleResolution`、运行时和打包器必须配套选择；仅让编辑器解析成功不代表 Node 或浏览器能够加载。发布包应同时测试 ESM/CJS 入口、`exports`、类型声明与源码映射。

## 来源

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/declaration-files/introduction.html)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/modules/reference)
- [TypeScript Handbook](https://www.typescriptlang.org/tsconfig/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/project-references.html)

访问日期：2026-07-16。
