---
title: Token Budget 的建立与分配
stage: intermediate
direction: ai
topic: context-engineering
---

# Token Budget 的建立与分配

Token Budget 是一次模型调用中对指令、用户输入、历史、检索证据、工具结果和输出预留的显式容量计划。它同时约束请求能否被模型接受、信息是否被截断、延迟、费用以及模型从长上下文中找到关键证据的难度。

## 前置知识与边界

前置阅读：

- [Tokenization、Context Window 与成本](../00-foundations/tokenization-context-cost.md)。
- [模型、Token、延迟与成本记录](../02-prompt/prompt-observability.md)。
- [指令与数据的边界](01-instruction-data-boundary.md)。

本文讨论应用层预算，不给任何模型写死通用窗口长度。模型标识、接口、输入类型和供应商规则共同决定实际上限，应从当前模型元数据与真实 Usage 获取。

## 预算方程

对一次调用，可以先写成：

```text
context_limit
>= protocol_overhead
 + fixed_instructions
 + current_user_input
 + conversation_history
 + retrieved_context
 + tool_results
 + reserved_output
 + safety_margin
```

其中：

- `context_limit`：特定模型与接口允许的总上下文。
- `protocol_overhead`：消息角色、内容项、工具 Schema 等编码开销。
- `fixed_instructions`：系统和任务规则。
- `reserved_output`：预期最大输出，不是平均输出。
- `safety_margin`：本地估算与服务端计数差异、动态字段和编码变化的余量。

如果接口把输出上限与输入窗口分别限制，还要同时满足两个约束，不能只减一个总数。

## 容量、费用和质量是三个问题

### 容量

请求是否超过硬上限。超限通常在生成前失败，或由某些框架静默截断。应用应自己决定删什么，不能依赖未知截断顺序。

### 费用与延迟

输入 Token、缓存 Token、输出 Token、推理 Token的计费方式可能不同。供应商报价和 Usage 字段必须按模型版本记录。更长输入通常增加传输、预填充和首 Token 延迟，但实际关系需要测量。

### 质量

“放得下”不表示“用得好”。长上下文可能引入无关、冲突和重复信息；相关证据在不同位置也可能表现不同。需要基于任务评估上下文长度、排序和召回，而不是把窗口填满。

## 预算分类

| 类别 | 是否可缩减 | 优先级 | 常用策略 |
|---|---|---|---|
| 安全与任务不变量 | 通常不可删 | 最高 | 简洁固定、版本化 |
| 输出 Schema | 不可随意删 | 最高 | 减少描述重复 |
| 当前用户任务 | 不可误删 | 高 | 超长输入另行解析 |
| 权限内关键事实 | 任务决定 | 高 | 结构化、保留来源 |
| 检索证据 | 可选择 | 中高 | rerank、阈值、去重 |
| 对话历史 | 可压缩 | 中 | 摘要、关键事实 |
| 示例 | 可选择 | 中 | 只保留区分决策的示例 |
| 工具结果 | 可裁剪 | 中 | 结构化字段、分页 |
| 调试说明 | 生产通常不发 | 低 | 留在 trace |

优先级不能仅按来源写死。对“总结当前上传文件”的任务，文件正文比历史闲聊重要；对“按刚才批准的参数执行”，授权记录和最新确认更重要。

## 建立预算配置

```json
{
  "modelProfile": "reader-large-context-v3",
  "contextLimit": 128000,
  "reservedOutput": 4000,
  "safetyMargin": 2000,
  "categories": {
    "instructions": {"hardMax": 6000, "priority": 100},
    "currentInput": {"hardMax": 12000, "priority": 95},
    "verifiedFacts": {"hardMax": 8000, "priority": 90},
    "history": {"hardMax": 18000, "priority": 60},
    "retrieval": {"hardMax": 70000, "priority": 80},
    "toolResults": {"hardMax": 8000, "priority": 70}
  }
}
```

数字只是一份任务配置示例，不是模型通则。预算配置要与完整模型标识、Prompt 版本和评估结果一起版本化。

## 计数策略

### 供应商 Tokenizer

模型提供官方 tokenizer 时，本地计数可用于发送前规划。必须使用与模型匹配的 encoding，并认识到消息封装、工具定义和多模态内容可能有额外开销。

### Token Counting API

部分平台提供服务端计数接口，能更接近真实请求格式，但增加一次调用、延迟和依赖。适合大文件、昂贵请求和调试，不一定适合每个低延迟请求。

### Usage 回填

最终以模型 API 返回的 Usage 为计费和观测依据。记录本地估算与实际输入 Token 的差值，差值持续扩大时更新计数器或安全余量。

### 字符估算

字符数只能作拒绝超大输入的粗筛。不同语言、代码、数字和空白的 Token 比例不同，不能用固定“四字符一 Token”作为跨模型准确预算。

## 一个确定性预算分配器

以下 JavaScript 不依赖具体 tokenizer。调用方传入已经计数的候选项：

```javascript
export function allocateContext({
  contextLimit,
  reservedOutput,
  safetyMargin,
  fixedTokens,
  candidates,
}) {
  for (const value of [
    contextLimit,
    reservedOutput,
    safetyMargin,
    fixedTokens,
  ]) {
    if (!Number.isSafeInteger(value) || value < 0) {
      throw new TypeError("token values must be non-negative integers");
    }
  }

  let remaining =
    contextLimit - reservedOutput - safetyMargin - fixedTokens;
  if (remaining < 0) {
    return {
      status: "fixed_content_overflow",
      selected: [],
      remaining,
    };
  }

  const ordered = candidates
    .map((item, index) => ({...item, stableIndex: index}))
    .sort((a, b) =>
      b.priority - a.priority ||
      a.tokens - b.tokens ||
      a.stableIndex - b.stableIndex
    );

  const selected = [];
  const excluded = [];

  for (const item of ordered) {
    if (!Number.isSafeInteger(item.tokens) || item.tokens < 0) {
      throw new TypeError(`invalid token count for ${item.id}`);
    }
    if (item.tokens <= remaining) {
      selected.push(item);
      remaining -= item.tokens;
    } else {
      excluded.push({...item, reason: "budget"});
    }
  }

  return {status: "ok", selected, excluded, remaining};
}
```

这个算法按优先级贪心，不保证全局最优。真实系统还要处理最小来源覆盖、同一文档最多若干片段、强制证据和相邻 chunk 合并。

## 为什么不能只从尾部截断

尾部截断可能删除：

- 当前用户刚提出的问题。
- 工具最新返回结果。
- JSON 的闭合结构。
- 代码块末尾的错误堆栈。
- 最近一次人工确认。

头部截断则可能删除任务规则和定义。正确方法是对内容项做语义选择，并在截断后重新构造完整结构。

## 应用案例一：文档问答

### 输入

模型配置：

```json
{
  "contextLimit": 32000,
  "reservedOutput": 2500,
  "safetyMargin": 1000
}
```

已计数内容：

| 类别 | Token | 说明 |
|---|---:|---|
| 指令与 Schema | 2800 | 固定 |
| 当前问题 | 240 | 固定 |
| 用户权限事实 | 160 | 固定 |
| 对话摘要 | 1800 | 候选 |
| 检索片段 A | 7200 | 直接回答 |
| 检索片段 B | 6800 | 直接回答 |
| 检索片段 C | 6400 | 背景重复 |
| 检索片段 D | 5200 | 低相关 |
| 引用 metadata | 900 | 固定 |

固定 Token 为 `4100`，候选空间：

```text
32000 - 2500 - 1000 - 4100 = 24400
```

### 分配步骤

1. A、B 是不同来源的直接证据，优先级最高。
2. 对话摘要只保留与问题相关的产品版本和日期。
3. C 与 A 重复，先做去重而不是硬塞入。
4. D 低于相关阈值，排除。
5. 最终候选为 A `7200`、B `6800`、压缩摘要 `700`。
6. 余量保留给实际消息封装差异，不为填满而加入 C。

### 输出预算

```json
{
  "fixedTokens": 4100,
  "selectedTokens": 14700,
  "reservedOutput": 2500,
  "safetyMargin": 1000,
  "estimatedTotal": 22300,
  "excluded": [
    {"id": "chunk-c", "reason": "duplicate"},
    {"id": "chunk-d", "reason": "below_relevance_threshold"}
  ]
}
```

### 验证

- API 返回的实际输入 Token 与估算差值进入指标。
- A/B 放在不同位置执行对照，测量答案与引用准确率。
- 加入 C 后若质量不升、延迟和费用上升，则继续排除。
- 输出达到 2500 上限时检查答案是否被截断。
- 无相关结果时不把剩余预算填充为低相关内容。

### 失败分支

若按 Top-K 固定取四段，C、D 会占用 11600 Token，并可能挤掉输出或关键历史。修复是把 Token 作为检索后选择约束，同时保留相关性、来源覆盖和去重规则。

## 应用案例二：代码修复助手

### 输入

用户要求修复测试失败。候选上下文：

- 当前失败测试：3100 Token。
- 堆栈：1200 Token。
- 被测函数：2400 Token。
- 直接依赖模块：5000 Token。
- 整个仓库 README：9000 Token。
- 最近 20 轮对话：11000 Token。
- 编码规范：1800 Token。

任务预算只有 18000 输入 Token。

### 分配策略

1. 保留当前任务、测试、堆栈和被测函数。
2. 从依赖模块按符号切片，只保留被调用函数，降至 1600 Token。
3. 编码规范只保留与修改语言相关部分，降至 500 Token。
4. 历史抽取用户明确约束和已尝试失败方案，降至 900 Token。
5. README 与任务无直接关系，排除。
6. 预留足够输出给 diff 和解释。

### 结果

| 内容 | Token |
|---|---:|
| 固定任务规则 | 1500 |
| 当前问题 | 300 |
| 测试 | 3100 |
| 堆栈 | 1200 |
| 被测函数 | 2400 |
| 依赖符号 | 1600 |
| 规范片段 | 500 |
| 历史事实 | 900 |
| 输入合计 | 11500 |

剩余容量用于输出和安全余量。模型不需要整个仓库才能提出局部修复。

### 验证

- 修复后运行真实测试，不以模型解释为成功。
- 比较“整个依赖文件”和“符号切片”两个版本的修复通过率。
- 记录首 Token 延迟、总时长和输入 Token。
- 检查是否因裁剪漏掉接口契约；若漏掉，把失败加入选择规则测试。

### 失败分支

若只保留最短内容，可能省 Token 却删除类型定义或调用约束。预算器需要优先级、依赖关系和最小完整单元，不能只做最短项装箱。

## 动态预算策略

### 按任务类型

- 抽取：输入预算高，输出 Schema 小。
- 长文生成：输出预留高，证据要精选。
- 代码修改：完整符号和测试优先。
- Agent：要为未来工具结果和步骤保留空间。
- 多模态：使用平台实际计数，不把图片按字符估算。

### 按模型能力

更大的窗口可允许更多证据，但不自动获得更高利用率。切换模型时同时运行：

- 长度分桶评估。
- 关键证据位置评估。
- 噪声和冲突评估。
- Token 估算误差。
- 质量、延迟与费用联合评估。

### 按风险

高风险任务宁可返回“证据不足”，也不以低相关内容填满预算。安全规则、权限事实和人工确认记录不可被摘要器静默删除。

## Prompt Cache 与预算

Prompt caching 可能减少重复前缀的费用或延迟，但缓存命中不增加模型上下文上限。为了缓存而把不相关固定内容放入请求，仍会增加上下文噪声。

缓存设计应记录：

- 可缓存前缀的版本。
- 命中 Token 与未命中 Token。
- 租户和数据隔离。
- 模型或工具 Schema 变化后的失效。
- 数据保留与平台兼容性。

## 调试指标

### 每次请求

- estimated input tokens。
- actual input tokens。
- cached input tokens。
- output tokens。
- 每个类别的 selected/excluded tokens。
- exclusion reason。
- context utilization ratio。
- 截断与摘要次数。

### 聚合指标

- 估算绝对误差和相对误差。
- 超限请求率。
- 输出截断率。
- 按长度分桶的任务成功率。
- 输入 Token 与首 Token 延迟相关性。
- 每个成功任务的费用，而非每次调用费用。

## 常见错误

### 把最大窗口当目标

最大窗口是硬边界，不是推荐填充量。信息相关性和组织质量更重要。

### 用字符切片破坏结构

字符切片可能破坏 UTF-8 之外的逻辑单元、JSON、Markdown 表格或代码语法。按结构单元切分并重新计数。

### 输出只预留平均值

长尾输出会被截断。根据产品上限、Schema 和任务复杂度设置最大值，并测试达到上限时的失败行为。

### 忽略工具 Schema

大量工具描述和参数 Schema 也占输入。只暴露当前步骤需要的工具，并记录其 Token。

### 只看单次费用

上下文过少可能增加重试、工具调用和人工修复。用完成一个真实任务的总质量、总延迟和总费用比较。

## 生产边界

- 模型限制从受控配置读取，不由客户端提交。
- 本地计数器版本与模型映射必须可更新。
- API 超限错误不能无限重试。
- 摘要失败时保留关键事实，不静默删除授权记录。
- 超长文件进入解析和检索流程，不直接全部塞入上下文。
- 预算日志脱敏，不能为了计数保存完整敏感输入。
- Prompt、模型、工具集合变化后重新校准固定开销。

## 综合练习：多来源研究报告

构建预算器，为十个来源生成带引用的研究报告。

验收标准：

- 配置明确总窗口、输出预留和安全余量。
- 每个内容项有真实 tokenizer 计数、来源、优先级和去重组。
- 至少保留三个独立来源，不能被单个超长文档占满。
- 关键证据位置在头、中、尾三种顺序下运行评估。
- 输出被截断时返回可恢复状态，不提交半截 JSON。
- 记录估算与 API Usage 差值。
- 比较短上下文、高相关上下文和填满窗口三种方案。
- 用报告正确性、引用准确率、延迟和总费用作联合决策。

## 来源

- [OpenAI tiktoken](https://github.com/openai/tiktoken)（访问日期：2026-07-17）
- [OpenAI Cookbook：How to count tokens with tiktoken](https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb)（访问日期：2026-07-17）
- [Lost in the Middle: How Language Models Use Long Contexts](https://arxiv.org/abs/2307.03172)（访问日期：2026-07-17）
- [OpenAI API：Data controls and prompt caching storage](https://platform.openai.com/docs/models/default-usage-policies-by-endpoint)（访问日期：2026-07-17）
