---
title: Persistent State：带版本、期限与隐私边界的客户端恢复
stage: intermediate
direction: frontend
tags:
  - architecture
  - persistence
  - web-storage
---

# Persistent State：带版本、期限与隐私边界的客户端恢复

持久化层保存的是某一时刻、某一 Schema 版本的可恢复快照，而不是自动同步的权威事实。恢复路径必须依次完成读取、解析、版本识别、迁移、校验和过期判断，任一步失败都应回退到安全默认值。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理 Web Storage、IndexedDB 与前端快照恢复；服务端数据库、离线优先同步协议另行讨论。

## 1. 定义、所有权与数据流

Persistent State 是跨页面重载或会话保留的客户端状态。持久化不是给所有 state 加 localStorage；它需要明确存储介质、序列化格式、版本迁移、过期、并发和隐私边界。

```mermaid
flowchart LR
    A0["内存状态"] --> A1
    A1["白名单投影"] --> A2
    A2["序列化与版本"] --> A3
    A3["存储介质"] --> A4
    A4["校验迁移"] --> A5
    A5["恢复内存"]
```

Persistent State 是经白名单投影、版本化和运行时校验后可跨重载恢复的客户端快照。存储介质只是最后一步，迁移、期限、冲突和隐私决定它能否安全恢复。

## 2. 关键机制

### 2.1 介质选择

localStorage 同步且小型，IndexedDB 异步并适合结构化大数据，cookie 随请求发送。

若边界缺失，大对象频繁写 localStorage 阻塞主线程。

验证：性能记录写入耗时和配额错误。

### 2.2 白名单

只投影恢复必要字段，不序列化整个 store。

若边界缺失，令牌、错误对象和临时 UI 被保存。

验证：扫描存储内容和字段清单。

### 2.3 schemaVersion

每个持久化文档含版本，迁移函数逐级转换。

若边界缺失，发布后旧数据使应用崩溃。

验证：用各历史 fixture 启动。

### 2.4 运行时校验

JSON.parse 只保证语法，不保证形状和取值。

若边界缺失，类型断言接受损坏数据。

验证：畸形 fixture 回退默认。

### 2.5 过期

expiresAt 按业务期限失效，使用绝对时间并考虑时钟偏差。

若边界缺失，永久草稿使用过期规则。

验证：虚拟时钟测试边界。

### 2.6 原子性

写临时记录后切换指针，或利用 IndexedDB transaction。

若边界缺失，写一半留下不可解析文档。

验证：中断写入后恢复旧快照。

### 2.7 多标签页

storage/BroadcastChannel 通知变化，需定义冲突策略。

若边界缺失，两个标签页最后写覆盖草稿。

验证：版本号或更新时间冲突提示。

### 2.8 序列化

Date、Map、BigInt、File、函数不能靠普通 JSON 无损恢复。

若边界缺失，恢复后 Date 变字符串导致比较错误。

验证：往返测试每个字段类型。

### 2.9 加密边界

客户端加密密钥若同在页面，不能抵抗 XSS；不要持久化高敏秘密。

若边界缺失，把 base64 当加密。

验证：威胁模型审查攻击者能力。

### 2.10 清理

登出、租户切换、过期和版本废弃都触发删除。

若边界缺失，共享设备暴露前一用户数据。

验证：端到端测试登出后存储为空。

## 3. 持久化文档而不是直接存 Store

文档至少包含 schemaVersion、savedAt、可选 ownerId 和 value。写入前从内存状态投影允许字段；读取先 JSON 解析，再做运行时 schema 校验，然后逐版本迁移。任何一步失败都回退安全默认并删除或隔离坏记录。TypeScript 类型断言不会验证磁盘中的旧 JSON。

## 4. 保存、恢复和失效顺序

1. 启动先确定当前用户与租户命名空间，禁止读取前一主体的 key。

2. 读取文档并限制尺寸，解析、校验版本和 expiresAt。

3. 迁移函数逐级执行且保持纯净，迁移失败不写回半成品。

4. 恢复内存后才建立跨标签订阅，避免初始化消息覆盖本地快照。

5. 登出、过期、未知版本和租户切换都清理记录与广播状态。

## 5. 应用案例一：编辑器草稿

1. 编辑器只投影 title、body、baseRevision，排除 File、AbortController 和上传 URL。

2. 空闲防抖保存到 IndexedDB transaction，保存状态有 pending/saved/failed。

3. 恢复时请求服务端 revision；相同则应用草稿，不同进入冲突界面。

4. 注入配额错误，编辑继续留在内存并显示未保存。

5. 关闭写入中页面后，重新打开只能得到完整旧文档或完整新文档。

结果：同版本可直接恢复，不同版本进入冲突选择。

失败分支：配额满时继续编辑并显示未保存状态，不能反复同步阻塞。

## 6. 应用案例二：界面偏好

1. 主题保存 system/light/dark 枚举而不是计算颜色。

2. v1 的 compact:boolean 迁移为 v2 density 枚举。

3. BroadcastChannel 消息带 revision 和 sourceTab，接收方不回播同一版本。

4. 两个标签并发修改时按产品规则提示或比较 revision，不能默认为最后写胜。

5. 未知主题值回退 system，并从 DOM class 白名单生成样式。

结果：重载前后稳定，跨标签页同步偏好。

失败分支：未知枚举回退 system，不把非法字符串注入 className。

## 7. TypeScript 核心实现

下面代码把编码、版本、迁移和过期策略封装成持久化协议。具体存储介质留在适配器中，以便用内存实现测试旧版本、损坏内容和容量错误。

```tsx
type PreferenceDoc = {
  schemaVersion: 2;
  savedAt: number;
  value: { theme: "system" | "light" | "dark"; density: "compact" | "comfortable" };
};

export function readPreferences(storage: Storage, now: number): PreferenceDoc["value"] {
  const fallback = { theme: "system", density: "comfortable" } as const;
  const raw = storage.getItem("preferences");
  if (!raw) return fallback;
  try {
    const doc: unknown = JSON.parse(raw);
    if (!doc || typeof doc !== "object") return fallback;
    const d = doc as Partial<PreferenceDoc>;
    if (d.schemaVersion !== 2 || !d.value || typeof d.savedAt !== "number") return fallback;
    if (now - d.savedAt > 180 * 24 * 60 * 60 * 1000) return fallback;
    return d.value;
  } catch {
    return fallback;
  }
}
```

从存储读取的数据不再受编译期类型保证，必须先解析、校验版本再迁移。测试要覆盖非法 JSON、未知版本、跨标签页写入和迁移中断，失败时安全回退而非带病使用。

## 8. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| 内存 | 只需当前页面 | 重载丢失但最安全简单 |
| localStorage | 小型、低频、同步偏好 | 阻塞、字符串、配额和 XSS 可读 |
| IndexedDB | 较大结构化数据与事务 | 异步 API、迁移和清理更复杂 |

只有刷新后仍有价值、可安全落盘且有明确过期规则的状态才应持久化。敏感凭据、可重新查询的远端缓存和瞬时交互状态不应因使用方便而进入 localStorage。

## 9. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 启动白屏 | JSON.parse 或字段访问未捕获 | 校验并回退 |
| 升级后值错 | 缺少版本迁移 | 历史 fixture 测试 |
| 输入卡顿 | 每键同步写 localStorage | 防抖或 IndexedDB |
| 共享设备泄密 | 登出未清理 | 按用户命名空间并删除 |
| 两标签覆盖 | 无冲突协议 | revision 比较和提示 |
| 配额错误循环 | 写失败仍高频重试 | 降级内存并通知 |
| 时间字段比较错 | Date 被序列化为字符串 | 显式 epoch 或 codec |
| XSS 读到 token | 高敏数据被持久化 | 使用 HttpOnly 会话并最小化存储 |

先读取原始存储值和版本号，再逐步运行解码、迁移、校验与过期判断，最后检查多标签页事件是否覆盖了更新版本。失败信号是启动崩溃、日期变字符串、旧字段静默丢失或敏感值可被脚本读取；用损坏夹具、迁移矩阵和安全扫描验证。

## 10. 性能、安全与运维边界

- 任何持久化 schema 都需版本、校验、迁移和安全默认。
- 写入捕获 QuotaExceededError 和不可用存储环境。
- 登录用户数据按用户和租户隔离命名空间。
- 登出清理存储、缓存和跨标签广播状态。
- 敏感令牌优先使用服务端会话与 HttpOnly cookie。
- 大数据使用 IndexedDB transaction 并做失败恢复。
- 后台标签页同步设置冲突和回环保护。
- 记录迁移失败率和配额错误，不记录正文内容。

生产验证至少记录一次正常路径和一次故障路径；对“Persistent State”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 11. 与其他架构模块集成

- Global State 只把白名单切片交给持久化适配器。
- Form State 的 dirty 草稿需要 revision 冲突处理。
- URL 显式值优先于本地偏好默认。
- Server State 缓存持久化要尊重用户隔离和 stale 规则。

集成时 store 仍是运行期所有者，持久化层只保存带版本的快照并在启动时恢复一次。服务端会话和查询缓存分别维护自己的真实性，不能用本地快照越过鉴权或覆盖新鲜远端数据。

## 12. 综合练习

实现版本化草稿存储，覆盖 v1→v2 迁移、过期、配额失败、双标签冲突和登出清理。

### 验收标准

- [ ] v1、v2、损坏和未来版本 fixture 均有结果。
- [ ] 配额失败不阻塞编辑且不无限重试。
- [ ] 双标签消息无回环并有冲突策略。
- [ ] 登出清理用户与租户命名空间。
- [ ] 存储扫描不含令牌、密码和文件句柄。

## 13. 存储故障矩阵

| 故障 | 恢复行为 | 证据 |
|---|---|---|
| JSON 被截断 | 捕获解析错误并回退默认 | 损坏 fixture 不白屏 |
| schemaVersion 过新 | 保留或隔离记录，不用旧代码猜字段 | 降级版本启动测试 |
| expiresAt 已过 | 删除并重新初始化 | 虚拟时钟边界测试 |
| QuotaExceededError | 降级内存，停止高频重试 | 写入次数与用户提示 |
| storage 被浏览器禁用 | 能力检测后继续运行 | 隐私模式测试 |
| 两标签同 revision | 进入冲突策略 | 并发自动化记录 |

客户端可读取的加密密钥也能被同源 XSS 使用，因此加密 localStorage 不能把访问令牌变成安全会话。高敏凭证优先由服务端会话和 HttpOnly cookie 管理。

## 14. IndexedDB 事务与升级

IndexedDB 的 schema 变化在 `versionchange` 事务中执行。升级回调不能等待任意网络请求；对象仓库与索引修改必须在升级事务存活期间完成。其他标签页持有旧连接会阻塞升级，因此应用应监听 `versionchange` 并关闭旧连接，同时向用户提示刷新。

草稿正文和元数据需要原子提交时放在同一 transaction。若正文成功而 revision 失败，恢复时无法判断草稿对应哪个服务端版本。测试在第二次 `put` 前主动 abort，重新读取应仍得到旧的完整文档。

大文档不要每个键击完整 JSON 序列化。记录一次 100KB、1MB 草稿的序列化和写入耗时，在空闲窗口防抖保存；页面隐藏时尝试最后一次短写，但不能承诺浏览器一定完成。UI 用最后成功 savedAt 表示保存证据，而不是在调用 `put` 时立即显示“已保存”。

## 15. 恢复结果的可观察状态

恢复不是同步布尔值。界面至少区分 `loading`、`none`、`restored`、`conflict` 和 `failed`：没有草稿可以直接开始；恢复成功显示保存时间；冲突让用户比较；读取失败允许忽略坏记录。

恢复期间禁止把空初始表单自动保存回存储，否则会在异步读取完成前覆盖旧草稿。先完成读取和迁移，再启用保存订阅。测试把读取延迟 500ms，确认期间没有 write 调用。

## 来源

- [HTML Living Standard：Web Storage](https://html.spec.whatwg.org/multipage/webstorage.html)（访问日期：2026-07-18）
- [MDN：Web Storage API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)（访问日期：2026-07-18）
- [Indexed Database API 3.0](https://www.w3.org/TR/IndexedDB/)（访问日期：2026-07-18）
- [MDN：BroadcastChannel](https://developer.mozilla.org/en-US/docs/Web/API/BroadcastChannel)（访问日期：2026-07-18）
