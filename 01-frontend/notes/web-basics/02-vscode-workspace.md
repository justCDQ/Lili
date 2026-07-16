# VS Code：文件、搜索、终端、插件与调试

## 是什么与为什么需要

VS Code 是以文件夹或工作区为边界的代码编辑器。Explorer 管文件，Search 跨文件检索，Terminal 运行 shell，Extensions 增加语言和工具支持，Run and Debug 控制调试器。统一在工作区内操作能让路径、任务、源码控制和调试配置一致。

## 关键特性与规则

- 用“打开文件夹”而非逐个打开文件；项目级配置存入 `.vscode/`。
- Search 支持普通文本、大小写、全词和正则，并可包含/排除 glob。
- 集成终端的初始目录通常是工作区根；终端使用系统已有 shell。
- 插件代码拥有较高权限，只安装可信发布者、检查权限与维护状态。
- 调试器依靠断点、调用栈、变量、监视表达式；配置通常在 `.vscode/launch.json`。

## 实际使用

1. `File > Open Folder` 打开项目。
2. `Cmd/Ctrl+Shift+F` 搜索 `TODO`；`Cmd/Ctrl+\`` 打开终端。
3. 在行号左侧设置断点，从 Run and Debug 启动；观察 Variables 和 Call Stack。
4. 将团队需要的插件写入 `.vscode/extensions.json` 的 `recommendations`，不要提交个人敏感设置。

```json
{"recommendations":["dbaeumer.vscode-eslint"]}
```

## 常见错误与边界

终端命令报错不等于编辑器报错，先确认终端工作目录和 shell。插件之间可能冲突；禁用后重载定位。浏览器 JavaScript 调试可能需要正确 source map。不要把 `.vscode/settings.json` 中的令牌提交到 Git。

## 补充知识

Workspace Trust 可限制不可信目录中任务、调试和插件的自动执行；Command Palette 是发现命令和快捷键的统一入口。

## 来源

- [VS Code：User interface](https://code.visualstudio.com/docs/editing/userinterface)
- [VS Code：Getting started](https://code.visualstudio.com/docs/editing/getting-started)
- [VS Code：Integrated terminal](https://code.visualstudio.com/docs/terminal/getting-started)
- [VS Code：Debugging](https://code.visualstudio.com/docs/editor/debugging)

访问日期：2026-07-16。
