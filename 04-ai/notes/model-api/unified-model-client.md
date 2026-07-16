---
type: ai-note
stage: junior
topic: unified-model-client
verified: 2026-07-16
tags: [ai, model-client, architecture, vendor-isolation]
---

# 统一模型 Client 与厂商隔离

## 是什么

统一模型 Client 是应用内部定义的接口，把业务需要的生成、流式、结构化、工具和 Usage 能力与具体供应商 SDK 分离。每个供应商 Adapter 负责请求转换、事件映射、错误分类和能力声明。

## 为什么需要

供应商在字段、角色、工具、错误、Usage、模型生命周期和数据控制上存在差异。业务代码直接依赖 SDK 会让切换、测试、监控和安全策略分散。统一接口集中共性，但不能消除真实差异。

## 关键特性

```ts
interface Capabilities {
  streaming: boolean;
  structuredOutput: "none" | "json" | "json-schema";
  toolCalling: boolean;
  inputModalities: Array<"text" | "image" | "audio">;
}

interface ModelClient {
  capabilities(model: string): Capabilities;
  generate(request: GenerateRequest): Promise<GenerateResult>;
  stream(request: GenerateRequest, signal?: AbortSignal): AsyncIterable<ModelEvent>;
}
```

- 统一结构保留模型、状态、结束原因、Usage、请求 ID、引用和工具事件。
- 错误至少映射为认证、权限、限流、超时、暂时服务、无效请求、内容拒绝、结构失败和未知。
- 供应商独有能力通过显式扩展或能力查询暴露，不用含义模糊的通用字段硬凑。
- 日志、超时、重试、预算和脱敏在统一层实施，业务规则仍在领域层。

## 实际怎么使用

1. 从真实业务用例定义最小接口，不从所有 SDK 字段求并集。
2. 为每个 Adapter 编写契约测试：文本、流、Schema、Tool、错误和 Usage。
3. 模型配置放注册表，记录供应商、完整模型标识、能力和生命周期。
4. 路由根据任务、质量、延迟、成本、区域和数据政策选择模型。
5. Fallback 只在语义允许时使用，并记录实际使用模型；结构或工具不兼容时不得静默切换。

## 常见错误与边界

- 抽象只返回字符串，丢失 Tool、引用、状态和 Usage。
- 假设所有供应商的 `temperature`、System、缓存和多轮语义相同。
- 自动 Fallback 到能力更弱或数据政策不同的模型，用户和审计均不可见。
- 在统一层加入大量业务 Prompt 和领域逻辑，Adapter 无法复用。
- 没有能力检测，运行时才发现模型不支持 Schema 或模态。

## 补充知识

统一 Client 不是必须从第一天支持多供应商。即使只有一家供应商，边界仍便于测试和升级；避免在没有迁移需求时构建过度复杂的动态路由平台。

## 来源

- [OpenAI API Reference](https://platform.openai.com/docs/api-reference)（访问日期：2026-07-16）
- [Anthropic API Reference](https://docs.anthropic.com/en/api)（访问日期：2026-07-16）
- [OpenTelemetry：Generative AI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/)（访问日期：2026-07-16）

