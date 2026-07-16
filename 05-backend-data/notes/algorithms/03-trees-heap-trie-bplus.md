# 二叉树、搜索树、堆、Trie 与 B+Tree

## 是什么

二叉树每节点最多两个子节点。平衡搜索树按键有序，查找/插删 O(log n)，退化树最坏 O(n)。二叉堆维持父子优先级，查看最值 O(1)、插入和弹出 O(log n)、建堆 O(n)。Trie 按字符串前缀逐字符走边，操作 O(L)，空间依总字符和分支。B+Tree 高分支、数据在叶节点并串联，适合磁盘页、范围扫描，操作约 O(log_B n)。

```go
h := &IntHeap{4, 1, 7}
heap.Init(h)             // O(n)
heap.Push(h, 2)          // O(log n)
min := heap.Pop(h).(int) // O(log n)
```

## 关键特性或规则

- 有序映射和范围：平衡搜索树。
- Top-K、调度、合并有序流：堆；Top-K 用大小 k 的堆为 O(n log k)、空间 O(k)。
- 前缀查询：Trie，但需评估内存。
- 数据库索引和范围扫描：B+Tree；应用层通常使用数据库，不自行实现。

## 常见错误与边界

二叉搜索树不是自动平衡。堆只保证根最优，不提供全局有序遍历。Trie 对 Unicode 应明确按字节、rune 或规范化字符。数据库索引高度低，但随机写会引发页分裂与写放大。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 实际怎么使用

先写输入规模 n、关键操作频率、内存上限和正确性不变量，再运行本文实现。为 n=0、1、重复值、最坏分布和大规模随机数据编写测试；用 JavaScript `performance.now()` 或 Go Benchmark 测量，并核对结果是否符合给出的时间与空间复杂度。

## 补充知识

Big-O 忽略常数、缓存局部性、分配和 I/O，不能直接预测实际延迟。生产选择还需比较标准库实现、最坏输入、稳定性、并发访问和数据是否超过内存。

## 来源

- [Go：container/heap](https://pkg.go.dev/container/heap)（访问日期：2026-07-16）
- [Open Data Structures：Binary Trees](https://opendatastructures.org/ods-java/6_Binary_Trees.html)（访问日期：2026-07-16）
- [PostgreSQL：B-Tree Indexes](https://www.postgresql.org/docs/current/btree.html)（访问日期：2026-07-16）
