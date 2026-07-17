# CPU、内存、磁盘与网络

## 学习目标

本文建立后端程序使用四类资源的基础模型：CPU 执行、内存保存工作集、存储持久化、网络跨主机传输。重点是用吞吐、延迟、利用率、饱和度和错误证据定位瓶颈。

## 1. 一次请求如何使用资源

典型请求会经历网络接收、解析、业务计算、缓存或数据库访问、编码和网络发送。每一阶段可能等待上一阶段，观察到的总延迟不是各资源利用率的简单结论。

```mermaid
flowchart LR
  A["网络接收"] --> B["CPU 解析与验证"]
  B --> C["内存中的对象/缓存"]
  C --> D["磁盘或远端存储"]
  D --> E["CPU 编码"]
  E --> F["网络发送"]
```

资源诊断至少区分：

- 利用率：资源忙碌的比例或已用容量。
- 饱和度：超过立即服务能力而排队的工作。
- 错误：失败、超时、重试、丢包、I/O 错误等。
- 吞吐：单位时间完成的请求、字节或操作数。
- 延迟：一个操作完成所需时间，应观察分布而非只看平均值。

平均值会掩盖尾延迟。至少同时观察 p50、p95/p99、最大值与样本量，并确认指标聚合窗口。

## 2. CPU

CPU 核执行指令。进程可有多个线程在不同核并行；调度器在可运行线程间分配时间。CPU 时间分用户态、内核态等类别，具体指标定义依操作系统。

CPU 利用率高不自动表示故障：批处理可有意用满核。问题出现在响应目标未达成、可运行队列持续增长或其他重要工作被挤压时。低 CPU 也不等于系统空闲，程序可能在等待磁盘、网络、锁或限流。

CPU 成本来源包括解析、序列化、压缩、加密、正则、复制、哈希、GC 扫描和内核协议处理。优化前用 CPU profile 找热点，不能只凭函数调用次数判断。

Go `pprof` CPU profile 采样运行栈，显示样本集中函数。采样比例是定位线索，不是精确逐指令计费；短任务应延长稳定负载或使用 benchmark。

## 3. 内存

内存为进程提供地址空间中的数据存储。虚拟内存让进程看到连续地址空间，操作系统按页把虚拟地址映射到物理内存或其他后备。RSS、虚拟大小、堆大小和容器统计口径不同，不能混用。

Go 堆由垃圾回收器管理。不再可达对象可被回收；仍可达但业务已不需要的对象仍占用内存，这类保留常被称为逻辑泄漏。短暂峰值、堆缓存和 GC 尚未归还操作系统不一定是泄漏。

需要同时观察：

- 当前存活堆和进程 RSS。
- 分配速率与每请求分配字节。
- GC 频率、CPU 开销与暂停。
- goroutine、缓存条目和队列长度。
- 容器/进程内存限制与 OOM 事件。

一次读取巨大文件会按输入规模分配；流式处理降低峰值，但增加状态管理与 I/O 次数。无界缓存和无界队列最终把流量问题转为内存问题，必须定义容量与淘汰/背压。

## 4. 磁盘与持久存储

存储设备以块处理数据，文件系统提供文件、目录和元数据。性能由介质、控制器、文件系统、缓存、I/O 模式和并发共同决定。

常见指标：吞吐量适合大顺序传输，IOPS 表示每秒操作数，延迟表示单次完成时间，队列/繁忙度反映饱和。随机小 I/O 与顺序大 I/O 即使字节数相同，成本也可能不同。

写系统调用成功可能只表示数据进入内核缓存。应用级耐久性需要明确 `fsync`/`fdatasync`、文件关闭、目录项更新和底层存储保证。原子重命名解决可见性切换，不自动等于断电后数据一定持久。

磁盘容量耗尽会使写入失败，也可能影响日志、临时文件、数据库和部署。inode 或配额也可能先耗尽。监控容量要留增长和恢复空间，不能等到 100%。

## 5. 网络

网络把字节从一个端点传到另一个端点。一次请求可能包含 DNS、建立连接、TLS 握手、发送、服务器排队处理、首字节返回与完整响应传输。

带宽是单位时间可传数据量，延迟是传输与处理等待，丢包会触发协议恢复，连接数和端口等也是有限资源。高带宽不能消除每次往返延迟；压缩降低字节数但增加双方 CPU。

网络调用不是本地函数：可能超时、连接重置、响应截断、请求已被服务端执行但客户端没收到结果。重试必须结合幂等性、截止时间、退避和重试预算，不能对所有错误立即无限重试。

观察网络阶段时分别记录 DNS、connect、TLS、首字节和总时长；只有总时长无法区分远端处理慢还是连接建立慢。响应状态成功也不保证业务内容正确，仍需协议验证。

## 6. 排队与背压

当到达速率持续高于服务速率，队列增长，等待延迟上升，最终耗尽内存或超时。增加队列容量只能延后失败，不能提高稳定吞吐。

背压让上游感知下游容量：限制并发、阻塞生产者、拒绝请求或丢弃低优先级工作。策略必须定义最大等待、拒绝错误和重试提示。

Little's Law 在稳定系统中关联平均在途工作量 L、平均到达率 λ 和平均停留时间 W：`L = λW`。使用前确认统计窗口稳定和单位一致。它可帮助检查“每秒 100 请求、平均 0.2 秒”对应约 20 个平均在途请求，但不描述尾部分布。

## 7. 性能测量原则

测量前定义目标，例如“10 KiB JSON 的 p99 < 50 ms，100 并发，无错误”。记录硬件、操作系统、Go 版本、数据规模、并发、预热和采样时间。

一次计时受缓存、调度、动态频率和后台任务影响。基准测试重复执行并报告每操作时间、分配等；生产 profile 在代表性负载下定位热点。先验证结果正确，再测性能。

一次只改变一个主要因素，并用相同输入比较。优化可能转移成本：缓存省 CPU/网络但增内存与一致性复杂度；压缩省网络和磁盘但增 CPU；批量写提升吞吐但增加单项等待和失败范围。

## 8. Go 的诊断工具

| 目标 | 工具或证据 | 注意 |
| --- | --- | --- |
| CPU 热点 | `runtime/pprof`、`net/http/pprof` | 使用代表性负载与足够时长 |
| 堆存活/分配 | heap、allocs profile | 区分 in-use 与累计分配 |
| 阻塞 | block profile | 需要适当采样设置 |
| 锁竞争 | mutex profile | 开销与采样率要评估 |
| 调度事件 | execution trace | 文件可能很大，限制窗口 |
| 微基准 | `go test -bench -benchmem` | 防止编译器消除与输入失真 |

指标先告诉“何时、影响多少”，profile/trace 再告诉“代码哪里”。日志适合离散失败上下文，不能替代时序指标。

## 9. 完整案例：整文件读取与流式处理

### 9.1 需求与输入

输入是 UTF-8 文本，每行一个事件。统计总行数、非空行和总字节。文件可能 512 MiB；单行上限 1 MiB。要求内存不随文件总大小线性增长。

整文件实现：

```go
func CountWhole(path string) (lines, nonEmpty, bytes int, err error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return 0, 0, 0, err
    }
    bytes = len(data)
    for _, line := range bytespkg.Split(data, []byte{'\n'}) {
        lines++
        if len(line) > 0 {
            nonEmpty++
        }
    }
    return lines, nonEmpty, bytes, nil
}
```

它至少保留完整 `data`，Split 还创建切片描述符；峰值与文件大小相关。并且尾随换行会使 Split 产生一个额外空段，行语义需明确。

### 9.2 流式实现

```go
package eventcount

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

type Result struct {
    Lines    int64
    NonEmpty int64
    Bytes    int64
}

func Count(reader io.Reader) (Result, error) {
    scanner := bufio.NewScanner(reader)
    scanner.Buffer(make([]byte, 64*1024), 1024*1024)
    var result Result
    for scanner.Scan() {
        result.Lines++
        result.Bytes += int64(len(scanner.Bytes()))
        if len(scanner.Bytes()) > 0 {
            result.NonEmpty++
        }
    }
    if err := scanner.Err(); err != nil {
        return Result{}, fmt.Errorf("scan events: %w", err)
    }
    return result, nil
}

func CountFile(path string) (Result, error) {
    file, err := os.Open(path)
    if err != nil {
        return Result{}, fmt.Errorf("open events: %w", err)
    }
    defer file.Close()
    return Count(file)
}
```

这里 `Bytes` 只统计行内容，不含换行分隔符，契约要如此说明。Scanner 复用缓冲区，`scanner.Bytes()` 只在下一次 Scan 前有效；若保存内容必须复制。

### 9.3 输入、步骤与输出

输入字节：

```text
start

ready
```

文件内容为 `start\n\nready\n`。Scanner 产生三条 token：`start`、空、`ready`。结果为 `Lines=3`、`NonEmpty=2`、`Bytes=10`。最后一个换行不产生额外第四行，这是 Scanner 的分词语义。

验证测试：

```go
func TestCount(t *testing.T) {
    got, err := Count(strings.NewReader("start\n\nready\n"))
    if err != nil { t.Fatal(err) }
    want := Result{Lines: 3, NonEmpty: 2, Bytes: 10}
    if got != want { t.Fatalf("got=%+v want=%+v", got, want) }
}
```

### 9.4 基准与 profile

```go
func BenchmarkCount(b *testing.B) {
    data := bytes.Repeat([]byte("event payload\n"), 10000)
    b.ReportAllocs()
    b.SetBytes(int64(len(data)))
    b.ResetTimer()
    for range b.N {
        got, err := Count(bytes.NewReader(data))
        if err != nil || got.Lines != 10000 {
            b.Fatalf("got=%+v err=%v", got, err)
        }
    }
}
```

执行：

```sh
go test ./...
go test -bench '^BenchmarkCount$' -benchmem -count=5 ./...
go test -bench '^BenchmarkCount$' -cpuprofile cpu.out -memprofile mem.out ./...
go tool pprof cpu.out
```

输入固定为 10000 行，每轮重新创建 Reader，不重新分配底层 data。基准验证结果，防止错误实现被当作更快。报告的 ns/op、MB/s、B/op 与 allocs/op 要连同机器和工具链记录。

仓库中的[可运行 Event Count 示例](../../examples/computer-systems/eventcount/)包含正常、超长行与 benchmark。

### 9.5 失败分支

构造超过 1 MiB 的单行，Scanner 返回 `token too long` 相关错误，Count 返回零值而非部分结果。调用者若需要部分统计，必须改变返回契约为 `{Partial Result, error}` 并明确消费者不得当成功。

文件总量可以很大，因为内存由最大 token 与固定状态控制；但 int64 计数最终仍有上限，极端长期流需要检查溢出。慢磁盘会使 CPU 利用率较低且总耗时高，应从 I/O 延迟和吞吐判断，不盲目优化循环。

若输入来自网络，必须增加 context/截止时间；`io.Reader` 接口本身没有统一取消方法，取消依赖具体连接或包装设计。

## 10. 四资源诊断案例

现象：p99 从 100 ms 升到 2 s，CPU 30%，内存稳定。不能因此判断“网络慢”。按阶段采集：请求队列等待、DNS/connect/TLS、数据库、文件 I/O、处理 CPU 和响应写入。

若并发队列从 20 增到 1000，数据库操作延迟从 20 ms 增到 1.5 s，数据库连接池全占用，则当前证据指向下游容量与排队。动作是限制并发、检查慢查询/数据库资源、设置截止时间，而不是先增加应用内存。

若 CPU 95%、队列增长、profile 显示 JSON 编码占 60% 样本，可从减少重复编码、调整响应结构或评估更合适编码入手，并比较端到端延迟与兼容性。

若 RSS 持续增长但 heap in-use 稳定，检查线程栈、映射、C 分配、文件缓存统计口径和内存归还，不能只称为 Go 堆泄漏。

## 11. 诊断清单

- CPU 高：确认按核口径、可运行队列、吞吐和 profile 热点。
- CPU 低但慢：检查 I/O、锁、队列、限流、远端和 goroutine 状态。
- 内存上涨：区分 RSS、heap in-use、累计分配、缓存与 goroutine。
- 磁盘慢：区分容量、IOPS、吞吐、延迟、队列和顺序/随机模式。
- 网络慢：拆 DNS、连接、TLS、首字节、下载与服务端时间。
- 重试放大流量：统计原请求与重试次数，设置预算、退避和幂等条件。
- 优化后尾延迟变差：比较完整分布、错误率和资源转移，不只看平均吞吐。

## 12. 练习

1. 创建流式案例和整文件版本，用 1 MiB、64 MiB 输入比较 B/op 与峰值堆。
2. 为 Count 增加 context 可取消的读取设计，说明如何中断具体 Reader。
3. 限制 10 个 worker 处理 1000 个任务，记录队列等待和执行时间。
4. 对同一数据测试无压缩与 gzip，比较 CPU 时间、字节数与总延迟。
5. 构造慢 Reader，每次只返回少量字节，验证 Count 不假设一次 Read 填满缓冲。

## 来源

- [Go 官方文档：Diagnostics](https://go.dev/doc/diagnostics)（访问日期：2026-07-17）
- [Go 官方文档：Profiling Go Programs](https://go.dev/blog/pprof)（访问日期：2026-07-17）
- [Go 标准库：runtime/pprof](https://pkg.go.dev/runtime/pprof)（访问日期：2026-07-17）
- [Linux Kernel Documentation：Memory Management](https://docs.kernel.org/admin-guide/mm/index.html)（访问日期：2026-07-17）
- [Linux Kernel Documentation：Block layer statistics](https://docs.kernel.org/admin-guide/iostats.html)（访问日期：2026-07-17）
