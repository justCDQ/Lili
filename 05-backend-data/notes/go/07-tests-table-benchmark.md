# Go：Test、Table-driven Test 与 Benchmark

## 是什么

单元测试用 `testing.T` 验证行为；表驱动测试将多组输入和期望统一运行；Benchmark 用 `testing.B` 在校准迭代中测量时间与分配。

```go
func TestParse(t *testing.T) {
    tests := []struct { name, in string; want int; wantErr bool }{
        {"ok", "42", 42, false}, {"bad", "x", 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.in)
            if (err != nil) != tt.wantErr || got != tt.want { t.Fatalf("got %d err=%v", got, err) }
        })
    }
}

func BenchmarkParse(b *testing.B) {
    for b.Loop() { _, _ = Parse("42") }
}
```

## 关键特性或规则

- 测可观察契约和边界，不绑定内部实现。
- 子测试名称稳定；并行子测试不得共享可变夹具。
- Benchmark 隔离初始化，报告分配，固定输入与环境，并用 `benchstat` 比较多次样本。
- 失败输出包含输入、得到值和期望值；清理用 `t.Cleanup`。
- 性能优化前先有 Profile 与基准，优化后保留回归基准。

## 常见错误与边界

微基准不代表真实负载；编译器可能消除无效工作；系统噪声和 CPU 调频影响结果。覆盖率仅表示语句执行，不表示断言充分。

## 为什么需要

这些语言、并发和诊断能力决定 Go 服务能否正确传播错误与取消、限制资源、避免竞争并用证据优化性能。只掌握语法而缺少生命周期和工具约束，会留下难以复现的并发与泄漏问题。

## 实际怎么使用

把本文示例放入独立包，补充正常、取消、超时、关闭、并发和失败测试，依次运行 `go test ./...`、`go test -race ./...` 与相关 Benchmark。并发示例还要在测试结束前确认所有 goroutine 可退出，并用 pprof 或 trace 记录证据。

## 补充知识

Go 的调度、垃圾回收和工具都会随版本演进，性能结论必须记录 Go 版本、平台、GOMAXPROCS 与负载。Race Detector、pprof 和 trace 回答的问题不同，应按症状选择而不是同时开启。

## 来源

- [Go：testing package](https://pkg.go.dev/testing)（访问日期：2026-07-16）
- [Go Wiki：TableDrivenTests](https://go.dev/wiki/TableDrivenTests)（访问日期：2026-07-16）
- [Go：Benchmark tools](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)（访问日期：2026-07-16）
