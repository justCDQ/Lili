---
type: ai-note
stage: junior
topic: solution-selection
verified: 2026-07-16
tags: [ai, prompt, rag, tools, workflow, architecture]
---

# Prompt、代码、RAG、Tool 与 Workflow 的选择

## 是什么

Prompt 控制模型推理时的任务表达；代码执行确定性规则；RAG 提供外部知识上下文；Tool 读取实时事实或执行动作；Workflow 用固定步骤、路由、并行和检查组合多个能力。

## 为什么需要

不同问题有不同根因。继续增加 Prompt 长度不能可靠解决权限、最新数据、计算、外部副作用和多步骤恢复。

## 选择规则

| 问题 | 优先方案 |
| --- | --- |
| 任务、语气、边界和输出说明不清 | Prompt |
| 必须严格执行的计算、校验、权限、状态转换 | 代码 |
| 需要私有、最新、可删除、可引用知识 | RAG |
| 需要查询数据库/API或执行外部动作 | Tool |
| 步骤可预测、需要重试/审批/恢复 | Workflow |
| 步骤无法预先确定且需要自主选择工具 | 受限 Agent，并设置停止和审批 |

## 实际怎么使用

1. 从失败样例定位问题层：格式、知识、计算、权限、行动还是流程。
2. 选择最简单且可验证的机制。
3. 定义输入、输出、失败和权限。
4. 在固定评测集上与基线比较。
5. 增加可观测性和回退，再进入生产。

示例：订单退款不能只靠 Prompt。模型可分类意图，Tool 查询订单，代码检查退款窗口和权限，Workflow 请求确认并执行幂等写入，最后模型解释结果。

## 常见错误与边界

- 用 Prompt 实现财务、权限和合规规则。
- 用 RAG 解决需要实时事务一致性的查询。
- 能用固定 Workflow 时直接引入自由 Agent。
- Tool 描述清楚就信任参数，不在服务端校验。
- 多层方案上线后没有分层评测，失败无法定位。

## 补充知识

方案可以组合，但每增加一层都会增加延迟、费用和故障模式。架构评审应记录为何简单方案不足。

## 来源

- [Anthropic：Building Effective Agents](https://www.anthropic.com/engineering/building-effective-agents)（访问日期：2026-07-16）
- [OpenAI：Prompt Engineering](https://developers.openai.com/api/docs/guides/prompt-engineering)（访问日期：2026-07-16）
- [Model Context Protocol：Tools](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)（访问日期：2026-07-16）
- [Lewis et al.：Retrieval-Augmented Generation](https://arxiv.org/abs/2005.11401)（访问日期：2026-07-16）

