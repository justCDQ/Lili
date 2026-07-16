# ps、top、htop、lsof、ss、curl、grep、awk 与 sed

## 是什么

ps 取进程快照，top/htop 交互观察资源，lsof 列进程打开文件，ss 查看 socket，curl 发协议请求；grep 筛选行，awk 按字段处理，sed 做流式转换。

## 为什么需要

组合这些工具可在不改代码时确认进程、监听、连接、HTTP 和日志证据。

## 关键特性或规则

先限定时间、PID、端口和固定字符串；脚本中使用稳定列格式；curl 区分 DNS/连接/TLS/首字节总超时；生产命令保留原始输出。

## 实际怎么使用

```sh
ps -eo pid,ppid,user,stat,%cpu,%mem,etime,cmd --sort=-%cpu | head
lsof -nP -p "$PID"
ss -lntp
curl -v --connect-timeout 3 --max-time 10 http://127.0.0.1:8080/health
grep -F 'request_id=' app.log | awk '{print $1,$NF}' | sed -n '1,20p'
```

## 常见错误与边界

grep 正则元字符会误匹配，固定文本用 -F；top 瞬时高 CPU 不代表长期；lsof/ss 可能需权限；curl 成功退出不代表业务状态成功，可用 --fail-with-body。

## 补充知识

命令输出受 locale 和版本影响；长期指标应由监控采集，现场命令用于验证。

## 来源

- [一手资料 1](https://man7.org/linux/man-pages/man1/ps.1.html)（访问日期：2026-07-16）
- [一手资料 2](https://man7.org/linux/man-pages/man8/ss.8.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man8/lsof.8.html)（访问日期：2026-07-16）
- [一手资料 4](https://curl.se/docs/manpage.html)（访问日期：2026-07-16）
- [一手资料 5](https://www.gnu.org/software/gawk/manual/gawk.html)（访问日期：2026-07-16）
