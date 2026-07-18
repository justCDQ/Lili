---
title: 标题、页码、来源与原文定位
stage: intermediate
direction: ai
topic: rag-parsing
---

# 标题、页码、来源与原文定位

RAG 的每个解析块必须能回答四个问题：它属于哪份来源、哪个版本、处于什么结构位置、用户如何回到原文。没有这些信息，检索结果无法可靠引用，解析错误也难以调试。标题、页码和 locator 不只是展示字段，而是分块、权限、版本与证据系统的基础。

## 前置知识与边界

- [异构文档格式解析](01-document-formats-and-parsing.md)
- 文档 block、revision 与 artifact 基础。
- [引用、证据视图与无法回答](../ai-ux/06-citations-evidence-no-answer.md)

本文处理解析层的结构与定位。检索后仍需验证主张与证据关系。

## Source identity

来源对象：

```json
{
  "sourceId": "src_31",
  "canonicalUri": "repo://policies/refund.md",
  "title": "退款政策",
  "publisher": "commerce-policy-team",
  "revision": "git:8c14d2a",
  "contentHash": "sha256:...",
  "mediaType": "text/markdown",
  "effectiveAt": "2026-07-01T00:00:00+08:00",
  "ingestedAt": "2026-07-17T09:20:00Z"
}
```

字段不能混用：

- `revision`：来源系统的版本。
- `contentHash`：具体字节内容。
- `effectiveAt`：业务生效时间。
- `ingestedAt`：进入索引时间。
- `updatedAt`：来源元数据更新时间。

同一 URL 的内容会变化，URL 不能代替 revision。

## Block identity

```json
{
  "blockId": "blk_019",
  "sourceId": "src_31",
  "sourceRevision": "git:8c14d2a",
  "type": "paragraph",
  "text": "标准退款期限为 14 天。",
  "headingPath": [
    {"id":"h1","level":1,"text":"退款政策"},
    {"id":"h2-cn","level":2,"text":"中国区"}
  ],
  "ordinal": 19,
  "locator": {
    "type": "text_range",
    "start": 810,
    "end": 823,
    "quoteHash": "sha256:..."
  }
}
```

`blockId` 是该解析 revision 内的稳定标识。内容变化后可创建新 ID，或保留逻辑 ID 同时增加版本；必须定义清楚。

## 标题层级

标题提供：

- 文档主题。
- 章节上下文。
- 分块边界。
- 过滤和导航。
- 引用的可读位置。

### Heading path

一个段落应携带完整路径：

```text
退款政策 > 中国区 > 特殊商品
```

只保存最近标题会丢失上层范围；把路径复制进正文会增加重复 Token。可以在 metadata 中保存，构建上下文时按需要注入。

### 层级跳跃

源文档可能从 h1 跳到 h3。解析器应保留实际 level 并记录 warning，不应无依据改成 h2。清洗层可以建立规范化 outline，但必须保留原级别。

### 视觉标题

PDF 没有 HTML heading 标签。可以依据：

- 字号、字重。
- 前后间距。
- 编号模式。
- 页面位置。
- 目录映射。
- 重复样式。

推断结果保存：

```json
{
  "type": "heading",
  "level": 2,
  "text": "适用范围",
  "derivation": "layout_inference",
  "confidence": 0.87,
  "styleCluster": "font-16-bold"
}
```

confidence 只在同一推断器版本内解释。标题层级要用标注集评估。

## 页码的三种含义

### 物理页索引

文件中的第几个 page object，通常从 1 展示。

### 印刷页码

页面上印出的 “1”“i”“附录 A-3”。封面和目录可能不计入。

### 查看器页码

PDF page labels 可以让查看器显示 `iii` 或 `A-1`。

locator 建议同时保存：

```json
{
  "pageIndex": 12,
  "pageLabel": "10",
  "printedPageText": "10"
}
```

用户点击“第 10 页”时必须明确是哪种。数据库内部使用 pageIndex，界面优先显示 pageLabel。

## PDF 坐标

```json
{
  "type": "pdf_region",
  "pageIndex": 12,
  "bbox": [72.5, 184.2, 510.8, 238.6],
  "coordinateSystem": "pdf_user_space",
  "pageRotation": 0,
  "quoteHash": "sha256:..."
}
```

坐标要声明：

- 原点。
- 单位。
- 页面旋转。
- crop box/media box。
- 是否归一化。

否则前端高亮会错位。跨行证据可使用多个 bbox，不要用一个大框覆盖中间无关文字。

## 文本范围

Markdown、HTML 和纯文本可用字符范围，但要定义单位：

- Unicode code point。
- UTF-16 code unit。
- UTF-8 byte。

JavaScript string offset 通常按 UTF-16 code unit，服务端语言可能不同。推荐协议明确 `offsetEncoding`：

```json
{
  "type": "text_range",
  "start": 810,
  "end": 823,
  "offsetEncoding": "utf-8-bytes",
  "quote": "标准退款期限为 14 天。",
  "quoteHash": "sha256:..."
}
```

使用 quote 与前后文可以在小改动后尝试重新锚定，但不能静默把引用迁移到不同语义。

## DOM 定位

HTML 的 CSS selector 或 XPath 可能随页面改版失效。组合：

```json
{
  "type": "dom_text",
  "headingId": "refund-window",
  "selector": "main article section:nth-of-type(3) p:nth-of-type(2)",
  "textPosition": {"start": 418, "end": 460},
  "quoteHash": "sha256:..."
}
```

优先使用来源提供的稳定 fragment ID。选择器只是辅助。重新抓取后用 quote hash 检查。

## 表格定位

```json
{
  "type": "table_cells",
  "sheet": "政策",
  "range": "B4:D7",
  "tableId": "refund_rules",
  "rowKeys": ["CN", "US"],
  "columnKeys": ["window_days", "status"]
}
```

只保存 `B4` 对插入行很脆弱。稳定表格 ID、主键和 revision 共同定位。PDF 表格则保存 page、cell bbox 与逻辑 row/column。

## 音视频与图像

- 音视频：起止时间码、转写 revision、speaker ID。
- 图像：图像 ID、bbox、页码、视觉描述 version。
- 幻灯片：slide index、shape ID、bbox。

模型生成的图像描述是派生内容，应关联原图 region，并标记 derivation。

## Locator version

```json
{
  "locatorVersion": 2,
  "parserVersion": "pdf-parser-4.1",
  "sourceRevision": "sha256:..."
}
```

解析器升级可能改变 block 与坐标。引用绑定旧 locator 时：

- 保留旧解析 artifact。
- 或运行迁移并验证 quote hash。
- 或标记 unavailable/stale。

不能让新解析结果无说明替换旧引用。

## 原文查看器

点击引用时：

1. 服务端验证当前用户有权查看 source revision。
2. 返回受控内容或短期 URL。
3. 查看器打开 page/heading/sheet。
4. 高亮 locator。
5. 展示前后上下文。
6. 若 hash 不匹配，显示“来源已变化”。

查看器不能只在前端隐藏无权限按钮；下载和 range 请求都要授权。

## Heading 与 chunk

分块时可以把 heading path 作为 metadata：

```json
{
  "chunkId": "chunk_8",
  "blockIds": ["blk_19", "blk_20"],
  "headingPath": ["退款政策", "中国区"],
  "sourceLocator": {
    "firstBlockId": "blk_19",
    "lastBlockId": "blk_20"
  }
}
```

检索展示时可生成：

```text
退款政策 / 中国区
标准退款期限为 14 天……
```

heading 本身也可以单独索引，但要避免在每个 chunk 重复造成 heading 关键词占比过高。

## 跨页段落

一个段落跨两页时保存多个 region：

```json
{
  "type": "multi_region",
  "regions": [
    {"pageIndex":8,"bbox":[72,700,520,760]},
    {"pageIndex":9,"bbox":[72,80,520,132]}
  ]
}
```

引用界面可打开第一页并提示跨页。不要因为分页把一句话拆成两个无上下文 chunk。

## 脚注

脚注需要双向关系：

```json
{
  "referenceBlockId": "blk_21",
  "noteBlockId": "footnote_3",
  "marker": "3"
}
```

正文检索命中时可以按预算附加相关脚注。不能把整页所有脚注拼到每个段落。

## 图注

图注与图像关联：

```json
{
  "figureId": "fig_4",
  "imageLocator": {"pageIndex":5,"bbox":[80,160,520,460]},
  "captionBlockId": "caption_4",
  "altDerivationId": "vision_18"
}
```

图注是原文，视觉描述是派生内容，两者不可混同。

## 来源链

同一文档可能经历：

```text
DOCX source
→ rendered PDF
→ OCR fallback
→ normalized blocks
→ chunk
→ retrieved evidence
→ answer citation
```

每层保存 `wasDerivedFrom`：

```json
{
  "entity": "chunk_8",
  "wasDerivedFrom": ["block_19", "block_20"],
  "source": "src_31",
  "sourceRevision": "8c14d2a"
}
```

调试时才能知道错误来自原文件、渲染、OCR、清洗还是分块。

## 完整案例一：PDF 政策问答

### 输入

80 页政策 PDF，封面后使用罗马数字目录，正文 pageIndex 7 显示印刷页 “1”，双栏且包含脚注。

### 解析

1. 保存 pageIndex 与 pageLabel。
2. 从目录和样式推断 heading hierarchy。
3. 双栏区域分别建立阅读顺序。
4. 跨栏脚注建立 note relation。
5. 每个 block 保存多个 bbox 和 quote hash。
6. chunk 保存 heading path 与 block 范围。

### 回答

用户问中国区退款期限。命中 block：

```json
{
  "headingPath": ["退款政策", "区域规则", "中国区"],
  "pageIndex": 19,
  "pageLabel": "13",
  "blockId": "blk_188"
}
```

界面显示“第 13 页”，内部打开 pageIndex 19 并高亮两个 bbox。

### 验证

- 50 个引用全部打开正确页。
- 高亮与原文字句一致。
- 双栏阅读顺序没有交叉。
- 脚注引用打开对应 note。
- page label 与查看器一致。

### 失败分支

解析器升级使 bbox 改变。旧回答仍使用旧 parse artifact；后台迁移只有在 quote hash 和上下文验证后才更新 locator，否则标为需要重新验证。

## 完整案例二：版本化技术文档

### 输入

Git 仓库中 Markdown API 文档，多个 release branch。标题 ID 在 v3 改名。

### 处理

1. source revision 使用 commit hash。
2. block 保存 path、heading path、UTF-8 byte range。
3. 内部链接解析到同 commit 的目标文件。
4. 每个发布文档建立 version metadata。
5. 用户问题包含产品 v2 时，metadata filter 只选 v2。

### 引用

引用固定 commit、path 与 line/offset。打开时展示 v2 文件，不自动跳到 main。

### 失败分支

用户问 v2，但检索只找到 v3。系统返回无相关结果或提示版本范围，不用 v3 回答后伪装成 v2。

### 验证

- commit 变化后旧链接仍可解析。
- 删除 heading 不影响历史 revision。
- UTF-8 offset 对中文和 emoji 正确。
- 跨文件链接不越过仓库权限。

## Locator 验证

批量检查：

- source/revision 存在。
- block range 在内容范围内。
- quote hash 匹配。
- page index 有效。
- bbox 在 page box 内。
- heading ID 唯一或有消歧。
- table range 有效。
- 用户可访问。

失败原因使用枚举：

```text
source_missing
revision_missing
locator_out_of_range
quote_mismatch
permission_denied
parser_artifact_expired
```

## 性能取舍

- 保存完整 bbox 增加存储，但支持高亮和布局调试。
- 只保存页码简单，但无法定位具体证据。
- 保存 quote 便于展示，但可能包含敏感文本；访问和日志要控制。
- 固定所有历史 parse artifact 可重现，但需要保留策略。
- locator 压缩不能丢失 offset unit 与版本。

## 可观测性

- heading 层级准确率。
- locator 生成覆盖率。
- 点击定位成功率。
- quote hash mismatch。
- 来源 revision 缺失。
- 页面/字符/表格 locator 分布。
- 解析升级后的迁移成功率。
- 权限拒绝率。
- 用户手动寻找原文的时间。

## 常见错误

### URL 等于来源版本

URL 内容会更新。保存 revision 和 hash。

### 所有页码都从 1 顺序显示

区分 pageIndex 与 pageLabel。

### 只保存字符 offset

声明 offset encoding，并保存 quote/hash。

### CSS selector 永久稳定

组合 stable ID、text position、quote 和 revision。

### Chunk 自己生成引用位置

引用应从组成 block 的 locator 聚合，不凭 chunk 文本在来源中模糊搜索。

### 查看器前端做权限

source、artifact 和短期 URL 都在服务端授权。

## 生产验收清单

- [ ] 每个 source 有 canonical ID、revision、hash 和媒体类型。
- [ ] 每个 block 有 source revision、类型、ordinal 和 locator。
- [ ] heading path 保留完整层级。
- [ ] PDF 区分 page index、label 与 printed text。
- [ ] bbox 声明坐标系、单位和旋转。
- [ ] text range 声明 offset encoding。
- [ ] HTML、表格、音视频使用适合的 locator。
- [ ] locator 关联 parser version。
- [ ] 跨页、脚注和图注保存关系。
- [ ] 点击引用重新授权并校验 hash。
- [ ] 来源更新不静默迁移旧引用。
- [ ] locator 有批量回归和可观测性。

## 集成练习

为 PDF、Markdown 和 XLSX 实现统一 locator：

1. PDF 使用 pageIndex、pageLabel、多 bbox 和 quote hash。
2. Markdown 使用 commit、path、heading path 和 UTF-8 byte range。
3. XLSX 使用 workbook revision、sheet、range 与 row key。
4. 每种格式生成 50 个标注证据。
5. 查看器点击后定位成功率达到团队设定门槛。
6. 修改来源后，旧 locator 要么打开历史 revision，要么明确 stale。
7. 解析器升级后运行迁移验证，不匹配项进入复核。
8. 跨租户 locator 请求不泄露 source 是否存在。

## 来源

- [W3C Web Annotation Data Model](https://www.w3.org/TR/annotation-model/)（访问日期：2026-07-17）
- [W3C PROV-O：The PROV Ontology](https://www.w3.org/TR/prov-o/)（访问日期：2026-07-17）
- [ISO：PDF 2.0（ISO 32000-2:2020）](https://www.iso.org/standard/75839.html)（访问日期：2026-07-17）
- [CommonMark Specification](https://spec.commonmark.org/)（访问日期：2026-07-17）
- [RFC 5147：URI Fragment Identifiers for text/plain](https://www.rfc-editor.org/rfc/rfc5147)（访问日期：2026-07-17）
