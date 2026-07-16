# Go：Goroutine、Channel、Select 与 Context

## 是什么

Goroutine 是由 Go Runtime 调度的并发函数执行；Channel 在 goroutine 间传递类型化值并提供同步；`select` 等待多个 Channel 操作；Context 沿调用链传递取消、截止时间和请求范围值。

```go
func fetch(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)
    for _, u := range urls {
        u := u
        g.Go(func() error { return get(ctx, u) })
    }
    return g.Wait()
}
```

## 关键特性或规则

- 启动 goroutine 时明确所有者、结束条件和错误去向。
- 无缓冲 Channel 发送与接收会同步；有缓冲只提供有限解耦，不是持久队列。
- 发送方负责关闭；关闭表示不再发送，不是广播任意状态。
- `select` 配合 `ctx.Done()` 避免永久阻塞；`default` 会变成非阻塞轮询，慎用。
- Context 作为首参数，不存入 Struct，不传 nil，不放可选业务参数。

## 常见错误与边界

从关闭 Channel 接收立即返回零值，使用 `v, ok` 区分。向关闭 Channel 发送或重复关闭会 panic。取消是协作式的，下游必须观察 Context。无并发上限的 goroutine fan-out 会耗尽连接和内存。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

原示例会为每个 URL 启动 goroutine，输入很大时缺少上限。使用 `errgroup.SetLimit` 限制并发，并给整批任务设置截止时间：

```go
func FetchAll(parent context.Context, urls []string, limit int) error {
    if limit < 1 { return fmt.Errorf("limit must be positive") }
    ctx, cancel := context.WithTimeout(parent, 10*time.Second)
    defer cancel()

    group, ctx := errgroup.WithContext(ctx)
    group.SetLimit(limit)
    for _, url := range urls {
        url := url
        group.Go(func() error {
            req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
            if err != nil { return fmt.Errorf("build %q: %w", url, err) }
            resp, err := http.DefaultClient.Do(req)
            if err != nil { return fmt.Errorf("get %q: %w", url, err) }
            defer resp.Body.Close()
            if resp.StatusCode / 100 != 2 { return fmt.Errorf("get %q: status %d", url, resp.StatusCode) }
            _, err = io.Copy(io.Discard, resp.Body)
            return err
        })
    }
    return group.Wait()
}
```

用 `httptest.Server` 建立成功、500、阻塞直到取消三类端点；断言最大在途请求不超过 `limit`，首个错误会取消其他请求，超时返回 `context deadline exceeded`。运行 `go test -race`，测试结束后检查服务器请求全部退出。

## 补充知识

Channel 建立的是特定发送/接收间的 happens-before 关系，不等于所有共享状态自动安全。`select` 在多个分支同时就绪时伪随机选择，不能依赖优先级。Context 的 Value 只用于跨 API 的请求范围元数据；业务可选参数应使用普通参数或配置 Struct。

## 来源

- [Go Blog：Concurrency Patterns — Context](https://go.dev/blog/context)（访问日期：2026-07-16）
- [Go Blog：Pipelines and cancellation](https://go.dev/blog/pipelines)（访问日期：2026-07-16）
- [Go：context package](https://pkg.go.dev/context)（访问日期：2026-07-16）
