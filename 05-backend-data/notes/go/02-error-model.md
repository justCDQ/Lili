# Go：错误包装、业务错误、系统错误与可重试错误

## 是什么

Go 用普通 `error` 值表示失败。错误包装用 `%w` 保留原因链；`errors.Is` 按身份/自定义规则判断，`errors.As` 提取类型。业务错误表示请求违反领域规则；系统错误来自资源、网络或程序环境；可重试性是特定操作与上下文下的策略，不是所有错误固有属性。

```go
var ErrInsufficientStock = errors.New("insufficient stock")

func reserve(ctx context.Context, id string) error {
    if id == "" { return fmt.Errorf("reserve: invalid id") }
    if err := store.Reserve(ctx, id); err != nil {
        return fmt.Errorf("reserve %q: %w", id, err)
    }
    return nil
}
```

## 关键特性或规则

- 错误信息增加操作与关键非敏感标识，不重复写“error”。
- 不用字符串匹配判断错误类别；用 sentinel、类型或接口。
- API 边界将内部错误映射为稳定错误码，不泄漏 SQL、路径和凭据。
- 重试前确认操作幂等、错误短暂、上下文未取消，并使用上限、退避和抖动。
- `context.Canceled` 与 `DeadlineExceeded` 通常直接向上传播。

## 常见错误与边界

包装过多会暴露内部抽象；公开 sentinel 会形成兼容承诺。业务冲突通常不靠原样重试恢复。panic 用于不可恢复的程序不变量，不作为普通错误流程。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

把本文示例放入独立包，补充正常、取消、超时、关闭、并发和失败测试，依次运行 `go test ./...`、`go test -race ./...` 与相关 Benchmark。并发示例还要在测试结束前确认所有 goroutine 可退出，并用 pprof 或 trace 记录证据。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go Blog：Working with Errors](https://go.dev/blog/go1.13-errors)（访问日期：2026-07-16）
- [Go：errors package](https://pkg.go.dev/errors)（访问日期：2026-07-16）
- [Go Wiki：CodeReviewComments / Error Strings](https://go.dev/wiki/CodeReviewComments#error-strings)（访问日期：2026-07-16）
