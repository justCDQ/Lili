# 终端、Shell 与 Git 基础

## 是什么

终端提供文本输入输出界面；Shell 解析命令、展开参数、连接管道并启动进程。Git 是内容寻址的分布式版本控制系统：提交指向树和父提交，分支是可移动引用，远端是另一组引用。

```sh
git status
git switch -c feat/parser
git add src/parser.go
git commit -m "feat: parse CSV records"
git fetch origin
git rebase origin/main
```

## 关键特性或规则

- 引号决定空格、通配符和变量是否展开；不可信输入不能拼进 Shell 命令。
- 管道传 stdout；退出码表示成功/失败；`&&` 仅在前一命令成功时继续。
- 提交应单一目的、可构建、消息说明意图；合并前同步远端并解决冲突。
- `fetch` 只更新远端跟踪引用，`pull` 通常还会合并或变基，`push` 修改远端。
- 不提交密钥；已提交密钥必须轮换，删除历史不足以撤销泄漏。

## 常见错误与边界

重写已共享历史会改变提交 ID，应协调。`.gitignore` 不会停止跟踪已纳入 Git 的文件。Shell 脚本默认错误处理容易遗漏，应显式检查失败。

## 为什么需要

后端程序直接依赖操作系统资源和开发工具链。理解这些对象的生命周期、权限和性能指标，才能从进程、文件、时间或资源层定位故障，并保证本地与 CI 的操作可复现。

## 实际怎么使用

按本文命令建立一个最小实验，先记录工作目录、用户、环境、进程和输入规模，再执行操作并保存退出码、日志与耗时。主动制造权限不足、路径错误、资源饱和或时区差异，确认诊断步骤能定位到具体层。

## 补充知识

容器和 CI 会改变命名空间、用户、工作目录、可用 CPU/内存和时区。任何依赖本机默认值的行为都应显式配置，并在目标运行环境再次验证。

## 来源

- [POSIX：Shell Command Language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html)（访问日期：2026-07-16）
- [Git：Git Internals](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)（访问日期：2026-07-16）
- [GitHub Docs：About remote repositories](https://docs.github.com/en/get-started/git-basics/about-remote-repositories)（访问日期：2026-07-16）
