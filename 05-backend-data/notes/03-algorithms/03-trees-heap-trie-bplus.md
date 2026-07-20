# 二叉树、搜索树、堆、Trie 与 B+Tree

## 学习目标

本文比较五类树结构的有序性、优先级、前缀与存储页语义，解释遍历和复杂度条件，并用最小堆完成 Top-K 案例。

## 1. 树的基本术语

树由节点和边组成。根没有父节点，叶没有子节点；深度是节点到根的边数，高度是节点到最深叶的最长边数。含 n 个节点的树有 n-1 条父子边。

二叉树每节点最多两个有方向的子节点 left/right，不保证键有序或平衡。完全二叉树除最后一层外填满，最后一层从左到右填充，适合数组表示堆。

遍历：前序根-左-右，适合复制结构；中序左-根-右，对二叉搜索树得到有序键；后序左-右-根，适合自底向上计算；层序用队列逐层访问。

递归遍历时间 O(n)，额外调用栈 O(h)。退化高度 h=n 时可能栈耗尽，大输入用显式栈并限制节点数。

## 2. 二叉搜索树

BST 不变量通常是左子树键小于节点、右子树键大于节点；重复键策略必须定义为计数、固定一侧或拒绝。

查找从根比较，每步进入一侧，成本 O(h)。平衡时 h=O(log n)，按升序插入普通 BST 可退化为链，h=O(n)。因此“BST 查找 O(log n)”只有平衡条件成立。

删除分三种：叶直接移除；单子节点用子节点替代；双子节点用中序后继/前驱替换，再删除该节点。每次变更都要保持全树不变量。

AVL、红黑树等平衡树通过旋转和额外信息约束高度。应用通常使用标准库或数据库实现，不手写生产平衡树。

## 3. 堆

二叉最小堆满足每个父节点优先级不大于子节点，根是全局最小。它不保证兄弟或不同子树整体有序，所以遍历数组不会得到排序结果。

完全二叉树可紧凑存在数组。0-based 索引：parent `(i-1)/2`，left `2i+1`，right `2i+2`。Push 在尾部插入后向上修复；Pop 把末项移到根后向下修复，均 O(log n)；Peek O(1)。

从任意数组自底向上 heapify 为 O(n)，不是逐个 Push 的 O(n log n)。

Go `container/heap` 要求类型实现 `sort.Interface` 的 Len/Less/Swap，加 Push/Pop。包通过这些方法维护堆；调用者应使用 `heap.Push/Pop`，不直接调用接收者 Push/Pop。

堆适合优先队列、定时任务、Top-K、Dijkstra 和 k 路合并。它不适合按任意键 O(1) 删除，除非额外维护索引并在交换时同步更新。

## 4. Trie

Trie 按键序列的单位逐层走边，节点表示前缀。查找/插入长度 L 的键为 O(L)，与已有键数量不直接线性相关；空间取决于所有不同前缀与每节点边表示。

```text
root
 ├─ c ─ a ─ t*
 │       └─ r*
 └─ d ─ o ─ g*
```

终止标记区分键 `car` 与前缀 `ca`。前缀查询先走到前缀节点，再遍历子树输出，成本还要加结果总大小。

字符串单位必须定义。按 UTF-8 字节实现紧凑但前缀可能落在码点中间；按 rune 避免切断码点但一个视觉字符仍可多码点。应先规范化大小写和 Unicode，并明确 Trie 保存的标准形式。

每节点用大数组存所有可能边查询快但浪费空间，用 map 节省稀疏分支但有哈希/分配成本，可使用压缩 Trie/radix tree 合并单一路径。

## 5. B-Tree 与 B+Tree

二叉树每节点少量键，存储在磁盘页时高度和随机 I/O 可能较大。B-Tree 类结构让一个节点保存多个有序键和多个子指针，提高分支因子，使高度约 O(log_B n)。

B+Tree 通常把记录或记录引用集中在叶节点，内部节点保存分隔键；叶节点按顺序链接，范围扫描定位起点后可顺序遍历叶页。不同数据库实现对页格式、重复键、并发和分裂细节不同。

PostgreSQL B-tree 索引支持可排序数据的相等与范围查询，并可用于某些 ORDER BY。文档称其结构为多层树和叶页；应用不应把通用教材 B+Tree 伪代码当 PostgreSQL 精确页实现。

插入满页可能分裂并更新父层，产生随机写、WAL 和写放大。删除可能留下可回收空间；具体 vacuum/页合并行为看数据库版本。索引不是免费排序副本。

复合索引 `(a,b)` 的利用与列约束、排序方向、选择性和实现优化有关。PG18 支持 B-tree skip scan 的特定计划机会，但只应在索引查询篇结合 EXPLAIN 说明，不把它写成所有数据库通则。

## 6. 结构选择

| 需求 | 候选 | 关键边界 |
| --- | --- | --- |
| 有序 map、范围内存查询 | 平衡搜索树 | 标准库支持、h 保证 |
| 只需当前最小/最大 | 堆 | 任意查找仍 O(n) |
| 保留 k 个最优项 | 大小 k 堆 | 比较方向、稳定 tie |
| 字符串前缀 | Trie/radix | Unicode 单位与空间 |
| 磁盘有序索引 | 数据库 B-tree | 页、事务、写成本 |

## 7. 完整案例：Top-K 最大值

### 7.1 思路

输入 n 个分数，返回最大的 k 个，按分数降序；同分按 ID 升序。维护容量 k 的“最差项在根”的最小优先堆。每来一项：堆未满则 Push；比根更好则替换根；否则丢弃。最终弹出并反转/排序。

定义“更差”：分数小更差；同分 ID 字典序大更差，因为最终希望 ID 小优先。

### 7.2 Go 实现

```go
package topk

import (
    "container/heap"
    "errors"
    "sort"
)

type Item struct { ID string; Score int }
type minHeap []Item

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
    if h[i].Score != h[j].Score { return h[i].Score < h[j].Score }
    return h[i].ID > h[j].ID
}
func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(value any) { *h = append(*h, value.(Item)) }
func (h *minHeap) Pop() any {
    old := *h
    last := old[len(old)-1]
    old[len(old)-1] = Item{}
    *h = old[:len(old)-1]
    return last
}

func better(a, b Item) bool {
    return a.Score > b.Score || (a.Score == b.Score && a.ID < b.ID)
}

func Select(items []Item, k int) ([]Item, error) {
    if k < 0 { return nil, errors.New("k must be non-negative") }
    if k == 0 { return []Item{}, nil }
    selected := make(minHeap, 0, min(k, len(items)))
    heap.Init(&selected)
    for _, item := range items {
        if item.ID == "" { return nil, errors.New("item id is empty") }
        if selected.Len() < k {
            heap.Push(&selected, item)
        } else if better(item, selected[0]) {
            selected[0] = item
            heap.Fix(&selected, 0)
        }
    }
    result := append([]Item(nil), selected...)
    sort.Slice(result, func(i, j int) bool { return better(result[i], result[j]) })
    return result, nil
}
```

### 7.3 输入、过程与输出

输入 `{a:90,b:70,c:90,d:80,e:95}`，k=3。堆先收 a,b,c；d 比根 b 好，替换；e 再替换当前最差。最终排序输出 `e:95,a:90,c:90`，同分 a 在 c 前。

复杂度：每项最多一次 O(log k) 堆操作，总 O(n log k)；最终 k 项排序 O(k log k)；空间 O(k)。k>=n 时仍做 n log n 量级并最终排序，可直接复制全量排序，复杂度同阶。

### 7.4 测试

```go
func TestSelect(t *testing.T) {
    input := []Item{{"a",90},{"b",70},{"c",90},{"d",80},{"e",95}}
    got, err := Select(input, 3)
    if err != nil { t.Fatal(err) }
    want := []Item{{"e",95},{"a",90},{"c",90}}
    if !reflect.DeepEqual(got, want) { t.Fatalf("got=%v want=%v", got, want) }
}
```

失败分支：k<0 返回错误；k=0 空结果；空 ID 失败且不返回部分选择；重复 ID 当前允许并作为不同记录，如业务要求唯一要先用 map 检查。

仓库中的[可运行 Top-K 示例](../../examples/algorithms/topk/)保存了堆实现与 tie-break 测试。

### 7.5 堆不变量验证

测试中可在每次操作后断言对所有 i>0，`Less(i,parent(i))` 为 false。只断最终 Top-K 可能漏掉中间修复 bug。benchmark 比较 n=1e6、k=10 的全排序与堆方案，记录内存。

## 8. Trie 案例要点

自动补全输入先 NFC + 明确大小写规则；节点保存 `map[rune]*node` 与 terminal。查前缀 O(L)，输出受结果数限制。必须设置最大结果、最大深度和总节点数，避免短前缀遍历整个词典。

删除键需取消 terminal，并从叶向上清理无子节点非终止节点。并发读写需要版本快照或锁，不能直接修改共享 map。

## 9. 数据库索引验证

对 `WHERE score >= 90 ORDER BY score,id LIMIT 20`，候选 `(score,id)` B-tree。用 PostgreSQL 18 真实数据运行 `EXPLAIN (ANALYZE, BUFFERS)`，同时测不同选择性与冷/热缓存。计划使用或不使用索引都是成本模型结果。

不要在应用实现 B+Tree 替代数据库索引；事务可见性、WAL、锁、恢复、页并发远超基础结构。

## 10. 调试清单

- BST 变慢：测高度和插入顺序，确认是否平衡实现。
- 堆输出“乱序”：堆只保证根，最终结果需排序。
- Top-K tie 不稳定：比较器加入确定性 ID，并保持堆“最差根”方向一致。
- Trie 查不到视觉相同词：比较码点和规范化策略。
- Trie 内存过大：统计节点/边、map 开销，评估 radix 压缩。
- 数据库索引未使用：看选择性、统计、查询形状、排序和缓存，不强制猜测。
- 页分裂写慢：检查键分布、填充、索引数量和目标数据库证据。

## 11. 练习

1. 为 Top-K 添加随机输入，与全排序前 k 项做属性对照测试。
2. 实现整数 BST 并输入升序 10000 项，记录高度与递归风险。
3. 用 `heap.Init` 与逐个 Push 建堆 benchmark，观察增长趋势。
4. 实现 rune Trie 的 Insert/Contains/Prefix，测试 NFC 等价输入。
5. 在 PostgreSQL 18 建复合索引，记录 equality/range/order 查询计划差异。

## 来源

- [Go 标准库：container/heap](https://pkg.go.dev/container/heap)（访问日期：2026-07-17）
- [Open Data Structures：Binary Trees](https://opendatastructures.org/ods-java/6_Binary_Trees.html)（访问日期：2026-07-17）
- [Open Data Structures：Binary Search Trees](https://opendatastructures.org/ods-java/6_2_BinarySearchTree_Unbala.html)（访问日期：2026-07-17）
- [PostgreSQL 18 文档：B-Tree Indexes](https://www.postgresql.org/docs/18/btree.html)（访问日期：2026-07-17）
- [PostgreSQL 18 文档：Index Types](https://www.postgresql.org/docs/18/indexes-types.html)（访问日期：2026-07-17）
