---
title: 流式断线、重复事件与不完整 JSON
stage: intermediate
direction: ai
topic: ai-ux
---

# 流式断线、重复事件与不完整 JSON

流式 AI 界面必须假设连接会断开、事件会重复、事件可能缺失，结构化内容也可能停在任意字节。正确实现依赖持久化 run、单调序号、可重放事件、增量解码与终态校验。把每次网络读取直接拼进文本或 `JSON.parse`，无法满足恢复与数据完整性。

## 前置知识与边界

- [AI 任务状态机与界面状态](01-ai-task-state-machine.md)
- [流式输出的停止、重试与继续生成](02-stream-stop-retry-continue.md)
- UTF-8、HTTP、SSE、JSON 和幂等性的基础知识。

本文处理传输和解析层。业务工具仍必须独立做权限、Schema、事务与副作用检查。

## 先区分四个对象

| 对象 | 说明 | 典型标识 |
|---|---|---|
| run | 一次持久化任务 | `runId` |
| transport connection | 一次浏览器到服务端的连接 | `connectionId` |
| event | run 中可重放的状态变化 | `eventId`、`sequence` |
| content fragment | 某个输出块的增量 | `blockId`、`delta` |

一次 run 可以经历多个连接。连接断开不能自动创建新 run；同一事件也可能在重连后再次传输。

## 事件协议

```json
{
  "protocolVersion": 1,
  "runId": "run_91",
  "eventId": "01JZ8K5NXD0X2VQ7S6N5",
  "sequence": 42,
  "type": "content.delta",
  "createdAt": "2026-07-17T10:20:31.418Z",
  "payload": {
    "blockId": "answer",
    "format": "text",
    "text": "新的内容"
  }
}
```

字段规则：

- `protocolVersion`：支持兼容检查，不能把未知字段语义静默猜测。
- `runId`：防止旧连接事件污染当前任务。
- `eventId`：全局或 run 内唯一，用于审计和去重。
- `sequence`：同一 run 内严格递增，表示逻辑顺序。
- `type`：受控枚举。
- `createdAt`：服务端产生时间，不用于代替顺序。
- `payload`：按事件类型校验。

客户端本地到达时间不能决定逻辑顺序。

## SSE 的恢复能力与边界

SSE 事件可以携带 `id`：

```text
id: 42
event: content.delta
data: {"runId":"run_91","sequence":42,"text":"新的内容"}

```

浏览器的 `EventSource` 在重新连接时可以发送 `Last-Event-ID`。服务端仍需要：

- 把 ID 映射到当前用户有权访问的 run。
- 从持久化事件日志补发后续事件。
- 设置有限的事件保留期。
- 对过旧的恢复点返回明确错误。
- 处理客户端刷新后自行携带 checkpoint 的情况。

原生 `EventSource` 只支持有限请求配置；若需要自定义认证头、POST 请求或精细取消，可以使用 fetch stream，并由应用自己维护 checkpoint。

## 断线状态

断线后界面不应立刻进入 failed：

```json
{
  "taskStatus": "streaming",
  "transportStatus": "reconnecting",
  "lastAppliedSequence": 42,
  "reconnectAttempt": 2,
  "nextRetryAt": "2026-07-17T10:20:35Z"
}
```

用户应看到：

- 已收到内容仍保留。
- “连接中断，正在恢复”。
- 任务可能仍在服务端运行。
- 停止操作是否仍可发送。
- 恢复失败后的查询或手动重试入口。

## 重连算法

```javascript
export async function resumeRun({
  api,
  runId,
  lastAppliedSequence,
  signal
}) {
  const snapshot = await api.getRun(runId, { signal });

  if (snapshot.status === "expired") {
    return { kind: "unrecoverable", reason: "event_history_expired" };
  }

  const events = await api.listEvents(runId, {
    after: lastAppliedSequence,
    signal
  });

  return {
    kind: "events",
    authoritativeStatus: snapshot.status,
    events
  };
}
```

恢复顺序：

1. 验证当前身份和 run 权限。
2. 查询权威 run 状态。
3. 请求 `lastAppliedSequence` 之后的事件。
4. 验证第一条是否正好是下一序号，或服务端提供可信快照。
5. 按顺序重放。
6. 再建立实时连接。
7. 实时事件和补发事件经过同一 reducer 去重。

若先建立实时连接再补历史，需要缓冲实时事件，否则可能先应用 sequence 51，再收到 43–50。

## 指数退避与抖动

```javascript
export function reconnectDelay(attempt, {
  baseMs = 500,
  capMs = 15000,
  random = Math.random
} = {}) {
  const ceiling = Math.min(capMs, baseMs * 2 ** attempt);
  return Math.floor(random() * ceiling);
}
```

重连还要限制：

- 最大连续尝试次数。
- 最大总等待时间。
- 页面离线时暂停。
- 服务端 `Retry-After`。
- 页面隐藏时的资源策略。
- 用户主动停止后不再自动连接。

网络恢复事件可以触发立即查询，但不能造成多个并发重连循环。

## 重复事件

重复产生于：

- 服务端至少一次投递。
- 客户端未确认前断线。
- SSE 自动重连。
- 网关或应用重试。
- 前端组件重复订阅。
- 历史补发与实时流重叠。

处理函数必须幂等：

```javascript
export function applyEvent(state, event) {
  if (event.runId !== state.runId) {
    return { ...state, protocolError: "wrong_run" };
  }

  if (event.sequence <= state.lastAppliedSequence) {
    return state;
  }

  if (event.sequence !== state.lastAppliedSequence + 1) {
    return {
      ...state,
      transportStatus: "needs_resync",
      gap: {
        expected: state.lastAppliedSequence + 1,
        received: event.sequence
      }
    };
  }

  const next = reduceKnownEvent(state, event);
  return { ...next, lastAppliedSequence: event.sequence };
}
```

“看到相同 text 就去重”是错误方法。模型可能有意重复词语，两个相同文本也可能属于不同 block。去重依据是事件身份与序号。

## 缺失与乱序

### 缺失

当前 sequence 10，收到 12，说明 11 未应用。客户端应：

- 暂停把任务标为 completed。
- 缓冲 12 或丢弃后重新拉取。
- 请求 10 之后的事件。
- 若历史已过期，请求服务端快照。
- 无法恢复时显示输出可能不完整。

### 乱序

HTTP 单连接通常保持字节顺序，但多路内部系统、补发和实时流合并仍可能乱序。应用协议不能依赖“通常不会发生”。

小范围缓冲：

```javascript
export function createOrderedBuffer(startSequence, apply) {
  let next = startSequence;
  const pending = new Map();

  return event => {
    if (event.sequence < next) return;
    pending.set(event.sequence, event);

    while (pending.has(next)) {
      const current = pending.get(next);
      pending.delete(next);
      apply(current);
      next += 1;
    }
  };
}
```

缓冲必须有限制。若 pending 数量或等待时间超限，触发 resync，防止攻击者发送巨大序号耗尽内存。

## 终态事件

`run.completed` 必须携带或可查询：

```json
{
  "type": "run.completed",
  "sequence": 88,
  "payload": {
    "finalOutputVersion": 3,
    "finalOutputHash": "sha256:...",
    "lastContentSequence": 87,
    "validationStatus": "passed"
  }
}
```

客户端只有在这些条件满足时标记 complete：

- 1 到 `lastContentSequence` 没有缺口。
- 最终输出 hash 与服务端 artifact 一致，或直接使用服务端最终 artifact。
- 输出通过声明的验证。
- 当前 run 尚未处于其他终态。

完成事件不应只是文本 delta 结束时的连接关闭。

## UTF-8 增量解码

网络 chunk 可以在多字节字符中间分割。错误示例：

```javascript
// 每块单独解码可能破坏跨块字符。
const text = chunks.map(chunk => new TextDecoder().decode(chunk)).join("");
```

正确使用流式解码：

```javascript
const decoder = new TextDecoder("utf-8", { fatal: true });
let text = "";

for await (const chunk of byteStream) {
  text += decoder.decode(chunk, { stream: true });
}

text += decoder.decode();
```

`fatal: true` 会在非法 UTF-8 时抛错，适合协议层明确失败；若选择替换字符策略，也要记录数据损坏而不是静默当成模型输出。

## 不完整 JSON 的三种来源

### 传输 framing 未完成

SSE 的某个 `data:` 事件尚未以空行结束。解析器要保留缓冲，不能提前发出事件。

### JSON 值跨 chunk

```text
chunk 1: {"name":"北
chunk 2: 京","count":2}
```

先完成字节解码，再按 framing 组装完整 JSON 文本。

### 模型输出本身截断

即使网络事件完整，模型内容可能是：

```text
{"items":[{"id":1},{"id":
```

这是完整收到的“不完整业务 JSON”。只有整体终态和 Schema 校验才能区分。

## 不要在每个 delta 上 `JSON.parse`

```javascript
// 错误：任意 delta 通常都不是完整 JSON。
onDelta(delta => {
  const value = JSON.parse(delta);
});
```

如果供应商提供结构化参数 delta，应先按 item/block 累积：

```javascript
export function createJsonAccumulator() {
  let source = "";
  let terminal = false;

  return {
    append(delta) {
      if (terminal) throw new Error("append_after_terminal");
      source += delta;
      return { byteLength: new TextEncoder().encode(source).byteLength };
    },
    finish() {
      terminal = true;
      return JSON.parse(source);
    },
    snapshot() {
      return source;
    }
  };
}
```

这个累积器只解决完整结束后的 JSON 语法，不解决 Schema、权限和业务有效性。

## 使用结构化事件代替裸 JSON 文本

更可靠的协议按字段或节点传输：

```json
{
  "type": "object.field.completed",
  "sequence": 9,
  "payload": {
    "objectId": "task_1",
    "field": "title",
    "value": "完成发布说明"
  }
}
```

服务端验证字段后才发出 completed。前端可以逐字段展示，但最终提交仍等待 `object.completed` 和整体 Schema。

取舍：

- 节点事件更易恢复和验证。
- 协议更复杂，需要稳定 ID 和版本。
- 对纯展示文本，字符串 delta 更简单。
- 对工具参数、表单和数据库命令，结构化节点更安全。

## JSON Schema 与业务校验

语法有效不等于可执行：

```json
{
  "action": "refund",
  "orderId": "ord_8",
  "amountCents": -100
}
```

验证层次：

1. UTF-8 有效。
2. framing 完整。
3. JSON 语法有效。
4. Schema 类型、必填字段和枚举有效。
5. 字段长度、金额范围和规范化有效。
6. 当前用户权限有效。
7. 订单状态和可退余额等业务不变量有效。
8. 幂等键与事务提交有效。

模型或供应商的 structured output 可以提高 Schema 命中率，但不能替代后四层。

## 部分内容的存储

建议分开：

```json
{
  "provisionalBuffer": {
    "lastSequence": 44,
    "text": "正在生成的文本",
    "publishable": false
  },
  "finalArtifact": {
    "version": 3,
    "hash": "sha256:...",
    "validation": "passed",
    "publishable": true
  }
}
```

临时缓冲可以用于恢复与展示，但不能覆盖已发布 artifact。完成时使用原子写入创建新 artifact version。

## 完整案例一：带引用的研究回答

### 输入

服务端检索资料并产生段落、引用与完成事件。客户端在 sequence 27 后断线。

### 服务端已有事件

```text
25 paragraph.delta
26 paragraph.completed
27 citation.completed
28 paragraph.delta
29 paragraph.completed
30 run.completed
```

### 恢复

1. 客户端保留 `lastAppliedSequence=27`。
2. 重连查询 run，权威状态为 completed。
3. 请求 `after=27`，得到 28–30。
4. 验证第一个序号是 28。
5. 应用段落并验证引用 source ID。
6. 使用 final artifact hash 校验本地聚合。
7. 将 provisional 切换为 complete。

### 重复分支

实时连接同时又收到 29 和 30。reducer 看到 sequence 小于等于 30，忽略，不重复段落和终态通知。

### 失败分支

事件保留期已过，服务端无法补发 28–30，但仍有最终 artifact。客户端放弃本地临时聚合，下载已授权 final artifact，并把 last sequence 设置为快照 checkpoint。若 final artifact 也不存在，显示“无法恢复完整结果”，不能把 sequence 27 的内容伪装成完成。

### 验证

- 网络面板中断后只有一个 run。
- 重复补发不造成重复段落。
- 缺失引用使验证失败。
- final hash 不一致时重新下载 artifact。
- 不同租户无法用 run ID 获取事件。

## 完整案例二：流式表单与工具调用

### 输入

模型协助填写会议安排：

```json
{
  "title": "产品评审",
  "attendees": ["a@example.com"],
  "start": "2026-07-20T09:00:00+08:00",
  "durationMinutes": 60
}
```

### 展示策略

- 字段通过独立 completed 事件进入可编辑表单。
- 原始 JSON 字符串只用于调试，不直接渲染为最终值。
- 所有字段完成后做 Schema 和时区校验。
- 用户确认时绑定字段 hash。
- 日历工具只接收服务端生成的类型化命令。

### 断线

断线发生在 attendees 字段的字符串中间。未完成 framing 不产生字段事件。恢复后补发完整 `field.completed`，页面只出现一个参会者。

### 重复

确认请求因超时重发。工具调用使用 `calendar:create:run_44:confirmation_2` 幂等键。查询发现会议已创建，返回原 event ID，不创建第二个会议。

### 失败分支

模型生成了语法正确但不存在的本地时间，例如落在时区夏令时跳跃区间。Schema 无法发现，时区业务校验拒绝并要求用户选择有效时刻。

### 验证

- 任意字节边界拆分 UTF-8，最终文本保持一致。
- 任意位置断开 JSON，不执行工具。
- 重放所有事件两次，表单值不重复。
- 缺少一个字段事件时，确认按钮保持禁用。
- 编辑字段后旧确认 hash 无效。

## 测试方法

### 字节切分测试

对一条包含中文、emoji 和转义字符的事件，在每个字节位置切分，逐块送入解码器，最终结果必须与原文一致。

### 事件故障注入

固定序列执行：

- 删除一个事件。
- 重复每个事件。
- 交换相邻事件。
- 完成事件提前。
- 旧 run 事件混入。
- sequence 跳到极大值。

预期结果应是恢复、拒绝或明确不完整，不能静默产生不同内容。

### JSON 截断测试

在每个字符位置截断结构化输出。除完整位置外都不得触发业务工具。完整位置仍需 Schema 和权限检查。

## 可观测性

记录：

- 连接次数、断线原因和恢复耗时。
- 每个 run 的最后生产、发送、应用 sequence。
- 重复、缺口、乱序和错误 run 事件数。
- 事件补发数量和快照恢复次数。
- UTF-8、framing、JSON、Schema、业务校验各层错误。
- provisional 与 final hash 不一致次数。
- 自动恢复成功率与事件过期率。

高基数 run ID 不应直接成为普通 metrics label，可进入 trace 或日志字段。

## 常见错误

### 连接关闭等于完成

网关、浏览器或网络都能关闭连接。只有权威终态事件或任务查询能证明完成。

### 用时间戳排序

分布式时钟可能偏移，同毫秒也可产生多个事件。使用 run 内 sequence。

### 用内容 hash 去重 delta

重复文本可能是合法输出。使用 event ID 和 sequence。

### 恢复后创建新 run

先恢复旧 run；只有明确失败并选择重试时才创建关联新 run。

### 补历史与实时订阅并行直接写 UI

会乱序。先建立恢复边界，或对两路事件统一排序缓冲。

### 半截 JSON 自动补括号

补到语法有效不等于恢复原意，尤其不能用于工具。保留最后验证节点或重新生成。

## 生产验收清单

- [ ] run 与 transport connection 有独立 ID 和状态。
- [ ] 每个事件有 protocol version、run ID、event ID 和 sequence。
- [ ] 事件可在有限保留期内补发。
- [ ] 恢复前重新授权 run。
- [ ] 重复事件处理幂等。
- [ ] 缺失和乱序触发 resync，不静默跳过。
- [ ] 缓冲有数量、字节和时间限制。
- [ ] 完成事件指明最后内容序号和最终 artifact。
- [ ] UTF-8 使用流式解码。
- [ ] framing 完成后才解析事件 JSON。
- [ ] 模型 JSON 只在终态后做整体解析。
- [ ] Schema、权限和业务校验分层执行。
- [ ] provisional 不覆盖 final artifact。
- [ ] 工具调用使用幂等键并可查询结果。
- [ ] 故障注入覆盖断线、重复、缺失、乱序和截断。

## 集成练习

实现一个流式生成任务表单：

1. 服务端事件日志保留 24 小时，run 内 sequence 严格递增。
2. 客户端保存 last applied sequence，刷新后恢复。
3. 将 UTF-8 事件随机切为 1–7 字节 chunk 传输。
4. 注入重复、缺失和乱序，客户端不得生成错误终态。
5. 表单字段只有收到 field completed 事件后才可确认。
6. 不完整 JSON 在任何截断位置都不能触发工具。
7. 最终 artifact 具有 version、hash 和 Schema 结果。
8. 跨租户恢复请求返回统一拒绝且不泄露 run 是否存在。

## 来源

- [WHATWG HTML：Server-sent events、event ID 与 Last-Event-ID](https://html.spec.whatwg.org/multipage/server-sent-events.html)（访问日期：2026-07-17）
- [WHATWG Encoding：TextDecoder 与流式解码](https://encoding.spec.whatwg.org/#interface-textdecoder)（访问日期：2026-07-17）
- [WHATWG Streams Standard](https://streams.spec.whatwg.org/)（访问日期：2026-07-17）
- [IETF RFC 8259：The JavaScript Object Notation Data Interchange Format](https://www.rfc-editor.org/rfc/rfc8259)（访问日期：2026-07-17）
- [JSON Schema 2020-12 Specification](https://json-schema.org/specification)（访问日期：2026-07-17）
