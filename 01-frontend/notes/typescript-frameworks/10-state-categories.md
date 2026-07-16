# 客户端、服务端、URL、持久化状态

## 是什么

客户端状态属于当前交互；服务端状态是远程数据缓存；URL 状态应可分享导航；持久化状态跨会话保存。分类决定所有权、失效和同步方式。

## 为什么需要

客户端临时状态、服务端缓存、URL 状态和持久化状态具有不同生命周期、所有者和一致性要求。分类错误会造成重复数据源、刷新行为不一致和难以解释的同步冲突。

## 关键特性与规则

远程数据保存获取时间、错误和失效策略；筛选分页放 URL；表单瞬时输入就近；持久化数据版本化并验证。

## 实际使用

```tsx
const url=new URL(location.href); url.searchParams.set('page','2'); history.pushState(null,'',url);
```

## 常见错误与边界

把服务端数据复制进全局 state 会产生双源；localStorage 不适合秘密；URL 不放敏感信息。

## 相关补充知识

可分享和可返回的筛选条件优先放 URL；服务器数据由查询缓存管理；输入草稿保留在局部或显式持久化。不要把可派生值复制进状态，也不要把敏感令牌放入 URL 或普通 Web Storage。

## 来源

- [React Documentation](https://react.dev/learn/choosing-the-state-structure)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/History_API)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)

访问日期：2026-07-16。
