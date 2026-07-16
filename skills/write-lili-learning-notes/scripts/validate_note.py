#!/usr/bin/env python3
"""Validate structural requirements for one or more Lili Markdown notes."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


FORBIDDEN = (
    "参考的写作方式",
    "经过联网搜索",
    "经过网络搜索",
    "本文核对的",
    "不把某一篇教程",
    "只要记住",
    "显而易见",
)


def validate(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")

    body = re.sub(r"\A---\n.*?\n---\n+", "", text, count=1, flags=re.S)
    if not body.startswith("# "):
        errors.append("缺少一级标题")
    if "## 来源" not in text:
        errors.append("缺少 ## 来源")
    if text.count("```") % 2:
        errors.append("代码围栏数量为奇数")

    frontmatter = text[len("---\n"):text.find("\n---", 4)] if text.startswith("---\n") else ""
    is_intermediate = bool(re.search(r"^stage:\s*intermediate\s*$", frontmatter, re.M))
    minimum_lines = 200 if is_intermediate else 120
    minimum_bytes = 12_000 if is_intermediate else 8_000
    line_count = len(text.splitlines())
    byte_count = len(text.encode("utf-8"))
    if line_count < minimum_lines:
        errors.append(f"正文过短：{line_count} 行，最低审计线为 {minimum_lines} 行")
    if byte_count < minimum_bytes:
        errors.append(f"内容过少：{byte_count} bytes，最低审计线为 {minimum_bytes} bytes")
    if re.search(r"\[\[[^\]]+\]\]", text):
        errors.append("包含 GitHub 无法解析的 Obsidian Wiki Link")
    if any(term in text for term in FORBIDDEN):
        found = [term for term in FORBIDDEN if term in text]
        errors.append(f"包含编辑过程说明：{', '.join(found)}")
    if re.search(r'data:image/svg\+xml,[^"\n]*\s', text):
        errors.append("SVG Data URL 包含未编码空白")

    source_text = text.split("## 来源", 1)[-1] if "## 来源" in text else ""
    sources = re.findall(r"^- \[[^]]+\]\(https?://[^)]+\)", source_text, re.M)
    if not 2 <= len(sources) <= 5:
        errors.append(f"来源数量应为 2–5，实际为 {len(sources)}")

    dates = re.findall(r"访问日期：\d{4}-\d{2}-\d{2}", source_text)
    if sources and len(dates) != len(sources):
        errors.append("每条来源都应包含访问日期")

    without_fences = re.sub(r"^```.*?^```\s*$", "", text, flags=re.M | re.S)
    without_inline = re.sub(r"`[^`\n]+`", "", without_fences)
    for raw in re.findall(r"\[[^]]*\]\(([^)]+)\)", without_inline):
        target = raw.split("#", 1)[0]
        if not target or re.match(r"^(https?://|mailto:|data:)", target, re.I):
            continue
        local = (path.parent / target.replace("%20", " ")).resolve()
        if not local.exists():
            errors.append(f"本地链接不存在：{raw}")

    for line_number, line in enumerate(text.splitlines(), start=1):
        if line.endswith((" ", "\t")):
            errors.append(f"第 {line_number} 行有尾随空白")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("paths", nargs="+", type=Path)
    args = parser.parse_args()
    failed = False

    for path in args.paths:
        errors = validate(path)
        if errors:
            failed = True
            for error in errors:
                print(f"{path}: {error}", file=sys.stderr)
        else:
            print(f"OK {path}")

    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
