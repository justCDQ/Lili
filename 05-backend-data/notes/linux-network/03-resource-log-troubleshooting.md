# CPU、内存、磁盘、网络与日志排查

## 是什么

资源排查是把症状映射到可测指标：CPU 使用/负载/调度，内存 RSS/page cache/swap/OOM，磁盘容量/inode/延迟/吞吐，网络丢包/重传/连接，日志事件与请求上下文。

## 为什么需要

服务慢或不可用可能由任何一层饱和，单看应用错误无法定位。

## 关键特性或规则

先确认时间窗和影响面，再看饱和、错误、延迟；比较基线；CPU 分 user/system/iowait；磁盘容量和 inode 都检查；日志用 request ID 串联。

## 实际怎么使用

```sh
uptime
vmstat 1 5
free -h
df -h; df -i
iostat -xz 1 5
ss -s
journalctl -k --since today | grep -Ei 'oom|error|reset'
```

## 常见错误与边界

Linux free 内存低不等于不足，available 更有意义；load average 包含不可中断任务；网络 ping 正常不代表应用端口正常；日志没有事件不等于没有故障。

## 补充知识

建立从 DNS→路由→端口→进程→依赖→资源→日志的固定 runbook。

## 来源

- [一手资料 1](https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html)（访问日期：2026-07-16）
- [一手资料 2](https://man7.org/linux/man-pages/man1/free.1.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man1/vmstat.1.html)（访问日期：2026-07-16）
- [一手资料 4](https://www.freedesktop.org/software/systemd/man/latest/journalctl.html)（访问日期：2026-07-16）
