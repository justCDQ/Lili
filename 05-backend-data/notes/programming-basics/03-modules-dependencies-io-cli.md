# 模块、包、依赖、输入输出与命令行参数

## 是什么

模块定义可版本化的代码单元与依赖边界；包组织同一职责的源文件；依赖是当前程序调用的外部模块。标准输入、输出、错误输出是进程的三个基础数据流；命令行参数为单次进程提供配置。

```go
func main() {
    fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    input := fs.String("input", "-", "input file or - for stdin")
    _ = fs.Parse(os.Args[1:])
    fmt.Fprintf(os.Stdout, "input=%s\n", *input)
}
```

Go 用 `go.mod` 声明模块路径与最低 Go 版本，`go.sum` 校验依赖内容。JavaScript ESM 使用显式 `import`/`export`，Node 项目以 `package.json` 和锁文件记录依赖。

## 关键特性或规则

- 正常数据写 stdout，诊断信息写 stderr；失败以非零退出码结束。
- 参数应有帮助文本、默认值、类型校验和互斥规则；密码不放命令行参数。
- 锁定依赖并审查升级；不要提交生成的依赖目录。
- 包应按能力和内聚性划分，避免循环依赖与万能工具包。

## 常见错误与边界

模块路径是导入身份，发布后不要随意修改。读取 stdin 时需考虑管道没有结束、输入过大和取消。CLI 输出若供脚本消费，应提供稳定机器格式并将日志分离。

## 为什么需要

这些基础决定程序如何表示数据、组织控制流、处理输入输出并报告失败。掌握它们才能明确函数契约、资源边界和可测试行为，而不是只让示例在单一输入下运行。

## 实际怎么使用

运行本文代码，并至少加入正常、空值、非法输入、边界规模和外部资源失败五类用例。先写预期输出或错误，再用测试固定；对文件和命令行示例同时检查 stdout、stderr、退出码、权限和大输入。

## 补充知识

同一逻辑在 JavaScript 与 Go 中会受到不同的数值范围、集合语义、复制方式和错误模型影响。跨语言或跨进程交换数据时，应把整数范围、空值、编码和错误格式写成明确契约。

## 来源

- [Go：Managing dependencies](https://go.dev/doc/modules/managing-dependencies)（访问日期：2026-07-16）
- [Node.js：ECMAScript modules](https://nodejs.org/api/esm.html)（访问日期：2026-07-16）
- [POSIX：Utility conventions](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap12.html)（访问日期：2026-07-16）
