---
type: ai-note
stage: beginner
topic: files-encoding-data-formats
verified: 2026-07-16
tags: [ai, files, encoding, csv, jsonl, data-cleaning]
---

# 文件、文本编码、CSV、JSONL 与数据清洗

## 是什么

文件是持久化字节序列；文本文件需要字符编码把字节解释为字符。UTF-8 是 AI 文本数据常用编码。CSV 用分隔字段表达表格，JSONL/NDJSON 每行保存一个独立 JSON 值，适合流式处理样例、模型输入输出和评估记录。数据清洗是把原始数据转换为可验证、可追踪的输入数据。

## 为什么需要

Prompt、文档、评测集、模型响应和 Trace 经常以文件保存。编码错误会破坏字符；格式解析错误会错位字段；不受控清洗会丢失证据或把测试集信息泄漏进训练/开发过程。

## 关键特性

### 文本编码

读写时显式指定 UTF-8，避免依赖操作系统默认编码。解码错误不应默认静默忽略，否则原文会在不知情时损坏。换行符在不同系统可能不同，比较和切分前应明确规范化策略。

### CSV

CSV 适合二维表格，但分隔符、引号、换行、空值和字符编码需要由解析器处理，不能简单使用字符串 `split(',')`。复杂嵌套数据不适合 CSV。

### JSONL

每行独立解析，便于追加、流式读取、错误定位和并行处理。一行必须是完整合法 JSON；缩进后的多行 JSON 不符合 JSONL 的逐行约定。文件结尾通常保留换行。

### 数据清洗

清洗包括去除确定性噪声、规范化编码和字段、去重、脱敏、验证 Schema、标记来源和隔离无效记录。原始数据应只读保留，清洗代码和产出版本可追踪。

## 实际怎么使用

```python
import json

with open("evals.jsonl", encoding="utf-8") as source:
    for line_number, line in enumerate(source, start=1):
        if not line.strip():
            continue
        try:
            record = json.loads(line)
        except json.JSONDecodeError as error:
            raise ValueError(f"invalid JSON at line {line_number}") from error

        if "input" not in record or "expected" not in record:
            raise ValueError(f"missing fields at line {line_number}")
```

建议数据目录：

```text
data/raw/          # 原始、只读、权限严格
data/processed/    # 清洗后的版本化数据
data/evals/        # 开发/回归评测集
schemas/           # 字段和约束
scripts/           # 可重复执行的处理代码
```

每条记录至少保留 `id`、来源、采集/生成时间、许可或权限、Schema 版本和处理版本。

## 常见错误与边界

- 使用默认编码，代码在另一操作系统读取失败。
- 用字符串分割 CSV，遇到被引号包围的逗号或换行时字段错位。
- 清洗后覆盖原始数据，无法审计删除和转换过程。
- 用同一批样例反复调整 Prompt 又称其为独立测试集，造成评估泄漏。
- 把真实用户敏感数据直接提交仓库或发送给未经批准的模型服务。

## 补充知识

大文件应使用流式读取，避免一次载入内存。数据版本变化可能导致评估分数变化，因此报告必须同时记录数据集版本和处理代码版本。

## 来源

- [Python `io`：Text Encoding](https://docs.python.org/3/library/io.html#text-encoding)（访问日期：2026-07-16）
- [Python `csv`](https://docs.python.org/3/library/csv.html)（访问日期：2026-07-16）
- [Python `json`](https://docs.python.org/3/library/json.html)（访问日期：2026-07-16）
- [JSON Lines](https://jsonlines.org/)（访问日期：2026-07-16）

