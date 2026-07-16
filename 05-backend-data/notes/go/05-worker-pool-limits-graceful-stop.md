# Go：Worker Pool、并发限制与优雅停止

## 是什么

Worker Pool 用固定数量 worker 消费任务，限制同时运行数并复用执行循环。并发限制保护 CPU、内存、连接池和下游配额。优雅停止停止接收新任务，等待在途任务完成，在截止时间后强制退出。

```go
func run(ctx context.Context, workers int, jobs <-chan Job) error {
    g, ctx := errgroup.WithContext(ctx)
    for i := 0; i < workers; i++ {
        g.Go(func() error {
            for {
                select {
                case <-ctx.Done(): return ctx.Err()
                case j, ok := <-jobs:
                    if !ok { return nil }
                    if err := handle(ctx, j); err != nil { return err }
                }
            }
        })
    }
    return g.Wait()
}
```

## 关键特性或规则

- 容量由资源预算和基准确定，不直接等于任务数。
- 队列有限；定义满载时阻塞、拒绝和超时策略。
- 每任务有 Context、幂等标识、失败分类和结果通道。
- 停止顺序：停止入口→关闭/停止生产→等待 worker→关闭资源；设置总截止时间。
- HTTP 服务使用 `Server.Shutdown`，并处理 SIGTERM/SIGINT。

## 常见错误与边界

进程内池不保证崩溃后任务保留。一个慢任务会占 worker，需任务级超时。直接关闭 jobs 需要唯一发送方协调；多个生产者可通过单独调度器关闭。停止期间仍应保留日志和指标。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

把本文示例放入独立包，补充正常、取消、超时、关闭、并发和失败测试，依次运行 `go test ./...`、`go test -race ./...` 与相关 Benchmark。并发示例还要在测试结束前确认所有 goroutine 可退出，并用 pprof 或 trace 记录证据。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go：os/signal](https://pkg.go.dev/os/signal)（访问日期：2026-07-16）
- [Go：http.Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)（访问日期：2026-07-16）
- [Go：errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)（访问日期：2026-07-16）
