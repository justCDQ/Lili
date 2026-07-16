# Vite 等开发与构建工具

## 是什么

开发服务器提供模块解析、转换和热更新；生产构建做依赖图分析、转换、拆包、压缩与资源输出。Vite 开发期基于原生 ESM，生产构建使用打包器。

## 为什么需要

开发服务器负责模块转换、热更新和源码调试，生产构建负责依赖解析、分块、压缩和静态资源输出。理解两套路径可避免“开发正常、生产失败”和客户端环境变量泄露。

## 关键特性与规则

区分 dev 与 production；配置 base、目标浏览器和入口；构建后预览并检查资源路径、source map、chunk 和缓存。

## 实际使用

```tsx
npm create vite@latest my-app -- --template vanilla-ts
npm run dev
npm run build
npm run preview
```

## 常见错误与边界

dev 正常不代表产物正常；preview 不是生产服务器；客户端环境变量会暴露；动态导入路径需静态可分析。

## 相关补充知识

生产预览不能完全替代真实 CDN、子路径、缓存头和服务器回退测试。插件可执行任意构建代码，应锁定版本、审查来源，并比较构建清单和产物体积变化。

## 来源

- [Vite Documentation](https://vite.dev/guide/)
- [Vite Documentation](https://vite.dev/guide/build.html)
- [Vite Documentation](https://vite.dev/guide/env-and-mode)

访问日期：2026-07-16。
