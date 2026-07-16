---
type: ai-note
stage: junior
topic: reliability-rate-limits-usage
verified: 2026-07-16
tags: [ai, cancellation, timeout, retry, rate-limit, usage]
---

# 取消、超时、有限重试、限流与 Usage

## 是什么

取消终止不再需要的工作；超时限制等待时间；重试再次尝试临时失败；限流限制单位时间或并发内的请求/Token；Usage 是供应商返回的实际资源用量。它们共同控制可靠性、容量和费用。

## 为什么需要

模型服务受网络、供应商容量和配额影响。无边界重试会放大故障和费用；不传播取消会让用户离开后任务继续运行；只按请求数限流会忽略长上下文和高输出请求。

## 关键特性

- 超时应分连接、首事件、空闲、单工具和总任务。
- 只对暂时性且幂等的失败重试；认证、权限、参数和 Schema 错误先修请求。
- 使用指数退避和随机抖动，遵守 `Retry-After` 或官方限流 Header。
- 重试次数必须是总次数的一部分，不能由 SDK、网关和业务层各自独立无限重试。
- 限流可按组织、项目、租户、用户、模型、请求和 Token 维度组合。
- Usage 可能只在完成事件返回，流中断时需按接口能力处理不完整计量。

## 实际怎么使用

```text
总尝试次数：3（首次 + 最多 2 次重试）
单次总超时：60 秒
退避：供应商 Retry-After，否则指数退避 + jitter
可重试：连接重置、部分 429、部分 5xx
不可重试：401、403、无效 Schema、业务拒绝
总任务上限：Token、费用、步骤、时长和并发
```

在请求开始前预留配额，完成后用实际 Usage 结算差额。Agent 每步都检查剩余预算。写 Tool 使用幂等键，避免模型请求重试导致重复副作用。

## 常见错误与边界

- 所有 `429` 固定等待一秒；不同限流窗口和配额需要读取官方信息。
- SDK 已重试，外层又重试，实际调用数呈乘法增长。
- 只限制请求数，不限制输入输出 Token、图片或工具调用。
- 用户取消后仅停止前端渲染，后端、队列和供应商请求未终止。
- 用估算 Token 代替结算 Usage，成本统计长期偏差。

## 补充知识

系统过载时优先拒绝或排队低优先任务，避免所有请求同时变慢。降级可切换模型、缩短上下文或关闭非必要工具，但需要评估质量和透明告知。

## 来源

- [OpenAI Docs：Rate Limits](https://platform.openai.com/docs/guides/rate-limits)（访问日期：2026-07-16）
- [Anthropic API：Rate Limits](https://docs.anthropic.com/en/api/rate-limits)（访问日期：2026-07-16）
- [AWS Architecture Blog：Exponential Backoff and Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)（访问日期：2026-07-16）
- [MDN：AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)（访问日期：2026-07-16）

