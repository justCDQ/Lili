---
title: 长任务的刷新、恢复与后台执行
stage: intermediate
direction: ai
topic: ai-ux
---

# 长任务的刷新、恢复与后台执行

长时间 AI 任务不能依赖一个浏览器标签页和一条持续连接存活。任务应在服务端持久化执行，页面通过 run ID 查询状态、补发事件并恢复 artifact。刷新、关闭页面、设备休眠和网络切换都不应自动创建第二个任务或丢失已经完成的步骤。

## 前置知识与边界

- [AI 任务状态机与界面状态](01-ai-task-state-machine.md)
- [流式断线、重复事件与不完整 JSON](03-disconnect-duplicate-events-partial-json.md)
- [Artifact 的差异、历史与恢复](05-artifact-diff-history-restore.md)
- 任务队列、幂等性和持久化基础。

本文讨论应用协议与产品行为，不规定具体队列技术。运行数小时或数天的任务还需要容量规划、供应商配额和运维响应。

## 长任务的定义

符合任一条件就应按长任务设计：

- 超过常见 HTTP、网关或函数超时。
- 包含多个模型、检索或工具步骤。
- 需要等待人工确认。
- 用户可能离开页面。
- 结果可在以后查看。
- 失败后应从 checkpoint 继续。
- 副作用必须可审计和去重。

“预计只需 20 秒”也可能因排队、限流和工具故障变长。可靠性不能建立在平均耗时上。

## 三层状态

### Run

用户可理解的一次目标：

```json
{
  "runId": "run_501",
  "type": "analyze_repository",
  "status": "running",
  "ownerId": "user_9",
  "tenantId": "tenant_2",
  "inputSnapshotId": "input_44",
  "currentStep": "scan_files",
  "progress": {
    "completedUnits": 83,
    "totalUnits": 214,
    "unit": "files"
  },
  "createdAt": "2026-07-17T08:00:00Z",
  "updatedAt": "2026-07-17T08:04:18Z",
  "version": 27
}
```

### Step

可重试或可恢复的工作单元：

```json
{
  "stepId": "step_scan_3",
  "runId": "run_501",
  "name": "scan_files",
  "status": "running",
  "attempt": 1,
  "idempotencyKey": "run_501:scan_files:batch_3",
  "checkpointRef": "checkpoint://run_501/scan/83",
  "leaseExpiresAt": "2026-07-17T08:05:00Z"
}
```

### Artifact

已经持久化的阶段或最终结果：

```json
{
  "artifactId": "art_88",
  "runId": "run_501",
  "kind": "repository_report",
  "status": "provisional",
  "version": 4,
  "contentHash": "sha256:..."
}
```

页面刷新恢复 run，不直接恢复某个 worker 进程。

## 创建任务

```http
POST /api/runs
Idempotency-Key: 1ac3c721-5315-4b76-a05f-f84a8eea9832
Content-Type: application/json

{
  "type": "analyze_repository",
  "inputSnapshotId": "input_44"
}
```

响应：

```http
HTTP/1.1 202 Accepted
Location: /api/runs/run_501
Content-Type: application/json

{
  "runId": "run_501",
  "status": "queued",
  "statusUrl": "/api/runs/run_501",
  "eventsUrl": "/api/runs/run_501/events"
}
```

`202 Accepted` 表示请求已接受处理，不表示成功完成。客户端保存 run ID 后才跳转到任务页面。

## Run URL

每个任务有稳定 URL：

```text
/projects/project_2/runs/run_501
```

刷新时：

1. 从路由读取 run ID。
2. 重新进行身份与项目授权。
3. 查询 run 快照。
4. 加载当前 artifact 与历史事件 checkpoint。
5. 建立事件连接。
6. 把 UI 映射到权威状态。

不要把完整 Prompt 或敏感参数塞入 URL。

## 服务端后台执行

请求处理器只做：

1. 认证与授权。
2. 输入 Schema 与配额检查。
3. 创建 input snapshot。
4. 用幂等键创建 run。
5. 入队初始 step。
6. 返回 run ID。

worker 负责实际执行。浏览器断开不取消 worker，除非用户发送明确取消命令。

## Lease 与心跳

worker 领取 step 时使用有限 lease：

```text
queued -> leased -> running -> completed
                   ├-> failed
                   └-> lease_expired -> queued
```

lease 过期允许其他 worker 恢复，但会产生至少一次执行。所有 step 必须：

- 使用幂等键。
- 在 checkpoint 后提交状态。
- 对外部副作用查询原操作。
- 防止旧 worker 在 lease 丢失后覆盖新结果。

可以用 fencing token：

```json
{
  "stepId": "step_scan_3",
  "leaseToken": 18,
  "workerId": "worker_7"
}
```

写结果时只有当前 token 允许提交。

## Checkpoint

checkpoint 是恢复所需的确定状态，不是模型自由文本“我做到了第 83 个文件”：

```json
{
  "runId": "run_501",
  "step": "scan_files",
  "inputSnapshotHash": "sha256:...",
  "completedFileIds": ["file_1", "file_2"],
  "nextCursor": "cursor_84",
  "partialArtifactVersion": 4,
  "toolOperations": {
    "index_batch_3": "completed"
  },
  "createdAt": "2026-07-17T08:04:18Z"
}
```

checkpoint 需要原子关联：

- 已完成单元。
- 对应 artifact version。
- 外部操作状态。
- 输入版本。
- 代码或 workflow 版本。

否则恢复可能重复或漏处理。

## 进度

### 可确定进度

当总数已知，显示：

```text
已处理 83 / 214 个文件
```

这不等于完成时间 39%。不同文件成本可能差异很大。

### 不确定进度

模型推理、搜索深度或 agent 步数未知时，显示阶段：

```text
正在检索资料
正在比较证据
正在生成报告
正在运行验证
```

不要让进度条按时间自动爬到 95%。

### 重新估计

发现更多文件时 total 从 214 变为 260，界面应说明“发现新的处理项”，而不是让百分比倒退却不解释。

## 页面刷新恢复

客户端持久化最少状态：

```json
{
  "runId": "run_501",
  "lastAppliedSequence": 112,
  "viewPreferences": {
    "activeTab": "progress"
  }
}
```

不应把权威 task status、完整输出或权限保存在 localStorage 并直接信任。它们从服务端恢复。

### 恢复 reducer

```javascript
export async function hydrateRun(api, runId, checkpoint) {
  const snapshot = await api.getRun(runId);
  const artifact = snapshot.artifactId
    ? await api.getArtifact(snapshot.artifactId, snapshot.artifactVersion)
    : null;
  const events = await api.getEvents(runId, {
    after: checkpoint?.lastAppliedSequence ?? 0
  });

  return {
    snapshot,
    artifact,
    events
  };
}
```

若快照版本晚于事件历史，使用服务端定义的 snapshot sequence，避免把旧事件重复应用到新 artifact。

## 多标签页

同一 run 在两个标签页打开时：

- 两边都可只读观察。
- 取消、确认、编辑等命令仍由服务端串行化。
- 一个标签确认后，另一标签收到状态变化。
- 旧确认请求返回 already resolved。
- artifact 编辑使用 base version。

不要仅用浏览器内全局变量防重复，它无法覆盖其他设备。

## 通知

用户离开页面后可以通过站内通知、邮件或系统通知告知完成，但需要：

- 用户明确的通知偏好。
- 不在通知正文泄露敏感结果。
- 链接打开后重新授权。
- 通知去重。
- completed、failed、needs confirmation 使用不同文案。
- 不承诺浏览器推送必达。

通知只是提示，run 状态仍在服务端。

## 取消长任务

取消流程：

1. 用户发送取消命令。
2. run 进入 `cancelling`。
3. 停止创建新 step。
4. 向当前 worker 发送取消信号。
5. 等待不可中断步骤到安全边界。
6. 标记未完成 artifact。
7. run 进入 cancelled 或 cancellation_failed。

工具副作用可能无法撤销。界面展示：

- 已完成步骤。
- 已取消步骤。
- 仍在进行或无法撤销的动作。
- 是否提供补偿。

## 暂停与恢复

暂停不是取消：

- 暂停保留 checkpoint，停止调度新 step。
- 运行中的 step 在安全边界结束。
- 恢复检查输入、权限、workflow 和依赖是否仍有效。
- 长时间暂停后重新确认外部状态。

如果系统没有真正的暂停能力，不应把“断开页面”称为暂停。

## 版本升级

长任务运行期间 workflow 代码可能部署新版本。run 记录：

```json
{
  "workflow": {
    "name": "repository_analysis",
    "version": "2026-07-17.3"
  },
  "modelConfigVersion": "cfg_28",
  "promptVersion": "prompt_14"
}
```

恢复策略：

- 继续旧版本 worker，直到完成。
- 在兼容迁移后升级。
- 明确失败并要求从新版本重启。

不能用新代码任意解释旧 checkpoint。

## 失败分类

### Step 可重试

临时网络、限流、worker 崩溃。恢复同一 step，受预算与幂等保护。

### Run 需要用户输入

附件权限失效、需要确认参数。进入 waiting，不消耗 worker。

### Run 永久失败

输入损坏、workflow 不兼容、策略禁止。显示具体修复方向。

### 状态未知

外部工具超时后不知道是否完成。先查询 operation ID，不能直接重试。

## 完整案例一：大型仓库分析

### 任务

分析 20,000 个文件，建立索引，生成架构报告并运行引用校验。

### 步骤

```text
snapshot_repository
enumerate_files
parse_batches
build_index
retrieve_architecture_evidence
generate_report
validate_citations
publish_artifact
```

每个 parse batch 具有稳定 batch ID 和 checkpoint。

### 刷新

用户在 `parse_batches` 关闭页面。worker 继续。两小时后打开 run URL：

1. 查询 run 为 `generate_report`。
2. 加载进度快照：20,000/20,000 parsed。
3. 加载 provisional report v2。
4. 从 snapshot sequence 之后订阅。
5. 页面不重新上传仓库、不创建新 run。

### Worker 故障

worker 在 batch 71 写 artifact 后、提交 step 状态前崩溃。新 worker 读取 checkpoint，发现 batch operation ID 已完成且 artifact 包含 batch hash，因此只补提交状态，不重复索引。

### 输入变化

分析基于 commit `abc123`。仓库后来更新不影响当前报告的可重现性；界面显示输入 commit，并提供“对新 commit 重新分析”，这会创建关联新 run。

### 验证

- 页面关闭不改变 run。
- 同一创建幂等键返回同一 run。
- worker 崩溃不重复批次。
- 报告引用都绑定 commit。
- current repository 与输入 snapshot 区分。

## 完整案例二：批量图片处理与人工确认

### 任务

处理 500 张商品图，生成替代文本，低质量项进入人工审核，最后批量发布。

### 状态

- 生成可并行。
- 每张图有独立 item status。
- 人工审核期间 run 是 waiting_for_review。
- 发布是单独高风险批次。

### 页面恢复

刷新后读取：

```json
{
  "generated": 500,
  "passed": 462,
  "needsReview": 38,
  "published": 0
}
```

界面只加载当前分页的 38 个审核项，不下载所有图像。

### 发布确认

人员完成审核后，系统生成发布 confirmation，绑定：

- 500 个 item version 的 manifest hash。
- 目标 catalog revision。
- 发布渠道。

若确认后又编辑一个替代文本，manifest hash 改变，旧确认失效。

### 部分失败

发布 500 项时 493 成功、7 失败。run 状态为 partially_completed，不显示“发布完成”。重试只针对失败 item，使用原 item idempotency key 查询结果。

### 验证

- 关闭审核页面不丢 ownership。
- 过期 lease 允许其他审核员接管前先释放所有权。
- 通知不包含商品敏感信息。
- 重试不重复发布 493 项。
- 完成统计可追溯到 item 状态。

## 数据保留

长任务会产生大量：

- 事件。
- checkpoint。
- 临时 artifact。
- 工具日志。
- 模型输入输出。

定义：

- 活动任务保留期。
- 终态事件压缩。
- 临时 blob 删除。
- 审计记录最小字段。
- 用户删除传播。
- 备份恢复后 tombstone 重放。

删除 run 记录不能让仍在队列的 worker 继续访问已撤销数据。

## 安全与隔离

- run ID 使用不可预测值，但不可预测性不替代授权。
- 每次 status、event、artifact 请求重新检查租户和项目。
- worker 使用最小权限短期凭证。
- checkpoint 不保存 Secret。
- 队列消息携带资源 ID，不携带完整敏感正文。
- 输入 snapshot 的访问权限撤销时，运行按策略停止或转审核。
- 下载 artifact 使用短期受控 URL。

## 可观测性

指标：

- queue latency、run duration 和各 step duration。
- active、waiting、stalled、cancelling 数。
- heartbeat 丢失与 lease 过期。
- step 重试和重复副作用拦截。
- checkpoint 恢复成功率。
- 页面刷新恢复成功率和事件补发量。
- 长时间无进展 run。
- 用户取消延迟。
- 通知发送与打开，不把打开等同于用户已知晓。
- partial completion 与人工队列等待。

stalled 检测应基于预期心跳、阶段和截止时间，不只看总耗时。

## 常见错误

### 浏览器连接承担任务生命周期

刷新就失败。创建持久化 run 并交给 worker。

### 页面加载时重新 POST

刷新创建重复任务。路由携带 run ID；创建请求使用幂等键。

### 只保存百分比

无法恢复实际工作。保存已完成 item ID、cursor、artifact version 和 operation ID。

### Worker lease 过期后仍写结果

使用 fencing token 拒绝旧 worker。

### 通知显示敏感输出

只提示状态并链接回受授权页面。

### 失败后从头开始

按 checkpoint 和幂等 step 恢复；若必须重启，明确原因与额外成本。

### 所有 item 成功前显示 complete

支持 partially completed，并列出失败集合与安全重试。

## 生产验收清单

- [ ] 长任务有持久化 run ID 和稳定 URL。
- [ ] 创建请求有幂等键并返回 202/status URL。
- [ ] 浏览器关闭不取消后台任务。
- [ ] run、step 和 artifact 独立建模。
- [ ] worker 使用 lease、heartbeat 和 fencing token。
- [ ] checkpoint 包含输入、进度、artifact 与工具状态。
- [ ] 页面刷新从服务端权威快照恢复。
- [ ] 多标签确认和取消由服务端串行化。
- [ ] 进度只展示真实阶段或可测单位。
- [ ] 取消等待安全边界并展示未撤销副作用。
- [ ] workflow 与 Prompt 版本随 run 保存。
- [ ] 自动重试受预算、幂等与错误类型限制。
- [ ] 部分成功不是 completed。
- [ ] 通知不泄露内容且打开后重新授权。
- [ ] 数据保留和用户删除覆盖临时对象与队列。

## 集成练习

实现一个可处理 1,000 个文档的长任务：

1. POST 创建接口使用幂等键，返回 run URL。
2. 文档按 20 个一批，每批有 step ID、lease 和 checkpoint。
3. 杀死 worker 后由新 worker 从 checkpoint 恢复。
4. 刷新页面只查询旧 run，并从 last sequence 补事件。
5. 两个标签同时点击取消，只产生一个取消流程。
6. 生成最终报告前验证全部 batch 与输入 snapshot hash。
7. 注入 10 个 item 失败，run 显示 partially completed 并只重试失败集合。
8. 记录 queue、step、恢复、取消与 artifact 发布的可观测数据。

## 来源

- [IETF RFC 9110：HTTP Semantics，202 Accepted 与条件请求](https://www.rfc-editor.org/rfc/rfc9110)（访问日期：2026-07-17）
- [WHATWG HTML：Server-sent events 与 Last-Event-ID](https://html.spec.whatwg.org/multipage/server-sent-events.html)（访问日期：2026-07-17）
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)（访问日期：2026-07-17）
- [OpenTelemetry：Semantic conventions for generative AI systems](https://opentelemetry.io/docs/specs/semconv/gen-ai/)（访问日期：2026-07-17）
- [PostgreSQL：Explicit Locking](https://www.postgresql.org/docs/current/explicit-locking.html)（访问日期：2026-07-17）
