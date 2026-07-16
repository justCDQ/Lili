---
type: ai-note
stage: beginner
topic: structured-output-validation
verified: 2026-07-16
tags: [ai, structured-output, json-schema, validation]
---

# Structured Output、Schema 与运行时校验

## 是什么

Structured Output 要求模型按预定义结构返回数据。Schema 描述允许的类型、字段、必填项和约束。运行时校验器在程序收到实际数据后检查其是否符合 Schema。三者解决格式和接口契约问题，不保证内容事实正确。

## 为什么需要

业务代码需要稳定字段来写数据库、调用 Tool、渲染 UI 和运行自动化。自然语言格式容易出现缺字段、字段名变化、类型错误和额外说明。Schema 将可接受结构显式化，运行时校验防止无效数据继续传播。

## 关键特性

- JSON Schema 需要声明使用的方言，例如 2020-12；供应商 Structured Output 通常只支持其中一个子集。
- `type`、`properties`、`required`、`items`、`enum` 和 `additionalProperties` 是常用约束。
- `title`、`description`、`examples` 等注解帮助模型和维护者理解，不自动施加业务验证。
- `format` 在一些实现中只是注解，是否强制校验取决于验证器配置。
- 模型可能拒绝任务、输出不完整或被安全策略截断，这些状态不能伪装成符合 Schema 的业务对象。

## 实际怎么使用

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "title": { "type": "string", "minLength": 1 },
    "priority": { "type": "string", "enum": ["low", "medium", "high"] },
    "evidence": { "type": "array", "items": { "type": "string" } }
  },
  "required": ["title", "priority", "evidence"],
  "additionalProperties": false
}
```

处理流程：

1. 使用供应商支持的 Structured Output 配置发送 Schema。
2. 检查响应状态、拒绝和不完整原因。
3. 获取完整结构化值。
4. 使用独立运行时验证器校验。
5. 再做业务验证，例如 ID 是否存在、日期是否允许、证据是否支持结论。
6. Schema 版本与 Prompt、模型和评估集一起记录。

## 常见错误与边界

- 只要求“返回 JSON”，没有使用 Schema 或运行时验证。
- 把 TypeScript 类型断言当成验证；类型断言不会检查实际网络数据。
- 使用供应商不支持的 Schema 关键字，或假设所有验证器对 `format` 行为相同。
- Schema 过度复杂，包含大量深层 Union 和可选字段，降低可维护性。
- 结构正确后直接执行写操作，没有业务规则、权限和人工确认。

## 补充知识

Tool Calling 也依赖参数 Schema，但 Tool 服务端必须重新验证。Schema 变更应考虑向后兼容；新增必填字段通常是破坏性变化。可以从一份权威 Schema 生成类型和文档，减少多处定义漂移。

## 来源

- [JSON Schema Specification](https://json-schema.org/specification)（访问日期：2026-07-16）
- [JSON Schema：Creating Your First Schema](https://json-schema.org/learn/getting-started-step-by-step)（访问日期：2026-07-16）
- [JSON Schema：Type-specific Keywords](https://json-schema.org/understanding-json-schema/reference/type)（访问日期：2026-07-16）
- [OpenAI API：Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)（访问日期：2026-07-16）

