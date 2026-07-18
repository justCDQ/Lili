---
title: 状态所有权与单一数据源：让每个事实只有一个权威写入者
stage: intermediate
direction: frontend
tags:
  - architecture
  - state-ownership
  - single-source
---

# 状态所有权与单一数据源：让每个事实只有一个权威写入者

状态所有权说明谁可以权威修改一个事实；单一数据源要求同一事实不在多个可写位置独立维护。组件可以拥有缓存、投影和草稿，但必须能指出其来源、同步方向和冲突规则。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理前端状态归属、派生与同步；数据库主从和事件溯源不在本文范围。

## 1. 定义、所有权与数据流

状态所有权说明谁可以权威修改一个事实；单一数据源要求同一事实不在多个可写位置独立维护。组件可以拥有缓存、投影和草稿，但必须能指出其来源、同步方向和冲突规则。

```mermaid
flowchart LR
    A0["权威事实"] --> A1
    A1["受控命令"] --> A2
    A2["唯一写入者"] --> A3
    A3["派生快照"] --> A4
    A4["多个消费者"]
```

状态所有权定义唯一的权威写入者；单一数据源禁止同一事实在多个可写位置独立演进。缓存、草稿和投影可以存在，但必须标明来源、方向与冲突处理。

## 2. 关键机制

### 2.1 权威性

服务端记录、URL、表单草稿或组件交互各有不同权威源。

若边界缺失，把所有东西都称为全局真相。

验证：为每个字段标记 owner。

### 2.2 生命周期

所有者至少活到最后一个消费者不再需要该事实。

若边界缺失，状态放得过低导致卸载丢失。

验证：导航和卸载测试。

### 2.3 写权限

消费者通过命令请求变更，不直接改共享对象。

若边界缺失，任意写入破坏不变量。

验证：公开动作白名单。

### 2.4 派生值

可由其他字段纯计算的值不另存。

若边界缺失，items 与 total 漂移。

验证：属性测试比较派生。

### 2.5 草稿与已提交

草稿属于表单，已提交事实属于服务端。

若边界缺失，背景刷新覆盖编辑。

验证：dirty 时进入冲突策略。

### 2.6 缓存

缓存有来源、时间和失效规则，不成为第二权威。

若边界缺失，修改缓存后假定服务端成功。

验证：mutation 对账。

### 2.7 受控边界

父组件拥有 value，子组件只发 onChange。

若边界缺失，父子同时维护副本。

验证：Props 变化与事件序列测试。

### 2.8 身份

稳定 id 连接多个投影，数组索引不承担身份。

若边界缺失，重排后选择错对象。

验证：插入删除重排测试。

### 2.9 同步

必须定义单向源→投影，双向同步需要冲突协议。

若边界缺失，两个 Effect 互相回写。

验证：数据流图检查环。

### 2.10 重置

登出、切租户和资源删除明确销毁所有投影。

若边界缺失，旧用户状态残留。

验证：边界端到端测试。

## 3. 建立状态清单与所有权表

给每个状态记录权威源、写命令、生命周期、消费者、是否派生、是否持久化和冲突规则。`filter` 若由 URL 拥有，组件只解析 location；`draft.name` 由表单拥有，服务端背景快照不能直接 reset；`isEmpty` 可由 items.length 派生，不应另存。只有写出这张表，团队才能区分“共享”与“权威”。

## 4. 消除双写的迁移顺序

1. 记录两个副本的全部写入路径和读消费者。

2. 选择能覆盖完整生命周期且有权执行不变量的一方作为 owner。

3. 把另一方改为纯派生或只读缓存，停止反向 Effect。

4. 将变更改成语义命令并迁移消费者。

5. 用前进后退、背景刷新、乱序响应和重置测试证明没有回弹。

## 5. 应用案例一：筛选表格

1. URL 拥有 query/category/page，解析结果生成 query key。

2. 表格本地只拥有 hover 和列宽，不初始化一份 filters。

3. 修改筛选执行导航并重置 page，popstate 后重新解析。

4. 查询缓存拥有结果快照，不复制到全局 store。

5. 刷新、分享和前进后退得到相同控件与请求。

结果：刷新和返回恢复同一查询，不存在 filter useState 副本。

失败分支：若本地筛选与 URL 同时更新，快速后退会产生回弹。

## 6. 应用案例二：编辑资料

1. 服务端 Profile 是基线，表单 draft 是未提交意图。

2. dirty=false 时背景新版本可替换基线；dirty=true 时进入 conflict。

3. 提交携带 baseVersion，409 返回 latest。

4. 用户选择覆盖、放弃或逐字段合并，不能由 Effect 自动决定。

5. 成功响应同时成为新基线并清除 dirty。

结果：提交携带 version，409 后让用户比较。

失败分支：简单 useEffect(reset(serverData)) 会销毁未提交输入。

## 7. TypeScript 核心实现

下面代码只实现本主题的核心契约；网络、DOM 或存储副作用留在调用边界。

```tsx
type ProfileState =
  | { kind: "clean"; server: Profile }
  | { kind: "editing"; baseVersion: number; draft: Profile }
  | { kind: "conflict"; baseVersion: number; draft: Profile; latest: Profile };
type Profile = { id: string; version: number; name: string };
export function receive(state: ProfileState, latest: Profile): ProfileState {
  return state.kind === "clean"
    ? { kind: "clean", server: latest }
    : { kind: "conflict", baseVersion: state.baseVersion, draft: state.draft, latest };
}
```

类型检查用于排除结构错误，运行时仍需校验外部输入、测试时序并执行安全约束。

## 8. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| 局部拥有 | 消费者同一子树且生命周期短 | 跨页会丢失 |
| 提升到共同祖先 | 少数兄弟需同步 | 祖先可能重渲染 |
| 外部权威源 | 跨树、跨页或远端事实 | 需订阅和冲突规则 |

选择应以所有权、生命周期、订阅范围和失败成本为依据。引入库不能替代这些判断；库只提供实现机制。

## 9. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 值来回跳 | 是否两个可写副本 | 选择 owner 并删除反向 Effect |
| 总数不对 | 是否保存派生值 | 用 selector 计算 |
| 刷新覆盖草稿 | 是否混淆基线与草稿 | 版本冲突状态 |
| 重排选错行 | 是否用 index 身份 | 稳定 id |
| 登出残留 | 是否缺少 reset | 销毁用户命名空间 |
| 测试依赖顺序 | 是否共享单例 | 每例新建 owner |
| 缓存假成功 | 是否跳过 mutation | 服务端对账 |
| 返回键失效 | URL 是否只是初始化源 | 持续从 location 派生 |

调试顺序是：确认输入事实，再检查所有者和转换，随后检查订阅与渲染，最后检查异步资源。跳过前序证据直接增加 Effect，通常会制造第二个状态源。

## 10. 性能、安全与运维边界

- 为关键状态建立 ownership 表。
- 服务端权威字段不得仅靠客户端限制。
- 派生 selector 要纯且可测试。
- 异步写入带版本或关联 ID。
- 登出与租户切换销毁投影。
- 日志区分草稿、缓存和已提交事实。
- 避免 Effect 维持双向镜像。
- 性能优化不能复制第二份可写状态。

生产验证至少记录一次正常路径和一次故障路径；对“状态所有权与单一数据源”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 11. 与其他架构模块集成

- URL 拥有导航事实。
- Query cache 拥有远端快照。
- Form 拥有未提交字段。
- Interaction state 留在局部组件。

集成时先画出事实所有者，跨边界只传递稳定契约。不要为了减少一层调用而复制同一事实。

## 12. 综合练习

为企业后台列出 20 个状态的所有者、生命周期、消费者、持久化和冲突策略，并消除至少三个双写。

### 验收标准

- [ ] 至少20个状态均标明 owner 与生命周期。
- [ ] 删除三个可由事实派生的可写字段。
- [ ] URL 与本地筛选不再双向 Effect 同步。
- [ ] 背景刷新不覆盖 dirty 表单。
- [ ] 登出与资源删除清理全部投影。

## 13. 派生值与缓存投影

派生不等于每次都昂贵重算。可以用 memoized selector 缓存计算，但缓存没有独立写入口：输入引用不变时复用结果，输入变化时重新计算。总价、已选数量、是否为空和按钮 disabled 常属于派生值。

远端查询缓存不同于纯派生，它保存异步快照并有 stale 与失效规则，但权威写入者仍是服务端。乐观更新是有回滚和对账的预测，不会把缓存提升为数据库真相。

审查发现两个 state 名字相同还不够，要比较语义。例如 `selectedOrderId` 可以由 URL 拥有，而拖拽期间的 `previewOrderId` 属于局部交互；它们不是重复事实。相反，`filters` 与 `searchParams` 若能无损互转并互相回写，通常是同一事实的双写。

## 14. 所有权冲突审查实例

一个订单后台同时出现以下字段：

| 字段 | 权威所有者 | 允许的投影 | 冲突处理 |
|---|---|---|---|
| `status` | 服务端订单记录 | Query cache、状态标签 | mutation 后按服务端 version 对账 |
| `statusFilter` | URL | 筛选控件值、query key | 非法值规范化后 replace |
| `editNote` | 当前表单 | 离开确认、草稿白名单 | 背景刷新不覆盖 dirty 值 |
| `activeRowId` | 表格交互 Root | roving tabindex | 行删除时选择后继 |
| `canRefund` | 服务端授权投影 | 按钮可见性 | 敏感操作再次由服务端授权 |

`status` 和 `statusFilter` 名字相近但不是同一事实；前者描述单个订单，后者描述集合查询。`editNote` 与订单 `note` 描述同一业务字段的不同时间版本，因此需要 baseVersion 和提交对账。

### 14.1 用删除测试识别派生字段

尝试删除 `isEmpty`、`selectedCount`、`canSubmit` 的写入代码，改为从 items、selection 和 form validity 计算。如果功能与性能都满足，它们不应拥有独立 state。性能不足时添加 memoized selector，而不是恢复第二个写入口。

### 14.2 用时间线识别竞态

记录 `t0` 加载 version 3、`t1` 用户编辑、`t2` 背景收到 version 4、`t3` 用户提交的完整时间线。每个节点写出谁能修改什么。若 t2 的 Effect 直接 reset 表单，它越过草稿所有者；正确模型在 t2 保存 latest 并进入 conflict。

### 14.3 用重置测试检查边界

执行登出、租户切换、删除订单和浏览器后退。每个事件都应有明确的销毁清单。仅清空全局 store 而保留 Query cache 或 localStorage，仍会让下一主体看到旧投影。

## 15. 所有权评审问题

评审一个新增 state 时要求提交者回答：哪个事件能修改它、删除组件后是否仍需要、能否由已有事实计算、刷新是否恢复、与服务端冲突谁获胜。答案必须落到具体对象，不能只写“放全局方便使用”。

若多个模块都要发写命令，仍可有单一所有者：命令入口集中验证并发布快照。单一数据源限制的是权威写入位置，不限制消费者数量，也不要求所有数据进入同一个 store。

## 来源

- [React：Sharing State Between Components](https://react.dev/learn/sharing-state-between-components)（访问日期：2026-07-18）
- [React：Choosing the State Structure](https://react.dev/learn/choosing-the-state-structure)（访问日期：2026-07-18）
- [React：You Might Not Need an Effect](https://react.dev/learn/you-might-not-need-an-effect)（访问日期：2026-07-18）
