# CPU、内存、磁盘、网络与日志排查

资源排查不是寻找一个“高数值”，而是把用户症状与同一时间窗内的利用率、饱和度、错误和延迟关联，定位最先限制请求推进的环节。

## 1. 建立可验证的排查模型

先明确开始/结束时间和时区、受影响比例、接口与实例、最近变更、正常基线、延迟目标。一次请求可能依次等待入口队列、CPU、锁、内存回收、磁盘、网络与下游；总延迟是这些阶段的结果。

```mermaid
flowchart LR
    A["用户症状"] --> B["限定时间与影响面"]
    B --> C["检查错误与饱和"]
    C --> D["按请求链分段计时"]
    D --> E["提出一个假设"]
    E --> F["用第二种证据验证"]
    F --> G["缓解并持续观察"]
```

利用率表示忙碌比例；饱和表示排队或无容量；错误表示失败事件；延迟表示工作完成时间。高利用率不一定故障，低利用率也不排除串行瓶颈、配额或锁。

## 2. CPU 与调度

CPU 时间常分 user、system、idle、iowait、steal。user 是用户态执行；system 是内核态；iowait 是 CPU 空闲且至少有未完成 I/O 的时间，不等于磁盘设备忙碌率；虚拟机 steal 表示 vCPU 等待宿主调度。

```sh
uptime
vmstat 1 5
ps -eo pid,stat,psr,nlwp,%cpu,%mem,cmd --sort=-%cpu | sed -n '1,21p'
```

`load average` 是 1、5、15 分钟指数衰减平均，Linux 计入可运行任务和不可中断睡眠任务，不能直接解释为 CPU 百分比。必须结合可用 CPU、`vmstat` 的 `r`、`b` 和 CPU 状态判断。容器看到的 CPU 数和 quota 也可能不同。

典型模式：

- user 高且运行队列增长：计算、序列化、压缩或忙循环；用 profile 找热点。
- system 高：频繁 syscall、网络包处理、上下文切换或内核路径。
- CPU 低但延迟高：检查锁、下游、配额、I/O 和队列。
- 单核满而总 CPU 低：串行热点或单线程事件循环阻塞。

## 3. 内存、page cache 与 swap

Linux 会把空闲内存用于 page cache。`free` 的 `available` 是在不 swap 的情况下可供新应用使用的估算，比单看 `free` 更有意义。

```sh
free -h
vmstat 1 5
grep -E 'MemAvailable|SwapTotal|SwapFree|Dirty|Writeback' /proc/meminfo
cat "/proc/$PID/smaps_rollup"
```

RSS 是当前驻留页总量，会包含共享页，不能跨进程简单相加；PSS 按共享者分摊共享页，更适合估算进程份额；匿名内存通常是 heap/stack 等，文件映射与 page cache 可被回收。swap in/out 的持续活动和延迟比“启用了 swap”更有诊断价值。

内存故障包括工作集超过物理或 cgroup 限制、无界缓存/队列、分配突增、内核 slab 增长、脏页回写压力、OOM kill。容器内应读 cgroup v2：

```sh
cat /sys/fs/cgroup/memory.current
cat /sys/fs/cgroup/memory.max
cat /sys/fs/cgroup/memory.events
```

`memory.max` 可能为 `max`。`memory.events` 的 `oom` 与 `oom_kill` 可区分触发与实际杀死。不要通过 `drop_caches` 当日常修复；它改变系统状态、可能制造 I/O 峰值且不解决应用增长。

## 4. 磁盘容量、inode 与 I/O

磁盘问题至少分文件系统容量、inode、设备延迟/队列、吞吐和错误：

```sh
df -hT
df -i
findmnt -T /var/lib/lili
iostat -xz 1 5
journalctl -k --since '-30 min' --no-pager | grep -Ei -- 'I/O error|reset|read-only|ext4|xfs|nvme'
```

`df` 观察文件系统；`du` 汇总可见目录项，两者不一致常见于已删除仍打开文件、挂载覆盖、权限或快照。inode 耗尽时即使仍有字节空间也不能创建文件。

`iostat` 中 `r/s`、`w/s` 是操作速率，`rkB/s`、`wkB/s` 是吞吐，`await` 是提交到完成的平均时间，队列和 `%util` 需结合设备并行能力解释；NVMe 的 `%util=100` 不必然等同已达吞吐极限。第一次报告常是自启动累计值，应看后续区间。

生产中不要在繁忙文件系统运行无范围 `du /`。优先限定挂载点，使用低优先级并在副本或低峰执行。不要在未确认恢复策略时 `fsck` 已挂载文件系统。

## 5. 网络路径与 socket

从本机到依赖分解为 DNS、路由、连接、TLS、请求发送、服务端处理、首字节和响应下载：

```sh
getent ahosts api.example.com
ip route get 192.0.2.10
ss -s
ss -tin dst 192.0.2.10
ip -s link
curl --silent --show-error --output /dev/null \
  --connect-timeout 3 --max-time 10 \
  --write-out 'code=%{response_code} dns=%{time_namelookup} conn=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer} total=%{time_total}\n' \
  https://api.example.com/health
```

丢包、重传、接口 drop、连接拒绝、超时具有不同含义。ping 使用 ICMP，可能被过滤或走不同策略；ping 成功不证明 TCP 端口或应用正常。`tcpdump` 会读取真实流量和敏感数据，必须限定主机、端口、包数、snaplen 与保存权限，并取得授权。

macOS 用 `route -n get ADDRESS`、`netstat`/`lsof` 等；没有 Linux `ip`、`ss` 和相同 `/proc`。

## 6. 日志：事件证据而非真相全集

日志应包含时间戳和时区、severity、service/version、instance、request/trace ID、操作、结果、耗时和结构化错误类型。不得记录密码、完整 token、私钥、支付卡或不必要个人数据。

```sh
journalctl -u lili-api.service \
  --since '2026-07-17 10:20:00' --until '2026-07-17 10:26:00' \
  --output=short-iso-precise --no-pager
```

日志缺失可能是未执行该路径、采样、级别过滤、缓冲未刷、进程被杀、磁盘满、时钟偏差或查询条件错误，不能直接证明事件没发生。日志量激增本身可耗尽 I/O 或额度，应有速率限制与保留策略。

## 7. 固定排查顺序

1. 确认用户可见症状和错误率，而非只看主机面板。
2. 锁定时间窗、实例、版本和变更。
3. 检查进程/容器重启、OOM、磁盘只读、内核错误等强信号。
4. 查看 CPU、内存、磁盘、网络的饱和与错误，和正常基线比较。
5. 用 trace/request ID 分解请求阶段。
6. 一次提出一个可证伪假设，并用独立证据验证。
7. 缓解后持续观察用户指标，保留回滚路径。

## 8. 完整案例：p99 上升且 CPU 不高

### 输入

- 10:20 后 `/orders` p99 从 180 ms 升到 3 s，错误率 8%。
- 总 CPU 35%，内存曲线平稳。
- 10:15 发布新版本；服务依赖 PostgreSQL。

### 步骤

1. 按版本和实例切分，确认只有新版本实例异常。
2. 查 OOM、重启、磁盘错误，均无强信号。
3. `vmstat` 显示运行队列和 swap 正常；`iostat` 显示应用磁盘正常。
4. trace 显示 2.7 s 位于数据库调用；连接池等待占 2.4 s。
5. 数据库本身查询耗时约 80 ms，连接数已到新版本池上限。
6. 对比配置发现每实例池上限从 40 错写为 4，并且请求并发约 30。
7. 回滚该配置，观察连接池 wait、p99 和错误率同时恢复。

### 输出与验证

根因是应用连接池排队，不是 CPU 或磁盘。修复后 15 分钟内 p99 回到 190 ms、池等待接近零、错误率恢复基线；压测以同样并发复现错误配置，并验证正确配置。

### 失败分支

如果数据库调用总时长高但池等待低，继续拆查询执行、锁等待、网络 RTT 和结果读取；不要仅扩大连接池，因为它可能把数据库压垮。若只有单实例异常且版本相同，检查该实例 cgroup、节点、网络路径与连接状态。

## 9. 常见误判

- “内存用了 95%”没有说明 available、回收、swap、限制和工作集。
- “load 8 很高”没有说明 CPU 数量、不可中断任务和延迟目标。
- “磁盘 util 100%”没有结合设备类型、队列、await 和业务吞吐。
- “网络慢”没有拆 DNS/connect/TLS/TTFB/download。
- “没有错误日志”没有检查采样、磁盘、进程死亡与查询窗口。
- 先重启会销毁现场证据；紧急缓解仍应先保留最小快照并记录动作。

## 10. 练习与完成标准

1. 对一个本地服务采集 5 个间隔样本，说明首个累计样本为何不可与后续直接比较。
2. 构造磁盘空间、inode、已删除打开文件三种情形的安全测试目录，并分别找到证据。
3. 为一次 HTTP 调用画出 DNS 到响应完成的阶段，并记录每段耗时。
4. 完成标准：结论至少有两种相互独立证据；每个命令限定范围；缓解后验证用户指标并说明失败分支。

## 来源

- [Linux Kernel：/proc 文件系统](https://docs.kernel.org/filesystems/proc.html)（访问日期：2026-07-17）
- [Linux Kernel：Memory Management](https://docs.kernel.org/admin-guide/mm/index.html)（访问日期：2026-07-17）
- [Linux Kernel：Block layer statistics](https://docs.kernel.org/admin-guide/iostats.html)（访问日期：2026-07-17）
- [procps-ng：vmstat(8)、free(1)](https://man7.org/linux/man-pages/man8/vmstat.8.html)（访问日期：2026-07-17）
- [systemd：journalctl](https://www.freedesktop.org/software/systemd/man/latest/journalctl.html)（访问日期：2026-07-17）
