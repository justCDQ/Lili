# 响应式更新与渲染模型

## 是什么

框架把状态变更映射为 UI 更新。React 重新调用组件并协调树；Vue 跟踪响应式依赖；Svelte 编译更新逻辑。理解模型能判断何时读取状态、何时产生副作用。

## 为什么需要

框架的响应式与渲染模型决定状态读取如何建立依赖、更新何时批处理、组件何时重新计算以及 DOM 如何提交。理解模型才能定位陈旧闭包、过度渲染和更新时序问题。

## 关键特性与规则

渲染阶段必须纯；状态更新通常批处理；DOM 提交后再执行需要真实节点的副作用；派生数据优先计算或 memo。

## 实际使用

```tsx
// React: 更新函数避免陈旧快照
setCount(c=>c+1);
// Vue: ref 变更触发依赖更新
count.value++;
```

## 常见错误与边界

把框架 state 当同步可变变量会读到旧快照；在 render 中更新 state 可无限循环；不同框架细节不能互套。

## 相关补充知识

React 的状态快照、Vue 的依赖追踪和 Svelte 的编译转换实现不同，不能照搬优化规则。性能判断应以 Profiler 和用户任务为证据，先修复错误依赖和不稳定标识，再考虑缓存。

## 来源

- [React Documentation](https://react.dev/learn/render-and-commit)
- [React Documentation](https://react.dev/learn/state-as-a-snapshot)
- [Vue Documentation](https://vuejs.org/guide/extras/reactivity-in-depth.html)

访问日期：2026-07-16。
