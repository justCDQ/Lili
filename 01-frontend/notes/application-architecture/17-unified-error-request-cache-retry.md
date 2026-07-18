---
title: 统一错误、请求、缓存与重试：建立端到端失败语义
stage: intermediate
direction: frontend
tags:
  - architecture
  - errors
  - http
  - retry
---

# 统一错误、请求、缓存与重试：建立端到端失败语义

统一模型把请求生命周期、错误分类、缓存身份和重试决策连接起来。统一不等于一个万能 fetch；每层保留自己的语义，并通过稳定错误联合、trace id 和策略接口协作。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理浏览器请求平台和查询层；服务端内部熔断与消息重投不展开。

## 1. 定义、所有权与数据流

统一模型把请求生命周期、错误分类、缓存身份和重试决策连接起来。统一不等于一个万能 fetch；每层保留自己的语义，并通过稳定错误联合、trace id 和策略接口协作。

```mermaid
flowchart LR
    A0["领域请求"] --> A1
    A1["HTTP 客户端"] --> A2
    A2["响应归一"] --> A3
    A3["缓存策略"] --> A4
    A4["重试决策"] --> A5
    A5["UI 恢复"]
```

统一请求模型连接传输、响应校验、错误分类、缓存和重试，但各层仍有明确职责。HTTP client 不决定业务缓存，查询层不猜授权语义，页面不按错误字符串分支。

## 2. 关键机制

### 2.1 请求标识

每次逻辑操作有 correlation/idempotency 标识。

若边界缺失，重复提交无法对账。

验证：日志串联。

### 2.2 错误分类

transport、unauthorized、validation、conflict、rate-limit 分开。

若边界缺失，统一 throw Error 丢语义。

验证：穷尽处理。

### 2.3 超时取消

超时策略产生 AbortSignal，用户取消单独识别。

若边界缺失，超时后请求仍运行。

验证：慢服务测试。

### 2.4 状态码

HTTP 状态先解析，再映射领域错误。

若边界缺失，只判断 response.ok 文本。

验证：fixture 覆盖 401/409/429/5xx。

### 2.5 响应校验

成功 JSON 也需 schema 校验。

若边界缺失，字段缺失进入 UI。

验证：畸形响应测试。

### 2.6 缓存 key

由规范输入生成。

若边界缺失，用户或筛选串缓存。

验证：key 工厂。

### 2.7 重试条件

仅幂等或具幂等键且错误可重试。

若边界缺失，POST 自动重复扣款。

验证：故障注入。

### 2.8 退避抖动

指数退避加随机抖动并尊重 Retry-After。

若边界缺失，客户端同步重试雪崩。

验证：虚拟时钟。

### 2.9 可观测性

记录耗时、状态、错误 kind、重试次数，不记秘密。

若边界缺失，只有 console error。

验证：trace 对账。

### 2.10 UI 恢复

错误模型决定登录、字段修正、重试或联系支持。

若边界缺失，所有错误只显示再试。

验证：交互分支测试。

## 3. 一次请求的可验证管线

领域适配器构造 method、URL 和 body；HTTP client 注入会话、trace id、超时与 AbortSignal；响应先按状态读取，再校验 payload schema；适配器把传输结果映射为领域联合；查询层按 key 缓存或失效；UI 根据 error kind 提供登录、修正、等待或重试。每一步保留关联 ID。

## 4. 运行顺序与边界

1. 认证刷新使用 single-flight：并发 401 只启动一次刷新，其余等待同一结果。

2. 刷新成功只重放满足安全条件的原请求；刷新失败统一结束会话。

3. 可重试 GET 使用指数退避、抖动和上限，并尊重 Retry-After。

4. 写请求只有带服务端幂等语义时才能在结果未知后重试。

5. 缓存按用户隔离，登出同时清理查询与正在等待的刷新队列。

## 5. 应用案例一：订单提交

1. 提交订单生成一次 idempotencyKey，所有重试复用同一键。

2. 客户端超时后先进入 unknown，不立即宣称失败。

3. 用同一键查询或重试，服务端返回原 orderId。

4. 422 字段错误和409业务冲突都不自动重试。

5. 日志用 traceId 与幂等键哈希对账，不记录购物车敏感字段。

结果：服务端返回同一订单，UI 不重复。

失败分支：无幂等能力时超时必须进入结果未知状态并查询确认。

## 6. 应用案例二：列表限流

1. 列表 GET 收到429，解析 Retry-After 秒或日期。

2. 保留旧缓存并显示最后更新时间。

3. 到期后加小幅抖动刷新，多个标签避免同时唤醒。

4. 第二次429增加等待但不超过页面重试预算。

5. 用户手动重试仍受全局限流时间约束。

结果：用户仍能浏览旧数据且看到更新时间。

失败分支：忽略 Retry-After 会加重限流。

## 7. TypeScript 核心实现

下面代码只实现本主题的核心契约；网络、DOM 或存储副作用留在调用边界。

```tsx
type AppError =
  | { kind: "network"; retryable: true }
  | { kind: "unauthorized"; retryable: false }
  | { kind: "conflict"; retryable: false }
  | { kind: "rateLimit"; retryable: true; retryAfterMs: number };
export function retryDelay(attempt: number, error: AppError, random: number): number | null {
  if (!error.retryable || attempt >= 3) return null;
  if (error.kind === "rateLimit") return error.retryAfterMs;
  return Math.min(30_000, 500 * 2 ** attempt) * (0.5 + random * 0.5);
}
```

类型检查用于排除结构错误，运行时仍需校验外部输入、测试时序并执行安全约束。

## 8. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| 页面自行 fetch | 极小原型 | 行为不一致 |
| 共享 HTTP client | 统一传输与错误 | 不应拥有业务缓存 |
| Query 层 | 缓存、失效、重试 | 需正确 key 与策略 |

选择应以所有权、生命周期、订阅范围和失败成本为依据。引入库不能替代这些判断；库只提供实现机制。

## 9. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 重复扣款 | POST 自动重试 | 幂等键 |
| 401 重试风暴 | 认证错误可重试 | 刷新令牌单航班 |
| 429 更严重 | 忽略 Retry-After | 延迟与抖动 |
| 缓存串用户 | key/实例未隔离 | 会话清理 |
| 错误不可恢复 | 只保留 message | 判别联合 |
| 取消显示失败 | 未区分 Abort | 单独 cancelled |
| 响应字段崩溃 | 无 schema | 边界校验 |
| 日志无法串联 | 无 trace id | 请求上下文 |

调试顺序是：确认输入事实，再检查所有者和转换，随后检查订阅与渲染，最后检查异步资源。跳过前序证据直接增加 Effect，通常会制造第二个状态源。

## 10. 性能、安全与运维边界

- 重试预算有上限。
- 非幂等写入不盲重试。
- 401 刷新使用 single-flight。
- 429 尊重服务端时间。
- 日志脱敏 header 和 body。
- 缓存按用户隔离。
- schema 错误告警但不渲染脏数据。
- 监控成功率、尾延迟、重试放大率。

生产验证至少记录一次正常路径和一次故障路径；对“统一错误、请求、缓存与重试”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 11. 与其他架构模块集成

- 适配器归一第三方错误。
- Query cache 承担缓存。
- State Machine 管理结果未知流程。
- UI 错误边界处理渲染异常而非业务错误。

集成时先画出事实所有者，跨边界只传递稳定契约。不要为了减少一层调用而复制同一事实。

## 12. 综合练习

实现带超时、取消、错误联合、响应校验、幂等重试、429 退避和 trace 日志的请求平台。

### 验收标准

- [ ] 错误联合覆盖网络、认证、校验、冲突、限流和取消。
- [ ] 并发401只执行一次刷新且失败只登出一次。
- [ ] 订单超时复用幂等键并处理unknown。
- [ ] 429尊重Retry-After和重试预算。
- [ ] 日志可按trace关联且已脱敏。

## 12. 认证刷新 single-flight 状态机

状态为 idle、refreshing、failed。第一个 401 从 idle 进入 refreshing 并保存刷新 Promise；后续 401 订阅同一结果。成功后回到 idle，等待请求各自决定是否重放；失败进入 failed，取消队列并触发一次登出。

如果刷新请求本身也经过会自动刷新的拦截器，会形成递归。刷新端点必须显式跳过该逻辑。重放前还要检查原 AbortSignal、请求体是否可重放以及幂等条件。

观测指标包括原始请求数、重试请求数、retry amplification、401 并发数、刷新次数、429 等待时间和取消率。成功率提高但放大率暴涨不是健康改进。

## 13. HTTP 响应读取顺序

不能对所有响应直接 `response.json()`。204 没有正文；错误正文可能是 `application/problem+json`、普通 JSON、文本或空。client 先检查 status 与 Content-Type，在最大尺寸限制内读取，再交给对应 schema。

成功状态也可能返回畸形对象。schema 失败映射为 protocol error，记录 endpoint、status、trace id 和字段路径，UI 显示稳定错误；不能把未经验证的对象传入缓存。

401 表示当前凭证不足，不一定都能刷新；403 通常是已认证但无权限；404 对详情可能是不存在，对搜索接口也可能代表路由配置错误。HTTP client 保留传输分类，领域适配器根据接口契约增加业务含义。

## 14. 超时、取消与结果未知

用户取消、页面卸载、客户端超时和浏览器离线都是不同原因。它们都可表现为 AbortError，但适配层应在 abort reason 或本地上下文中区分，决定是否提示和是否重试。

GET 超时通常可以在预算内重试；创建订单超时可能已在服务端成功。若服务端支持幂等键，同键重试安全取得同一结果；否则进入 unknown 并提供状态查询。换新键自动重试会把一次意图变成两次订单。

总超时预算包含每次请求和退避等待。例如用户只容忍 10 秒，不能执行三次各 10 秒再加退避。重试函数接收剩余预算，下一次预计无法完成时停止。

## 15. 缓存错误与旧数据策略

后台刷新失败时，已有可用快照不必清空。查询状态可以同时是 `data` 存在、`isFetching` false、最近刷新 error 存在。UI 显示旧数据和更新时间，并提供重试；初次加载失败则没有可展示快照。

认证错误应清理或隔离用户缓存；限流错误保留旧值到 Retry-After；schema 错误不能继续标记新响应为成功。错误缓存的持续时间也要明确，避免每次 render 立即请求同一个确定失败的资源。

## 16. 失败注入脚本的验收矩阵

测试服务依次返回：连接重置、延迟到超时、401 后刷新成功、十个并发 401、刷新失败、422 字段错误、409 冲突、429 带 `Retry-After: 2`、503 两次后成功、200 但 JSON schema 错误。

每个场景断言实际请求次数、等待时间、最终错误 kind、缓存内容和用户恢复入口。十个并发 401 只能产生一个 refresh；429 在两秒前不重发；schema 错误不能进入成功缓存；取消请求不显示红色业务失败。

使用虚拟时钟验证退避，使用可预测 random 输入验证抖动范围。真实端到端测试再确认浏览器 AbortSignal、网络面板请求数和服务端 trace 能对齐。

## 来源

- [HTTP Semantics RFC 9110](https://www.rfc-editor.org/rfc/rfc9110)（访问日期：2026-07-18）
- [MDN：AbortSignal](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal)（访问日期：2026-07-18）
- [TanStack Query：Query Retries](https://tanstack.com/query/latest/docs/framework/react/guides/query-retries)（访问日期：2026-07-18）
- [OpenTelemetry：HTTP semantic conventions](https://opentelemetry.io/docs/specs/semconv/http/)（访问日期：2026-07-18）
