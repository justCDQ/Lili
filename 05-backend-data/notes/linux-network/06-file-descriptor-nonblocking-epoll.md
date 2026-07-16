# File Descriptor、阻塞/非阻塞、I/O Multiplexing 与 epoll

## 是什么

文件描述符是进程文件描述表中的整数索引，指向打开文件描述。阻塞 I/O 等待就绪；O_NONBLOCK 让操作立即返回 EAGAIN；I/O 多路复用等待多个 fd；epoll 是 Linux 可扩展就绪通知接口。

## 为什么需要

高并发网络服务需在有限线程上管理大量连接，并正确处理部分读写和资源上限。

## 关键特性或规则

fd 在进程内有范围且可复用；close 后旧整数可能指新资源；非阻塞读写需处理 EAGAIN 与部分结果；epoll 支持 level/edge triggered。

## 实际怎么使用

```sh
ulimit -n
ls -l /proc/$PID/fd | wc -l
ss -antp
# Go net/http 默认通过 runtime netpoll 使用平台多路复用
GODEBUG=schedtrace=1000 ./server
```

## 常见错误与边界

fd 泄漏会达 EMFILE；edge-trigger 未读到 EAGAIN 会丢失后续通知；可读不保证完整消息；并发 close/use 会竞态。

## 补充知识

socket、pipe、eventfd 等都可被 epoll 监视；应用层仍需 framing、超时与背压。

## 来源

- [一手资料 1](https://man7.org/linux/man-pages/man2/open.2.html)（访问日期：2026-07-16）
- [一手资料 2](https://man7.org/linux/man-pages/man2/fcntl.2.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man7/epoll.7.html)（访问日期：2026-07-16）
- [一手资料 4](https://go.dev/src/runtime/netpoll.go)（访问日期：2026-07-16）
