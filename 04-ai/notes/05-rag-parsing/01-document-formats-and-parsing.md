---
title: PDF、Word、Markdown、HTML、表格与扫描件解析
stage: intermediate
direction: ai
topic: rag-parsing
---

# PDF、Word、Markdown、HTML、表格与扫描件解析

RAG 文档解析的目标不是“得到一些文字”，而是把异构文件转换为可追溯的结构化文档：正文、标题层级、列表、表格、图片、页码、来源和定位关系。不同格式保存的信息不同，因此不能用同一个纯文本提取器处理所有文件。

## 前置知识与范围

- 文件 MIME、字符编码与压缩包基础。
- HTML DOM、Markdown 和表格基础。
- RAG 的解析→清洗→分块流程。

本文说明六类文件的解析策略、验证与失败边界。具体库只是实现选择，必须用自己的文档集评估。

## 统一中间表示

```json
{
  "documentId": "doc_91",
  "revision": "sha256:...",
  "mediaType": "application/pdf",
  "parser": {
    "name": "parser_service",
    "version": "2026.07.3"
  },
  "blocks": [
    {
      "id": "block_1",
      "type": "heading",
      "level": 1,
      "text": "退款政策",
      "locator": {
        "page": 1,
        "bbox": [72, 86, 430, 112]
      }
    }
  ],
  "warnings": [],
  "quality": {}
}
```

核心字段：

- `documentId/revision`：区分来源与内容版本。
- `mediaType`：不要只信文件扩展名。
- `parser/version`：支持回归和重新解析。
- `blocks`：保留类型、顺序与层级。
- `locator`：能回到原文。
- `warnings`：显式保存丢失、OCR 和不支持内容。
- `quality`：按格式记录可计算指标。

## 解析管线

```mermaid
flowchart LR
    U["上传"] --> V["类型、大小、恶意文件检查"]
    V --> R["格式路由"]
    R --> P["格式解析器"]
    P --> N["规范化中间表示"]
    N --> Q["质量检查"]
    Q --> C["清洗与去重"]
    C --> K["分块与索引"]
```

解析器的输出在进入索引前必须经过质量门。失败文件不能静默产生空文档。

## 输入检查

处理前检查：

- 文件大小、页数、解压后大小和压缩层级。
- 声明扩展名、HTTP Content-Type 与魔数。
- 密码保护或加密状态。
- 宏、嵌入对象和外部关系。
- 恶意软件扫描结果。
- 租户、来源与访问权限。
- 上传是否完整，计算内容 hash。

解析服务应在隔离进程中运行，并设置 CPU、内存、时间和输出上限。文档解析器处理攻击者可控的复杂二进制文件，不能与主业务进程共享无限权限。

## PDF

### PDF 保存的是绘制指令

PDF 页面由文本、字体、图形和图像对象组成，不必包含语义阅读顺序。屏幕上相邻的字可能在内容流中分散；多栏页面可能按错误顺序提取。

### 三类 PDF

- 文本型：含可选择文字。
- 扫描型：页面主要是图像，需要 OCR。
- 混合型：文字层、扫描图和注释并存。

### 文本块重建

解析器通常依据 glyph 坐标组合：

```text
glyph -> word -> line -> paragraph -> region -> page
```

需要处理：

- 字符间距与词间空格。
- 行基线与旋转。
- 双栏和浮动侧栏。
- 连字符断行。
- 脚注与页眉页脚。
- 字体映射缺失。
- ligature 和不可见文字层。

### 阅读顺序

按 y、x 坐标排序只适合简单单栏。多栏页面需要布局区域识别。验证样本应包含：

- 单栏。
- 双栏。
- 表格跨页。
- 图注。
- 竖排或旋转文本。
- 脚注。

### PDF 表格

表格可能没有真实单元格对象。解析器依据线条、对齐和文本位置推断。结果应保留：

```json
{
  "type": "table",
  "page": 4,
  "columns": 3,
  "rows": [
    {
      "cells": [
        {"text":"地区","rowSpan":1,"colSpan":1},
        {"text":"期限","rowSpan":1,"colSpan":1},
        {"text":"备注","rowSpan":1,"colSpan":1}
      ]
    }
  ]
}
```

若结构置信不足，可同时保存表格区域图像与线性文本，标记 `table_structure_uncertain`。

### PDF 安全

不要执行：

- JavaScript actions。
- 附件。
- 外部链接自动请求。
- 嵌入媒体。

解析所需字体和对象也要限制资源使用。

## Word

现代 `.docx` 通常是 ZIP 包装的 Office Open XML 部件。正文、样式、关系、图片、批注和页眉页脚分布在不同 XML 文件。

### 应保留

- 段落与 run 文本。
- heading 样式层级。
- 编号和项目符号列表。
- 表格、合并单元格。
- 超链接目标。
- 图片及图注。
- 脚注、尾注。
- 页眉页脚。
- 批注和修订状态，按产品策略。

### 样式不是视觉字体大小

标题优先从段落样式和 outline level 获取。只按字号猜标题会把封面大字、强调段落误判。

### 修订

文档可能包含插入、删除与修订。产品必须选择：

- 接受修订后的视图。
- 原始视图。
- 同时保存修订关系。

不能把删除文本和当前文本无标记拼接。

### 页码问题

Word 是流式排版格式，页码取决于字体、纸张、打印机和渲染环境。仅解析 XML 通常无法得到与用户 Word 客户端完全相同的页码。需要页码时：

1. 使用固定渲染环境转为 PDF。
2. 记录渲染器、字体和纸张配置。
3. 从渲染 PDF 获取页码定位。

不要凭段落数伪造页码。

### 外部关系

DOCX 可引用外部图片或链接。解析时默认不请求外部资源，避免 SSRF、隐私泄露和不可重现结果。

## Markdown

Markdown 是纯文本，但方言影响结构：

- CommonMark。
- GitHub Flavored Markdown。
- Obsidian 扩展。
- Frontmatter。
- 数学、Mermaid 或自定义容器。

### 解析为 AST

不要只按 `#` 和空行切分。使用与目标方言一致的 parser，保留：

- heading depth。
- paragraph。
- list 与嵌套。
- blockquote。
- fenced code 与 info string。
- link、image。
- table 扩展。
- source range。

### Frontmatter

YAML frontmatter 是文档元数据，不一定进入正文：

```yaml
---
title: 退款政策
effective_date: 2026-07-01
region: CN
---
```

字段要按 Schema 允许列表解析，不能让任意 frontmatter 覆盖 tenant、权限或系统字段。

### 链接

相对链接根据文档所在仓库路径解析，保存：

- 原始 target。
- 解析后的内部资源 ID。
- revision。
- 是否存在。

不要在离线索引时自动抓取任意外链。

### 代码块

代码内容不应当作普通说明文字清洗。保存语言、原始换行和行号范围。检索时可以建立代码专用字段或索引。

## HTML

### 先解析 DOM

HTML 不能用正则可靠提取结构。解析器按 HTML 规范构建 DOM，然后选择正文。

### 正文与页面框架

常见非正文：

- navigation。
- cookie banner。
- footer。
- related articles。
- 评论。
- 广告。
- 隐藏模板。

可以使用语义元素、站点规则和正文密度，但必须在目标站点评估。通用 readability 算法可能删除 API 参数表和侧栏警告。

### 保留结构

将：

- `h1`–`h6` 转为 heading。
- `p` 转为 paragraph。
- `ul/ol` 保留列表。
- `table` 保留 cell。
- `pre/code` 保留代码。
- `figure/figcaption` 关联。
- `a` 保存目标和 anchor text。

### 动态页面

HTTP 返回的 HTML 可能没有 JavaScript 渲染后的正文。选择：

- 使用网站提供的结构化 API。
- 用受控浏览器渲染。
- 抓取静态页面。

浏览器渲染成本高且风险更大，要限制网络来源、脚本、下载和执行时间。登录态内容必须遵守权限和数据使用规则。

### HTML 安全

索引文本不能把页面中的指令当系统指令。原始 HTML 渲染到调试页面前要 sanitization。外部资源默认不加载。

## 表格文件

包括 CSV、TSV、XLSX 和导出表。

### CSV/TSV

需要确定：

- 字符编码。
- 分隔符。
- 引号与转义。
- 换行。
- 是否有表头。
- 空值表示。
- 日期和数字 locale。

CSV 是记录格式，不自带可靠类型。`00123` 可能是编号，不应自动变成 123；日期 `01/02/2026` 含糊。

### XLSX

需要处理：

- workbook 与 sheet。
- cell value 与显示格式。
- 公式及缓存结果。
- 合并单元格。
- 隐藏行列和 sheet。
- named range。
- comments。
- 外部链接。

产品要决定索引公式字符串、计算结果还是两者。解析器通常不应执行不可信宏或自动刷新外部链接。

### 表格到文本

不要简单逐行拼接：

```text
CN 14 active
US 30 active
```

应关联表头：

```text
地区: CN | 退款天数: 14 | 状态: active
地区: US | 退款天数: 30 | 状态: active
```

同时保留 sheet、row 和 cell locator，回答可以回到单元格。

### 宽表和大表

- 只选业务相关列。
- 保留主键。
- 分页或按行组 chunk。
- 汇总信息与明细分开。
- 限制公式、样式和空白区域造成的虚假 used range。

## 扫描件与 OCR

### OCR 管线

```text
page image
→ orientation detection
→ deskew/denoise
→ layout detection
→ text recognition
→ reading order
→ language/post-processing
→ confidence and locator
```

### OCR 不能只保存字符串

保存：

```json
{
  "text": "退款期限为14天",
  "page": 2,
  "bbox": [118, 240, 508, 271],
  "ocr": {
    "engine": "ocr_4.2",
    "language": ["zh-Hans"],
    "confidence": 0.91
  }
}
```

confidence 的尺度由引擎定义，不能跨引擎直接比较。关键数字、姓名和 ID 可使用规则或人工抽检。

### 常见 OCR 错误

- `0/O`、`1/l/I`。
- 小数点和负号。
- 中文相似字。
- 表格列错位。
- 页边批注进入正文。
- 图章遮挡。
- 旋转、模糊和低对比。
- 手写内容。

不要用语言模型静默“纠正”后覆盖原 OCR。可以保存候选修正及来源，让关键字段经确定性或人工验证。

### 文本层与 OCR 冲突

扫描 PDF 可能已有质量差的隐藏文字层。策略：

1. 抽样比较文字层与页面图像 OCR。
2. 按页选择更高质量结果。
3. 保存采用的路径与质量指标。
4. 不把两层文本重复拼接。

## 多模态元素

图表、流程图和截图中的信息不会自动出现在文本提取中。选择：

- 提取图注和 alt。
- 保存图片 region 与页码。
- 对重要图像运行视觉模型或专用图表解析。
- 结果标记为派生描述。
- 与原图 locator 关联。

视觉描述可能遗漏数值，不应在高精度数据问答中替代表格源数据。

## 规范化

跨格式统一：

- Unicode 正规化策略。
- 换行。
- 空白。
- heading level。
- list item。
- table cell。
- language tag。
- source locator。

不要过早：

- 全部转小写。
- 删除标点。
- 合并所有换行。
- 删除代码缩进。
- 把表格变成无表头词袋。

这些操作会丢失后续检索与引用所需结构。

## 完整案例一：产品手册知识库

### 输入

- 120 页双栏 PDF。
- 同内容 DOCX 源文件。
- Markdown 更新日志。
- XLSX 参数表。

### 处理

1. 计算 hash 和访问范围。
2. DOCX 提取标题层级与表格。
3. PDF 提供稳定页码和坐标。
4. 建立 DOCX block 与 PDF page locator 的映射。
5. Markdown 按 release heading 解析。
6. XLSX 按型号主键转结构化行。
7. 去重时保留权威来源优先级，不删除定位。
8. 对型号、单位与版本做一致性检查。

### 输出

回答中的功能说明引用 PDF 页码；参数值同时引用 XLSX sheet/cell。若 DOCX 与 PDF revision 不同，标记冲突，不把两者拼成一条事实。

### 验证

- 抽查 30 个标题层级。
- 双栏页面阅读顺序正确。
- 20 个参数单元格值和单位一致。
- 版本更新日志可定位 heading。
- 相同正文没有重复索引两次。

### 失败分支

PDF 的某些字体无 Unicode 映射，提取为乱码。质量门阻止该页进入索引，尝试 OCR 或使用 DOCX 文本，同时保留 PDF page 作为引用定位。

## 完整案例二：财务扫描件

### 输入

扫描 PDF 包含发票、表格、印章和手写批注。

### 处理

1. 每页分类：发票、说明、空白。
2. OCR 中英文与数字。
3. 用布局模型识别字段和表格。
4. 金额字段通过币种、格式和合计规则检查。
5. 原图 bbox 与 OCR 值关联。
6. 低 confidence 或合计不一致进入人工复核。
7. 只有复核值进入业务数据库；原 OCR 可进入受限检索。

### 验证

- 发票号逐字符准确率。
- 金额 exact match。
- 行项目结构准确率。
- 空白页不产生 chunk。
- 人工修正保留原值、修正值和操作者。

### 失败分支

OCR 把 `8,000.00` 读成 `3,000.00`，但合计校验失败。系统不能让语言模型按上下文猜测金额，应显示页面区域并请求复核。

## 格式路由

```javascript
export function selectParser(file) {
  const routes = new Map([
    ["application/pdf", "pdf"],
    [
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "docx"
    ],
    ["text/markdown", "markdown"],
    ["text/html", "html"],
    ["text/csv", "csv"],
    [
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "xlsx"
    ]
  ]);

  const parser = routes.get(file.detectedMediaType);
  if (!parser) {
    return { accepted: false, reason: "unsupported_media_type" };
  }

  return { accepted: true, parser };
}
```

路由后仍要在解析器内部验证文件签名和结构。

## 质量指标

按格式记录：

- 解析成功率。
- 空文本页/块比例。
- 字符乱码比例。
- heading、list、table 结构准确率。
- 阅读顺序错误率。
- OCR 字符/词错误率。
- locator 打开成功率。
- 警告与人工复核率。
- 处理延迟、峰值内存和成本。

不要只记录“文件已完成”。

## 常见错误

### 所有文件 `extractText()`

丢失结构与定位。使用格式路由和统一中间表示。

### PDF 文本顺序等于阅读顺序

对多栏、侧栏和图注抽样验证。

### DOCX 页码来自 XML

页码依赖渲染。需要固定渲染到 PDF。

### HTML 用正则去标签

会错误处理脚本、实体和嵌套。使用 HTML parser 和正文选择。

### CSV 自动类型推断

会破坏前导零、日期和大整数。使用列 Schema。

### OCR 修正覆盖原值

保留原 OCR、候选修正、证据图像和验证状态。

## 生产验收清单

- [ ] 类型路由基于检测结果，不只扩展名。
- [ ] 解析进程有资源隔离和限制。
- [ ] 输出统一 block、locator、revision 和 parser version。
- [ ] PDF 评估阅读顺序、表格和字体。
- [ ] DOCX 处理样式、修订、关系和页码边界。
- [ ] Markdown 使用目标方言 AST。
- [ ] HTML 使用 DOM 和受控正文选择。
- [ ] 表格保留表头、主键、sheet/row/cell。
- [ ] OCR 保存 bbox、引擎、语言和 confidence。
- [ ] 低质量结果不静默进入索引。
- [ ] 外部资源、宏和脚本不自动执行。
- [ ] 每种格式有独立质量样本和指标。

## 集成练习

建立一个包含六种格式的解析基准：

1. 每种至少 20 个真实文件，并包含困难样本。
2. 统一输出 block、结构、locator、warning 和 quality。
3. 为标题、表格、页码、阅读顺序与关键值建立人工标注。
4. 运行两个解析器版本并比较。
5. 任一格式低于质量门时阻止发布索引。
6. 调试页面可以从 block 打开原文位置。
7. 解析进程超时、内存超限和密码文件产生明确错误。
8. 所有结果绑定文件 hash、parser version 和处理时间。

## 来源

- [ISO：PDF 2.0（ISO 32000-2:2020）](https://www.iso.org/standard/75839.html)（访问日期：2026-07-17）
- [ECMA-376：Office Open XML File Formats](https://ecma-international.org/publications-and-standards/standards/ecma-376/)（访问日期：2026-07-17）
- [CommonMark Specification](https://spec.commonmark.org/)（访问日期：2026-07-17）
- [WHATWG HTML Living Standard](https://html.spec.whatwg.org/multipage/)（访问日期：2026-07-17）
- [IETF RFC 4180：Common Format and MIME Type for CSV Files](https://www.rfc-editor.org/rfc/rfc4180)（访问日期：2026-07-17）
