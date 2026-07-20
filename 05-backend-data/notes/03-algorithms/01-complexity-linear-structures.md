# 复杂度、数组、链表、栈、队列与双端队列

## 学习目标

本文建立算法成本分析方法，并逐项解释线性数据结构的布局、操作、不变量和边界。完成后应能根据输入规模与操作比例选择结构，而不是只背复杂度表。

## 1. 输入规模与基本操作

复杂度描述资源消耗随输入规模增长的数量级。分析前先定义 n：数组元素数、字符串字节数、图顶点数都可能是不同变量。若算法同时处理 n 条记录和每条平均长度 L，应写出 O(nL) 等多变量表达。

时间复杂度通常计关键基本操作，空间复杂度区分输入存储与额外空间。递归调用栈、临时数组、哈希表和输出都要按题目约定计入。

## 2. O、Ω、Θ 与情况分类

Big-O 给渐近上界，Big-Ω 给下界，Big-Θ 表示上下界同阶。工程文章常用 O 表示“数量级”，但仍应注明最坏、平均、期望或摊还。

```text
3n + 20       属于 Θ(n)
n² + 100n     属于 Θ(n²)
log₂n 与 log₁₀n 只差常数因子，均为 Θ(log n)
```

最坏情况对攻击输入和延迟上限重要；平均情况必须声明输入分布；期望复杂度可能来自随机哈希或随机选择；摊还复杂度把一串操作总成本平均，不承诺单次上限。

动态数组 append 常为摊还 O(1)：大多数写入常数时间，偶尔扩容分配并复制 O(n)。n 次 append 总复制量在几何扩容下为 O(n)，但某一次仍会造成延迟峰值。

## 3. 常见增长阶

| 复杂度 | n 翻倍时的趋势 | 典型操作 |
| --- | --- | --- |
| O(1) | 近似不变 | 数组已知索引访问 |
| O(log n) | 增加常数量级 | 有序数组二分 |
| O(n) | 约翻倍 | 扫描全部元素 |
| O(n log n) | 略高于翻倍 | 通用比较排序 |
| O(n²) | 约四倍 | 两两比较 |
| O(2ⁿ) | 指数增长 | 无充分剪枝的子集枚举 |

Big-O 忽略常数、缓存、分配、分支预测和 I/O。O(n) 连续数组扫描可能比 O(log n) 的多次磁盘随机访问更快；复杂度用于排除不可扩展方案，实际延迟用基准和 profile 验证。

## 4. 数组与动态数组

数组把同类型元素连续存储，可通过基地址加索引偏移定位，因此随机访问 O(1)。连续布局有缓存局部性。中间插入/删除要移动后续元素，O(n)。

Go 数组 `[N]T` 长度属于类型且赋值复制全部元素。切片 `[]T` 是底层数组视图，包含指针、长度和容量；多个切片可能共享底层存储。

```go
values := make([]int, 0, 4)
values = append(values, 10, 20, 30)
middle := values[1:]
middle[0] = 99 // values[1] 也变为 99
```

切片索引 O(1)，尾部 append 摊还 O(1)，中间插入需移动 O(n)。预知近似元素数时提供容量可减少扩容，但容量提示不是上限。

删除含指针元素后，尾部旧引用若仍留在底层数组，可能延长对象存活。可在缩短前把不再需要位置清零，或使用 `slices.Delete` 并理解当前标准库语义。

## 5. 链表

链表节点保存值和下一个节点引用；双向链表还保存前驱。已持有目标节点及必要邻接信息时插入/删除为 O(1)，但按第 i 个位置查找 O(n)。

```go
type Node[T any] struct {
    Value T
    Next  *Node[T]
}
```

“链表插入 O(1)”不包括找到插入位置的成本。若每次先从头找第 i 项，总体仍 O(n)。链表节点额外保存指针并分散分配，缓存局部性通常弱于切片。

稳定节点地址、频繁已知节点摘除、拼接链段时链表有价值。通用顺序集合通常先选择切片，除非测量和操作模式支持链表。

空链表、头尾删除、单节点、自环和节点属于哪个列表都是实现不变量。Go 标准库 `container/list` 使用内部哨兵/根结构，节点不能被多个列表同时拥有。

## 6. 栈

栈遵循后进先出。核心操作 Push、Pop、Peek 均应 O(1)。Go 常用切片末尾作为栈顶：

```go
stack = append(stack, value)
value := stack[len(stack)-1]
stack = stack[:len(stack)-1]
```

Pop 前检查非空。栈适合括号匹配、表达式求值、撤销、显式 DFS 和函数调用。递归隐式使用调用栈，深度由输入控制时可能栈增长或耗尽，显式栈便于设置上限。

栈不保证按值优先级，不能替代堆。撤销栈还要保存逆操作或不可变快照，并处理内存上限。

## 7. 队列

队列遵循先进先出，Enqueue 在尾部，Dequeue 从头部。Go 若用 `slice = slice[1:]`，单次不搬移元素，但长期可能保留大底层数组；定期压缩或用环形缓冲。

JavaScript `Array.shift()` 通常需要移动索引，不能在大 BFS 中假设 O(1)。可保存 head 索引，消费到阈值后一次 slice，获得摊还 O(1)。

队列用于 BFS、到达顺序任务和缓冲。队列必须定义容量、空/满语义、是否线程安全和阻塞策略。数据结构队列与可靠消息队列不是同一层。

## 8. 环形队列与双端队列

环形队列使用固定数组与 head、size，把逻辑索引映射为 `(head+i) % capacity`。满时拒绝、覆盖或扩容必须明确。

双端队列允许两端 push/pop，四种操作期望 O(1)。可用环形数组实现，适合滑动窗口、0-1 BFS 与工作窃取等模式。

环形不变量：`0 <= size <= capacity`；head 指向首元素；tail 可由 `(head+size)%capacity` 算出；空与满不能只用 head==tail 区分，需 size 或额外标志。

## 9. 完整案例：泛型环形队列

```go
package ringqueue

import "errors"

var ErrFull = errors.New("queue is full")
var ErrEmpty = errors.New("queue is empty")

type Queue[T any] struct {
    items []T
    head  int
    size  int
}

func New[T any](capacity int) (*Queue[T], error) {
    if capacity < 1 { return nil, errors.New("capacity must be positive") }
    return &Queue[T]{items: make([]T, capacity)}, nil
}

func (q *Queue[T]) Len() int { return q.size }
func (q *Queue[T]) Cap() int { return len(q.items) }

func (q *Queue[T]) Push(value T) error {
    if q.size == len(q.items) { return ErrFull }
    tail := (q.head + q.size) % len(q.items)
    q.items[tail] = value
    q.size++
    return nil
}

func (q *Queue[T]) Pop() (T, error) {
    if q.size == 0 {
        var zero T
        return zero, ErrEmpty
    }
    value := q.items[q.head]
    var zero T
    q.items[q.head] = zero
    q.head = (q.head + 1) % len(q.items)
    q.size--
    return value, nil
}
```

### 9.1 输入、步骤与输出

容量 3，依次 Push A、B、C；Pop 得 A；Push D 时 tail 环绕到索引 0；再 Pop 得 B、C、D。物理数组可能是 `[D,B,C]`，逻辑顺序仍 B,C,D。

```go
func TestWrapAround(t *testing.T) {
    q, _ := New[string](3)
    for _, value := range []string{"A", "B", "C"} {
        if err := q.Push(value); err != nil { t.Fatal(err) }
    }
    got, _ := q.Pop()
    if got != "A" { t.Fatalf("first=%q", got) }
    if err := q.Push("D"); err != nil { t.Fatal(err) }
    for _, want := range []string{"B", "C", "D"} {
        got, err := q.Pop()
        if err != nil || got != want { t.Fatalf("got=%q err=%v want=%q", got, err, want) }
    }
}
```

每次 Push/Pop 只执行固定次数索引、赋值和取模，为 O(1)；存储 O(capacity)。清零弹出槽避免引用元素长期滞留。

### 9.2 失败分支

容量 0 时 New 返回错误。满队列再次 Push 返回 ErrFull 且不覆盖旧数据。空队列 Pop 返回 T 零值与 ErrEmpty，调用者必须检查 error，不能用零值判断空，因为零值可能是合法元素。

该实现不是并发安全的；并发 Push/Pop 要由调用者锁保护或设计专门同步队列。添加锁会改变阻塞语义但不改变渐近复杂度。

仓库中的[可运行 Ring Queue 示例](../../examples/algorithms/ringqueue/)保存了环绕与 FIFO 测试。

## 10. 结构选择表

| 操作模式 | 首选候选 | 需要验证 |
| --- | --- | --- |
| 大量按索引读、顺序扫描 | 切片/数组 | 扩容、共享底层数组 |
| 已知节点频繁摘除 | 双向链表 | 查找成本与缓存 |
| LIFO | 切片栈 | 最大深度、清零 |
| FIFO 有界缓冲 | 环形队列 | 满策略与并发 |
| 两端操作 | deque | 实现与容量策略 |

## 11. Benchmark 方法

分别 benchmark Push/Pop 混合，而不是只测创建。固定容量并防止每轮测试状态不同。报告 ns/op、B/op、allocs/op，并确认结果正确。

比较切片队列与环形队列时使用相同序列、容量和元素类型。若数据很小，简单切片可能更快；若长时间消费，保留底层数组和搬移成本才明显。

## 12. 调试清单

- 复杂度结论不一致：重新定义 n 和被计数操作。
- append 延迟尖峰：检查扩容、复制和 GC，用预分配对比。
- 切片修改互相影响：打印 len/cap 并检查共享底层数组。
- 队列内存不降：检查 head 前元素引用和大数组保留。
- 环形顺序错：逐步记录 head、size、tail 与物理槽。
- 空值混淆：API 返回 `(value, ok/error)`，不使用零值哨兵。
- 链表反而慢：包括查找、分配和缓存成本做端到端 benchmark。

## 13. 练习

1. 为 Queue 增加 Peek，保证空队列不改变状态。
2. 实现自动扩容版本，证明 n 次 Push 的摊还 O(1)。
3. 构造含指针元素，验证 Pop 清零后对象可回收。
4. 用显式栈实现深树 DFS，设置最大节点数防止无界输入。
5. benchmark JS `shift()` 与 head-index 队列，检查不同 n 的增长趋势。

## 来源

- [Go 语言规范：Array 与 Slice types](https://go.dev/ref/spec#Array_types)（访问日期：2026-07-17）
- [Go 官方博客：Go Slices—usage and internals](https://go.dev/blog/slices-intro)（访问日期：2026-07-17）
- [Go 标准库：container/list](https://pkg.go.dev/container/list)（访问日期：2026-07-17）
- [Open Data Structures：Array-Based Lists](https://opendatastructures.org/ods-java/2_Array_Based_Lists.html)（访问日期：2026-07-17）
- [Open Data Structures：Linked Lists](https://opendatastructures.org/ods-java/3_Linked_Lists.html)（访问日期：2026-07-17）
