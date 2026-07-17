# Go：竞态条件、数据竞争、死锁与 Goroutine 泄漏

## 学习目标

本文区分逻辑竞态和数据竞争，解释 happens-before、锁与 channel 死锁、goroutine 泄漏的根因，并用可取消流水线展示正确生命周期。

## 1. 竞态条件与数据竞争

竞态条件表示结果依赖不可控的执行顺序。数据竞争是更精确的内存模型概念：两个 goroutine 并发访问同一内存位置，至少一个写，且访问未由同步关系排序。

数据竞争一定是程序错误；没有数据竞争仍可有逻辑竞态。例如两个请求各自在锁内读取库存 1，解锁后决定购买，再分别在锁内减 1：每次内存访问受锁保护，却因检查与修改未放在同一临界区而超卖。

```go
mu.Lock()
available := stock > 0
mu.Unlock()
if available {
    mu.Lock()
    stock--
    mu.Unlock()
}
```

正确做法是在一个临界区内检查并修改，保护的是“不小于零”的不变量，而不只是变量本身。

## 2. happens-before

Go 内存模型规定同步操作建立顺序。例如 Mutex 的 Unlock 与之后成功 Lock、channel 发送与对应接收、关闭 channel 与因关闭返回零值的接收等建立特定同步关系。

goroutine 启动前的写对新 goroutine 可见，但 goroutine 完成没有自动同步；必须通过 channel、WaitGroup、锁等等待。`time.Sleep` 只延迟，不建立需要的同步契约。

无竞争程序具有可按某种顺序交错解释的保证。不要依赖“机器字写入看起来原子”替代内存顺序；使用 sync/atomic 的类型或锁。

## 3. Race Detector

`go test -race ./...` 构建插桩二进制，运行时记录冲突访问。报告包含两个访问栈及 goroutine 创建栈，先找共享内存和所有权，再修同步。

Race Detector 只发现实际执行路径。一次通过不证明所有路径无竞争；提高并发测试覆盖，在代表负载运行 `-race` 构建。它显著增加时间和内存，不用于性能数字。

不要用 sleep、扩大 channel buffer 或跳过测试“修复”报告。若报告为第三方库，仍需升级、隔离或报告，不能忽略共享状态风险。

## 4. 常见竞争模式

循环变量捕获在现代 Go 的特定循环语义已改进，但闭包仍可能捕获之后被修改的外部变量。最清楚方式是把任务值作为参数或在循环体建立局部所有权。

共享 map 并发读写、错误变量被多个 goroutine 赋值、计数器非原子 `n++`、懒初始化无同步、复用 buffer 后交给异步消费者，都是常见竞争。

slice header 与底层数组是不同内存。两个 goroutine append 同一 slice 可能同时改 header/数组；即使各持有不同 slice header，共享重叠底层区域的写仍竞争。

## 5. 死锁

死锁表示一组参与者互相等待，无法推进。常见条件包括互斥资源、持有并等待、不可抢占、循环等待。破坏其中一项可避免特定锁死。

锁顺序死锁：goroutine A 持 Lock1 等 Lock2，B 持 Lock2 等 Lock1。全局规定锁顺序并在所有路径遵守；更好是缩小共享状态，避免嵌套锁。

Mutex 不可重入，同一 goroutine二次 Lock 会等待自身。持锁调用未知 callback 可能回调同对象或长期阻塞，因此锁内复制必要状态，解锁后调用。

## 6. Channel 死锁

无缓冲发送需要接收者，有缓冲发送在满时等待。以下会永久阻塞：当前 goroutine 在启动接收者前发送无缓冲 channel；所有发送者等待满 buffer 而消费者等待另一资源；range 等待从未关闭的 channel。

只有能证明不再有发送的协调者关闭 channel。接收方随意关闭会让并发发送 panic；多个发送者可由 WaitGroup 等待后单一 goroutine close。

nil channel 的发送/接收永久阻塞。在 select 中可用 nil 禁用 case，但未初始化错误也会造成泄漏。

Go runtime 能在所有 goroutine 都无法推进等情形报告 fatal deadlock，但部分 goroutine 死锁时其他 goroutine（如 HTTP server）仍运行，runtime 不会替应用检测。

## 7. Goroutine 泄漏

goroutine 泄漏是 goroutine 已无业务价值却因阻塞或循环仍存活。它会保留栈、引用对象、timer、连接和下游资源。

常见根因：发送结果但调用者提前返回；读取永不关闭的 channel；外部 I/O 无 timeout/context；ticker 未 Stop；后台循环没有 shutdown；生产者无限快、消费者退出。

每次 `go f()` 前回答：谁拥有它、如何成功结束、如何失败、如何取消、谁等待结束、错误去哪里。无法回答就不应启动。

## 8. Buffer 不是生命周期修复

给结果 channel 加 buffer 1 可让单个 goroutine 在调用者超时后完成发送，适合明确“最多一个结果”。但任务仍需取消外部工作，数量大时 buffer 只延迟饱和。

无界 goroutine 等同隐藏队列。用并发限制、有限队列与过载策略保护资源。缓冲大小由最大在途与内存预算决定，不用“足够大”掩盖。

## 9. Context 取消

context 取消是协作信号。阻塞发送、接收和外部调用要 select/传递 context：

```go
select {
case out <- value:
case <-ctx.Done():
    return ctx.Err()
}
```

创建 WithCancel/Timeout 的调用者负责调用 cancel，释放 timer 和子引用。不要把 Context 存进长期 struct；它属于一次调用树。

取消后函数应停止启动新工作，并让已有工作有界退出。清理操作是否使用原已取消 context 要判断；必要时创建有严格短 deadline 的独立清理 context。

## 10. 完整案例：可取消平方流水线

### 10.1 契约

输入 int slice，generator 逐项发送，square 计算平方，collector 收集。任一阶段观察 context；超过 `MaxInt` 平方返回错误；所有 goroutine 在函数返回前结束。

```go
package pipeline

import (
    "context"
    "errors"
    "math"
)

func generate(ctx context.Context, values []int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, value := range values {
            select {
            case out <- value:
            case <-ctx.Done(): return
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) (<-chan int, <-chan error) {
    out := make(chan int)
    errCh := make(chan error, 1)
    go func() {
        defer close(out)
        defer close(errCh)
        for {
            select {
            case <-ctx.Done(): return
            case value, ok := <-in:
                if !ok { return }
                squared, err := checkedSquare(value)
                if err != nil { errCh <- err; return }
                select {
                case out <- squared:
                case <-ctx.Done(): return
                }
            }
        }
    }()
    return out, errCh
}

func checkedSquare(value int) (int, error) {
    if value == math.MinInt { return 0, errors.New("square overflow") }
    absolute := value
    if absolute < 0 { absolute = -absolute }
    if absolute != 0 && absolute > math.MaxInt/absolute {
        return 0, errors.New("square overflow")
    }
    return value * value, nil
}

func Run(parent context.Context, values []int) ([]int, error) {
    ctx, cancel := context.WithCancel(parent)
    defer cancel()
    numbers := generate(ctx, values)
    squared, errs := square(ctx, numbers)
    result := make([]int, 0, len(values))
    for value := range squared { result = append(result, value) }
    if err := <-errs; err != nil { return nil, err }
    if err := ctx.Err(); err != nil { return nil, err }
    return result, nil
}
```

`checkedSquare` 先处理无法取绝对值的 MinInt，再用除法边界检查乘法。错误发送到容量 1 的 errCh 后返回；buffer 防止 Run 因先等待 out 关闭而与错误发送互锁。

### 10.2 输入、输出与验证

输入 `[2,-3,4]`，各阶段顺序发送，输出 `[4,9,16]`。测试精确结果并用 `go test -race -count=100` 执行。

父 context 预先 cancel 时，Run 应返回 context.Canceled，不能返回空成功。输入含 MaxInt 时 checkedSquare 返回 overflow，square 关闭 out/errCh，Run 返回错误，defer cancel 让仍发送的 generator 退出。

### 10.3 生命周期验证

使用 WaitGroup 显式注入/记录阶段结束比比较 `runtime.NumGoroutine` 更可靠。goroutine 总数会受 testing/runtime 影响。可让 generate/square 接受 done callback，测试等待 done channel 有 deadline。

## 11. 死锁诊断

goroutine profile 展示每个栈的等待点；`go test -timeout` 失败会打印栈。mutex/block profile 用于争用/阻塞采样，不能证明循环等待图。

记录锁所有权和顺序，画 wait-for graph。不要只扩大 timeout。若持锁 I/O，先重构资源边界。

## 12. 泄漏诊断

观察 goroutine 数持续趋势与 profile 栈分组。大量栈停在同一 channel send、HTTP RoundTrip 或 timer 能定位所有权。采集前后使用相同负载，排除正常池/GC 辅助 goroutine。

生产 pprof 必须限制访问；profile 可能泄露路径、URL 和运行结构。

## 13. 调试清单

- race 报告：读两次访问栈与创建栈，找共享状态。
- 无 race 但结果错：检查多步不变量是否原子。
- 测试超时：看所有 goroutine 栈，不先加 sleep。
- channel send 堆积：谁接收、接收者是否提前返回、是否观察 context。
- goroutine 数上涨：按栈聚合，检查 I/O deadline 和 ticker。
- close panic：是否多个关闭者或仍有发送者。
- 修复后：运行 race、重复并发测试和取消/失败分支。

## 14. 练习

1. 把最终 pipeline 落为测试包，用 WaitGroup 证明三类结束路径。
2. 构造无数据竞争的库存逻辑竞态并修复临界区。
3. 构造两个锁反序死锁，用统一排序获取修复。
4. 写生产者在消费者首项返回后的泄漏，再加 context select。
5. 采集 goroutine/block/mutex profile，分别说明证据边界。

## 来源

- [Go Memory Model](https://go.dev/ref/mem)（访问日期：2026-07-17）
- [Go 官方文档：Data Race Detector](https://go.dev/doc/articles/race_detector)（访问日期：2026-07-17）
- [Go 官方博客：Pipelines and cancellation](https://go.dev/blog/pipelines)（访问日期：2026-07-17）
- [Go 标准库：context](https://pkg.go.dev/context)（访问日期：2026-07-17）
- [Go 标准库：runtime/pprof](https://pkg.go.dev/runtime/pprof)（访问日期：2026-07-17）
