# 复杂度、数组、链表、栈、队列与双端队列

## 是什么

时间复杂度描述输入规模 n 增长时操作次数的数量级；空间复杂度描述额外内存。Big-O 给上界，常同时说明平均/最坏情况和摊还成本。数组随机访问 O(1)，末尾追加摊还 O(1)，中间插删 O(n)；链表已知节点插删 O(1)，按索引查找 O(n)。栈后进先出，队列先进先出，双端队列两端 O(1) 插删。

```js
class Queue {
  #items = []; #head = 0;
  push(x) { this.#items.push(x); }
  shift() {
    if (this.#head === this.#items.length) return undefined;
    const x = this.#items[this.#head++];
    if (this.#head > 1024 && this.#head * 2 > this.#items.length) {
      this.#items = this.#items.slice(this.#head); this.#head = 0;
    }
    return x;
  }
}
```

## 关键特性或规则

- 需要缓存友好随机访问：数组/Slice。
- 只在两端操作：环形队列/双端队列，避免 JS `Array.shift()` 每次搬移 O(n)。
- 括号匹配、DFS、撤销：栈。
- BFS、任务按到达顺序处理：队列。
- 链表仅在节点插删和稳定引用确有价值时使用，通常有指针与缓存成本。

## 常见错误与边界

摊还 O(1) 不代表每次 O(1)；动态数组扩容会复制。递归调用栈深度受运行时限制。Go Slice 删除元素时若保留大底层数组，可能造成内存滞留。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 实际怎么使用

先写输入规模 n、关键操作频率、内存上限和正确性不变量，再运行本文实现。为 n=0、1、重复值、最坏分布和大规模随机数据编写测试；用 JavaScript `performance.now()` 或 Go Benchmark 测量，并核对结果是否符合给出的时间与空间复杂度。

## 补充知识

Big-O 忽略常数、缓存局部性、分配和 I/O，不能直接预测实际延迟。生产选择还需比较标准库实现、最坏输入、稳定性、并发访问和数据是否超过内存。

## 来源

- [Go Blog：Slices](https://go.dev/blog/slices-intro)（访问日期：2026-07-16）
- [MDN：Indexed collections](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Indexed_collections)（访问日期：2026-07-16）
- [Open Data Structures](https://opendatastructures.org/)（访问日期：2026-07-16）
