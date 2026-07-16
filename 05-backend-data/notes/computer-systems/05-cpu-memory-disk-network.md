# CPU、内存、磁盘与网络

## 是什么

CPU 执行指令和计算；内存保存进程正在使用的数据，低延迟但易失；磁盘持久保存数据，访问通常比内存慢；网络在主机间传输数据，受带宽、延迟、丢包和远端状态影响。

## 实际怎么使用

```go
start := time.Now()
data, err := os.ReadFile("large.json")
log.Printf("read bytes=%d duration=%s err=%v", len(data), time.Since(start), err)
```

先测量再归因：CPU 看使用率和 Profile；内存看存活对象、分配率和 GC；磁盘看吞吐、IOPS、延迟与容量；网络看 DNS、连接、吞吐、重传和远端延迟。区分资源利用率、饱和度和错误。

## 常见错误与边界

- CPU 高不一定异常，持续饱和且队列增长才说明容量不足。
- 内存泄漏是不可达/不再需要对象仍被引用，峰值高不等于泄漏。
- 磁盘写入成功可能只进入缓存；耐久性依赖刷新和文件系统语义。
- 网络不是可靠的本地函数调用，会超时、重复、乱序或部分失败。
- 优化一个资源可能增加另一个资源成本，例如压缩省网络但耗 CPU。

## 为什么需要

后端程序直接依赖操作系统资源和开发工具链。理解这些对象的生命周期、权限和性能指标，才能从进程、文件、时间或资源层定位故障，并保证本地与 CI 的操作可复现。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 补充知识

容器和 CI 会改变命名空间、用户、工作目录、可用 CPU/内存和时区。任何依赖本机默认值的行为都应显式配置，并在目标运行环境再次验证。

## 来源

- [Go：Diagnostics](https://go.dev/doc/diagnostics)（访问日期：2026-07-16）
- [Linux kernel：Memory management](https://docs.kernel.org/admin-guide/mm/index.html)（访问日期：2026-07-16）
- [Linux kernel：Block layer statistics](https://docs.kernel.org/admin-guide/iostats.html)（访问日期：2026-07-16）
