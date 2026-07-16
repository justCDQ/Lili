#!/usr/bin/env python3
"""Audit a group of Lili notes and report structural and copy-quality failures."""

from __future__ import annotations

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path

from validate_note import validate


NOTE_ROOTS = (
    "01-frontend/notes",
    "02-product/notes",
    "03-interaction-design/notes",
    "04-ai/notes",
    "05-backend-data/notes",
)


def collect_paths(inputs: list[Path]) -> list[Path]:
    paths: set[Path] = set()
    for item in inputs:
        if item.is_dir():
            paths.update(item.rglob("*.md"))
        elif item.suffix == ".md":
            paths.add(item)
    return sorted(path for path in paths if path.name != "README.md")


def prose_paragraphs(text: str) -> list[str]:
    without_frontmatter = re.sub(r"\A---\n.*?\n---\n+", "", text, count=1, flags=re.S)
    without_fences = re.sub(r"^```.*?^```\s*$", "", without_frontmatter, flags=re.M | re.S)
    paragraphs = re.split(r"\n\s*\n", without_fences)
    normalized: list[str] = []
    for paragraph in paragraphs:
        paragraph = re.sub(r"\[[^]]+\]\([^)]+\)", "", paragraph)
        paragraph = re.sub(r"[`*_>#|\-]", "", paragraph)
        paragraph = re.sub(r"\s+", "", paragraph)
        if len(paragraph) >= 120:
            normalized.append(paragraph)
    return normalized


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("paths", nargs="*", type=Path)
    parser.add_argument(
        "--all",
        action="store_true",
        help="audit the five direction note roots from the current repository",
    )
    args = parser.parse_args()

    inputs = list(args.paths)
    if args.all:
        inputs.extend(Path(root) for root in NOTE_ROOTS)
    if not inputs:
        parser.error("provide note paths/directories or --all")

    paths = collect_paths(inputs)
    if not paths:
        print("没有找到 Markdown 学习笔记", file=sys.stderr)
        return 2

    failed = False
    passed = 0
    paragraph_owners: dict[str, list[Path]] = defaultdict(list)

    for path in paths:
        errors = validate(path)
        if errors:
            failed = True
            for error in errors:
                print(f"ERROR {path}: {error}")
        else:
            passed += 1

        for paragraph in set(prose_paragraphs(path.read_text(encoding="utf-8"))):
            paragraph_owners[paragraph].append(path)

    repeated_groups = 0
    for paragraph, owners in paragraph_owners.items():
        unique_owners = sorted(set(owners))
        if len(unique_owners) < 2:
            continue
        repeated_groups += 1
        failed = True
        preview = paragraph[:80]
        joined = ", ".join(str(path) for path in unique_owners)
        print(f"ERROR 重复长段落 [{preview}…]: {joined}")

    print(
        f"审计 {len(paths)} 篇：单篇校验通过 {passed}，"
        f"失败 {len(paths) - passed}，跨文章重复长段落 {repeated_groups}"
    )
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
