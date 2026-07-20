# 同步、异步、并发、并行与队列

## 学习目标

本文区分调用完成方式和任务执行关系，解释 goroutine、channel、锁、取消、有限队列和 worker pool，并实现一个能停止、限流、返回结果的 Go 并发处理器。

## 1. 同步与异步

同步调用在操作得到结果或失败前不把控制权交还调用者。异步调用先返回一个句柄、Promise、任务 ID 或注册回调，完成结果在之后传递。

同步/异步描述接口的完成方式，不直接描述底层是否并行。异步文件 API 可能由线程池执行阻塞系统调用；同步计算也可在内部并行后等待所有分支。

```text
同步：caller -> call -> 等待 -> result/error -> caller 继续
异步：caller -> submit -> handle -> caller 继续 ... completion(result/error)
```

异步接口必须定义结果所有权、错误传播、取消、超时和进程重启后状态。只“启动 goroutine 不管”会丢失错误并产生生命周期泄漏。

## 2. 并发与并行

并发表示多个任务生命周期重叠，可以在一个执行资源上交错；并行表示同一时刻在多个执行资源上运行。并发用于组织独立等待和响应性，并行可能提高 CPU 工作吞吐。

并发不自动更快。创建任务、调度、同步、缓存竞争和内存占用都有成本。对单个很小 CPU 操作，并发可能更慢；对大量外部 I/O，并发可隐藏等待，但会增加下游负载。

Go goroutine 是由 Go runtime 调度的并发执行单位，不等于固定操作系统线程。`GOMAXPROCS` 影响同时执行 Go 代码的处理器数，但阻塞系统调用、cgo 等还有更细行为。

## 3. 数据竞争与同步

多个 goroutine 并发访问同一变量，至少一个写且没有适当同步，会形成数据竞争。Go 内存模型定义哪些同步事件建立 happens-before。无竞争程序才有可理解的顺序一致性保证。

保护共享状态的常见方式：

- `sync.Mutex`：同一时刻一个持有者访问临界区。
- `sync.RWMutex`：允许多个读者或一个写者；不保证总比 Mutex 快。
- channel：在 goroutine 间传值并建立特定同步关系。
- `sync/atomic`：对受支持的单变量原子操作，适合小型状态。
- 所有权隔离：数据只由一个 goroutine 修改，其他通过消息请求。

“只读时不加锁”只有在没有并发写、对象已安全发布且内部不变时成立。map 只要可能并发写，所有相关访问都必须按设计同步。

## 4. channel 的语义

channel 是带元素类型的通信通道。无缓冲 channel 的发送与对应接收同步；有缓冲 channel 在缓冲未满时发送可完成，在为空时接收等待。

```go
jobs := make(chan Job, 16)
jobs <- job
received := <-jobs
```

关闭 channel 表示不会再有新值。接收方可继续读取缓冲值，之后得到元素零值和 `ok=false`；`range channel` 直到关闭且排空。向已关闭 channel 发送或再次关闭会 panic。

通常由发送方或能证明所有发送结束的协调者关闭。接收方不知道是否还有其他发送者，不应随意关闭。

nil channel 的发送和接收永久阻塞，select 中 nil channel case 永不选择；这可动态禁用 case，也可能造成泄漏。

## 5. select、取消与超时

`select` 等待多个 channel 操作；多个 case 同时可进行时伪随机选择，不保证业务优先级。`default` 使 select 非阻塞，但循环中无等待的 default 可能忙等占满 CPU。

Go 用 context 在调用树传播截止时间、取消和请求范围值。接受 context 的函数通常把它作为第一个参数，不保存到长期结构；调用者创建并调用 cancel 释放计时器资源。

```go
select {
case result := <-results:
    return result, nil
case <-ctx.Done():
    return Result{}, ctx.Err()
}
```

取消是协作式的：任务必须在等待点或循环中观察 `ctx.Done()`，下游 API 也要接收 context。调用 cancel 不会强行终止任意 goroutine。

## 6. 队列

队列保存等待处理的工作，解耦生产速率与消费时刻。队列可在进程内存、数据库或消息系统中实现，可靠性语义不同。

设计必须定义：

- 容量上限，按条数和/或字节。
- 满时阻塞、拒绝、丢弃还是降级。
- 顺序保证及分区范围。
- 投递语义：至多一次、至少一次或特定去重协议。
- 处理超时、重试次数、退避和死信。
- 任务确认时点和进程崩溃后的恢复。
- 监控深度、最老任务年龄、吞吐和失败率。

内存 channel 在进程崩溃时丢失缓冲任务，不是持久队列。至少一次队列可能重复交付，消费者必须幂等；“exactly once” 通常只在限定系统边界和事务协议下成立。

## 7. worker pool 与并发上限

worker pool 用固定数量 worker 消费任务，限制并发。数量选择受 CPU 核、任务 CPU/I/O 比、数据库连接池、远端配额、内存和延迟目标影响，应压测而非使用固定公式。

生产者把任务发送到有界 channel。当满时阻塞形成背压，或显式返回 overloaded。把每个请求无界启动 goroutine 只是把队列隐藏到调度器和内存。

结果必须关联任务 ID。只共享一个错误变量会竞态；每个任务通过结果 channel 返回，或由协调者独占聚合。

## 8. WaitGroup、锁与生命周期

WaitGroup 等待一组任务结束。计数必须在启动 goroutine 前增加，每个任务确保 Done。WaitGroup 不传错误、不取消任务；需要这些能力时组合 context 和结果 channel，或使用合适的高层结构。

锁临界区保持短，不在持锁时执行无界网络调用或把锁泄露给回调。`defer Unlock` 可保护多返回路径，但热点极短循环要基准确认开销，正确性优先。

死锁可能来自锁顺序循环、无人发送/接收 channel、等待自身、丢失 Done 或所有 goroutine 都阻塞。超时能避免无限等待，但不能修复不变量；要分析 goroutine dump 和资源所有权。

## 9. 完整案例：有限并发文本处理器

### 9.1 契约

输入是任务切片，每项含 ID 与 Text。最多 `workers` 个任务并发；队列容量等于 worker 数；结果按输入顺序返回；第一个失败取消未完成任务。Text 为 `FAIL` 时返回业务错误，其他转大写。

```go
package workerpool

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "sync"
)

type Job struct { ID, Text string }
type Result struct { ID, Output string }
type indexedJob struct { index int; job Job }
type indexedResult struct { index int; result Result; err error }

func process(ctx context.Context, job Job) (Result, error) {
    select {
    case <-ctx.Done():
        return Result{}, ctx.Err()
    default:
    }
    if job.ID == "" { return Result{}, errors.New("job id is empty") }
    if job.Text == "FAIL" { return Result{}, fmt.Errorf("job %s failed", job.ID) }
    return Result{ID: job.ID, Output: strings.ToUpper(job.Text)}, nil
}

func Run(ctx context.Context, jobs []Job, workers int) ([]Result, error) {
    if workers < 1 { return nil, errors.New("workers must be positive") }
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    jobCh := make(chan indexedJob, workers)
    resultCh := make(chan indexedResult, len(jobs))
    var wg sync.WaitGroup
    for range workers {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobCh {
                result, err := process(ctx, item.job)
                resultCh <- indexedResult{item.index, result, err}
                if err != nil { cancel() }
            }
        }()
    }

    go func() {
        defer close(jobCh)
        for index, job := range jobs {
            select {
            case jobCh <- indexedJob{index, job}:
            case <-ctx.Done(): return
            }
        }
    }()
    go func() { wg.Wait(); close(resultCh) }()

    results := make([]Result, len(jobs))
    completed := 0
    var firstErr error
    for item := range resultCh {
        completed++
        if item.err != nil && firstErr == nil { firstErr = item.err }
        if item.err == nil { results[item.index] = item.result }
    }
    if firstErr != nil { return nil, firstErr }
    if err := ctx.Err(); err != nil { return nil, err }
    if completed != len(jobs) {
        return nil, fmt.Errorf("completed %d of %d jobs", completed, len(jobs))
    }
    return results, nil
}
```

### 9.2 正常输入与输出

输入 `[{1,go},{2,http},{3,sql}]`、workers=2。两个 worker 生命周期重叠，完成顺序可能不同；indexedResult 带原索引，最终输出仍为 `GO, HTTP, SQL`。

```go
func TestRun(t *testing.T) {
    jobs := []Job{{"1", "go"}, {"2", "http"}, {"3", "sql"}}
    got, err := Run(context.Background(), jobs, 2)
    if err != nil { t.Fatal(err) }
    want := []Result{{"1", "GO"}, {"2", "HTTP"}, {"3", "SQL"}}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got=%v want=%v", got, want) }
}
```

步骤：验证并发度；创建派生 context；启动两个 worker；生产者有限队列发送；每项返回带索引结果；等待 worker 后关闭 resultCh；聚合并验证数量。

### 9.3 失败分支

输入第二项 Text=`FAIL`。处理该项的 worker返回错误并 cancel；尚未发送或开始的任务可停止，已进入 process 的任务可能完成。Run 丢弃部分结果并返回首个业务错误。

这段实现仍会让 worker 对 jobCh 中已缓冲项调用 process，process 立即检查取消并返回 context error。firstErr 可能因并发顺序成为 context canceled 而不是业务错误：如果必须稳定返回根业务错误，需要单独的 cause channel 或 `context.WithCancelCause` 并设计优先级。测试不应假设当前调度顺序。

外部 context 在开始前已取消时，生产者不发送，worker 退出后 Run 返回 ctx error。workers=0 立即返回参数错误，不建立 channel。

### 9.4 验证竞争与泄漏

运行 `go test -race -count=100 ./...` 检查实际路径竞争；在测试前后直接比较 goroutine 数并不可靠，可使用超时、阻塞 profile 和专用泄漏检测。所有 goroutine 都有终止条件：jobCh 被关闭、resultCh 在 Wait 后关闭、生产者观察 context。

仓库中的[可运行 Worker Pool 示例](../../examples/service-data-basics/workerpool/)保存了顺序结果、并发处理和失败取消测试。

## 10. 持久队列补充

真实任务系统要在“业务事务提交”和“消息发布”之间处理双写。若数据库提交成功而 publish 失败，任务丢失；反序也会发布未提交数据。transactional outbox 把业务变更与待发布行写在同一数据库事务，再由发布器重试发送。

消费者在副作用完成后确认消息；崩溃可能导致重复投递。用任务 ID 唯一约束、幂等状态机或事务把去重记录与副作用放在同一一致性边界。

## 11. 调试清单

- goroutine 持续增长：查没有接收者的发送、未取消 I/O、未关闭输入和永久 nil channel。
- 队列深度上涨：比较到达与完成速率，增加容量不能修复长期过载。
- 结果顺序变化：用索引/ID 重排，不依赖完成顺序。
- race 报告：读两个访问栈，明确共享变量所有权和同步关系。
- 超时后仍占资源：下游是否接受/检查 context，是否有不可取消阻塞。
- channel close panic：确认唯一关闭者以及所有发送者已结束。
- 重试风暴：限制次数、指数退避加抖动、全局预算并只重试可恢复错误。

## 12. 练习

1. 用 `context.WithCancelCause` 改造案例，稳定返回第一个业务失败原因。
2. 给 process 加可控阻塞，测试外部 deadline 能让 Run 在期限内结束。
3. 增加最大排队时间，满时返回 overloaded，而不是无限阻塞生产者。
4. 实现同一 job ID 幂等写入，重复投递不重复产生副作用。
5. 比较 workers=1、CPU 核数和 100 的基准，记录吞吐、p99 和内存而非只看总时间。

## 来源

- [Go Memory Model](https://go.dev/ref/mem)（访问日期：2026-07-17）
- [Go 语言规范：Go statements](https://go.dev/ref/spec#Go_statements)（访问日期：2026-07-17）
- [Go 语言规范：Select statements](https://go.dev/ref/spec#Select_statements)（访问日期：2026-07-17）
- [Go 标准库：context](https://pkg.go.dev/context)（访问日期：2026-07-17）
- [Go 标准库：sync](https://pkg.go.dev/sync)（访问日期：2026-07-17）
