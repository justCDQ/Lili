---
type: ai-note
stage: beginner
topic: model-adaptation-methods
verified: 2026-07-16
tags: [ai, pretraining, sft, preference-optimization, prompt, rag, fine-tuning]
---

# Pretraining、SFT、偏好优化、Prompt、RAG 与 Fine-tuning

## 是什么

- Pretraining：在大规模数据上训练基础模型，学习通用表示和生成能力。
- SFT（Supervised Fine-Tuning）：使用输入与目标输出样例继续训练，塑造任务行为和格式。
- 偏好优化：使用人类或模型偏好数据，让模型更倾向被选中的响应；具体算法包括 RLHF、DPO 等，不能互相等同。
- Prompt：在推理时提供任务、上下文、约束和示例，不更新模型参数。
- RAG：推理前从外部知识源检索内容并加入上下文，使回答可使用最新、私有或可引用资料。
- Fine-tuning：所有继续更新模型参数的训练方法的宽泛称呼，SFT 和偏好优化都可属于其中。

## 为什么需要区分

每种方法解决的问题、数据需求、费用、更新速度和风险不同。把所有质量问题都交给 Fine-tuning 会增加数据和运维成本，也不能自动解决最新知识、权限或引用。

## 选择规则

1. 任务和成功标准不清：先定义评估集，不先改模型。
2. 指令、格式或少量示例即可表达：先改 Prompt 和 Schema。
3. 规则确定且必须保证：使用代码校验或确定性逻辑。
4. 需要最新、私有、可删除和可引用知识：优先 RAG 或 Tool。
5. 需要执行外部动作或精确事实：使用 Tool，并在服务端校验权限。
6. 大量稳定样例显示特定行为、风格或任务模式无法靠前述方法达到：评估 Fine-tuning。
7. 人类偏好难以写成单一参考答案：在高质量偏好数据和可靠评测下考虑偏好优化。

## 实际怎么使用

建立基线顺序：

```text
基础模型 + 简单 Prompt
→ 明确 Prompt + Structured Output
→ 代码规则 / Tool
→ RAG + 引用
→ 更合适的模型
→ SFT / 偏好优化
```

每一步在同一评估集上比较质量、延迟、成本和安全。Fine-tuning 前确认数据许可、隐私、划分、版本、回滚和模型生命周期。

## 常见错误与边界

- 用 Fine-tuning 注入会频繁变化的事实；更新和删除困难。
- 认为 RAG 能保证正确；检索可能漏召回、召回过时或不相关内容，生成仍可能不受证据约束。
- 用 Prompt 要求模型执行确定性权限规则。
- 把供应商“微调”统一理解为同一种算法，不看训练目标和支持能力。
- 没有独立评估集就训练，无法判断改善还是记忆训练样例。

## 补充知识

方法可以组合，例如经过 SFT 的模型使用 RAG 和 Tool。但组合会增加故障点，应能定位问题发生在任务定义、模型、检索、上下文、工具还是验证阶段。

## 来源

- [Ouyang et al.：Training language models to follow instructions with human feedback](https://arxiv.org/abs/2203.02155)（访问日期：2026-07-16）
- [Lewis et al.：Retrieval-Augmented Generation](https://arxiv.org/abs/2005.11401)（访问日期：2026-07-16）
- [Rafailov et al.：Direct Preference Optimization](https://arxiv.org/abs/2305.18290)（访问日期：2026-07-16）
- [Hugging Face LLM Course](https://huggingface.co/learn/llm-course/)（访问日期：2026-07-16）

