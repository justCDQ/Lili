# Array、Slice、Map、Struct、Object、String 与文件

## 是什么

JavaScript `Array` 是可变长有序集合；`Map` 保存任意可比较键到值的映射；普通 `Object` 主要表示具名字段记录；`String` 是不可变的 UTF-16 码元序列。Go 数组长度属于类型，Slice 是底层数组的窗口，Map 是哈希映射，Struct 是固定字段的值类型。

```go
type User struct { ID int `json:"id"`; Name string `json:"name"` }

func load(path string) ([]User, error) {
    data, err := os.ReadFile(path)
    if err != nil { return nil, fmt.Errorf("read %s: %w", path, err) }
    var users []User
    if err := json.Unmarshal(data, &users); err != nil {
        return nil, fmt.Errorf("decode JSON: %w", err)
    }
    return users, nil
}
```

## 关键特性或规则

- Go Slice 含指针、长度、容量；切片共享底层数组，`append` 可能换数组。
- Go Map 未初始化时可读但写入会 panic；并发读写不安全。
- Go Struct 字段导出需首字母大写；JSON 标签控制序列化字段名。
- JavaScript Object 键是字符串或 Symbol；需要保持键类型时用 `Map`。
- 文件读取必须处理不存在、权限、编码、超大文件和部分写入；大文件用流式读取。

## 常见错误与边界

按 Unicode 字符处理文本时不要把 JS `length` 当用户可见字符数；Go String 是只读字节序列，遍历 `range` 解码 UTF-8 rune。写关键文件时使用临时文件、刷新并原子替换，避免进程中断留下半文件。

## 为什么需要

这些基础决定程序如何表示数据、组织控制流、处理输入输出并报告失败。掌握它们才能明确函数契约、资源边界和可测试行为，而不是只让示例在单一输入下运行。

## 实际怎么使用

运行本文代码，并至少加入正常、空值、非法输入、边界规模和外部资源失败五类用例。先写预期输出或错误，再用测试固定；对文件和命令行示例同时检查 stdout、stderr、退出码、权限和大输入。

## 补充知识

同一逻辑在 JavaScript 与 Go 中会受到不同的数值范围、集合语义、复制方式和错误模型影响。跨语言或跨进程交换数据时，应把整数范围、空值、编码和错误格式写成明确契约。

## 来源

- [Go Blog：Go Slices](https://go.dev/blog/slices-intro)（访问日期：2026-07-16）
- [Effective Go：Maps](https://go.dev/doc/effective_go#maps)（访问日期：2026-07-16）
- [Node.js：File system](https://nodejs.org/api/fs.html)（访问日期：2026-07-16）
- [Unicode Standard](https://www.unicode.org/standard/standard.html)（访问日期：2026-07-16）
