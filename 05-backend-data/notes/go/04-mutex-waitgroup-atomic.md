# Go：Mutex、RWMutex、WaitGroup 与 Atomic

## 是什么

Mutex 为共享状态建立互斥临界区；RWMutex 区分多个读者和单个写者；WaitGroup 等待一组 goroutine 完成；Atomic 对单个机器字或专用类型提供原子操作和内存同步。

```go
type Counter struct {
    mu sync.Mutex
    n  map[string]int
}
func (c *Counter) Add(k string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n[k]++
}
```

## 关键特性或规则

- 锁保护的是不变量，不只是变量；所有访问遵循同一锁协议。
- 临界区尽量短，不在持锁时执行未知回调、网络或慢 I/O。
- 不复制已使用的 Mutex/RWMutex/WaitGroup/Atomic 值。
- WaitGroup 在启动 goroutine 前增加计数，goroutine `defer Done()`，并确保不会负计数。
- Atomic 适合简单计数、标志或指针发布；多个字段一致性仍需锁。

## 常见错误与边界

RWMutex 不是必然更快，读临界区短或争用低时开销可能更高。Mutex 不可重入；同 goroutine 二次 Lock 会死锁。WaitGroup 不传递错误或取消，任务组可使用 `errgroup`。任何优化先用 Race Detector 与基准验证。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

把本文示例放入独立包，补充正常、取消、超时、关闭、并发和失败测试，依次运行 `go test ./...`、`go test -race ./...` 与相关 Benchmark。并发示例还要在测试结束前确认所有 goroutine 可退出，并用 pprof 或 trace 记录证据。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go：sync package](https://pkg.go.dev/sync)（访问日期：2026-07-16）
- [Go：sync/atomic](https://pkg.go.dev/sync/atomic)（访问日期：2026-07-16）
- [Go Memory Model](https://go.dev/ref/mem)（访问日期：2026-07-16）
