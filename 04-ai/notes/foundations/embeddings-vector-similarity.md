---
type: ai-note
stage: beginner
topic: embeddings-vector-similarity
verified: 2026-07-16
tags: [ai, embeddings, vectors, dot-product, cosine-similarity]
---

# Embedding、向量、点积与余弦相似度

## 是什么

Embedding 是把离散对象映射为固定维度数值向量的表示。向量的各维通常不能单独用人工语义解释，但整体几何关系可用于相似度、检索、聚类和分类。点积是对应维度乘积之和。余弦相似度是点积除以两个向量长度，衡量方向接近程度。

## 为什么需要

文本不能直接用传统数值索引比较语义。Embedding 将查询和文档放到同一向量空间，允许近邻检索。理解相似度才能正确选择索引度量、阈值和归一化方式。

## 关键特性

对于向量 `a`、`b`：

```text
dot(a,b) = Σ aᵢbᵢ
cos(a,b) = dot(a,b) / (||a|| ||b||)
```

- 点积同时受方向和长度影响。
- 余弦相似度忽略整体长度，值通常位于 -1 到 1；具体 Embedding 分布不保证覆盖整个区间。
- 若向量已做 L2 归一化，点积等于余弦相似度。
- 不同模型、维度或版本生成的向量通常不能直接混用。
- 相似表示统计相关，不保证事实相同、逻辑蕴含或满足权限。

## 实际怎么使用

1. 明确检索单位：标题、段落、页面还是对象字段。
2. 使用同一 Embedding 模型和预处理生成文档与查询向量。
3. 在向量数据库中配置与模型建议一致的距离度量。
4. 取 Top-K 后再做元数据权限过滤或在检索阶段强制过滤。
5. 用标注问题集评估 Recall@K，而不是凭几次搜索调整阈值。
6. 模型升级时新建索引或完整重算，并记录版本。

## 常见错误与边界

- 混用不同 Embedding 模型的向量。
- 将余弦距离和余弦相似度的排序方向弄反；数据库定义可能不同。
- 只检索向量相似内容，不结合关键词、实体、日期和权限过滤。
- 把高相似度当作答案正确或文档支持结论。
- 在向量写入后丢失原文、来源和版本，无法引用或删除。

## 补充知识

近似最近邻索引用召回率换取速度和内存。Hybrid Search 将关键词与向量信号结合；Reranker 对较小候选集做更精细排序。两者都需要真实查询集评估。

## 来源

- [Google ML Crash Course：Embeddings](https://developers.google.com/machine-learning/crash-course/embeddings)（访问日期：2026-07-16）
- [TensorFlow：Vector Embeddings](https://www.tensorflow.org/text/guide/word_embeddings)（访问日期：2026-07-16）
- [SciPy：Cosine Distance](https://docs.scipy.org/doc/scipy/reference/generated/scipy.spatial.distance.cosine.html)（访问日期：2026-07-16）

