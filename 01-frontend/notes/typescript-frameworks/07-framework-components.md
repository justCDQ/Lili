# 组件、Props、State、事件、条件、列表与生命周期

## 是什么

组件封装 UI 与逻辑；props 是父级输入；state 是组件记忆；事件处理用户动作；条件/列表从数据生成 UI；生命周期描述挂载、更新、卸载及外部同步。

## 为什么需要

组件把界面结构、输入、内部状态和事件封装为可组合边界。清晰边界能限制状态传播和重渲染范围，也便于独立测试、复用和无障碍检查。

## 关键特性与规则

render 保持纯；state 最小化且不重复 props 派生值；列表 key 使用稳定业务 ID；外部系统同步时清理订阅。

## 实际使用

```ts
function Counter({step=1}){const [n,setN]=useState(0);return <button onClick={()=>setN(x=>x+step)}>次数 {n}</button>}
```

## 常见错误与边界

直接修改 state 可能不触发更新；数组索引 key 在重排时串状态；Effect 不应用于纯派生计算。

## 相关补充知识

组件不应按视觉矩形机械拆分，应按数据所有权、复用、交互和变化频率决定边界。列表项需要稳定 key；生命周期副作用必须有清理，并避免在渲染阶段执行外部写操作。

## 来源

- [React Documentation](https://react.dev/learn)
- [React Documentation](https://react.dev/learn/state-a-components-memory)
- [React Documentation](https://react.dev/learn/synchronizing-with-effects)

访问日期：2026-07-16。
