# ESLint、Formatter、类型检查与 Git Hooks

## 是什么

linter 检查代码模式，formatter 统一格式，类型检查验证静态契约，Git hook 在 Git 事件触发脚本。它们是互补质量门。

## 为什么需要

Lint、格式化、类型检查和 Git Hooks 分别检查不同缺陷。明确各自边界并在 CI 重复关键检查，才能让本地反馈快速且合并标准一致。

## 关键特性与规则

配置版本化；本地与 CI 使用同命令；formatter 与 lint 规则避免重叠；hook 是便利层，CI 才是不可绕过门禁。

## 实际使用

```tsx
npx eslint .
npx prettier . --check
npx tsc --noEmit
# pre-commit 只跑快速检查，CI 跑完整检查
```

## 常见错误与边界

钩子可用 --no-verify 绕过；自动修复可能改变语义；仅检查 staged 文件不能证明全项目通过。

## 相关补充知识

Hook 可被跳过且受本机环境影响，不能作为唯一门禁。格式化只处理表现形式，Lint 规则需要控制误报，类型通过也不证明运行时输入和业务行为正确。

## 来源

- [ESLint Documentation](https://eslint.org/docs/latest/use/getting-started)
- [Prettier Documentation](https://prettier.io/docs/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/compiler-options.html)
- [Git Documentation](https://git-scm.com/docs/githooks)

访问日期：2026-07-16。
