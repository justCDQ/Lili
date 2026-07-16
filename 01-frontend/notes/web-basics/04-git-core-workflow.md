# Git：仓库、工作区、暂存区、提交、分支、合并与远端

## 是什么

Git 是分布式版本控制系统。仓库存储对象和引用；工作区是检出的可编辑文件；暂存区（index）是下一次提交的候选快照；提交记录快照、父提交和作者信息。分支是指向提交的可移动引用，合并组合历史，远端是另一个仓库的命名地址。

## 为什么需要

Git 提供可审查历史、并行开发、回退依据和跨设备同步。暂存区允许一次只提交一个逻辑变化，分支隔离实验。

## 实际使用

```sh
git init
git status
git add README.md
git diff --staged
git commit -m "docs: add readme"
git switch -c feature/nav
git switch main
git merge feature/nav
git remote add origin https://github.com/OWNER/REPO.git
git push -u origin main
git fetch origin
git pull --ff-only
```

提交前分别看 `git diff` 与 `git diff --staged`。解决冲突后删除冲突标记，测试，再 `git add` 和提交。

## 关键规则

- 提交的是暂存区，不是整个工作区。
- `origin/main` 是远端跟踪分支；本地 `main` 与它不是同一引用。
- `.gitignore` 不会停止追踪已纳入仓库的文件。
- commit 哈希标识内容和历史；改写历史会产生新哈希。

## 常见错误与边界

不要提交密码、令牌和私钥；删除后仍可能存在历史中。共享分支避免强制推送。合并冲突不是 Git 判断哪方正确，需根据意图处理。Git 不适合直接存大量频繁变化的二进制文件。

## 补充知识

`HEAD` 通常指向当前分支；detached HEAD 指向具体提交。`git restore`、`reset`、`revert`语义不同，执行前读官方文档并确认是否改写历史。

## 来源

- [Git：gitglossary](https://git-scm.com/docs/gitglossary)
- [Git：Everyday Git](https://git-scm.com/docs/giteveryday)
- [Git：gitignore](https://git-scm.com/docs/gitignore)

访问日期：2026-07-16。
