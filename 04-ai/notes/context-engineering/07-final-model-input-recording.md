---
title: 最终模型输入的记录与可重放性
stage: intermediate
direction: ai
topic: context-engineering
---

# 最终模型输入的记录与可重放性

最终模型输入记录描述一次推理实际使用了哪些指令、上下文项、工具 Schema、模型参数和版本。它用于解释“为什么包含这些内容”、定位线上问题和运行回归，但不能无条件保存完整 Prompt、用户数据、Secret 或模型私有推理。

## 前置知识与范围

前置阅读：

- [模型、参数、Token、延迟与成本记录](../prompt/prompt-observability.md)。
- [上下文权限与租户隔离](06-context-permission-tenant-isolation.md)。
- [上下文去重、过期与冲突](05-dedup-staleness-conflicts.md)。

本文讨论应用侧 trace。模型供应商的请求日志、服务端会话存储和数据保留行为应单独核对。

## 为什么普通应用日志不够

只记录用户问题和最终答案，无法回答：

- 使用了哪个 Prompt 版本。
- 检索到了哪些文档版本。
- 哪些候选因权限、过期或预算被排除。
- 工具定义是否占用了大量 Token。
- 模型看到的是原文、摘要还是裁剪内容。
- 上下文排序是否改变。
- 哪一步引入了其他租户数据。
- 输出校验前的内容与最终展示是否相同。

仅记录 SDK 请求对象又可能保存过多敏感内容，并绑定某一家接口。

## 三层记录

### 结构清单

长期保留的低敏记录：

- request ID、trace ID。
- 模型完整标识。
- Prompt、Schema、上下文策略版本。
- 内容项 ID、来源、版本、Token、选择理由。
- 工具名称和 Schema hash。
- 输出状态、Usage、延迟和费用。

### 脱敏预览

受限保存，用于调试：

- 内容项前后少量字符。
- 结构字段。
- PII/Secret 已替换。
- 访问受 RBAC 和审计控制。

### 加密原始快照

只在高价值故障、用户授权或短期评估中保存：

- 独立加密密钥。
- 严格 TTL。
- 更少人员可访问。
- 删除和法律保留流程。

默认应从结构清单开始，而不是默认保存全文。

## Trace 数据模型

```json
{
  "traceId": "4bf92f3577b34da6a3ce929d0e0e4736",
  "requestId": "req_20260717_991",
  "tenantIdHash": "hmac:tenant_42",
  "task": "support_answer",
  "model": {
    "provider": "provider-a",
    "modelId": "reader-2026-06-01",
    "parameters": {
      "maxOutputTokens": 1200
    }
  },
  "artifacts": {
    "promptVersion": "support-v12",
    "schemaVersion": "answer-v4",
    "contextPolicyVersion": "ctx-v7",
    "toolSetHash": "sha256:..."
  },
  "contextItems": [
    {
      "id": "chunk_91",
      "source": "runbook_12",
      "version": "v5",
      "trust": "untrusted",
      "estimatedTokens": 740,
      "position": 3,
      "decision": "selected"
    }
  ],
  "usage": {
    "inputTokens": 4210,
    "outputTokens": 381
  },
  "result": {
    "status": "validated",
    "finishReason": "completed"
  }
}
```

`tenantIdHash` 只能在确实需要聚合时使用，并用受控 HMAC 避免可枚举 ID 的普通哈希反查。

## 内容寻址

不保存全文时，可记录内容哈希：

```javascript
import {createHash} from "node:crypto";

export function sha256Utf8(value) {
  return createHash("sha256").update(value, "utf8").digest("hex");
}

export function contextManifest(items) {
  return items.map((item, index) => ({
    id: item.id,
    source: item.source,
    version: item.version,
    position: index,
    bytes: Buffer.byteLength(item.content, "utf8"),
    contentHash: sha256Utf8(item.content),
    trust: item.trust,
  }));
}
```

哈希证明重放时内容是否相同，不会让已经删除的内容自动恢复。若源数据按政策删除，trace 只保留允许的元数据。

## 记录上下文决策

候选项不仅记录 selected，也记录排除理由：

| decision | 含义 |
|---|---|
| `selected` | 进入最终输入 |
| `permission_denied` | 当前主体无权访问 |
| `expired` | 不在查询有效时间 |
| `duplicate` | 与已选项重复 |
| `conflict_suppressed` | 按受控权威规则排除 |
| `below_threshold` | 相关性不足 |
| `budget` | 超出 Token 预算 |
| `deleted` | 来源已撤回 |

权限拒绝项的 ID 可能本身敏感。普通调试用户只看聚合计数，安全审计角色才可看受控标识。

## 重放等级

### 结构重放

重建上下文项顺序和版本，检查选择算法。无需再次调用模型，成本最低。

### 模型重放

使用同一模型标识、参数和输入重新调用。云模型仍可能非确定，服务端实现也可能变化；重放用于比较，不承诺字节级相同。

### 影子重放

把历史输入发送给新模型或 Prompt，不影响用户。需要确认数据使用目的、保留和供应商政策允许。

### 业务重放

重新执行工具或外部动作风险最高。默认只重放只读模拟，写操作必须使用沙箱和新幂等域，不能复用生产副作用。

## 可重放包

```json
{
  "manifestVersion": 2,
  "requestId": "req_20260717_991",
  "modelProfile": "reader-2026-06-01",
  "promptArtifact": "sha256:prompt...",
  "schemaArtifact": "sha256:schema...",
  "toolArtifacts": [
    "sha256:tool-search-v3..."
  ],
  "contextArtifacts": [
    {
      "id": "chunk_91",
      "version": "v5",
      "hash": "sha256:content..."
    }
  ],
  "permissionsSnapshot": {
    "policyDecisionId": "pd_771",
    "filterHash": "sha256:filter..."
  },
  "replayPolicy": "model_only_no_tools"
}
```

包中不保存 bearer token。重放服务使用独立身份，并重新检查是否仍允许访问原始 artifact。

## 应用案例一：错误引用调查

### 现象

客服答案引用了已经失效的退款政策。普通日志只有问题和答案。

### Trace 提供的证据

- Prompt `support-v12`。
- context policy `ctx-v6`。
- chunk `policy-v3#refund` 被选择。
- 当前政策 `v5` 候选被标记 `budget`。
- v3 的 `validTo` metadata 在索引中为空。
- 模型输出引用 v3，Schema 校验通过。

### 诊断步骤

1. 问题不在模型“忘记”v5，因为 v5 没进入输入。
2. 文档更新流程没有把 v3 标为失效。
3. 预算器按片段长度先选 v3，暴露了第二个缺陷。
4. 修复索引状态传播和版本优先规则。
5. 用原 manifest 结构重放，确认 v3 被排除。
6. 用固定问答集运行新模型调用。

### 修复后输出

```json
{
  "selected": ["policy-v5#refund"],
  "excluded": [
    {"id": "policy-v3#refund", "reason": "expired"}
  ],
  "answerStatus": "validated",
  "citationVersion": "v5"
}
```

### 验证

- v3 的 tombstone 传播到全文和向量索引。
- 历史日期查询仍可按授权读取 v3。
- 当前日期查询的引用准确率回归通过。
- trace 不含客户完整问题和个人信息。

### 失败分支

若记录了最终拼接 Prompt 但没记录候选决策，只能看到 v3 在输入中，无法解释 v5 为何缺失。Manifest 必须覆盖选择过程。

## 应用案例二：跨租户泄露调查

### 现象

租户 42 的回答出现租户 17 的内部项目名。

### 紧急处理

1. 暂停相关功能或切换安全降级路径。
2. 保存受控安全证据，限制访问。
3. 撤销可能暴露的缓存和会话。
4. 检查同一 filter fingerprint 的其他请求。
5. 通知安全和隐私响应流程。

### Trace 分析

结构清单显示：

- 认证租户为 42。
- 检索 filter hash 对应的对象缺少 tenant 条件。
- `chunk_17` 的 source tenant 为 17。
- answer cache key 只包含 normalized query。

由此定位两个确定性缺陷：检索过滤遗漏和缓存键遗漏。

### 修复

- 检索适配器要求强制 tenant filter，缺失时拒绝。
- 缓存键加入安全上下文 fingerprint。
- 所有派生 artifact 保存 tenant。
- CI 加入跨租户 canary。
- 回溯 trace 识别可能受影响请求。

### 验证

- 以相同问题在两个租户运行，缓存条目不同。
- 检索器无法构造无 tenant 的生产查询。
- 普通运维界面只显示哈希租户标识。
- 安全角色访问证据有单独审计。

### 失败分支

如果完整 Prompt 被写入普通集中日志，调查系统本身会扩大泄露。安全事件取证和日常可观测性需要不同的数据访问级别。

## OpenTelemetry 映射

生成式 AI 语义约定可以统一模型、Token、操作与工具 span，但约定可能处于不同稳定级别，升级时应固定版本并测试导出字段。

建议 span 层级：

```text
http.request
└── context.assemble
    ├── permission.evaluate
    ├── retrieval.search
    ├── context.deduplicate
    └── context.allocate
└── gen_ai.generate
    ├── tool.execute
    └── output.validate
```

高基数字段和完整 Prompt 不应默认作为 span attribute。大内容使用受控 artifact 存储，只在 span 中保存引用。

## Trace Context 安全

W3C `traceparent` 用于跨服务关联，不承载用户、租户或业务秘密。外部请求带来的 trace header 需要校验；跨安全边界可选择重启 trace，同时用内部受控关联 ID 连接事件。

不要把以下内容写进 `tracestate`：

- 邮箱。
- 用户 ID 明文。
- 租户名称。
- Prompt。
- token 或 session ID。
- 文档标题。

## 采样策略

### 头部采样

请求开始时决定，成本低；不知道最终是否失败，可能漏掉稀有错误。

### 尾部采样

根据错误、延迟、成本或安全信号决定，信息更有价值；需要缓冲和更复杂基础设施。

### 强制保留

安全事件、Schema 失败、跨租户 canary 命中和高成本异常可提高采样，但仍遵守数据最小化。

### 评估样本

用于质量评估的样本应单独获得合法目的和访问控制，不能因为被 trace 采样就自动成为训练或评估数据。

## 脱敏策略

脱敏发生在导出前：

- Secret 使用检测器直接删除，不可逆掩码。
- 邮箱、电话等用类型占位符。
- 稳定关联需要 HMAC，不用无盐哈希。
- 文档正文默认只记录 hash 与 ID。
- 工具参数按字段 allowlist。
- 错误堆栈移除请求正文和 token。

脱敏器本身要有测试、版本和失败关闭策略。解析失败时宁可不导出敏感 payload。

## 观测指标

- trace coverage。
- selected/excluded item counts。
- local estimate vs API Usage。
- context assembly latency。
- permission denial rate。
- stale/duplicate/conflict rate。
- output validation failure rate。
- replay success rate。
- raw snapshot access count。
- retention deletion lag。

## 常见错误

### 记录“发送前 Prompt”但 SDK 又修改

最终接口可能增加工具、消息封装或托管会话状态。记录应用 manifest，并同时保存供应商 request ID 与 Usage。

### 把模型输出当完整解释

模型生成的理由不是上下文选择器的真实决策日志。选择理由由代码记录。

### 追求完全可复现

模型行为可能非确定。可重放目标是固定输入与配置、量化差异，而不是保证相同文本。

### 永久保存一切

这增加隐私、安全和成本风险。结构清单、短期快照与严格访问可以兼顾诊断。

## 生产验收

- 每次推理有 request ID 和 trace ID。
- Prompt、Schema、工具、模型和上下文策略都有不可变版本。
- 每个 selected item 有来源、版本、位置和 hash。
- 排除项有安全的理由统计。
- 完整内容默认不进入普通 trace。
- Secret 在导出前移除。
- 重放不执行生产写工具。
- trace retention 与源数据删除联动。
- 安全调查访问有审批与审计。
- 供应商 request ID 可用于对账，但不暴露给无关用户。

## 综合练习：可重放 RAG 请求

构建一个请求 manifest 和离线重放器。

验收标准：

- manifest 能重建上下文顺序和每项版本。
- Prompt、Schema、工具定义用内容 hash 寻址。
- 删除源文档后，普通重放不能绕过删除恢复内容。
- 模型重放与工具重放分离，写工具永不自动执行。
- 对同一输入比较旧/新模型的答案、引用、Token 和延迟。
- trace 中的 PII 和 Secret 检测为零。
- 权限缺陷和过期文档两个案例可从 trace 定位到具体步骤。
- 采样与保留策略经过故障和隐私验收。

## 来源

- [OpenTelemetry：Generative AI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/)（访问日期：2026-07-17）
- [OpenTelemetry：GenAI Observability](https://opentelemetry.io/blog/2026/genai-observability/)（访问日期：2026-07-17）
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)（访问日期：2026-07-17）
- [NIST Privacy Framework Core](https://www.nist.gov/system/files/documents/2021/05/05/NIST-Privacy-Framework-V1.0-Core-PDF.pdf)（访问日期：2026-07-17）
- [OpenAI API：Data controls](https://platform.openai.com/docs/models/default-usage-policies-by-endpoint)（访问日期：2026-07-17）
