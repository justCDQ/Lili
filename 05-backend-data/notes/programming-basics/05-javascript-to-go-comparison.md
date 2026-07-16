# 从 JavaScript 算法表达迁移到 Go

## 是什么

本文知识点是“从 JavaScript 算法表达迁移到 Go”所涉及的概念、运行规则和工程边界。以下内容以可执行示例说明其数据或控制语义。

## 为什么需要

JavaScript 便于快速表达算法；Go 的静态类型、值语义、显式错误和内存表示会迫使实现者明确输入边界。重写的目标不是逐行翻译，而是比较类型、整数范围、集合语义、分配和错误契约。

```js
function count(xs) {
  const m = new Map();
  for (const x of xs) m.set(x, (m.get(x) ?? 0) + 1);
  return m;
}
```

```go
func Count[T comparable](xs []T) map[T]int {
    m := make(map[T]int, len(xs))
    for _, x := range xs { m[x]++ }
    return m
}
```

两者平均时间均为 O(n)，额外空间 O(k)，k 为不同值数量。

## 关键特性或规则

- JS `Number` 与 Go `int/int64` 的范围和溢出行为。
- JS Array 动态混合类型；Go Slice 元素类型固定并共享底层数组。
- JS 异常；Go 通常返回 `error`，调用者显式处理。
- JS 对象引用；Go Struct 默认复制，指针共享可变状态。
- Map 键规则、迭代顺序和缺失值区分。

## 常见错误与边界

性能结论必须 Benchmark，不根据语法判断。Go Map 迭代顺序未指定；JS Map 保持插入顺序。Go `int` 宽度依平台，跨边界数据优先固定宽度类型并校验。

## 实际怎么使用

运行本文代码，并至少加入正常、空值、非法输入、边界规模和外部资源失败五类用例。先写预期输出或错误，再用测试固定；对文件和命令行示例同时检查 stdout、stderr、退出码、权限和大输入。

## 补充知识

同一逻辑在 JavaScript 与 Go 中会受到不同的数值范围、集合语义、复制方式和错误模型影响。跨语言或跨进程交换数据时，应把整数范围、空值、编码和错误格式写成明确契约。

## 来源

- [MDN：JavaScript data structures](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Data_structures)（访问日期：2026-07-16）
- [Go Specification](https://go.dev/ref/spec)（访问日期：2026-07-16）
- [Go：Tutorial on generics](https://go.dev/doc/tutorial/generics)（访问日期：2026-07-16）
