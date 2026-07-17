---
title: 上下文去重、过期与冲突处理
stage: intermediate
direction: ai
topic: context-engineering
---

# 上下文去重、过期与冲突处理

上下文选择不能只按相似度追加内容。重复片段会浪费 Token 并放大某个来源的权重；过期内容会让模型依据失效事实回答；冲突内容若没有来源、版本和时间，模型可能任意选择一个看似合理的说法。

## 前置知识与边界

前置阅读：

- [Token Budget 的建立与分配](03-token-budget-allocation.md)。
- [长对话摘要与关键事实保留](04-conversation-summary-key-facts.md)。
- [Embedding 与向量相似度](../foundations/embeddings-vector-similarity.md)。

本文处理进入模型前的上下文治理。它不能替代源系统的数据质量、事务一致性和访问控制。

## 三类问题的准确含义

### 重复

两个内容项表达同一事实或包含高度重叠原文。重复可能来自：

- 同一文件的相邻滑窗。
- PDF 页眉页脚。
- 同一政策的 HTML、PDF 和邮件副本。
- 对话原文与摘要同时出现。
- 工具重试返回相同结果。
- 多查询检索召回同一 chunk。

字符串相同只是重复的一种。改写后的副本和包含关系也会造成语义重复。

### 过期

内容超出它的有效时间或已被新版本取代。过期依据应由数据源定义：

- `valid_from` / `valid_to`。
- 文档版本状态。
- 数据库行版本。
- 缓存 TTL。
- 用户删除或撤销事件。
- 当前请求的快照时间。

“创建时间较旧”不必然过期；历史合同可能仍是指定日期的权威事实。

### 冲突

两个上下文项对同一可比较命题给出不能同时成立的值，例如同一政策版本的退款期限分别为 7 天和 30 天。不同时间、地区、产品版本或适用对象的值不一定冲突。

## 上下文项最小元数据

```json
{
  "id": "policy:refund:v5#section-3",
  "canonicalEntity": "refund-policy",
  "claimKeys": ["refund.window_days"],
  "source": "policy-service",
  "sourcePriority": 100,
  "version": 5,
  "contentHash": "sha256:...",
  "validFrom": "2026-07-01T00:00:00Z",
  "validTo": null,
  "observedAt": "2026-07-17T09:00:00Z",
  "tenantId": "tenant_42",
  "content": "标准商品可在签收后 30 天内申请退款。"
}
```

没有 `canonicalEntity` 和 `claimKeys`，系统很难判断两个不同措辞是否谈论同一事实。

## 处理顺序

```mermaid
flowchart LR
    A["授权候选"] --> B["有效期与删除过滤"]
    B --> C["精确重复去除"]
    C --> D["近重复聚类"]
    D --> E["命题与适用范围归一化"]
    E --> F["冲突检测"]
    F --> G["按策略选择或显式保留冲突"]
    G --> H["Token Budget 与排序"]
```

必须先过滤权限和失效内容。不能为了冲突分析，把当前用户无权读取的文档加载进内存或模型。

## 精确去重

对规范化后的内容计算哈希：

```javascript
import {createHash} from "node:crypto";

export function normalizeForExactDedup(text) {
  return text
    .normalize("NFC")
    .replace(/\r\n/g, "\n")
    .replace(/[ \t]+$/gm, "")
    .trim();
}

export function contentHash(text) {
  return createHash("sha256")
    .update(normalizeForExactDedup(text), "utf8")
    .digest("hex");
}
```

规范化不能随意删除标点、数字、大小写或空白，因为代码、金额和否定含义可能改变。哈希相同只能证明规范化结果相同，不能证明来源或权限相同；保留所有 provenance。

## 近重复检测

常用方法：

- token shingles + Jaccard：适合大段复制并有少量编辑。
- MinHash：近似比较大量集合。
- SimHash：适合文本指纹。
- Embedding 相似度：能发现语义改写，但可能把不同事实误合并。
- 文档结构：相邻滑窗共享相同 chunk ID 范围。

近重复阈值必须用标注数据校准。数字不同的两条政策可能向量非常相似，却恰好是必须保留的冲突。

### 合并原则

合并后至少保留：

- 所有来源 ID。
- 最高权威来源。
- 最新有效版本。
- 差异字段。
- 原始片段定位。
- 合并算法和阈值。

不要只保留“最长文本”；最长副本可能包含更多噪声或恶意注入。

## 过期判断

### 有效时间与系统时间

- 有效时间：事实何时在业务世界成立。
- 系统时间：系统何时记录或观察到该事实。

用户问“2025 年 12 月的退款规则”时，需要有效时间覆盖该日期的历史版本，而不是当前最新版本。

### TTL

TTL 适合允许短暂陈旧的缓存，不应被误用为业务真相：

- 股票报价可设秒级 TTL。
- 用户权限应在高风险操作前重新读取。
- 合同条款由版本和有效期决定，不是简单缓存 24 小时。

### 删除与撤销

删除事件应形成 tombstone，并传播到：

- 文档存储。
- chunk 表。
- 向量索引。
- 搜索索引。
- 缓存。
- 对话摘要和长期记忆。

仅从主表删除而向量库仍可召回，会产生“幽灵上下文”。

## 冲突识别

先把自然语言转为可比较命题：

```json
{
  "subject": "refund-policy",
  "predicate": "window_days",
  "object": 30,
  "qualifiers": {
    "region": "CN",
    "productType": "standard"
  },
  "validFrom": "2026-07-01",
  "sourceItemId": "policy:v5#3"
}
```

冲突条件至少需要：

1. 主体相同。
2. 谓词相同。
3. 适用限定重叠。
4. 有效时间重叠。
5. 值不能同时成立。

模型抽取的命题仍需 Schema 校验和抽样审核，不能作为唯一事实源。

## 冲突解决策略

| 策略 | 适用条件 | 风险 |
|---|---|---|
| 权威源优先 | 事实所有者明确 | 权威配置错误 |
| 更高版本优先 | 版本线性且已发布 | 草稿版本误入 |
| 有效时间匹配 | 查询有明确时间点 | 时间缺失 |
| 更窄适用范围优先 | 特例覆盖通则有制度依据 | 范围推断错误 |
| 多数来源 | 独立来源且同等权威 | 副本放大票数 |
| 显式无答案 | 无法可靠裁决 | 用户体验需解释 |

不能让模型按措辞自信度裁决法律、财务、权限或库存事实。

## 应用案例一：企业退款政策问答

### 候选输入

| ID | 内容 | 状态 |
|---|---|---|
| A | 标准商品 30 天退款 | policy-service v5，当前有效 |
| B | 标准商品 7 天退款 | PDF v3，已于 2026-07-01 失效 |
| C | 标准商品 30 天退款 | v5 的网页副本 |
| D | 定制商品不可退款 | policy-service v5，特例 |
| E | 所有商品 90 天退款 | 用户上传文件，无权威 |

问题：“中国区标准商品现在多久可以退款？”

### 处理步骤

1. 按租户和地区授权过滤。
2. 以当前时间过滤掉 B。
3. A 与 C 内容近重复，合并 provenance，A 为事实所有者。
4. D 的限定是定制商品，不与标准商品冲突。
5. E 可作为用户附件数据，但无权覆盖政策。
6. 最终上下文只保留 A 的内容、A/C 来源和 D 的适用边界。

### 输出

```json
{
  "answer": "中国区标准商品可在签收后 30 天内申请退款。",
  "status": "answered",
  "citations": [
    "policy-service:v5#section-3"
  ],
  "applicability": {
    "region": "CN",
    "productType": "standard",
    "asOf": "2026-07-17"
  }
}
```

### 验证

- 引用必须定位 A，而非无权威的 E。
- 查询日期改为 2026-06-15 时，应选择 v3。
- 商品类型改为定制时，应返回 D。
- A 暂时不可用时，不自动用 E 替代。
- 去重前后答案一致，输入 Token 降低。

### 失败分支

若按“更新时间最新”选择，用户刚上传的 E 会覆盖政策服务。修复是先定义语义权威，再考虑版本和观察时间。

## 应用案例二：事故响应上下文

### 输入

事故过程中收到：

- 09:00 监控：错误率 2%。
- 09:05 监控重试：同一测量重复三次。
- 09:10 值班员：错误率约 20%。
- 09:12 指标服务：错误率 18.7%，窗口 5 分钟。
- 09:14 旧缓存：错误率 2%，缓存年龄 14 分钟。
- 09:16 部署系统：版本 `2026.07.17-2` 已回滚。

### 处理

1. 用事件 ID 去除 09:05 的传输重试。
2. 保留 09:10 人工观察，但标记为估算。
3. 以指标服务的时间窗和采样定义作为数值来源。
4. 09:14 缓存超过事故看板允许的 60 秒，排除。
5. 回滚事件是部署事实，不自动推断错误率已恢复。
6. 上下文同时包含当前指标与回滚状态。

### 输出状态

```json
{
  "incident": "inc_771",
  "currentErrorRate": {
    "value": 0.187,
    "window": "5m",
    "observedAt": "2026-07-17T09:12:00Z",
    "source": "metrics-service"
  },
  "deployment": {
    "version": "2026.07.17-2",
    "state": "rolled_back",
    "observedAt": "2026-07-17T09:16:00Z"
  },
  "unknowns": [
    "回滚后错误率是否恢复"
  ]
}
```

### 验证

- 时间线保留原始事件，不覆盖历史值。
- 当前值必须满足 freshness SLA。
- 指标单位与窗口随值传递。
- 回滚后等待新指标，不让模型臆测恢复。
- 事件重放不会因重复消息产生三次相同事实。

### 失败分支

若按最后到达时间选择，09:14 的旧缓存可能成为“最新”值。修复是比较事实的 `observedAt` 和 freshness，而不是消息到达顺序。

## 实现一个选择器

```javascript
export function selectCurrent(items, queryTime) {
  const now = new Date(queryTime).getTime();

  const eligible = items.filter((item) => {
    if (item.deleted) return false;
    const start = new Date(item.validFrom).getTime();
    const end = item.validTo
      ? new Date(item.validTo).getTime()
      : Number.POSITIVE_INFINITY;
    return start <= now && now < end;
  });

  const byClaim = new Map();
  for (const item of eligible) {
    const key = `${item.canonicalEntity}:${item.claimKey}`;
    const existing = byClaim.get(key);
    if (
      !existing ||
      item.sourcePriority > existing.sourcePriority ||
      (
        item.sourcePriority === existing.sourcePriority &&
        item.version > existing.version
      )
    ) {
      byClaim.set(key, item);
    }
  }

  return [...byClaim.values()];
}
```

该函数只适用于已定义线性版本和 source priority 的数据。相同优先级、同版本却值不同应返回冲突，不应由遍历顺序静默覆盖。

## 调试与监控

### 请求级 trace

- 初始候选数。
- 权限过滤数。
- 过期和删除过滤数。
- 精确重复组与近重复组。
- 冲突命题。
- 最终选择理由。
- 每项 Token 和来源。

### 线上指标

- duplicate token ratio。
- stale retrieval rate。
- unresolved conflict rate。
- deleted content recall incidents。
- citation-to-current-version accuracy。
- 冲突导致无答案的比例。
- 重建索引完成时间。

## 常见错误

### 相似度高就合并

“退款 7 天”和“退款 30 天”语义向量很近，但值冲突。命题中的数字和限定必须保留。

### 最新写入覆盖一切

系统时间最新不代表业务有效时间最新，也不代表来源权威。

### 去重丢 provenance

副本可合并内容，但所有来源、权限和引用位置仍要保存。

### 冲突全部交给模型

模型能解释冲突，不能替代事实所有者。无法确定时应输出冲突和无答案状态。

### 删除只处理主库

派生索引、缓存和摘要继续保留内容，会违反用户控制和权限边界。

## 生产边界

- 版本、有效期和 source priority 由受控服务定义。
- 所有时间统一保存绝对时刻，并保留业务时区。
- 近重复阈值按语言和文档类型评估。
- 更新与删除使用幂等事件。
- 索引重建有进度、水位和对账。
- 发现冲突时记录双方，不破坏性覆盖。
- 引用指向内容版本，而不是可变化的普通 URL。
- 任何自动裁决都能解释规则并回放。

## 综合练习：多版本员工手册

构建一个能回答“某日期、某地区、某员工类型适用什么规则”的上下文选择器。

验收标准：

- 文档含版本、有效时间、地区、员工类型和 source priority。
- 同一版本的 PDF/HTML 副本合并但保留两个引用。
- 相似措辞中的数字冲突不会被近重复算法吞掉。
- 历史日期查询选择当时有效版本。
- 删除事件能从全文、向量、缓存和摘要中传播。
- 无法裁决时返回冲突来源而不是猜测。
- 固定评估覆盖重复、过期、迟到事件、特例和跨租户内容。
- 输出 Token、引用准确率和冲突检测准确率都有记录。

## 来源

- [RFC 9110：HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-17）
- [RFC 9111：HTTP Caching](https://www.rfc-editor.org/rfc/rfc9111.html)（访问日期：2026-07-17）
- [Lamport：Time, Clocks, and the Ordering of Events](https://lamport.azurewebsites.net/pubs/time-clocks.pdf)（访问日期：2026-07-17）
- [W3C PROV-O: The PROV Ontology](https://www.w3.org/TR/prov-o/)（访问日期：2026-07-17）
