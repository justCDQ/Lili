---
type: ai-note
stage: beginner
topic: programming-basics
verified: 2026-07-16
tags: [ai, programming, async, errors]
---

# AI 应用需要的编程基础

## 是什么

AI 应用的编程基础包括值与变量、函数、集合、模块、异常和异步。模型调用本质上是带有不确定结果、网络延迟和费用的远程操作；这些语言能力用于构造请求、验证响应、组合步骤和处理失败。

## 为什么需要

- 变量保存配置、请求参数、模型响应和中间状态。
- 函数隔离模型调用、校验、重试、检索和工具执行。
- 对象/字典表达 JSON 结构；数组/列表表达消息、文档和评测样例。
- 模块隔离厂商 SDK、业务规则和基础设施，降低替换成本。
- 异常处理网络失败、认证失败、限流、超时、格式错误和业务失败。
- 异步让程序在等待网络、文件和数据库时继续处理其他工作。

## 关键特性

### 数据与函数

外部输入不可信。即使语言具有静态类型，请求体、环境变量、文件和模型结果仍需运行时验证。函数应有明确输入、输出和失败方式，避免直接读取大量全局状态。

### 模块

建议至少分为：`model-client`、`schemas`、`prompts`、`domain`、`evals` 和入口层。业务代码不要散落厂商字段名，模型 Client 负责转换统一结构。

### 异常

语法错误发生在代码无法被解析时；异常发生在代码运行期间。应用应区分可重试错误和不可重试错误。认证失败、Schema 不合法通常不应自动重试；临时网络错误可能有限重试。

### 异步

JavaScript 的异步 API 通常返回 Promise；`async` 函数也总是返回 Promise。`await` 只等待当前异步流程，不等于阻塞整个运行时。相互独立的请求可受控并发，存在依赖的步骤必须按顺序执行。

## 实际怎么使用

```ts
type Extracted = { title: string; tags: string[] };

async function extract(text: string): Promise<Extracted> {
  if (!text.trim()) throw new Error("input is empty");

  const raw = await modelClient.generate({
    task: "extract_metadata",
    input: text,
  });

  return extractedSchema.parse(raw);
}
```

这段代码体现：输入检查、函数边界、异步模型调用、统一 Client 和运行时 Schema 校验。生产代码还要增加超时、取消、Usage、日志和错误分类。

## 常见错误与边界

- 捕获异常后返回空对象，导致调用方无法区分真实空结果与系统失败。
- 对所有失败无限重试，造成费用增加和重试风暴。
- 在模块导入时立即调用模型，导致测试、工具脚本和构建过程产生副作用。
- 对一批请求直接 `Promise.all` 而不限制并发，触发速率限制。
- 把 Prompt 字符串、Schema、业务判断和网络调用写在一个函数中，无法单独测试。

## 补充知识

初级阶段应同时学习日志、单元测试和依赖注入。模型是外部依赖，测试时应使用固定响应或受控测试服务；评估模型质量时再调用真实模型。

## 来源

- [MDN JavaScript Guide](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide)（访问日期：2026-07-16）
- [MDN：使用 Promise](https://developer.mozilla.org/en-US/docs/Learn_web_development/Extensions/Async_JS/Promises)（访问日期：2026-07-16）
- [MDN：JavaScript Modules](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)（访问日期：2026-07-16）
- [Python：Errors and Exceptions](https://docs.python.org/3/tutorial/errors.html)（访问日期：2026-07-16）

