# 文件、目录、路径、扩展名、文本编码与压缩包

## 是什么

文件是带名称的持久数据；目录用于组织文件和子目录。路径是定位资源的字符串：绝对路径从文件系统根开始，相对路径从当前工作目录或当前文档开始。扩展名是文件名末尾用于提示格式的部分，不能保证内容真实类型。文本编码定义字符与字节的映射；Web 文本通常使用 UTF-8。压缩包把一个或多个文件封装并可压缩，ZIP、tar.gz 是常见格式。

## 为什么需要

HTML 的 `href`、`src`，命令行参数、构建配置都依赖路径。错误编码会产生乱码；隐藏扩展名会造成 `index.html.txt` 等错误；解压不可信文件可能发生路径穿越或覆盖。

## 关键规则

- POSIX 使用 `/`；Windows 常用 `\`，URL 始终使用 `/`。
- `.` 是当前目录，`..` 是父目录；路径是否区分大小写取决于文件系统。
- Web 项目使用小写、连字符、无空格文件名；不要依赖本机大小写不敏感特性。
- HTML 声明 `<meta charset="utf-8">`，文件本身也必须保存为 UTF-8。
- 扩展名不是 MIME 类型；服务器通过 `Content-Type` 告知浏览器资源类型。

## 实际使用

```text
site/
├── index.html
└── assets/
    └── logo.svg
```

`index.html` 中使用 `./assets/logo.svg`；`assets` 内文件引用首页可用 `../index.html`。ZIP 解压前先查看清单，确认没有绝对路径、`../` 或意外可执行文件。

## 常见错误与边界

不要在 HTML 中写本机绝对路径；部署后不可访问。重命名只改大小写在部分系统不会生效，可用中间名过渡。BOM 通常不必添加。压缩与归档不同：tar 主要归档，gzip 压缩单一字节流。

## 补充知识

URL 百分号编码与文件编码不是同一机制；Git 记录路径和内容，不记录空目录；可用 `.gitkeep` 约定保留空目录，但它不是 Git 特性。

## 来源

- [MDN：Dealing with files](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Environment_setup/Dealing_with_files)
- [MDN：Character encoding](https://developer.mozilla.org/en-US/docs/Glossary/Character_encoding)
- [MDN：MIME type](https://developer.mozilla.org/en-US/docs/Glossary/MIME_type)

访问日期：2026-07-16。
