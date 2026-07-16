# Node.js、包管理器、package.json、依赖与脚本

## 是什么

Node.js 是基于 V8 的 JavaScript 运行时。包管理器解析、下载并记录项目依赖。`package.json` 是包的清单，包含名称、版本、脚本、依赖和模块类型；锁文件记录已解析的精确依赖图。`dependencies` 是运行所需依赖，`devDependencies` 是开发和构建所需依赖。

## 为什么需要

前端工具链、开发服务器、测试器和构建器通常运行于 Node.js。清单与锁文件让团队和 CI 能重建接近一致的依赖环境。

## 实际使用

```sh
node --version
npm --version
npm init -y
npm install lodash
npm install --save-dev eslint
npm run test
npm ci
```

```json
{
  "private": true,
  "scripts": {"dev":"vite", "test":"node --test"},
  "dependencies": {"lodash":"^4.17.21"},
  "devDependencies": {"vite":"^7.0.0"}
}
```

## 关键规则与边界

- 提交 `package.json` 和所选包管理器的锁文件，不提交 `node_modules`。
- 同一项目避免混用 npm、pnpm、Yarn 锁文件。
- `npm run` 会把本地 `node_modules/.bin` 加入 PATH。
- 语义版本范围不保证升级无缺陷；升级后必须测试。
- 安装脚本可执行代码，依赖不是天然可信。

## 常见错误与边界

全局安装项目工具会导致版本漂移；优先项目内依赖。不要手改锁文件。`npm install` 可更新锁文件，CI 通常使用 `npm ci` 做冻结安装。不要把前端可见环境变量当秘密。

## 补充知识

Node 的 `engines` 可声明支持版本但不必然强制；可配合版本管理器。定期使用审计工具，但漏洞报告仍需结合代码是否实际可达判断。

## 来源

- [Node.js：Introduction](https://nodejs.org/en/learn/getting-started/introduction-to-nodejs)
- [npm：package.json](https://docs.npmjs.com/cli/configuring-npm/package-json)
- [npm：npm scripts](https://docs.npmjs.com/cli/using-npm/scripts)
- [npm：npm ci](https://docs.npmjs.com/cli/commands/npm-ci)

访问日期：2026-07-16。
