---
type: ai-note
stage: beginner
topic: request-metadata-usage
verified: 2026-07-16
tags: [ai, model, tokens, latency, cost, observability]
---

# 模型标识、输入输出、Token、延迟与费用记录

## 是什么

一次模型实验的最小记录包括：供应商、API/SDK 版本、模型完整标识、请求参数、Prompt/Schema 版本、输入、结构化输出、状态、Token Usage、延迟、费用和请求 ID。它们共同描述“运行了什么”，只保存最终文本无法复现或比较。

## 为什么需要

- 模型别名可能指向更新版本，行为会随时间变化。
- 输入、输出和缓存 Token 的计费可能不同，费用不能只按字符估算。
- 首 Token 延迟影响交互反馈，总时长影响任务完成和资源占用。
- 质量回归需要知道变更来自模型、Prompt、数据、参数还是应用代码。

## 关键字段

```text
timestamp / environment / feature
provider / endpoint / sdk_version / model_requested / model_returned
prompt_version / schema_version / dataset_version
parameters / tool_versions / retrieval_version
request_id / status / stop_reason / error_category
input_tokens / cached_tokens / output_tokens / total_tokens
time_to_first_event / total_latency / queue_or_tool_time
estimated_cost / currency / pricing_version
quality_result / safety_result
```

输入输出可能包含个人数据、商业信息或 Secret。详细 Trace 的访问、保留和脱敏策略必须先于采集设计。

## 实际怎么使用

用单调时钟测量持续时间，避免系统时间调整影响。费用根据供应商当日价格和实际 Usage 计算，同时保存 `pricing_version` 或价格表日期。

```ts
const started = performance.now();
const response = await client.generate(request);
const latencyMs = performance.now() - started;

record({
  modelRequested: request.model,
  modelReturned: response.model,
  requestId: response.requestId,
  usage: response.usage,
  latencyMs,
  promptVersion: "extract-v3",
  schemaVersion: "metadata-2",
});
```

生产看板至少按功能、模型、版本、用户/租户（使用不可逆标识）、状态和分位数聚合。不要只看平均延迟；P50、P95、P99 能显示尾部问题。

## 常见错误与边界

- 只记录请求使用的别名，不记录响应返回的实际模型标识。
- Streaming 只测总时长，不测首事件或首可用内容时间。
- 流中断时假设已收到最终 Usage；部分接口只有完成事件包含完整 Usage。
- 将估算 Token 当作供应商结算 Usage，或忽略缓存、推理、图片/音频等计费项。
- 把完整 Prompt、文档和模型输出写入普通日志，扩大敏感数据范围。

## 补充知识

质量、延迟和成本必须放在同一实验表中。只优化其中一个可能使另外两个退化。供应商价格、模型生命周期和数据保留政策具有时效性，应链接官方页面并定期复查。

## 来源

- [OpenAI API：Responses 与 Usage](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [OpenAI API：Streaming Events](https://platform.openai.com/docs/api-reference/responses-streaming)（访问日期：2026-07-16）
- [Anthropic API：Messages](https://docs.anthropic.com/en/api/messages)（访问日期：2026-07-16）
- [OpenTelemetry：Semantic Conventions for Generative AI](https://opentelemetry.io/docs/specs/semconv/gen-ai/)（访问日期：2026-07-16）

