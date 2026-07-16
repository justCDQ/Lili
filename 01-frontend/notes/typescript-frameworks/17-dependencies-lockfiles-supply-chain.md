# 依赖选择、升级、锁文件与供应链

## 是什么

依赖把外部包引入运行时或工具链；锁文件固定解析图；供应链风险包括恶意包、接管、安装脚本和漏洞。

## 为什么需要

直接与传递依赖会进入构建、开发或生产信任边界。锁文件、冻结安装、升级策略和来源审查共同保证可复现性，并降低恶意包、接管和未审查更新风险。

## 关键特性与规则

选择看维护、许可、体积、API 和替代成本；提交单一锁文件；升级分批并跑测试；最小化发布/CI 权限。

## 实际使用

```tsx
npm install
npm ci
npm outdated
npm audit
```

## 常见错误与边界

只看下载量不代表安全；audit 报告需判断可达性；忽略锁文件会产生不可重复安装；盲目 major 更新会破坏 API。

## 相关补充知识

漏洞扫描结果要结合依赖是否到达生产和代码路径是否可达判断。减少依赖、限制安装脚本、检查维护状态和许可证，通常比只在告警后升级更有效。

## 来源

- [npm Documentation](https://docs.npmjs.com/cli/commands/npm-ci)
- [npm Documentation](https://docs.npmjs.com/cli/commands/npm-audit)
- [npm Documentation](https://docs.npmjs.com/about-lockfiles)
- [npm Documentation](https://docs.npmjs.com/threats-and-mitigations)

访问日期：2026-07-16。
