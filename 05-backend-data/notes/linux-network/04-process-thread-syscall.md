# 进程、线程、用户态、内核态与系统调用

## 是什么

进程是资源与隔离容器；同进程线程共享地址空间和文件描述符但各有栈与调度状态。用户态代码权限受限，通过系统调用进入内核态请求 I/O、内存、进程等服务。

## 为什么需要

并发、崩溃隔离、系统 CPU 和阻塞问题都依赖这一执行边界。

## 关键特性或规则

系统调用有返回值和 errno；线程共享数据需要同步；一次 syscall 不等于一次物理 I/O；调度器在可运行线程间切换。

## 实际怎么使用

```sh
strace -f -tt -T -p "$PID"
cat /proc/$PID/status
ls /proc/$PID/task
```

## 常见错误与边界

把 goroutine 等同 OS 线程会误判成本；system CPU 高可能是频繁 syscall/网络/锁；strace 有开销且可能包含敏感数据。

## 补充知识

Go runtime 把 goroutine 多路复用到 OS 线程，阻塞 syscall 与网络 poller 由运行时协调。

## 来源

- [一手资料 1](https://man7.org/linux/man-pages/man2/syscalls.2.html)（访问日期：2026-07-16）
- [一手资料 2](https://man7.org/linux/man-pages/man7/pthreads.7.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man5/proc_pid_status.5.html)（访问日期：2026-07-16）
