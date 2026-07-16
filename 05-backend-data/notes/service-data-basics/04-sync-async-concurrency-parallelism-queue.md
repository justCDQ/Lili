# 同步、异步、并发、并行与队列

## 是什么

同步调用在结果可用前不返回；异步调用把完成通知延后。并发是多个任务的生命周期重叠，并行是同一时刻在多个执行资源上运行。队列缓冲生产者和消费者的速率差，保存待处理工作。

```go
jobs := make(chan int, 8)
var wg sync.WaitGroup
for i := 0; i < 2; i++ {
    wg.Add(1)
    go func() { defer wg.Done(); for id := range jobs { process(id) } }()
}
for id := 0; id < 10; id++ { jobs <- id }
close(jobs); wg.Wait()
```

## 关键特性或规则

- 异步不自动并行；并发不自动更快。
- 队列必须有限容量，并定义满时阻塞、拒绝或丢弃策略。
- 任务要处理取消、超时、重复、失败重试和结果归属。
- 并行度受 CPU、I/O、远端配额和内存约束；必须测量。
- 共享可变状态需要同步或所有权转移。

## 常见错误与边界

无界 goroutine/Promise 会耗尽资源；关闭 Channel 由发送方协调，向已关闭 Channel 发送会 panic；仅靠进程内队列不能耐受进程崩溃。需要可靠任务时使用持久队列和幂等消费者。

## 为什么需要

这些概念组成一次服务请求从寻址、传输、处理到持久化的最小闭环。缺少其中任一层的明确契约，都会让连接失败、协议错误、并发问题或数据不一致难以定位。

## 实际怎么使用

运行本文 Go 服务或数据库示例，使用 curl 发出正常、非法方法、错误 JSON、超大正文和并发请求。逐层记录 DNS/地址、端口、请求 Header、状态码、日志、数据变化和错误恢复，并为核心处理函数添加测试。

## 补充知识

本地成功只验证单进程与本机网络条件。进入容器、代理或远端数据库后，还要显式处理超时、连接池、取消、重试、幂等、事务边界和敏感日志。

## 来源

- [Go：Memory Model](https://go.dev/ref/mem)（访问日期：2026-07-16）
- [Effective Go：Concurrency](https://go.dev/doc/effective_go#concurrency)（访问日期：2026-07-16）
- [MDN：Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise)（访问日期：2026-07-16）
