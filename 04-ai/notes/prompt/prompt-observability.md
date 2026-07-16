---
type: ai-note
stage: junior
topic: prompt-observability
verified: 2026-07-16
tags: [ai, prompt, observability, tokens, latency, cost]
---

# Prompt 的模型、参数、Token、延迟与成本记录

## 是什么

Prompt 可观测性把一次调用与模型、Prompt 版本、参数、上下文组成、Token Usage、延迟、费用、状态和质量反馈关联。它服务于回归分析、容量、成本和故障排查。

## 为什么需要

同一 Prompt 在不同模型、参数和上下文下行为不同。只记录最终文本无法判断退化来自 Prompt、模型别名、检索、输入增长还是供应商异常。

## 关键特性与规则

- 一次任务需要稳定 trace 标识，并关联模型调用、检索、Tool、Grader 和最终状态。
- 模型请求标识与实际返回标识应分别记录，避免别名更新无法追踪。
- Token、延迟和费用必须带计量口径与价格日期，历史报表不能用新价格无说明回算。
- 日志默认最小化敏感内容；Prompt、输出和文件内容使用分级权限、脱敏和保留期。
- Streaming 同时记录首事件延迟、总延迟、完成状态和 Usage 可用性。

## 实际怎么使用

每次调用记录：

```text
trace_id / request_id / feature / tenant_hash
provider / model_requested / model_returned
prompt_id / prompt_version / schema_version
temperature / max_output / tool_set / retrieval_version
input_tokens / cached_tokens / output_tokens
first_event_ms / total_ms / tool_ms
status / stop_reason / error_category
estimated_cost / pricing_date
quality_label / user_feedback
```

敏感 Prompt 内容和用户数据单独按权限保存，普通指标日志只保留版本、哈希、长度和分类。使用 P50/P95/P99 分位数，不只看平均值。

## 常见错误与边界

- 把完整 Prompt 和 Secret 写入全员可查日志。
- 记录模型别名但不记录响应实际模型。
- 费用按当前价格回算旧请求，历史报表变化。
- Streaming 不记录首事件延迟和不完整状态。
- 只监控成本，不把成本与质量和任务完成率关联。

## 补充知识

Trace 应覆盖检索、模型、Tool 和 Grader，但采样率、保留期和脱敏必须符合数据政策。OpenTelemetry 的生成式 AI 语义约定仍会演进，使用时锁定版本。

## 来源

- [OpenTelemetry：Generative AI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/)（访问日期：2026-07-16）
- [OpenAI API：Responses Usage](https://platform.openai.com/docs/api-reference/responses)（访问日期：2026-07-16）
- [Anthropic API：Messages](https://docs.anthropic.com/en/api/messages)（访问日期：2026-07-16）
