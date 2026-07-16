---
type: ai-note
stage: beginner
topic: streaming-cancellation-timeout-errors
verified: 2026-07-16
tags: [ai, streaming, cancellation, timeout, errors]
---

# Streaming、取消、超时与错误展示

## 是什么

Streaming 在模型生成过程中持续发送事件或增量。取消由客户端或服务端停止仍在进行的工作。超时是应用允许某一阶段或整个任务占用的最长时间。错误展示把系统失败转换为用户可理解、可恢复且不泄露内部信息的状态。

## 为什么需要

模型生成可能持续数秒到数分钟。流式输出缩短首个可见结果时间，但增加事件顺序、部分结果和中断处理复杂度。没有取消和超时，断开的页面、后台 Worker 或 Agent 仍可能继续消耗资源和费用。

## 关键特性

### Streaming 状态

至少区分 `connecting`、`streaming`、`completed`、`cancelled`、`failed` 和 `incomplete`。事件可能包含文本、Tool 参数、引用、Usage、错误和完成状态，不能把所有增量直接拼成字符串。

### 取消

取消信号要向下传播到 HTTP 请求、模型 SDK、Tool、数据库查询和 Worker。客户端断开不保证供应商端立即停止计费，必须以具体 API 文档为准。

### 超时

区分连接超时、首事件超时、流空闲超时、单个 Tool 超时和总任务超时。仅设置一个总超时难以定位问题，也可能错误终止仍在稳定输出的长任务。

### 错误展示

用户需要知道发生了什么、已完成内容是否保留、能否重试、重试会不会重复写操作。内部日志保留请求 ID 和错误类别，UI 不显示堆栈、Secret 或供应商原始敏感信息。

## 实际怎么使用

```ts
const controller = new AbortController();
const timeout = setTimeout(() => controller.abort("total timeout"), 60_000);

try {
  for await (const event of client.stream(request, controller.signal)) {
    applyEvent(state, event); // 按事件类型和序号处理
  }
  assertCompleted(state);
} catch (error) {
  showRecoverableState(classifyError(error));
} finally {
  clearTimeout(timeout);
}
```

保存策略：文本草稿可以保留并标记“不完整”；结构化结果只有完成并通过 Schema 校验后才能进入业务数据库；写 Tool 需要幂等键和明确确认。

## 常见错误与边界

- 把收到最后一个网络 Chunk 当作模型成功完成，没有检查完成状态。
- 取消 UI 更新但没有取消后端请求，费用仍继续产生。
- 对半截 JSON 直接 `JSON.parse`；结构化流必须按协议等待完整事件或使用增量解析器。
- Stream 失败后自动重试写操作，产生重复副作用。
- 显示“网络错误”覆盖限流、权限、内容拒绝、超时和服务端故障，用户无法恢复。

## 补充知识

真实链路可能经过 CDN、反向代理和应用网关，它们可能缓冲或关闭长连接。部署前应验证心跳、空闲超时、压缩、代理缓冲和移动网络切换。

## 来源

- [OpenAI API：Streaming Events](https://platform.openai.com/docs/api-reference/responses-streaming)（访问日期：2026-07-16）
- [OpenAI Realtime API：response.cancel](https://platform.openai.com/docs/api-reference/realtime-client-events/response/cancel)（访问日期：2026-07-16）
- [MDN：AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)（访问日期：2026-07-16）
- [MDN：Server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)（访问日期：2026-07-16）

