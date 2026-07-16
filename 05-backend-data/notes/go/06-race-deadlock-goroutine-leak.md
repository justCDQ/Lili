# Go：Race Condition、Deadlock 与 Goroutine Leak

## 是什么

竞态条件是结果依赖不可控执行顺序；数据竞争是并发访问同一内存且至少一次写、没有正确同步。死锁是参与者循环等待而无法推进。Goroutine Leak 是 goroutine 因永久阻塞或缺少取消路径而不再有用但仍存活。

```go
func first(ctx context.Context, in <-chan string) (string, error) {
    select {
    case v, ok := <-in:
        if !ok { return "", io.EOF }
        return v, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

## 实际怎么使用

- 共享状态使用锁/原子或单一所有者；运行 `go test -race ./...`。
- 固定锁顺序，不在持锁时等待 Channel 或外部调用。
- 每个 goroutine 明确退出条件、取消信号和阻塞操作上限。
- 测试取消、超时、Channel 关闭和部分失败；观察 goroutine Profile。
- Leak 检查前后 goroutine 数只是线索，需 Profile 确认栈。

## 常见错误与边界

Race Detector 只发现实际执行路径，且有运行开销；无数据竞争不代表无逻辑竞态。Go Runtime 只能检测某些全局停滞，部分 goroutine 死锁时进程仍可能运行。增加缓冲可能隐藏而非修复泄漏。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go：Data Race Detector](https://go.dev/doc/articles/race_detector)（访问日期：2026-07-16）
- [Go Memory Model](https://go.dev/ref/mem)（访问日期：2026-07-16）
- [Go Blog：Pipelines and cancellation](https://go.dev/blog/pipelines)（访问日期：2026-07-16）
