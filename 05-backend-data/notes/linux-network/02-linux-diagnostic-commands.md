# ps、top、htop、lsof、ss、curl、grep、awk 与 sed

现场诊断的目标是用最小影响取得可重复证据，逐层回答“哪个进程、哪个资源、哪个端点、哪个请求、哪个时间窗发生了什么”。

## 1. 诊断前的安全边界

先记录主机、容器或网络命名空间、时区、开始时间、影响对象和命令。读取不等于零风险：`strace`、抓包、完整环境变量、命令行和日志可能泄漏密码、令牌或用户数据；高频采样和无边界递归搜索会消耗 CPU、I/O 与终端内存。

安全规则：

- 先限定 PID、端口、文件、时间窗和最大输出行数。
- 默认不使用 `sudo`；权限不足时说明缺失证据，再申请最小权限。
- 原始证据保存到权限为 `0600` 的受控目录，分享前脱敏。
- 不用 `kill`、`rm`、`truncate`、`sysctl -w` 等改变状态的动作代替诊断。
- 命令参数使用 `--` 终止选项，并给变量加双引号。

## 2. ps：进程快照

`ps` 读取某一时刻的进程信息。GNU/procps `ps` 可选择字段和排序：

```sh
ps -eo pid,ppid,user,stat,nlwp,%cpu,%mem,rss,etime,cmd --sort=-%cpu | sed -n '1,21p'
ps -L -p "$PID" -o pid,tid,psr,stat,%cpu,wchan:24,comm
```

| 字段 | 含义与边界 |
|---|---|
| `PID/PPID` | 进程与父进程 ID；会复用，需同时核对启动时间和命令 |
| `STAT` | `R` 可运行、`S` 可中断睡眠、`D` 不可中断睡眠、`T` 停止、`Z` zombie；附加字符有独立含义 |
| `%CPU` | procps 常表示生命周期内 CPU 时间比例，不等于最近一秒瞬时值 |
| `RSS` | 当前驻留物理内存的近似值，不等于独占内存或泄漏量 |
| `ETIME` | 启动后经过时间，用于排除刚重启进程 |
| `WCHAN` | 睡眠时等待的内核符号，受权限和符号可用性影响 |

macOS 的 BSD `ps` 不支持 GNU `--sort` 和部分字段，可用 `ps -axo pid,ppid,user,state,%cpu,%mem,etime,command | sort -k5 -nr`，但排序首行和字段语义需单独处理。

## 3. top 与 htop：连续采样

`top` 周期采集进程和系统状态。先观察负载、CPU 各状态、内存和任务数，再按 PID/线程定位。交互快捷键会因实现不同而变化，批处理更适合留证：

```sh
top -b -n 3 -d 1 -p "$PID" > /tmp/top-sample.txt
```

第一轮 CPU 差值可能覆盖工具启动前较长窗口，比较后续轮次。`us` 是用户态 CPU，`sy` 是内核态，`id` 是空闲，`wa` 是 CPU 空闲时至少有未完成 I/O 的时间比例；`wa` 不是“磁盘利用率”。

`htop` 提供树、线程和筛选界面，但输出不适合稳定脚本解析。诊断记录应保留选择条件与时间，而不是只截一张无法复现的图。

## 4. lsof：从进程到打开对象

`lsof` 枚举进程打开的文件，包括普通文件、目录、pipe 和 socket：

```sh
lsof -nP -p "$PID" | sed -n '1,80p'
lsof -nP -a -p "$PID" -iTCP
lsof +L1 | sed -n '1,40p'
```

`-nP` 避免 DNS 和服务名解析，减少延迟和歧义；`-a` 表示条件相与，否则多个选择项可能是相或。`+L1` 查找链接数小于 1 的已打开文件，可定位“文件已删除但磁盘未释放”。完整枚举可能需要 root，但提升权限也会看到更多敏感路径。

macOS 也常有 `lsof`；Linux `/proc/$PID/fd` 是另一条直接路径，macOS 没有兼容的 `/proc`。

## 5. ss：监听、连接与 TCP 内部状态

`ss` 是 Linux iproute2 工具，读取 socket 状态：

```sh
ss -lntp
ss -nt state established '( sport = :8080 or dport = :8080 )'
ss -tin dst 192.0.2.10
ss -s
```

- `-l` 只看监听，`-a` 看监听和非监听。
- `-n` 保留数值地址/端口，避免名称解析。
- `-t/-u/-x` 分别选择 TCP、UDP、Unix socket。
- `-p` 显示进程，需要足够权限。
- `-i` 展示 TCP 内部信息，如 RTT、拥塞窗口、重传；字段依内核版本和状态而异。

只看到 `0.0.0.0:8080` 表示该网络命名空间内所有本地 IPv4 地址监听，不证明防火墙、路由、容器端口映射或外部负载均衡可达。macOS 没有 `ss`，可用 `netstat -anv`、`lsof -nP -iTCP -sTCP:LISTEN`，参数不可直接照搬。

## 6. curl：分解请求阶段

`curl` 可以验证 DNS、连接、TLS、HTTP 状态与响应体。探活命令必须同时限制连接和总时长：

```sh
curl --fail-with-body --silent --show-error \
  --connect-timeout 3 --max-time 10 \
  --output /tmp/health-body.txt \
  --write-out 'code=%{response_code} remote=%{remote_ip} dns=%{time_namelookup} connect=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer} total=%{time_total}\n' \
  https://example.com/health
```

时间单位是秒。`time_connect`、`time_appconnect` 等是从请求开始累计的时间，不应直接相加；阶段耗时需做差。`--fail-with-body` 让 HTTP 400 及以上返回失败退出码并保留响应体，但仍需验证业务 JSON。`-v` 会输出 header，Authorization/Cookie 可能泄漏；公开证据优先使用 `--trace-config ids,time` 等受控方式并脱敏。

`--resolve example.com:443:127.0.0.1` 可绕过 DNS同时保留 URL 主机名、SNI 和 Host，用于区分 DNS 与服务问题。它不是生产修复，测试完应删除。

## 7. grep：按行筛选

```sh
grep -F -- 'request_id=abc-123' app.log
grep -nE -- 'timeout|connection reset' app.log | sed -n '1,50p'
grep -RIl --exclude='*.gz' -- 'deprecated_key' ./config
```

`-F` 把模式当固定字符串，避免 `.`、`[` 等正则元字符意外扩展；`-E` 使用扩展正则；`-n` 显示行号；`-I` 跳过二进制文件。模式或文件名可能以 `-` 开头，因此使用 `--`。对超大压缩日志、活跃写入文件或敏感目录要先限定范围。

## 8. awk：按记录与字段计算

`awk` 默认按行读记录、按空白分字段，变量和数值转换是动态的：

```sh
awk '$9 ~ /^[0-9]+$/ {count[$9]++} END {for (code in count) print code, count[code]}' access.log | sort -n
```

真实访问日志若包含带空格的引号字段，简单 `$9` 可能不是状态码。稳定分析优先输出 JSON 日志并用 JSON 解析器；若使用 awk，要声明日志格式、locale 和字段分隔符。

```sh
LC_ALL=C awk -F '\t' 'NR > 1 && $3+0 >= 1 {sum += $3; n++} END {if (n) printf "mean_ms=%.2f\n", sum/n; else print "no_data"}' metrics.tsv
```

`LC_ALL=C` 固定排序、字符类和小数点行为。必须处理空输入，避免除零或输出误导性的 0。

## 9. sed：选择与流式变换

```sh
sed -n '20,40p' app.log
sed -E 's/(Authorization: Bearer )[A-Za-z0-9._-]+/\1[REDACTED]/g' trace.txt
```

默认 sed 会打印每个处理后的 pattern space，`-n` 配合 `p` 只输出选择内容。不同实现的就地编辑参数差异明显：GNU 常用 `sed -i.bak`，BSD/macOS 常见 `sed -i ''`。诊断期间不要对原始日志就地修改；输出到新文件并限制权限。

## 10. 工具组合与退出状态

管道中默认退出状态常来自最后一个命令，前面的 `grep` 或 `curl` 失败可能被掩盖。Bash 脚本使用：

```sh
set -o pipefail
grep -F -- 'request_id=abc-123' app.log | sed -n '1,20p'
```

`grep` 退出 0 表示匹配、1 表示无匹配、2 表示错误；“无匹配”不一定是脚本故障。需要按业务语义显式分支。不要解析为人类展示而设计、受版本和 locale 影响的列；自动化长期采集优先 `/proc`、结构化 API、监控 exporter 或工具的机器格式。

## 11. 完整案例：本机健康检查间歇超时

### 输入

- 服务应监听 `127.0.0.1:8080`。
- 监控在 10:20–10:25 报告 2 秒超时。
- 不能重启、抓取完整请求体或查看其他用户进程。

### 步骤

1. `date -Ins` 记录时区和开始时间。
2. `ss -lntp 'sport = :8080'` 确认监听地址、PID；无进程字段时不立即 sudo。
3. `ps -p PID -o pid,ppid,user,stat,%cpu,%mem,rss,etime,cmd` 核对身份与是否刚重启。
4. 连续三次运行带 `--write-out` 的 curl，保存状态码、connect、TTFB 和 total，不保存含用户数据响应体。
5. 若 connect 快而 TTFB 慢，检查该 PID 的线程/CPU；若 connect 本身慢或拒绝，检查 backlog、监听和重启日志。
6. 用 `journalctl -u service --since '10:20' --until '10:26' --no-pager` 限定日志，再以固定 request ID 关联。

### 输出

证据显示 connect 约 1 ms，TTFB 在故障请求为 2 s，进程存在且 CPU 不高；应用日志同一 request ID 显示下游数据库超时。结论是应用等待下游，不是 DNS 或监听失败。

### 验证

在测试环境注入同样下游延迟，curl 的 connect 保持稳定而 TTFB 上升；移除注入后恢复。监控应分别记录 connect、TTFB 和业务状态。

### 失败分支

若 `ss` 无监听，继续检查 unit 状态、网络命名空间和启动日志；不要因为 curl `connection refused` 就修改防火墙。若只有外部访问失败而 loopback 成功，从入口逐层检查目标 IP、路由、防火墙、容器发布和代理，不把本机成功当端到端成功。

## 12. 练习与完成标准

1. 启动一个仅监听 loopback 的本地 HTTP 服务，使用 `ss`、`lsof`、`ps` 建立 PID—端口—可执行文件证据链。
2. 用 curl 输出各阶段累计时间，并正确计算 DNS、TCP、TLS、TTFB 阶段差值。
3. 构造 TSV 指标，用 awk 输出样本数和均值，对空文件输出 `no_data`。
4. 完成标准：所有命令有范围、超时、副作用说明；Linux 与 macOS 替代命令不混用；失败退出状态有分支。

## 来源

- [procps-ng：ps(1)、top(1)](https://man7.org/linux/man-pages/man1/ps.1.html)（访问日期：2026-07-17）
- [iproute2：ss(8)](https://man7.org/linux/man-pages/man8/ss.8.html)（访问日期：2026-07-17）
- [lsof manual](https://lsof.readthedocs.io/en/stable/manpage/)（访问日期：2026-07-17）
- [curl command line manual](https://curl.se/docs/manpage.html)（访问日期：2026-07-17）
- [GNU grep、gawk 与 sed manuals](https://www.gnu.org/software/grep/manual/grep.html)（访问日期：2026-07-17）
