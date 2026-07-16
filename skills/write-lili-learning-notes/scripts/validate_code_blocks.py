#!/usr/bin/env python3
"""Parse executable fenced blocks in Lili notes with available local tools."""

from __future__ import annotations

import argparse
import json
import re
import shutil
import subprocess
import sys
from pathlib import Path


FENCE = re.compile(r"^```([A-Za-z0-9_+-]+)\s*\n(.*?)^```\s*$", re.M | re.S)


def run_parser(command: list[str], source: str) -> str | None:
    result = subprocess.run(
        command,
        input=source,
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode == 0:
        return None
    return (result.stderr or result.stdout).strip()


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("paths", nargs="+", type=Path)
    parser.add_argument("--node", type=Path, help="Node.js executable for JS syntax checks")
    args = parser.parse_args()

    node = str(args.node) if args.node else shutil.which("node")
    bash = shutil.which("bash")
    failed = False
    checked = {"json": 0, "javascript": 0, "shell": 0}
    skipped_typescript = 0

    paths: list[Path] = []
    for item in args.paths:
        if item.is_dir():
            paths.extend(sorted(item.rglob("*.md")))
        elif item.suffix == ".md":
            paths.append(item)

    for path in sorted(set(paths)):
        text = path.read_text(encoding="utf-8")
        for index, match in enumerate(FENCE.finditer(text), start=1):
            language = match.group(1).lower()
            source = match.group(2)
            line = text.count("\n", 0, match.start()) + 1
            error: str | None = None

            if language in {"json", "jsonc"}:
                if language == "jsonc":
                    continue
                try:
                    json.loads(source)
                    checked["json"] += 1
                except json.JSONDecodeError as exc:
                    error = str(exc)
            elif language in {"js", "javascript", "mjs"}:
                if not node:
                    print(f"SKIP {path}:{line}: 找不到 Node.js")
                    continue
                error = run_parser([node, "--input-type=module", "--check", "-"], source)
                checked["javascript"] += 1
            elif language in {"bash", "sh", "shell"}:
                if not bash:
                    print(f"SKIP {path}:{line}: 找不到 bash")
                    continue
                error = run_parser([bash, "-n"], source)
                checked["shell"] += 1
            elif language in {"ts", "typescript", "tsx"}:
                skipped_typescript += 1

            if error:
                failed = True
                print(f"ERROR {path}:{line} block {index} ({language}): {error}")

    print(
        "代码围栏检查："
        f"JSON {checked['json']}，JavaScript {checked['javascript']}，"
        f"Shell {checked['shell']}；TypeScript 待编译器复验 {skipped_typescript}"
    )
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
