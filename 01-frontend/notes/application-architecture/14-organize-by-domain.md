---
title: 按业务领域组织模块：让变化沿用例而不是文件类型聚合
stage: intermediate
direction: frontend
tags:
  - architecture
  - domain-modules
  - modularity
---

# 按业务领域组织模块：让变化沿用例而不是文件类型聚合

按领域组织把同一业务能力的 UI、用例、类型、测试和适配器放在可识别模块中。目录反映业务语言和依赖边界，而不是把所有 components、hooks、services 堆成全局技术层。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理前端模块切分、公开入口和所有权；完整领域驱动设计战术模式不是前置要求。

## 1. 定义、所有权与数据流

按领域组织把同一业务能力的 UI、用例、类型、测试和适配器放在可识别模块中。目录反映业务语言和依赖边界，而不是把所有 components、hooks、services 堆成全局技术层。

```mermaid
flowchart LR
    A0["应用入口"] --> A1
    A1["领域模块公开 API"] --> A2
    A2["用例与组件"] --> A3
    A3["领域端口"] --> A4
    A4["基础设施适配"]
```

按业务领域组织把同一能力的组件、用例、类型、端口和测试放进一个可识别模块。模块通过公开入口协作，目录反映业务变化与所有权，不是把技术文件简单换个文件夹。

## 2. 关键机制

### 2.1 领域识别

按订单、账单、身份等业务能力及变化共振切分。

若边界缺失，按页面截图随意切模块。

验证：变更历史聚类。

### 2.2 公开入口

模块只从 index 导出允许依赖的契约。

若边界缺失，深层导入内部文件。

验证：lint 禁止越界。

### 2.3 内部凝聚

类型、用例和测试靠近共同变化原因。

若边界缺失，技术目录让一次功能跨十处。

验证：统计变更文件跨度。

### 2.4 共享内核

只放稳定跨域值对象和协议。

若边界缺失，shared 变成无主杂物箱。

验证：每项有 owner 与消费者。

### 2.5 跨域协作

通过公开命令、事件或只读投影。

若边界缺失，直接读另一个域 store。

验证：依赖图审计。

### 2.6 路由组合

应用层组合领域路由，不让领域反向依赖 app。

若边界缺失，模块导入根路由。

验证：构建层级检查。

### 2.7 数据边界

每个域拥有 query key 和 DTO 映射。

若边界缺失，全局 api.ts 理解全部业务。

验证：适配器归属清晰。

### 2.8 测试边界

领域用例可用内存端口独立运行。

若边界缺失，只能启动整站测试。

验证：模块测试速度与隔离。

### 2.9 团队所有权

目录与评审责任一致。

若边界缺失，跨团队共享文件频繁冲突。

验证：CODEOWNERS 统计。

### 2.10 渐进迁移

从高变更用例建立新边界，旧代码经兼容入口调用。

若边界缺失，一次性搬目录制造风险。

验证：逐模块依赖收敛。

## 3. 从变更历史识别领域边界

抽取最近十个需求，记录每次同时修改的规则、组件和接口。退款资格、退款理由和退款状态经常共同变化，应归入 refunds；通用 Button 因品牌和无障碍变化而修改，属于 shared-ui。页面相邻不代表同领域，复用同一个 DTO 也不代表共享所有权。

## 4. 一个领域模块的内部方向

1. domain 公开 types、commands、queries 和必要 UI，内部 DTO 与 reducer 不导出。

2. application 用例编排领域规则并依赖本模块定义的端口。

3. infrastructure 适配 HTTP 或 SDK 并把 DTO 转领域对象。

4. ui 调用用例、消费领域投影，不直接拼接 API URL。

5. 模块 index 只导出承诺兼容的入口，内部文件使用相对导入避免 barrel 回环。

## 5. 应用案例一：退款领域

1. 建立 refunds/application、domain、infrastructure、ui 和 index。

2. 把 canRefund 的证据与规则移入退款域，而不是订单 Badge。

3. Order 页面只从 refunds 公开入口导入 RefundPanel。

4. RefundGateway 的内存实现运行资格和提交用例测试。

5. 改变后端退款 DTO 时只修改适配器与其契约 fixture。

结果：规则修改集中在退款模块。

失败分支：若 refunds 导入 orders 内部 reducer，边界仍未成立。

## 6. 应用案例二：身份与权限

1. identity 提供当前主体与会话生命周期。

2. authorization 把 capability 作为稳定只读查询，不暴露角色字符串。

3. 订单域声明 order.refund 等能力并决定业务操作组合。

4. 应用层在登录/登出时协调两个域和缓存清理。

5. 直接调用服务端接口仍执行授权，前端 capability 只控制体验。

结果：页面不散落角色字符串。

失败分支：前端能力判断不替代服务端授权。

## 7. TypeScript 核心实现

下面代码只实现本主题的核心契约；网络、DOM 或存储副作用留在调用边界。

```tsx
// refunds/index.ts
export { RefundPanel } from "./ui/RefundPanel";
export { requestRefund } from "./application/requestRefund";
export type { RefundGateway, RefundReason } from "./application/ports";

// 其他模块只能从公开入口导入：
import { RefundPanel } from "../refunds";
```

类型检查用于排除结构错误，运行时仍需校验外部输入、测试时序并执行安全约束。

## 8. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| 按技术类型 | 小型原型、模块少 | 增长后一次功能跨目录 |
| 按页面 | 页面高度独立 | 共享领域规则易复制 |
| 按领域 | 业务持续演进、多人协作 | 需维护公开 API 和跨域协议 |

选择应以所有权、生命周期、订阅范围和失败成本为依据。引入库不能替代这些判断；库只提供实现机制。

## 9. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 一次需求改十目录 | 是否按技术类型分散 | 建立领域垂直切片 |
| shared 无人维护 | 是否缺 owner | 移回领域或指定内核 |
| 跨域深层导入 | 缺少 exports 规则 | 公开入口加 lint |
| 模块互相调用 store | 协作协议不清 | 命令或投影 |
| 重复 DTO | 基础设施边界散落 | 每域统一映射 |
| 循环依赖 | 双向业务知识 | 提取协调用例 |
| 测试只能 E2E | 模块依赖具体设施 | 端口与内存适配 |
| 迁移停滞 | 新旧没有桥接 | 兼容入口与指标 |

调试顺序是：确认输入事实，再检查所有者和转换，随后检查订阅与渲染，最后检查异步资源。跳过前序证据直接增加 Effect，通常会制造第二个状态源。

## 10. 性能、安全与运维边界

- 模块入口控制可见 API。
- 领域类型不直接等于后端 DTO。
- 共享内核保持小且有版本责任。
- 路由和组合根位于应用层。
- 权限由服务端最终执行。
- 构建输出检查跨模块 chunk 和包体。
- 变更日志记录公开契约。
- 依赖可视化进入 CI。

生产验证至少记录一次正常路径和一次故障路径；对“按业务领域组织模块”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 11. 与其他架构模块集成

- 单向依赖规则执行模块边界。
- 依赖倒置隔离基础设施。
- ADR 记录争议切分。
- 循环检测防止边界腐化。

集成时先画出事实所有者，跨边界只传递稳定契约。不要为了减少一层调用而复制同一事实。

## 12. 综合练习

把按 components/hooks/services 组织的订单退款功能迁移为领域模块，并保持旧入口可用。

### 验收标准

- [ ] 退款规则及适配器集中在 refunds。
- [ ] 其他域只能从公开 index 导入。
- [ ] 订单与库存协作位于上层协调器且无环。
- [ ] shared 每项有明确 owner 和两个真实消费者。
- [ ] 替换退款 HTTP DTO 不改页面和领域用例。

## 13. 跨领域协作的三种形式

同步查询适合读取稳定投影，例如订单域读取“当前主体是否具有 refund capability”；它不应取得 authorization 内部 store。命令适合明确请求另一个域执行动作，并返回领域结果。事件适合已经发生的事实通知多个未知消费者，但需要事件版本、幂等和失败观测。

订单创建后立即预留库存同时需要两个域时，把用例放在 application orchestration：它依赖 OrderPort 与 InventoryPort。不要让 orders 导入 inventory 内部函数，同时又让 inventory 回调 orders。上移协调器可以保持两个领域无环。

`shared` 只接受稳定且确有多个 owner 共同承诺的内容。为了消除 import 错误把 `OrderStatus` 搬进 shared，会丢失订单所有权。共享内核的每个导出都需要负责人、消费者清单和兼容策略。

## 14. 目录与公开契约示例

```text
src/
  app/
    routes.tsx
    composition-root.ts
  domains/
    refunds/
      application/
        request-refund.ts
        ports.ts
      domain/
        refund.ts
        eligibility.ts
      infrastructure/
        http-refund-gateway.ts
      ui/
        refund-panel.tsx
      index.ts
  shared/
    ui/
    platform/
```

`refunds/index.ts` 只导出 `RefundPanel`、`requestRefund` 和调用方需要的公开类型。`http-refund-gateway.ts`、后端 DTO 与内部状态枚举不导出。退款模块内部也不从自己的 index 回导，否则 barrel 可能形成初始化环。

应用组合根导入 `RefundGateway` 端口与 HTTP 实现，创建用例后传给页面。若希望领域模块完全不知道 React，可让页面适配用例结果；若领域 UI 是团队承诺的公开能力，也可从模块入口导出，但仍不能让 domain 目录依赖 React。

## 15. 切分过细与过粗的信号

只有一个类型和一个函数、每次都与另一个模块同时修改的“领域”，可能切得过细。包含十几个互不相关用例、需要多个团队共同批准任何修改的模块，可能过粗。

可记录三个月数据：

- 一次需求跨多少模块；
- 哪些文件经常共同修改；
- 跨模块公开 API 的变更次数；
- 循环依赖与例外数量；
- 每个模块的负责人和独立测试耗时。

数据不能自动给出领域答案，但能暴露边界与真实变化不一致。重切边界时用兼容入口迁移消费者，先阻止新增深层导入，再逐步移动实现，避免一次性目录搬迁与功能修改混在同一提交。

## 16. 模块测试与发布边界

退款模块的应用测试使用内存 `RefundGateway` 覆盖资格通过、证据不足、重复请求和取消。HTTP 适配器使用契约 fixture 覆盖 DTO 转换。`RefundPanel` 只测试可访问交互和用例结果，不在组件测试里重复后端规则。

若领域模块作为独立 package 发布，`package.json#exports` 只列公开入口，类型声明也不能泄漏内部 DTO。若仍在单仓源码中，ESLint/Nx 边界规则承担同一职责。两种形式都要用一个故意深层导入的失败 fixture 证明约束真实生效。

跨模块事件需要名字、载荷 schema、版本和幂等语义。`refund.completed.v1` 表示已经发生的事实；消费者失败不能回滚退款本身。需要原子业务结果时应使用同步协调用例，而不是广播事件后假定所有消费者成功。

领域拆分也影响 bundle：只有进入退款页面才需要的复杂编辑器可以从退款公开入口动态加载；共享基础按钮不应在多个领域重复打包。构建 metafile 用于检查边界是否造成意外公共 chunk。

## 来源

- [TypeScript：Modules](https://www.typescriptlang.org/docs/handbook/2/modules.html)（访问日期：2026-07-18）
- [Node.js：Package entry points](https://nodejs.org/api/packages.html#package-entry-points)（访问日期：2026-07-18）
- [Nx：Enforce Module Boundaries](https://nx.dev/features/enforce-module-boundaries)（访问日期：2026-07-18）
- [Vite：Features](https://vite.dev/guide/features.html)（访问日期：2026-07-18）
