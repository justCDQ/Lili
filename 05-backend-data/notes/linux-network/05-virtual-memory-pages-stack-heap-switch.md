# 虚拟内存、Page、Stack、Heap 与上下文切换

## 是什么

虚拟内存为进程提供独立地址空间并由页表映射物理页。stack 保存调用帧和局部状态，heap 管理动态对象；上下文切换保存并恢复执行状态，让 CPU 在任务间调度。

## 为什么需要

内存占用、page fault、栈溢出、分配和并发调度成本都需要这些概念。

## 关键特性或规则

RSS 是驻留物理内存，VSZ 是虚拟地址范围；匿名页与文件映射行为不同；minor/major page fault 成本不同；stack 通常按线程分配。

## 实际怎么使用

```sh
pmap -x "$PID" | tail -n 20
cat /proc/$PID/smaps_rollup
/usr/bin/time -v ./server
```

## 常见错误与边界

虚拟地址大不等于物理占用大；heap 不等于 Go 所有活对象；上下文切换数量本身需结合吞吐与延迟解释。

## 补充知识

page cache 可复用文件数据；容器内存限制可能早于宿主物理耗尽触发 OOM。

## 来源

- [一手资料 1](https://www.kernel.org/doc/html/latest/admin-guide/mm/concepts.html)（访问日期：2026-07-16）
- [一手资料 2](https://kernel.org/doc/html/next/mm/process_addrs.html)（访问日期：2026-07-16）
- [一手资料 3](https://man7.org/linux/man-pages/man5/proc_pid_smaps.5.html)（访问日期：2026-07-16）
