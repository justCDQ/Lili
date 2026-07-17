# ESLint、Formatter、类型检查与 Git Hooks

质量门把可自动判定的约束放进本地和 CI：formatter 统一文本格式，ESLint 检查语法与代码模式，TypeScript 证明静态契约，测试验证行为，Git hooks 提供快速反馈，CI 是最终权威。它们检查不同问题，不能互相替代。

## 1. 门禁流水线

```mermaid
flowchart LR
    A["编辑器保存"] --> B["format"]
    B --> C["pre-commit: changed files"]
    C --> D["CI install frozen lockfile"]
    D --> E["format check"]
    E --> F["lint"]
    F --> G["typecheck"]
    G --> H["tests"]
    H --> I["production build"]
```

hook 可被跳过或环境损坏，保护分支必须要求 CI 结果。

## 2. Formatter

Formatter 处理缩进、换行、引号等可机械决定的布局，减少 review 噪声。配置应固定版本并在仓库根统一：

```json
{
  "scripts": {
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  }
}
```

生成文件、dist、coverage 和大型快照放 ignore。不要同时启用与 formatter 冲突的 ESLint 样式规则。

## 3. ESLint Flat Config

当前 ESLint 使用 `eslint.config.js/mjs` flat config：

```js
import js from "@eslint/js";
import { defineConfig } from "eslint/config";
import globals from "globals";

export default defineConfig([
  { ignores: ["dist/**", "coverage/**"] },
  {
    files: ["**/*.{js,mjs,cjs}"],
    languageOptions: { globals: globals.browser },
    extends: [js.configs.recommended],
    rules: {
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "no-warning-comments": ["warn", { terms: ["FIXME"], location: "anywhere" }],
    },
  },
]);
```

配置数组按匹配和顺序合并；不同目录的 Node、browser、test globals 分开定义。`warn` 默认不使进程失败，CI 若要求零 warning 使用 `--max-warnings 0`。

### 3.1 Type-aware lint

typescript-eslint 可使用类型信息检查浮动 Promise、不安全 any 等，但启动更慢且依赖 TypeScript Compiler API。TypeScript 7.0 无 API，当前工具链可能需要 `@typescript/typescript6`。不要强制把 lint 的 `typescript` 模块指向 TS7 后假设兼容。

把快速语法 lint 和 type-aware lint 分层；CI 执行完整，pre-commit 只跑受影响文件的安全子集。

## 4. TypeScript 门

```json
{
  "scripts": {
    "typecheck": "tsc --noEmit --pretty false"
  }
}
```

TS7 默认变化不能代替显式配置。CI 固定 `strict`、target、module/resolution、types、exact optional 和 unchecked index。框架模板另跑专用检查器。

类型错误不能用批量 `as any`、`@ts-ignore` 或关闭 strict 消除。确有外部类型缺陷时建立局部 adapter、运行时验证和追踪 issue；`@ts-expect-error` 附原因并在错误消失时失败。

## 5. Git Hooks

pre-commit 目标是 1–5 秒内捕获格式和明显 lint：

```json
{
  "scripts": {
    "lint": "eslint . --max-warnings 0",
    "check": "pnpm format:check && pnpm lint && pnpm typecheck && pnpm test && pnpm build"
  },
  "lint-staged": {
    "*.{js,jsx,ts,tsx}": ["prettier --write", "eslint --max-warnings 0"],
    "*.{json,css,md,yaml,yml}": ["prettier --write"]
  }
}
```

不要在 hook 中自动修改未暂存部分而不理解部分暂存语义。大型 monorepo 的 typecheck/test 使用受影响项目图，但 CI 定期全量验证。

## 6. CI 配置原则

- 锁定运行时与包管理器；
- frozen lockfile 安装；
- 最小权限和无 secret 的 PR 环境；
- cache key 包含 lockfile、平台和工具版本；
- 失败立即非零退出；
- 产物、测试报告和 sourcemap 按敏感级别保存；
- 同一 commit 只构建一次，部署已验证 artifact。

```yaml
steps:
  - run: corepack enable
  - run: pnpm install --frozen-lockfile
  - run: pnpm format:check
  - run: pnpm lint
  - run: pnpm typecheck
  - run: pnpm test --run
  - run: pnpm build
```

这是步骤片段，不包含特定 CI 平台 action。真实流水线还应固定 action SHA、超时和并发取消。

## 7. 规则选择

规则必须能说明缺陷类型和修复方式。优先：

- 未处理 Promise；
- React Hook 规则；
- 无用变量/导入；
- 禁止危险 eval 和不安全 DOM API；
- 测试中 focused/disabled case；
- 包边界和 server/client import；
- 无障碍 JSX 规则作为补充。

lint 无法证明运行时授权、完整无障碍和业务正确。高误报规则会诱导 disable，应试运行、修复基线、记录例外，再升为 error。

## 8. 例外治理

```ts
// @ts-expect-error 第三方声明缺少已存在的 runtime overload；移除条件：vendor#123
vendorCall(options);
```

ESLint disable 指定单条规则和最小范围，附实际原因。定期统计 suppressions；禁止无说明 `eslint-disable`、`@ts-nocheck` 和全局规则关闭。

## 9. 完整案例：引入质量门

现有项目有 1200 个 lint warning，直接 `--max-warnings 0` 会阻塞所有交付。

实施：

1. 固定 formatter，单独提交纯格式变更；
2. 启用推荐规则，按规则统计；
3. 自动修复安全项；
4. 将真实高风险规则设 error；
5. 对遗留目录建立明确临时 override 与到期责任人；
6. 新增代码不允许增加 warning；
7. 分批清零后 CI `--max-warnings 0`；
8. 加入 TS7 typecheck；typescript-eslint type-aware 链按 TS6 API 兼容验证；
9. pre-commit 跑 staged，CI 跑全量。

输出：本地快速门、完整 CI 门、例外清单和指标。验证故意加入格式错、floating Promise、类型错、失败测试和 build-only 错，五类都在正确阶段失败。

失败分支：hook 跑完整 E2E 导致开发者频繁跳过；共享 config 未锁版本使本地/CI 结果漂移；自动 `eslint --fix` 改行为而未 review；类型工具被强制 TS7 后插件崩溃。

## 10. 调试

- `eslint --print-config file.ts` 查看最终规则；
- `eslint --debug` 查看文件匹配和插件；
- `tsc --showConfig` 查看最终 TS 配置；
- 本地使用 CI 同一 Node/pnpm 复现；
- 清 cache 检查是否陈旧；
- 分离 lint、typecheck、test、build 确认责任；
- 遇到只在 CI 失败检查大小写、换行、时区、环境变量和未跟踪文件。

## 11. Monorepo 与生成代码

Monorepo 的门禁应按项目依赖图执行：底层公共类型变化会让全部消费者 typecheck，独立文档变更不必重跑所有浏览器矩阵。受影响计算属于加速层；main 分支和定时任务仍做全量检查，以捕获依赖图配置遗漏。

生成代码不直接手改。CI 先运行生成器，再用 `git diff --exit-code` 验证仓库产物没有过期；生成器本身固定版本并测试。若生成文件排除 lint/typecheck，消费它的公共入口仍要通过编译和契约测试。

共享 ESLint/tsconfig 包按普通依赖发布与锁定。升级共享规则应单独 PR，展示新增错误数量和迁移方式，不能让浮动 workspace 配置在不同分支产生不同结果。

## 12. 常见错误

1. formatter 和 lint 重复争夺样式。
2. 把 warning 当“以后处理”，数量无限增长。
3. hook 是唯一门禁。
4. 仅 lint staged，主分支已有错误永远不发现。
5. 用类型断言消灭错误。
6. CI 安装浮动依赖或不使用 lockfile。
7. cache 未包含 lockfile，复用错误依赖。
8. 不验证 TS7 与 type-aware lint 的 TS6 API 边界。

## 13. 练习

为 TypeScript 框架项目建立质量门。验收：

1. flat config 区分 browser/node/test；
2. formatter 与 ESLint 无冲突；
3. staged hook 5 秒内完成；
4. CI 执行 format/lint/typecheck/test/build；
5. warnings 为零；
6. suppression 有原因、范围和追踪项；
7. 注入五类错误均被对应门捕获；
8. TS7 CLI 和需要的 TS6 lint 工具均复验。

## 来源

- [ESLint：Configure ESLint](https://eslint.org/docs/latest/use/configure/)（访问日期：2026-07-17）
- [ESLint：CLI](https://eslint.org/docs/latest/use/command-line-interface)（访问日期：2026-07-17）
- [typescript-eslint：Typed Linting](https://typescript-eslint.io/getting-started/typed-linting/)（访问日期：2026-07-17）
- [TypeScript TSConfig Reference](https://www.typescriptlang.org/tsconfig/)（访问日期：2026-07-17）
- [lint-staged Documentation](https://github.com/lint-staged/lint-staged)（访问日期：2026-07-17）
