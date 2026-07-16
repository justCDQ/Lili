# 哈希表、冲突、扩容与负载因子

## 是什么

哈希表用哈希函数把键映射到桶。不同键落到同一桶是冲突，常用链地址或开放寻址处理。负载因子通常是元素数/桶数；过高增加冲突，扩容会分配新表并重新放置元素。

平均查找、插入、删除 O(1)，最坏 O(n)；空间 O(n)。前提是哈希分布良好且负载受控。

```go
func Frequencies(xs []string) map[string]int {
    m := make(map[string]int, len(xs))
    for _, x := range xs { m[x]++ }
    return m
}
```

## 关键特性或规则

- 需要按键快速查找、去重、计数、缓存时使用。
- 键必须有稳定相等关系；Go Map 键必须 `comparable`，Slice/Map/Function 不能作键。
- 预知规模时提供容量提示以减少扩容，但这不是硬容量限制。
- 不依赖 Go Map 迭代顺序；需要有序输出时提取并排序键。
- 不并发读写普通 Go Map；用锁、所有权单 goroutine 或合适并发结构。

## 常见错误与边界

哈希表不适合范围查询或有序遍历。攻击者控制键时需考虑哈希碰撞 DoS。简化 HashMap 实现应测试同键覆盖、冲突链、删除、扩容后可达性和负载阈值。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 实际怎么使用

先写输入规模 n、关键操作频率、内存上限和正确性不变量，再运行本文实现。为 n=0、1、重复值、最坏分布和大规模随机数据编写测试；用 JavaScript `performance.now()` 或 Go Benchmark 测量，并核对结果是否符合给出的时间与空间复杂度。

## 补充知识

Big-O 忽略常数、缓存局部性、分配和 I/O，不能直接预测实际延迟。生产选择还需比较标准库实现、最坏输入、稳定性、并发访问和数据是否超过内存。

## 来源

- [Go Specification：Map types](https://go.dev/ref/spec#Map_types)（访问日期：2026-07-16）
- [Effective Go：Maps](https://go.dev/doc/effective_go#maps)（访问日期：2026-07-16）
- [Open Data Structures：Hash Tables](https://opendatastructures.org/ods-java/5_Hash_Tables.html)（访问日期：2026-07-16）
