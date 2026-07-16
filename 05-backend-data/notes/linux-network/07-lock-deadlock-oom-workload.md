# Lock、Deadlock、OOM、CPU 密集与 I/O 密集

## 是什么

锁保护共享临界区；死锁是参与者循环等待且无法推进；OOM 是可用内存/限制无法满足分配；CPU 密集工作受计算吞吐约束，I/O 密集工作主要等待外部设备或网络。

## 为什么需要

它们决定并发模型、超时、容量和故障恢复方式。

## 关键特性或规则

固定锁顺序并缩短临界区；不要持锁做慢 I/O；内存设上限和有界队列；CPU 工作限制并发接近可用核，I/O 工作仍需连接/内存上限。

## 实际怎么使用

```sh
go test -race ./...
go tool pprof -http=:0 cpu.pprof
journalctl -k | grep -i oom
cat /sys/fs/cgroup/memory.events
```

## 常见错误与边界

锁超时不自动恢复不变量；OOM killer 选择受评分和 cgroup 影响；无限加 goroutine 不能提升饱和资源吞吐。

## 补充知识

mutex profile、block profile 和 execution trace 可区分锁等待、调度与 I/O 等待。

## 来源

- [一手资料 1](https://go.dev/doc/articles/race_detector)（访问日期：2026-07-16）
- [一手资料 2](https://pkg.go.dev/runtime/pprof)（访问日期：2026-07-16）
- [一手资料 3](https://www.kernel.org/doc/html/latest/mm/oom.html)（访问日期：2026-07-16）
- [一手资料 4](https://go.dev/doc/diagnostics)（访问日期：2026-07-16）
