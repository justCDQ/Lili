---
type: ai-note
stage: junior
topic: schema-in-prompt-workflow
verified: 2026-07-16
tags: [ai, prompt, json-schema, validation]
---

# Prompt 工作流中的 JSON Schema 与运行时校验

## 是什么

Prompt 工作流把业务输出契约表示为 JSON Schema，并在模型返回后用独立验证器检查。Structured Output 负责引导或约束生成，运行时校验负责判断收到的数据，业务校验负责判断数据能否用于当前系统。

## 为什么需要

自然语言中的“严格返回 JSON”不可作为程序契约。即使供应商声明严格 Schema，响应仍可能拒绝、取消、不完整，且供应商只支持 JSON Schema 子集。

## 关键特性与规则

- Schema 方言、版本和供应商支持子集必须明确，不能假设所有关键字都执行相同语义。
- Structured Output 约束生成，独立 Validator 判断结构，业务校验判断数据能否安全使用。
- `default`、`format` 等关键字在不同验证器中的执行行为可能不同，应通过测试确认。
- Tool 参数、模型输出和持久化数据都要在信任边界重新校验。
- Schema 变更需要兼容策略；新增必填字段或修改类型通常会破坏旧消费者。

## 实际怎么使用

1. 先从业务消费者定义最小字段和约束。
2. 声明 Schema 方言和版本，关闭不需要的额外字段。
3. 在 Prompt 中解释字段语义、缺失处理和证据规则，不重复粘贴另一套冲突结构。
4. 调用 Structured Output。
5. 检查完成/拒绝状态，解析后运行独立 Validator。
6. 检查 ID、权限、时间、枚举来源和跨字段业务规则。
7. 失败时保留错误路径和响应 ID，不自动把无效值改成看似合法数据。

## 常见错误与边界

- TypeScript 类型和 JSON Schema 分别维护，长期漂移。
- 使用 `default` 期待验证器自动填字段；JSON Schema 的注解不必执行赋值。
- 依赖 `format` 必然强校验；具体实现可能只把它当注解。
- Schema 通过就认为事实准确。
- 结构失败时用正则修补后直接入库，掩盖模型或接口问题。

## 补充知识

Schema 需要版本策略。新增可选字段通常兼容，修改类型或新增必填字段可能破坏旧客户端。Tool 参数也遵循同样原则，并由工具服务端重新验证。

## 来源

- [JSON Schema Specification](https://json-schema.org/specification)（访问日期：2026-07-16）
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12)（访问日期：2026-07-16）
- [OpenAI：Structured Outputs](https://platform.openai.com/docs/guides/structured-outputs)（访问日期：2026-07-16）
