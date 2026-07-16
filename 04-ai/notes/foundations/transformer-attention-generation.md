---
type: ai-note
stage: beginner
topic: transformer-attention-generation
verified: 2026-07-16
tags: [ai, transformer, attention, qkv, positional-encoding, autoregressive]
---

# Transformer、Attention、Q/K/V、位置与自回归生成

## 是什么

Transformer 是以 Attention 和前馈网络为核心的序列模型架构。Self-Attention 让序列中每个位置根据其他可见位置计算加权信息。输入表示通过不同线性投影生成 Query、Key 和 Value：Query 与 Key 决定权重，权重再聚合 Value。位置信息用于区分 Token 顺序。自回归语言模型按条件概率一次生成下一个 Token，并把已生成 Token 作为后续上下文。

## 为什么需要

这些机制解释了模型为何受上下文、顺序、Token 数和解码策略影响，也解释了生成输出为何不是直接查询事实表。应用工程不需要从零训练 Transformer，但需要理解上下文限制和生成误差来源。

## 关键特性

Scaled dot-product Attention 的核心形式：

```text
Attention(Q,K,V) = softmax(QKᵀ / √dₖ)V
```

- Self-Attention 的 Q/K/V 来自同一序列的不同投影。
- Multi-head Attention 使用多组投影并行计算，再合并结果。
- 因果 Mask 阻止自回归模型在训练/生成某位置时查看未来 Token。
- 原始 Transformer 使用正弦/余弦位置编码；现代模型也使用可学习、旋转或其他位置方法，具体取决于模型。
- 自回归生成在每一步产生下一 Token 分布，再按贪心、采样或搜索策略选择 Token。
- Attention 权重不是可靠的自然语言解释或事实依据。

## 实际怎么使用

应用层面的直接影响：

1. 指令、数据和示例的顺序会影响可见上下文和注意分配，需通过评测确认。
2. 长文本要保留结构、来源和相关片段，不能仅靠一次塞满窗口。
3. Streaming 展示逐 Token/事件生成结果，但中间内容可能被后续内容修正。
4. 温度和采样控制输出分布，不修复缺失知识、错误上下文和权限问题。
5. 要求事实可靠时使用检索、工具、引用和校验，不把语言流畅度当作真实性。

## 常见错误与边界

- 把 Attention 解释为模型“理解”或可直接读取的推理过程。
- 认为模型逐字从训练数据复制；生成基于参数化概率分布，仍可能记忆或复现训练片段，但机制不能简单等同数据库检索。
- 认为更大 Context Window 自动解决长期记忆和事实更新。
- 用温度 `0` 承诺完全确定和永远正确。
- 把 Q/K/V 当作固定的人类可读“问题、键、答案”；它们是学习得到的向量投影。

## 补充知识

推理时 KV Cache 保存先前 Token 的 Key/Value，减少重复计算，但占用显存并影响长上下文服务成本。不同模型可能采用稀疏、滑动窗口或其他 Attention 变体。

## 来源

- [Vaswani et al.：Attention Is All You Need](https://arxiv.org/abs/1706.03762)（访问日期：2026-07-16）
- [NeurIPS：Attention Is All You Need](https://papers.nips.cc/paper/7181-attention-is-all-you-need)（访问日期：2026-07-16）
- [Hugging Face LLM Course](https://huggingface.co/learn/llm-course/en/chapter1/4)（访问日期：2026-07-16）
- [Google：Introduction to Large Language Models](https://developers.google.com/machine-learning/resources/intro-llms)（访问日期：2026-07-16）

