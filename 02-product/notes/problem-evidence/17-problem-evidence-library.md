# 问题证据库：来源、版本、权限与可追溯性

问题证据库把原始证据、可证伪主张、分析判断、冲突、决策和后续验证连接起来。它不是功能愿望池，也不是网页全文备份。目标是让任何结论能追溯到来源和条件，让新证据能定位受影响的判断，并在产品、政策或权限变化后重新评估。

## 一、分离四种对象

### Evidence：证据

可定位的观察、文档声明、行为统计、日志结果或复现记录。证据保留原始条件，不直接等于结论。

### Claim：主张

由证据支持或反对的可证伪陈述。一个问题可有“存在”“范围”“原因”“后果”多个主张。

### Problem：问题

目标角色在特定任务与条件下遭受的可观察损失。问题引用主张，不复制所有证据正文。

### Decision：决策

保留、删除、改进、验证或不行动的记录，包含当时证据、风险、负责人、回滚和复审。

分离后，同一证据可支持多个主张；主张被推翻不会删除原始记录；决策可以解释“当时为什么合理”。

## 二、目录结构

适合 Obsidian、GitHub 和 VS Code 的最小结构：

```text
evidence-library/
├── README.md
├── schema/
│   ├── evidence.schema.json
│   ├── claim.schema.json
│   ├── problem.schema.json
│   └── decision.schema.json
├── evidence/
│   └── 2026/07/EV-202607-0001.md
├── claims/
│   └── CLM-IMPORT-0001.md
├── problems/
│   └── PRB-IMPORT-0001.md
├── decisions/
│   └── DEC-202607-0001.md
├── reviews/
│   └── 2026-07.md
└── attachments-private/
    └── README.md
```

公共仓库不存敏感附件。`attachments-private` 可只放说明，真实受限证据留在获批系统，用权限内指针引用。

## 三、稳定 ID 与不可覆盖历史

ID 不使用标题或顺序号作为唯一语义：

```text
EV-202607-0001  证据
CLM-IMPORT-0001 主张
PRB-IMPORT-0001 问题
DEC-202607-0001 决策
```

标题可修改，ID 不变。合并同义问题时设置 `merged_into`，不删除旧记录；状态变化追加历史，不覆盖旧值。纠正错误时记录谁、何时、为什么改，以及原值。

## 四、可用的 Evidence schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.org/schema/evidence.schema.json",
  "title": "Problem Evidence Record",
  "type": "object",
  "required": [
    "id",
    "title",
    "source",
    "dates",
    "scope",
    "content",
    "permissions",
    "quality",
    "status",
    "review"
  ],
  "properties": {
    "id": {"type": "string", "pattern": "^EV-[0-9]{6}-[0-9]{4}$"},
    "title": {"type": "string", "minLength": 1},
    "source": {
      "type": "object",
      "required": ["type", "origin", "locator", "independence_group"],
      "properties": {
        "type": {
          "enum": [
            "public_review", "issue", "discussion", "official_document",
            "status_incident", "workflow_observation", "dogfooding",
            "support_secondary", "behavioral_analysis", "system_log",
            "controlled_test", "dataset", "policy"
          ]
        },
        "origin": {"type": "string"},
        "locator": {"type": "string"},
        "archived_locator": {"type": ["string", "null"]},
        "independence_group": {"type": "string"},
        "author_role": {"type": ["string", "null"]},
        "license_or_terms": {"type": ["string", "null"]}
      },
      "additionalProperties": false
    },
    "dates": {
      "type": "object",
      "required": ["occurred_or_published", "captured_at"],
      "properties": {
        "occurred_or_published": {"type": ["string", "null"], "format": "date-time"},
        "captured_at": {"type": "string", "format": "date-time"},
        "source_updated_at": {"type": ["string", "null"], "format": "date-time"},
        "valid_from": {"type": ["string", "null"], "format": "date"},
        "valid_until": {"type": ["string", "null"], "format": "date"}
      },
      "additionalProperties": false
    },
    "scope": {
      "type": "object",
      "required": ["task", "product_version", "environment"],
      "properties": {
        "role": {"type": ["string", "null"]},
        "task": {"type": "string"},
        "product_version": {"type": ["string", "null"]},
        "plan_or_tier": {"type": ["string", "null"]},
        "region": {"type": ["string", "null"]},
        "environment": {"type": "object"},
        "sample": {"type": ["object", "null"]}
      },
      "additionalProperties": false
    },
    "content": {
      "type": "object",
      "required": ["fact", "summary", "limitations"],
      "properties": {
        "fact": {"type": "string"},
        "summary": {"type": "string"},
        "minimal_excerpt": {"type": ["string", "null"]},
        "limitations": {"type": "array", "items": {"type": "string"}},
        "counter_observations": {"type": "array", "items": {"type": "string"}}
      },
      "additionalProperties": false
    },
    "permissions": {
      "type": "object",
      "required": ["classification", "allowed_use", "contains_personal_data"],
      "properties": {
        "classification": {"enum": ["public", "internal", "restricted"]},
        "allowed_use": {"type": "array", "items": {"type": "string"}},
        "contains_personal_data": {"type": "boolean"},
        "redaction_applied": {"type": "array", "items": {"type": "string"}},
        "retention_until": {"type": ["string", "null"], "format": "date"},
        "owner": {"type": "string"},
        "withdrawal_or_deletion": {"type": ["string", "null"]}
      },
      "additionalProperties": false
    },
    "quality": {
      "type": "object",
      "required": ["directness", "reproducibility", "coverage", "freshness"],
      "properties": {
        "directness": {"enum": ["direct", "indirect", "unknown"]},
        "reproducibility": {"enum": ["reproduced", "reproducible_steps", "not_reproduced", "not_applicable"]},
        "coverage": {"type": "string"},
        "freshness": {"enum": ["current", "review_due", "expired"]},
        "data_quality_checks": {"type": "array", "items": {"type": "string"}}
      },
      "additionalProperties": false
    },
    "relationships": {
      "type": "object",
      "properties": {
        "supports_claims": {"type": "array", "items": {"type": "string"}},
        "counters_claims": {"type": "array", "items": {"type": "string"}},
        "conflicts_with": {"type": "array", "items": {"type": "string"}},
        "derived_from": {"type": "array", "items": {"type": "string"}},
        "duplicates": {"type": "array", "items": {"type": "string"}}
      },
      "additionalProperties": false
    },
    "status": {"enum": ["active", "superseded", "expired", "withdrawn", "deleted_source"]},
    "review": {
      "type": "object",
      "required": ["last_reviewed_at", "next_review_at", "reviewer"],
      "properties": {
        "last_reviewed_at": {"type": "string", "format": "date"},
        "next_review_at": {"type": "string", "format": "date"},
        "reviewer": {"type": "string"},
        "review_reason": {"type": "string"}
      },
      "additionalProperties": false
    },
    "history": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["changed_at", "changed_by", "field", "from", "to", "reason"],
        "properties": {
          "changed_at": {"type": "string", "format": "date-time"},
          "changed_by": {"type": "string"},
          "field": {"type": "string"},
          "from": {},
          "to": {},
          "reason": {"type": "string"}
        }
      }
    }
  },
  "additionalProperties": false
}
```

实际项目可拆小，但不能删掉来源、时间、范围、权限、质量、关系、状态和复审这些核心层。

## 五、Claim schema 的关键字段

```json
{
  "id": "CLM-IMPORT-0001",
  "statement": "Windows桌面端v3.2导入空日期CSV的失败率高于v3.1",
  "claim_type": "problem_frequency",
  "scope": {
    "role": "管理员",
    "task": "导入CSV",
    "versions": ["3.1", "3.2"],
    "platform": "Windows",
    "time_window": "2026-07-01/2026-07-14"
  },
  "falsifier": "固定输入与环境的版本对照无差异",
  "supports": ["EV-202607-0001", "EV-202607-0002"],
  "counters": ["EV-202607-0003"],
  "conflicts": [],
  "facts": ["v3.2观察失败率6%，v3.1为1.2%"],
  "inference": "版本变化与失败上升相关",
  "hypotheses": ["日期规范化顺序改变"],
  "confidence": {
    "level": "medium",
    "reasons": ["行为数据与受控复现一致"],
    "limits": ["只覆盖Windows"]
  },
  "status": "supported_with_limits",
  "review_on": "2026-07-31",
  "history": []
}
```

事实、推断和假设不写在同一字符串里。`supports` 与 `counters` 都不可为空时静默删除。

## 六、Problem 与 Decision

问题记录包含：目标任务、可观察损失、当前工作流、替代方案、频率/严重度/范围/可信度、相关主张、未决问题和状态。

决策记录包含：

- 决策对象与类型；
- 决策日期、负责人和参与角色；
- 当时引用的 claim/evidence 版本；
- 选择与未选择方案；
- 预期结果、主指标与守护；
- 风险、停止、回滚和复审；
- 后续结果和是否需要撤销。

证据后来更新，不自动改写历史决策；而是触发复审并追加新决策。

## 七、来源和时间的双时间模型

至少区分：

- `occurred_or_published`：事件发生或页面发布；
- `captured_at`：研究者获得证据；
- `source_updated_at`：来源后来编辑；
- `valid_from/until`：规则或结论适用；
- `last_reviewed/next_review`：证据库维护。

网页访问日期不能代替事件发生时间。未知日期保留 null 并写精度，不猜测。

## 八、权限、脱敏和删除

来源公开不等于可以复制全文。记录许可或条款，只保存最小摘录和自己的摘要。内部或受限材料只保留获批系统指针，不进入公共仓库。

权限状态变化时：

1. 标记证据 `withdrawn` 或更新允许用途；
2. 删除无权保留的副本和派生产物；
3. 评估依赖主张是否仍能由其他证据支持；
4. 追加历史，不保留被要求删除的敏感原值；
5. 触发相关决策复审。

数据删除权优先于“完整历史”的技术偏好。

## 九、反证与冲突关系

反证是相对于主张的证据，不是“坏数据”。冲突记录包括：双方 ID、冲突字段、是否同口径、可能解释、区分测试、负责人和状态。

```text
CONFLICT-07
claim: CLM-IMPORT-0001
evidence_a: EV-202607-0001（行为失败上升）
evidence_b: EV-202607-0009（总体支持工单稳定）
alignment: 人群与分母不同，暂不可直接比较
next_test: 按版本筛选工单并核对支持入口变化
status: open
```

不要为了保持“可信度高”删除不一致。

## 十、过期、复审与状态

触发过期的条件：

- 产品、套餐、政策或数据 schema 变化；
- 来源删除或重要编辑；
- 证据超出预定时效；
- 关键反证出现；
- 权限和保留期限变化；
- 决策结果与预期相反。

状态建议：`hypothesis`、`collecting`、`supported_with_limits`、`contradicted`、`resolved_for_version`、`expired`、`merged`、`archived`。`resolved_for_version` 比永久 `resolved` 更准确。

## 十一、完整案例：批量导入无失败清单

### 1. 新建问题

问题 PRB-IMPORT-0001：管理员批量导入部分失败时，只看到总数，无法定位失败行，导致人工逐行比较。

### 2. 输入证据

- EV-0001：公开评论，v3.1，用户自述；
- EV-0002：Issue 最小复现，维护者确认；
- EV-0003：行为分析，部分失败后 30 分钟内重复导入增加；
- EV-0004：官方文档，声明只显示成功/失败总数；
- EV-0005：反证，小文件全成功路径无需清单；
- EV-0006：受控任务，100行中3行失败后人工定位耗时18分钟。

每条有来源、版本、时间、权限和独立组。

### 3. 拆分主张

CLM-01：v3.1 部分失败时没有行级清单。高可信。

CLM-02：缺少清单造成重复导入。中可信，行为关联但意图未知。

CLM-03：所有导入都需要行级清单。被 EV-0005 限定，只在部分失败时成立。

### 4. 冲突与推理

支持工单数量不高，但支持入口仅企业套餐可用，不与行为重复直接冲突。保留该限制。方案请求“增加下载按钮”改写为需求：部分失败后，管理员需要定位每条失败及原因，并能安全重试失败子集。

### 5. 决策

DEC-0001：先原型验证行级错误文件和失败子集重试；主指标为部分失败任务 24 小时完成率，守护为重复写入和数据泄露；5% 灰度，严重权限事件立即回滚。

### 6. 新版本复审

v3.2 发布错误文件后：

- 旧证据不删除；
- PRB 状态改为 `resolved_for_version: 3.2`；
- 新建 EV 记录回归与行为结果；
- v3.1 自托管用户仍保持 active；
- 设置 60 天后复查支持和重复导入。

### 7. 失败分支

- 若原始评论删除，标记 `deleted_source`，不保留超许可全文；
- 若行为事件 schema 发现重复计数，撤回 EV-0003 并重新评估 CLM-02；
- 若错误文件包含敏感原始值，停止灰度并修复权限/最小化；
- 若 v3.2 只对云端发布，自托管范围不能标解决；
- 若反证说明小文件无需清单，保持条件化需求；
- 若 schema 迁移失败，保留备份并验证 ID 关系完整。

## 十二、维护工作流

### 新证据进入

检查重复、权限、最小化、来源和时间；创建 EV；关联或新建 Claim；不直接新建功能。

### 每周

处理冲突、链接失效和权限到期；审查新增高风险证据。

### 每月

复审到期主张、问题状态和未决测试；验证随机记录能回到来源。

### 发布后

按版本创建新证据，更新范围化状态，比较决策指标和守护。

### schema 变化

版本化 schema，提供迁移脚本和变更说明；先在副本验证，检查必填、枚举、关系和历史，再替换。记录 schema_version。

## 十三、Obsidian、GitHub 与 VS Code

Markdown 正文便于解释，YAML frontmatter 或 JSON 便于校验。稳定 ID 可用普通 Markdown 链接，避免只有 Obsidian 能解析的专有嵌入。Git 提交记录不能替代字段内的业务变更原因，但可以提供审计补充。

在 CI 中检查：schema、唯一 ID、引用存在、来源日期、复审到期、公共记录不得含 restricted 分类、无孤立 Claim、决策必须引用证据版本。

## 十四、完成检查与练习

- Evidence、Claim、Problem、Decision 分离；
- ID 稳定，合并不删除历史；
- 来源、事件/采集日期、版本环境齐全；
- 权限、敏感级别、最小化、保留和删除可执行；
- 事实、推断、假设分字段；
- 支持、反证、冲突和独立组可查询；
- 可信度含理由和限制；
- 过期、复审、状态和历史存在；
- 决策引用当时证据并有回滚；
- 公共仓库不包含受限原文。

练习：按本章 schema 建立 10 条 Evidence、3 条 Claim、1 条 Problem 和 1 条 Decision。加入一条反证、一次来源删除和一次版本解决。完成标准是从决策能追溯到原始证据，从证据能找到权限与复审，并能在不覆盖历史的情况下处理冲突。

## 来源

- [W3C：PROV Overview](https://www.w3.org/TR/prov-overview/)（访问日期：2026-07-17）
- [JSON Schema：Draft 2020-12](https://json-schema.org/draft/2020-12)（访问日期：2026-07-17）
- [GitHub Docs：About issue and pull request templates](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates)（访问日期：2026-07-17）
- [ICO：A guide to the data protection principles](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/data-protection-principles/a-guide-to-the-data-protection-principles/)（访问日期：2026-07-17）
