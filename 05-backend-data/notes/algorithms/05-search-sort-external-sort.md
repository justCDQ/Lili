# 二分、快速、归并、稳定排序与外部排序

## 是什么

二分查找要求单调条件或有序数据，时间 O(log n)、空间可 O(1)。快速排序平均 O(n log n)、最坏 O(n²)，通常原地但不稳定。归并排序 O(n log n)，额外空间 O(n)，稳定且适合链表/外部合并。稳定排序保持相等键的原相对顺序。

```go
i := sort.Search(len(xs), func(i int) bool { return xs[i] >= target })
if i < len(xs) && xs[i] == target { /* found */ }

sort.SliceStable(users, func(i, j int) bool {
    return users[i].Score < users[j].Score
})
```

## 关键特性或规则

- 查找阈值、首个满足条件位置：二分，但先证明谓词单调。
- 内存数组通用排序：使用标准库，不自行手写快排。
- 多键排序需要稳定性时，明确次序或使用稳定排序。
- 数据超过内存：分块读入排序为 runs，再用最小堆 k 路归并；I/O 是主要成本。

## 常见错误与边界

比较器必须满足严格弱序，否则结果不可靠。浮点 NaN、大小写、区域排序需明确。二分中点计算要避免整数溢出。外部排序需预算临时磁盘、文件句柄、失败恢复和清理。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 实际怎么使用

先写输入规模 n、关键操作频率、内存上限和正确性不变量，再运行本文实现。为 n=0、1、重复值、最坏分布和大规模随机数据编写测试；用 JavaScript `performance.now()` 或 Go Benchmark 测量，并核对结果是否符合给出的时间与空间复杂度。

## 补充知识

Big-O 忽略常数、缓存局部性、分配和 I/O，不能直接预测实际延迟。生产选择还需比较标准库实现、最坏输入、稳定性、并发访问和数据是否超过内存。

## 来源

- [Go：sort package](https://pkg.go.dev/sort)（访问日期：2026-07-16）
- [Go：slices package](https://pkg.go.dev/slices)（访问日期：2026-07-16）
- [PostgreSQL：Resource Consumption / work_mem](https://www.postgresql.org/docs/current/runtime-config-resource.html)（访问日期：2026-07-16）
