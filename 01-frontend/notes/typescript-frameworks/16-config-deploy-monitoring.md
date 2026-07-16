# 环境变量、配置、构建、部署与错误监控

## 是什么

配置区分环境行为；构建生成不可变产物；部署发布产物；错误监控采集异常、版本、环境和上下文，支持回归定位。

## 为什么需要

配置、构建、部署与监控连接源码和生产运行状态。明确构建时变量、运行时配置、发布版本和错误关联，才能复现环境差异并安全回滚。

## 关键特性与规则

秘密只在服务端；产物关联 commit/release 和 source map；部署支持回滚；监控脱敏、采样并设告警。

## 实际使用

```tsx
// Vite 客户端只读取允许前缀
const apiBase=import.meta.env.VITE_API_BASE;
window.addEventListener('error',reportError);
window.addEventListener('unhandledrejection',e=>reportError(e.reason));
```

## 常见错误与边界

前端环境变量构建后公开；上传 source map 需访问控制；只收集错误不设告警等于不可观测。

## 相关补充知识

进入浏览器包的环境变量都可被用户读取，Secret 必须保留在服务端。监控事件应包含 release、source map、路由和关联 ID，同时执行脱敏、采样和保留期策略。

## 来源

- [Vite Documentation](https://vite.dev/guide/env-and-mode)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Window/error_event)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Window/unhandledrejection_event)

访问日期：2026-07-16。
