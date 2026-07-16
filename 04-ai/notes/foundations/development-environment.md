---
type: ai-note
stage: beginner
topic: development-environment
verified: 2026-07-16
tags: [ai, cli, git, package-management, environment]
---

# 命令行、Git、包管理、虚拟环境与环境变量

## 是什么

命令行用于执行程序、安装依赖、运行测试和自动化任务。Git 保存文件快照和历史。包管理器安装并锁定第三方依赖。虚拟环境隔离 Python 项目的包。环境变量向进程传入配置，常用于 API 地址、模型名和 Secret 引用。

## 为什么需要

AI 项目依赖 SDK、Schema、数据处理和评估工具。没有版本和环境管理，同一代码可能因依赖差异无法复现；把密钥写入代码又会造成泄露。

## 关键特性

### 命令行

程序从当前工作目录运行，路径可以是绝对路径或相对路径。命令的退出码 `0` 通常表示成功，非零表示失败，CI 会据此判断步骤结果。标准输出用于正常结果，标准错误用于诊断。

### Git

工作区是正在编辑的文件；暂存区选择下一次提交的内容；提交保存一个可追踪快照。分支是指向提交的可移动引用，远端是另一个 Git 仓库的命名地址。数据集、模型结果或 Secret 是否进入 Git 必须单独判断。

### 包与版本

JavaScript 项目使用 `package.json` 描述依赖和脚本，锁文件固定实际解析版本。Python 使用项目元数据或 requirements 文件记录依赖，`venv` 为项目创建隔离环境。锁文件应提交，虚拟环境目录和 Secret 文件不提交。

### 环境变量

环境变量属于进程环境。它解决配置注入，不等于完整的 Secret 管理系统；环境变量可能出现在进程检查、崩溃报告或错误日志中。

## 实际怎么使用

TypeScript 项目：

```bash
npm init -y
npm install <official-model-sdk>
git init
git add package.json package-lock.json
git commit -m "chore: initialize ai experiment"
```

Python 项目：

```bash
python3 -m venv .venv
source .venv/bin/activate
python -m pip install <official-model-sdk>
python -m pip freeze > requirements.txt
```

环境变量只在运行时读取，并在启动阶段检查：

```ts
const apiKey = process.env.MODEL_API_KEY;
if (!apiKey) throw new Error("MODEL_API_KEY is required");
```

`.gitignore` 至少排除 `.env`、`.venv/`、本地缓存和包含敏感输入的实验产物。提供 `.env.example` 时只写变量名和非敏感示例。

## 常见错误与边界

- 只提交依赖清单但不提交锁文件，导致安装结果漂移。
- 把 `.env`、API Key、真实用户输入或完整模型 Trace 提交到 Git。
- 认为从 Git 历史删除当前文件即可撤销 Secret；已泄露 Secret 必须立即吊销和轮换。
- 使用系统全局 Python 环境安装所有包，项目间版本互相影响。
- 未记录运行时版本、依赖版本和模型标识，实验无法复现。

## 补充知识

生产环境优先使用云平台或组织的 Secret Manager，将短期凭据和最小权限结合。依赖升级要运行测试与评估，因为 SDK 结构变化和默认参数变化都可能影响行为。

## 来源

- [Git User Manual](https://git-scm.com/docs/user-manual)（访问日期：2026-07-16）
- [Git Reference](https://git-scm.com/docs)（访问日期：2026-07-16）
- [Python Packaging：pip 与 venv](https://packaging.python.org/en/latest/guides/installing-using-pip-and-virtual-environments/)（访问日期：2026-07-16）
- [Node.js：Environment Variables](https://nodejs.org/api/environment_variables.html)（访问日期：2026-07-16）

