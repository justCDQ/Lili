# Go：Slice、Map、Struct、Interface、Pointer、Error、Package、Module 与 Generic

## 是什么

Slice 是数组视图；Map 是无序键值结构；Struct 组合固定字段；Interface 以方法集合描述行为并由类型隐式实现；Pointer 保存地址但不支持指针运算；`error` 是内置接口；Package 是编译与可见性单元；Module 定义版本化依赖边界；泛型以类型参数复用对多种类型相同的算法。

```go
type ID interface{ ~int64 | ~string }
type Repository[T any, K ID] interface { Get(context.Context, K) (T, error) }
type User struct { ID int64; Name string }

func Clone[T any](in []T) []T { return append([]T(nil), in...) }
```

## 关键特性或规则

- Slice 传递后元素共享；需要隔离就复制。区分 nil Slice 与空 Slice 的序列化契约。
- Map 读取用 `v, ok := m[k]` 区分缺失与零值；并发写需同步。
- 小接口由使用方定义，通常 1–3 个方法；不要为每个 Struct 预先造接口。
- Pointer Receiver 用于修改接收者或避免大对象复制；不要返回指向短命资源的无效语义引用。
- 每个错误都处理或显式返回；包名简短且职责内聚。
- 泛型用于真正的类型无关算法，不用于替代清晰接口。

## 常见错误与边界

Interface 值由动态类型和值组成，装有 nil 指针的 Interface 本身可能不为 nil。Map 迭代顺序未指定。Module 主版本 v2+ 进入模块路径。Go 的零值应尽量可用，但 Mutex 等值复制后不可使用。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

把本文示例放入独立包，补充正常、取消、超时、关闭、并发和失败测试，依次运行 `go test ./...`、`go test -race ./...` 与相关 Benchmark。并发示例还要在测试结束前确认所有 goroutine 可退出，并用 pprof 或 trace 记录证据。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go Specification](https://go.dev/ref/spec)（访问日期：2026-07-16）
- [Effective Go](https://go.dev/doc/effective_go)（访问日期：2026-07-16）
- [Go：Modules Reference](https://go.dev/ref/mod)（访问日期：2026-07-16）
- [Go：Generics Tutorial](https://go.dev/doc/tutorial/generics)（访问日期：2026-07-16）
