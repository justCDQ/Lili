---
title: 开发本地 MCP Server
stage: intermediate
direction: ai
topic: mcp
---

# 开发本地 MCP Server

本地 MCP Server 通常由 Host 作为 stdio 子进程启动。最小实现要完成 initialize、声明 capabilities、注册 primitives、验证输入、限制本机资源，并保证 stdout 只有协议消息。本文使用 TypeScript SDK 的高层 Server 结构说明实现；具体 SDK 版本应固定在 lockfile，并按该版本导出 API 编译。

## 前置知识与目标

前置阅读：

- [MCP Host、Client 与 Server](01-host-client-server.md)。
- [使用 MCP Inspector 调试 Server](03-inspector-debugging.md)。

目标 Server：

- Resource：读取当前项目 roadmap。
- Tool：按关键词搜索 Markdown 标题。
- Prompt：生成学习笔记复核消息。
- stdio transport。
- 只读 workspace。
- 有 Schema、错误、上限、日志和测试。

## 项目结构

```text
lili-mcp-server/
├── package.json
├── package-lock.json
├── tsconfig.json
├── src/
│   ├── index.ts
│   ├── config.ts
│   ├── paths.ts
│   ├── resources.ts
│   ├── tools.ts
│   └── prompts.ts
└── test/
    ├── protocol.test.ts
    ├── paths.test.ts
    └── tools.test.ts
```

入口只做配置、注册和连接。文件访问、Schema 与业务逻辑可单测。

## 依赖与 runtime

package 固定版本，不使用宽泛 `latest` 进入生产：

```json
{
  "name": "lili-local-mcp",
  "private": true,
  "type": "module",
  "engines": {
    "node": ">=22"
  },
  "scripts": {
    "build": "tsc -p tsconfig.json",
    "start": "node dist/index.js",
    "test": "node --test dist-test/**/*.test.js"
  },
  "dependencies": {
    "@modelcontextprotocol/sdk": "1.24.3",
    "zod": "4.1.13"
  },
  "devDependencies": {
    "typescript": "7.0.2"
  }
}
```

版本号示例必须在实现时以官方 SDK 当前稳定版和 lockfile 为准；API 与包版本一起验证。文章不要求读者复制未经安装验证的版本组合。

## 配置

Host 传入允许 workspace root：

```json
{
  "command": "node",
  "args": ["/absolute/path/lili-mcp-server/dist/index.js"],
  "env": {
    "LILI_WORKSPACE_ROOT": "/absolute/path/Lili",
    "LOG_LEVEL": "info"
  }
}
```

Server 启动时：

- 环境变量存在。
- `realpath`。
- 是目录。
- 不位于禁止系统目录。
- 不从 Tool arguments 覆盖。

Secret 不需要就不传。

## stdio 入口

高层结构：

```typescript
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

const server = new McpServer({
  name: "lili-project",
  version: "1.0.0"
});

registerResources(server);
registerTools(server);
registerPrompts(server);

const transport = new StdioServerTransport();
await server.connect(transport);
```

SDK 处理 initialize 与 JSON-RPC framing。应用仍负责 primitive 语义、安全和异常。

stdout 禁止：

```typescript
console.log("server started");
```

日志：

```typescript
process.stderr.write(
  JSON.stringify({ level: "info", event: "server_started" }) + "\n"
);
```

日志不含 Tool arguments 原文。

## 路径边界

```typescript
import path from "node:path";
import { realpath } from "node:fs/promises";

export async function resolveInsideRoot(
  root: string,
  relativePath: string
): Promise<string> {
  if (path.isAbsolute(relativePath)) {
    throw new Error("absolute_path_not_allowed");
  }
  const rootReal = await realpath(root);
  const candidate = await realpath(path.resolve(rootReal, relativePath));
  const prefix = rootReal.endsWith(path.sep) ? rootReal : rootReal + path.sep;
  if (candidate !== rootReal && !candidate.startsWith(prefix)) {
    throw new Error("path_outside_workspace");
  }
  return candidate;
}
```

边界：

- `../`。
- absolute。
- symlink。
- case sensitivity。
- 文件不存在的写路径另用父目录 canonicalization。
- root 本身可否读取按 Tool 决定。

本 Server 只读 `.md`，并显式拒绝：

- `.git`。
- `.env`。
- hidden Secret。
- 超过 1MB。

## Resource

URI 使用 Server scheme：

```text
lili://roadmap/ai
```

注册逻辑读取固定相对路径，不从 URI 直接映射任意路径：

```typescript
server.resource(
  "ai-roadmap",
  "lili://roadmap/ai",
  async (uri) => {
    const text = await readAllowedMarkdown("04-ai/roadmap.md");
    return {
      contents: [{
        uri: uri.href,
        mimeType: "text/markdown",
        text
      }]
    };
  }
);
```

Resource 返回 source revision/hash 可放 `_meta` 或内容 wrapper，便于缓存与引用。

## Tool

Tool 只搜索已允许 Markdown 标题，不接受 path：

```typescript
import { z } from "zod";

server.tool(
  "search_note_headings",
  "在当前 Lili workspace 的 Markdown 标题中搜索关键词；只读，不返回正文。",
  {
    query: z.string().trim().min(1).max(100),
    limit: z.number().int().min(1).max(20).default(10)
  },
  async ({ query, limit }) => {
    const matches = await searchHeadings(query, limit);
    return {
      content: [{
        type: "text",
        text: JSON.stringify({ matches })
      }],
      structuredContent: { matches }
    };
  }
);
```

`searchHeadings`：

- `rg` 可作为受控子进程，但 arguments 数组固定，不用 shell。
- 或 Node 遍历允许目录。
- 限制文件数、大小、总时间。
- 返回相对 path、heading、line。
- path 再次确保 root。

## 输出 Schema

若 SDK 注册 API 支持 outputSchema，声明：

```json
{
  "type": "object",
  "properties": {
    "matches": {
      "type": "array",
      "maxItems": 20,
      "items": {
        "type": "object",
        "properties": {
          "path": {"type": "string"},
          "line": {"type": "integer", "minimum": 1},
          "heading": {"type": "string"}
        },
        "required": ["path", "line", "heading"],
        "additionalProperties": false
      }
    }
  },
  "required": ["matches"],
  "additionalProperties": false
}
```

Server 在返回前也验证，不只依赖类型系统。

## Prompt

Prompt arguments：

- direction enum。
- level enum。

返回 message：

```typescript
server.prompt(
  "review_learning_note",
  "生成学习笔记复核任务模板",
  {
    direction: z.enum(["frontend", "product", "interaction", "ai", "backend"]),
    level: z.enum(["beginner", "junior", "intermediate"])
  },
  ({ direction, level }) => ({
    messages: [{
      role: "user",
      content: {
        type: "text",
        text: `复核 ${direction} 的 ${level} 笔记：检查概念、实例、失败边界与来源。`
      }
    }]
  })
);
```

Prompt 是用户可选择模板，不包含 workspace 正文，也不要求绕过 Host policy。

## 错误

业务错误返回安全结构：

```json
{
  "code": "query_too_broad",
  "retryable": false,
  "message": "请提供更具体的标题关键词。"
}
```

内部 `EACCES`、绝对路径和 stack 写脱敏 stderr。未知错误不把 stack 放 Tool content。

## 初始化与 Capability

SDK 根据注册内容声明 tools/resources/prompts。需要验证：

- 真的支持 listChanged 才声明。
- 没实现 subscribe 不声明。
- 没实现 Tasks 不声明。
- Server instructions 简短且不含安全承诺。

高层 SDK 默认行为仍用 Inspector 查看 initialize response。

## SDK 升级

升级官方 SDK 时先读取 release notes，并运行：

- TypeScript compile。
- lifecycle contract。
- primitive list snapshot。
- Schema valid/invalid fixtures。
- Inspector smoke。
- cancellation/timeout。
- stdout framing。

SDK 可能新增协议能力，但 Server 不应在未实现业务语义时自动暴露。lockfile、构建 artifact 与 Host catalog 一起发布，失败可回到上一版本。

## 优雅退出

处理：

- stdin EOF。
- SIGTERM/SIGINT。
- active request。
- 子进程。
- temp files。

只读请求可取消；不要因为断开继续扫描大型目录。设置 AbortSignal/timeout 传到搜索。

## 应用案例一：Obsidian/IDE 笔记助手

Host 连接 Server：

- Resource 打开 roadmap。
- Tool 搜标题。
- Prompt 复核笔记。

用户“找 AI 中所有权限主题”：

- 模型调用 search。
- Server 返回标题/路径。
- Host 再按用户选择读取文件。

Server 不自动读取所有笔记进模型，降低数据暴露。

### 失败

query=`"."` 可能匹配全部。Server 限制最小有效字符、结果和扫描时间，返回 broad query。

## 应用案例二：Git branch 切换

Resource 内容改变。简单 Server 不支持 subscription：

- Host 每次 read 得到新 hash。
- citation 绑定旧 Git revision。
- 不虚报 listChanged。

扩展版本可监听 `.git/HEAD`，但文件 watcher：

- 去抖。
- 不读取 Secret。
- 发送 resource updated。
- Host 重新 read。

### 失败

watcher 洪泛不能阻塞协议；设置队列上限，合并同 URI。

## 应用案例三：受控诊断工具

添加 `get_repository_status`：

- 运行固定 `git status --porcelain=v1`。
- `spawn("git", args, {shell:false})`。
- cwd 为 root。
- timeout 2 秒。
- max stdout 256KB。
- 不支持任意 Git arguments。

这仍是只读诊断，但子进程环境最小化。

## 协议测试

用内存或 stdio Client：

1. initialize。
2. capabilities。
3. tools/list。
4. resources/list/read。
5. prompts/list/get。
6. valid tool call。
7. invalid args。
8. unknown primitive。
9. cancellation。
10. shutdown。

断言 JSON-RPC ID 和 output Schema。

## 安全测试

- `../`。
- absolute path。
- symlink outside。
- `.env`。
- 2MB file。
- 100k files。
- regex-like query。
- null byte。
- Unicode path。
- Tool result injection heading。
- stdout log。
- child timeout。
- environment Secret。

测试 Server 进程看到的环境变量 allowlist。

## Inspector 验收

```sh
npm run build
npx @modelcontextprotocol/inspector node dist/index.js
```

检查：

- initialize。
-三类 list。
- resource。
- tool valid/invalid。
- prompt arguments。
- stderr。

CLI smoke 可放 CI，但单元安全测试仍独立。

## 打包与 Host 配置

- 绝对入口路径。
- lockfile。
- `npm ci`。
- build artifact hash。
- executable allowlist。
- 不使用 `npx` 每次下载 Server 作为生产启动。
- Server 升级需 catalog diff。

本地安装也有供应链风险。审查依赖、最小权限和签名/来源。

## 观测

stderr JSON：

- event。
- request ID。
- method/tool。
- duration。
- result code。
- scanned file count。

脱敏：

- query 可 hash/截断。
- path 使用 workspace-relative。
- 不记录内容。
- 不记录 env。

## 故障恢复

- malformed protocol：进程退出，Host 标 unavailable。
- transient file busy：有限重试 read。
- permission change：返回 forbidden。
- search timeout：Tool isError，不重启整个 Server。
- Server crash：Host 可重启并 initialize。
- pending read 可重试；未来写 Tool 需幂等。

## 综合练习

实现本文 Server：

1. stdio。
2. Resource。
3. read Tool。
4. Prompt。
5. path canonicalization。
6. output validation。
7. Inspector/CLI。
8. 12 个安全负例。

### 验收标准

- `tsc` 与测试通过。
- stdout 只有 MCP。
- root 不能被 arguments 覆盖。
- Tool 不执行 shell string。
- Resource URI 不映射任意 path。
- input/output 有上限。
- Prompt 不提升为 policy。
- crash/timeout 可隔离。
- Host 配置可重现。

## 来源

- [MCP TypeScript SDK 官方仓库](https://github.com/modelcontextprotocol/typescript-sdk)（访问日期：2026-07-18）
- [Build an MCP server](https://modelcontextprotocol.io/docs/develop/build-server)（访问日期：2026-07-18）
- [MCP Tools 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)（访问日期：2026-07-18）
- [MCP Resources 2025-11-25](https://modelcontextprotocol.io/specification/2025-11-25/server/resources)（访问日期：2026-07-18）
