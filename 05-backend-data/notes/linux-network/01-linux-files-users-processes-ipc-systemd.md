# 文件、权限、用户、进程、Signal、Pipe、Socket 与 systemd

## 是什么

Linux 把普通文件、目录、设备和部分内核接口暴露为文件。权限由 owner/group/other 的读写执行位和 ACL 控制；进程拥有 PID、凭据、地址空间与文件描述符；signal 是异步通知；pipe/socket 是进程通信端点；systemd unit 描述服务生命周期。

## 为什么需要

这些对象决定服务能否读取配置、绑定端口、通信、启动和停止，也是最常见的部署故障边界。

## 关键特性或规则

目录执行位控制进入/遍历；最小权限运行服务；SIGKILL 不能捕获清理；pipe 单向字节流；Unix socket 受文件权限保护；systemd 用 ExecStart、User、Restart、TimeoutStopSec 管理服务。

## 实际怎么使用

```sh
id
stat /srv/app/config.yaml
chmod 640 /srv/app/config.yaml
chown app:app /srv/app/config.yaml
kill -TERM "$PID"
systemctl status myapp.service
journalctl -u myapp.service --since '10 min ago'
```

## 常见错误与边界

chmod 777 掩盖所有权问题；删除仍被进程打开的文件不会立即释放空间；SIGTERM 后服务仍需限时优雅退出；unit 修改后需 daemon-reload。

## 补充知识

/proc 提供进程与内核视图；umask 影响新文件默认权限；capabilities 可拆分 root 特权。

## 来源

- [一手资料 1](https://man7.org/linux/man-pages/man7/inode.7.html)（访问日期：2026-07-16）
- [一手资料 2](https://man7.org/linux/man-pages/man7/credentials.7.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man7/signal.7.html)（访问日期：2026-07-16）
- [一手资料 4](https://man7.org/linux/man-pages/man7/pipe.7.html)（访问日期：2026-07-16）
- [一手资料 5](https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html)（访问日期：2026-07-16）
