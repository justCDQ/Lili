---
type: ai-note
stage: beginner
topic: sdk-vs-http
verified: 2026-07-16
tags: [ai, sdk, http, api]
---

# SDK 与原始 HTTP 模型调用

## 是什么

原始 HTTP 调用由应用直接构造 URL、Header 和 JSON 请求体，再解析状态码与响应。SDK 是供应商或社区提供的语言库，将这些协议细节封装为类型、函数、迭代器和错误对象。两者访问同一远程能力，但抽象层和升级责任不同。

## 为什么两种都要理解

- SDK 适合业务开发，减少重复的认证、序列化和流解析代码。
- HTTP 是最终接口事实来源；新能力可能先出现在 API 参考中，SDK 可能尚未支持。
- 调试代理、网关、超时、错误体和兼容性时需要查看真实请求与响应。
- 迁移供应商或构建统一 Client 时，必须区分通用语义和 SDK 专有便利方法。

## 关键特性

### SDK

- 提供请求/响应类型、自动 JSON 编解码和流式事件迭代。
- 可能自动读取环境变量、设置重试或附加遥测；这些默认行为必须阅读文档确认。
- SDK 版本与 API/模型版本不是同一概念。升级 SDK 不会自动固定模型行为。
- 便利字段可能丢失原始事件细节，必要时保留响应 ID、Header 和原始 Usage。

### HTTP

- 请求由方法、URL、Header 和内容组成，响应由状态码、Header 和内容组成。
- 可用 `curl` 形成最小可复现案例，排除应用框架干扰。
- 必须自行处理认证、超时、错误分类、重试、流分片、Schema 和日志脱敏。

## 实际怎么使用

学习顺序：

1. 按官方 Quickstart 使用官方 SDK 完成一次非流式请求。
2. 记录 SDK 实际发送的模型、请求字段和返回结构。
3. 用 `curl` 按 API Reference 重做同一请求。
4. 对比文本、结束状态、Usage、请求 ID 和错误结构。
5. 在业务中封装自己的 `ModelClient`，只暴露业务需要的统一能力。

```ts
interface ModelClient {
  generate(request: GenerateRequest): Promise<GenerateResult>;
  stream(request: GenerateRequest, signal?: AbortSignal): AsyncIterable<ModelEvent>;
}
```

统一接口不应假装所有供应商完全相同。无法稳定映射的推理参数、缓存、内置工具和数据保留选项应作为显式能力声明或供应商扩展。

## 常见错误与边界

- 只复制 Quickstart，不读取模型生命周期、错误、限流和数据控制文档。
- 在每个业务函数直接创建 SDK Client，无法集中配置超时、日志和权限。
- 把 SDK 自动重试与应用重试叠加，实际请求次数超过预期。
- 用 HTTP 成功状态代替模型任务成功；仍需检查拒绝、不完整、Tool Call 和 Schema。
- 自建抽象只保留文本，丢失引用、Usage、结束原因和工具事件。

## 补充知识

生产排障应保存供应商请求 ID，并允许在脱敏条件下查看统一结构和原始错误。契约测试可以在 SDK 升级时检查统一 Client 的字段映射是否变化。

## 来源

- [OpenAI Developer Quickstart](https://platform.openai.com/docs/quickstart)（访问日期：2026-07-16）
- [OpenAI API Reference](https://platform.openai.com/docs/api-reference)（访问日期：2026-07-16）
- [Anthropic API：Getting Started](https://docs.anthropic.com/en/api/getting-started)（访问日期：2026-07-16）
- [RFC 9110：HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-16）

