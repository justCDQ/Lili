---
title: Global State：按共享范围和订阅粒度建立应用级状态
stage: intermediate
direction: frontend
tags:
  - architecture
  - global-state
  - external-store
---

# Global State：按共享范围和订阅粒度建立应用级状态

Global State 不是“很多组件都用”的同义词，而是生命周期和所有权确实跨越多个远距离子树的客户端状态。应用级 store 应提供明确命令、选择器、实例隔离和销毁策略。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理客户端全局状态和外部 store；服务端缓存、URL、表单与持久化只说明边界。

## 1. 定义、所有权与数据流

Global State 不是“很多组件都用”的同义词，而是生命周期和所有权确实跨越多个远距离子树的客户端状态。应用级 store 应提供明确命令、选择器、实例隔离和销毁策略。

```mermaid
flowchart LR
    A0["领域命令"] --> A1
    A1["store reducer"] --> A2
    A2["不可变快照"] --> A3
    A3["selector"] --> A4
    A4["订阅组件"]
```

Global State 是生命周期跨越多个远距离子树的客户端事实。外部 store 的核心不是全局变量，而是可隔离实例、不可变快照、语义命令和细粒度订阅。

## 2. 关键机制

### 2.1 提升条件

会话、跨页面购物车或全局工作区选择确实跨树共享。

若边界缺失，为了避免传两层 Props 就全局化。

验证：画消费者和生命周期图。

### 2.2 快照

store 暴露不可变快照，更新后引用变化。

若边界缺失，原地修改导致 selector 看不到更新。

验证：冻结测试和引用断言。

### 2.3 命令

组件发送语义动作，不任意 set 任意字段。

若边界缺失，调用方可以制造非法组合。

验证：公开动作白名单。

### 2.4 selector

消费者订阅最小派生切片，结果身份要稳定。

若边界缺失，每个组件订阅整个 store。

验证：Profiler 检查提交范围。

### 2.5 一致性

React 外部 store 通过 useSyncExternalStore 提供一致快照。

若边界缺失，自写 Effect 订阅可能 tearing。

验证：并发渲染测试与官方适配。

### 2.6 实例范围

浏览器会话、测试和 SSR 请求各有独立实例。

若边界缺失，模块单例在 SSR 串用户。

验证：并行请求隔离测试。

### 2.7 归一化

大量实体按 id 存储，关系用 id 引用。

若边界缺失，嵌套复制难以原子更新。

验证：实体删除后引用完整性测试。

### 2.8 派生数据

可计算结果通过 selector 得出，不持久保存重复字段。

若边界缺失，total 与 items 漂移。

验证：属性测试比较派生结果。

### 2.9 中间件

日志、持久化和 devtools 是边界插件，不改变 reducer 语义。

若边界缺失，中间件吞事件或记录秘密。

验证：脱敏和故障旁路测试。

### 2.10 销毁

动态 workspace 或微前端 store 需取消订阅和定时器。

若边界缺失，切换租户后旧状态残留。

验证：destroy 后监听器数为零。

## 3. React 外部 Store 的一致性契约

`useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot?)` 要求缓存快照：底层数据未变化时 `getSnapshot` 必须返回同一引用。subscribe 返回取消函数。服务端渲染使用的 `getServerSnapshot` 必须与 hydration 的初始值一致，否则客户端会重新解释页面。store 实例由 Context 注入可以解决依赖范围，但 Context value 应是稳定实例，不是每次变化的完整快照。

## 4. Store 实例与会话生命周期

1. 浏览器应用启动时为当前会话创建实例；SSR 每个请求创建，不导出携带用户数据的模块单例。

2. 组件通过 selector 订阅所需切片；selector 返回数组或对象时保持结构共享或使用相等比较。

3. 领域动作在 store 内保持购物车数量、工作区身份等不变量，组件不能任意 set 整个对象。

4. 登出先停止请求和订阅，再清空用户状态、查询缓存与持久化命名空间。

5. 测试每例创建 store 并在结束时销毁，监听器和定时器数量归零。

## 5. 应用案例一：多工作区选择

1. 定义 workspaceId 的来源：显式 URL 优先，本地默认只在 URL 缺失时生效。

2. 切换 workspace 作为单个动作，原子更新选择并触发租户缓存清理。

3. 两个组件分别订阅 id 与 capability，名称变化不应重渲染只读 capability 消费者。

4. 并行创建两个 SSR store，断言不同用户快照互不可见。

5. 登出后重新登录另一用户，旧 workspace 不得恢复。

结果：切换时清空租户级缓存并重建依赖容器。

失败分支：SSR 不能使用进程全局 store；两个请求的 workspace 不得互见。

## 6. 应用案例二：购物车

1. 购物车按 productId 保存数量，禁止小于零；总数和总价均为 selector。

2. add、remove、setQuantity 是唯一写入口。

3. 价格只用于预览，提交时服务端重新定价。

4. 模拟一百个行组件，修改一行时测量提交数量。

5. 删除商品后所有派生 selector 同时一致，不另存 total。

结果：多个入口共享同一购物车，只有受影响 selector 更新。

失败分支：价格最终由服务端结算，客户端全局状态不能决定应付金额。

## 7. TypeScript 核心实现

下面代码只实现本主题的核心契约；网络、DOM 或存储副作用留在调用边界。

```tsx
type Listener = () => void;
type Cart = Readonly<Record<string, number>>;

export function createCartStore(initial: Cart = {}) {
  let snapshot = initial;
  const listeners = new Set<Listener>();
  return {
    getSnapshot: () => snapshot,
    subscribe(listener: Listener) {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },
    setQuantity(id: string, quantity: number) {
      const next = { ...snapshot };
      if (quantity <= 0) delete next[id]; else next[id] = quantity;
      snapshot = next;
      listeners.forEach((listener) => listener());
    },
  };
}
```

类型检查用于排除结构错误，运行时仍需校验外部输入、测试时序并执行安全约束。

## 8. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| Props/Context | 共享范围局部、更新低频 | 边界清楚但深层传递增加 |
| 外部 store | 跨树且需选择订阅 | 实例和生命周期管理更复杂 |
| 服务端缓存 | 数据由远端拥有 | 不能承载未提交客户端流程 |

选择应以所有权、生命周期、订阅范围和失败成本为依据。引入库不能替代这些判断；库只提供实现机制。

## 9. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 全页重渲染 | 是否订阅整个对象 | 使用稳定 selector |
| SSR 用户串数据 | 是否导出单例 | 每请求创建实例 |
| 状态无法重置 | 是否缺少会话边界 | 登出与租户切换显式 reset |
| 派生值不一致 | 是否重复保存 | 改为 selector |
| DevTools 泄密 | 中间件记录什么 | 字段脱敏和生产禁用 |
| 更新丢失 | 是否基于旧闭包 | 原子 reducer/函数更新 |
| 测试相互污染 | 是否复用 store | 每例创建并销毁 |
| 缓存复制进 store | 是否双份远端实体 | 移回查询缓存 |

调试顺序是：确认输入事实，再检查所有者和转换，随后检查订阅与渲染，最后检查异步资源。跳过前序证据直接增加 Effect，通常会制造第二个状态源。

## 10. 性能、安全与运维边界

- 全局 store 不保存访问令牌明文或服务端权威金额。
- SSR、测试与多租户必须隔离实例。
- selector 保持纯净并避免每次返回新集合。
- 记录动作类型但脱敏载荷。
- 大型实体更新测量提交数量和选择器耗时。
- 登出清空用户级状态与相关持久化。
- 动态模块卸载时清理订阅和定时器。
- 版本升级为持久化切片提供迁移。

生产验证至少记录一次正常路径和一次故障路径；对“Global State”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 11. 与其他架构模块集成

- URL 的显式选择优先于全局默认。
- Server State 不复制进全局 store。
- Persistent State 只保存允许跨重载的切片。
- Context 可注入 store 实例，而不是广播整个变化对象。

集成时先画出事实所有者，跨边界只传递稳定契约。不要为了减少一层调用而复制同一事实。

## 12. 综合练习

实现 SSR 可隔离的多工作区 store，包含 selector、登出重置、租户切换和订阅泄漏测试。

### 验收标准

- [ ] SSR 两个请求的 store 实例隔离。
- [ ] selector 只提交受影响消费者并有 Profiler 证据。
- [ ] 派生总数未在快照中重复保存。
- [ ] 登出销毁订阅、定时器和用户状态。
- [ ] 服务端权威价格与权限不由 store 决定。

## 13. 高频订阅的测量

一次动作先记录旧快照引用，再执行命令，随后统计被通知监听器和实际 React commit。所有监听器收到通知不必然导致所有组件提交，但订阅整个 store 会让 selector 失去隔离价值。用 Profiler 比较“订阅整个购物车”和“按 productId 订阅”两版：证据应包含行数、动作次数、commit 数和耗时，而不是只写“感觉更快”。

外部事件进入 store 时还要处理批量与次序。WebSocket 连续更新使用服务端 version 丢弃旧事件；不能用到达顺序假设生成顺序。时间旅行或 DevTools 重放不得重复发送真实支付、分析等副作用。

## 14. 外部事件与批量一致性

WebSocket 推送、BroadcastChannel 和 Service Worker 都可能在 React 事件之外修改 store。适配层先校验事件 schema，再按实体 version 应用；version 小于等于当前值的消息丢弃。一次消息修改多个字段时使用一个原子命令发布一次快照，不能逐字段通知消费者看到中间状态。

离线恢复后先取得服务端增量或完整快照，再重放仍允许的本地意图。购物车数量可以按产品规则合并，权限和结算金额不能用本地最后写胜出。合并失败应成为显式 conflict 状态。

测试创建两个订阅者：一个选择 `workspaceId`，一个选择 `capabilities.canRefund`。只更新工作区显示名称时，两者都不应提交；切换工作区时两者按新快照各提交一次。随后调用 unsubscribe 和 destroy，发送外部事件不得再触发监听器。

## 15. Store API 的最小公开面

消费者只获得 `getSnapshot`、`subscribe`、领域动作和必要 selector。调试替换快照、批量导入和内部 listener 数量属于测试辅助入口，不发布给业务组件。否则任何页面都能绕过 `setQuantity` 直接写负数。

选择器按业务含义命名，例如 `selectCartLine(id)` 和 `selectCanRefund(orderId)`，而不是把内部 state 路径当 API。重构归一化结构时，只要选择器结果语义不变，消费者无需迁移。

错误动作在开发环境立即抛出，在生产环境记录带动作类型的诊断并保持旧快照。不能发布半更新对象让一部分消费者看到新 workspace、另一部分仍看到旧 capability。

## 来源

- [React：useSyncExternalStore](https://react.dev/reference/react/useSyncExternalStore)（访问日期：2026-07-18）
- [Redux：Style Guide](https://redux.js.org/style-guide/)（访问日期：2026-07-18）
- [Redux Toolkit：Usage Guide](https://redux-toolkit.js.org/usage/usage-guide)（访问日期：2026-07-18）
- [React：Choosing the State Structure](https://react.dev/learn/choosing-the-state-structure)（访问日期：2026-07-18）
