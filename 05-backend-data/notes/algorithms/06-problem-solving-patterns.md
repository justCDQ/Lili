# 双指针、滑动窗口、分治、贪心、回溯与动态规划

## 学习目标

本文说明六类解题模式各自成立的结构条件、状态与不变量，避免按题目关键词套模板。完整案例用动态规划计算带约束的最小成本，并给出可恢复选择路径。

## 1. 先定义问题

算法前写清输入、输出、规模、允许修改、重复值、排序、数值范围和失败语义。再写暴力方案与复杂度，识别哪种重复工作或单调性质可优化。

每个优化必须有正确性依据：循环不变量、交换论证、归纳、最优子结构或状态覆盖。只有样例通过不能证明算法。

## 2. 双指针

双指针用两个位置协调扫描。指针可同向、相向或快慢，关键是每次移动能安全排除一部分搜索空间。

有序数组两数和：left 指最小、right 指最大。和太小移动 left，因为固定 left 时更小 right 只会更小；和太大移动 right。每指针最多走 n 次，O(n)。无序输入先排序 O(n log n) 并保留原索引，或用哈希表 O(n) 期望时间。

快慢指针可原地去重或检测链表环。Floyd 环检测中 slow 每次一步、fast 两步；有环会在环内相遇。要找入口还需第二阶段证明，不能把相遇点直接当入口。

双指针不自动 O(n)：若一个指针每轮回退或重新扫描，总成本可能 O(n²)。分析所有移动次数。

## 3. 滑动窗口

滑动窗口维护连续区间 `[left,right]` 的增量状态。适用条件是右边界前进加入元素，违反约束时左边界单调前进移除元素，状态可高效更新。

最长至多 k 个不同值：map 保存窗口计数，加入 right；不同数超 k 时移除 left 直到合法；更新最大长度。每元素最多进入/离开一次，O(n)，map 空间 O(k+1)。

```go
func LongestAtMostK(values []string, k int) int {
    if k <= 0 { return 0 }
    counts := make(map[string]int)
    left, best := 0, 0
    for right, value := range values {
        counts[value]++
        for len(counts) > k {
            old := values[left]
            counts[old]--
            if counts[old] == 0 { delete(counts, old) }
            left++
        }
        best = max(best, right-left+1)
    }
    return best
}
```

“最短子数组和至少 S”的普通正数窗口依赖加入不减、移除不增。允许负数后单调性消失，需前缀和+单调队列等其他算法。

固定长度窗口更简单：加新项、减离开项。数值累加要检查溢出。

## 4. 分治

分治把问题拆成规模更小、通常独立的子问题，递归解决后合并。归并排序：两个 n/2 子问题加 O(n) 合并，递推 `T(n)=2T(n/2)+O(n)=O(n log n)`。

二分查找也是分治/减治，每次只保留一个一半子问题，`T(n)=T(n/2)+O(1)=O(log n)`。

设计要有基例、规模严格缩小、合并正确且不会重复处理。子问题大量重叠时直接分治可能指数级，记忆化/DP 可复用结果。

并行分治只有子问题足够大才值得启动任务；小粒度调度和合并成本会超过收益。设置阈值并 benchmark。

## 5. 贪心

贪心每步选择当前看起来最优的决策且不回退。正确需要证明局部选择能扩展到全局最优，常用交换论证、cut property 或 matroid 结构。

区间调度最大不重叠数量：按结束时间升序，每次选择第一个与已选兼容的区间。交换证明可把任一最优解第一项替换为最早结束项，不减少后续空间，递归成立。

按开始最早、持续最短或价值最大都可能失败，必须给反例。0/1 背包按价值/重量贪心不保证最优，分数背包才成立。

贪心实现仍需确定 tie-breaker、空输入和区间端点是闭/开。`[start,end)` 中 end==next.start 可兼容。

## 6. 回溯

回溯遍历选择树：选择、递归、撤销。用于组合、排列、约束满足。最坏通常指数/阶乘，必须有输入上限和剪枝。

```text
search(state):
    if state complete: record
    for choice in candidates(state):
        if choice violates constraint: continue
        apply(choice)
        search(state)
        undo(choice)
```

不变量是每次递归 state 精确表示当前路径；undo 必须恢复进入前状态。忘记复制结果会让所有答案引用同一可变切片。

剪枝只能删除不可能产生合法/更优答案的分支。启发式顺序可更早找到解但不改变最坏复杂度。需要一个解时找到后传播停止信号。

## 7. 动态规划

DP 适用于重叠子问题和最优子结构。步骤：定义状态含义；写转移；基例；计算顺序；答案位置；复杂度=状态数×每状态转移成本。

自顶向下记忆化贴近递归，但有调用栈和缓存键成本；自底向上明确顺序，易做空间压缩。空间压缩前确认转移只依赖哪些历史层，且不会覆盖仍需值。

DP 值若表示“不可达”，选择哨兵要防加法溢出。用 `ok`、足够大但可检查的 INF 或分离布尔可达状态。

恢复具体方案需保存 choice/previous，或从 DP 表反推。只计算最优值不自动得到路径。

## 8. 模式对比

| 模式 | 结构条件 | 常见失败 |
| --- | --- | --- |
| 双指针 | 移动可排除搜索区 | 未排序却用相向规则 |
| 滑动窗口 | 连续区间、边界单调 | 负数破坏和单调 |
| 分治 | 子问题缩小、可合并 | 重叠导致重复指数计算 |
| 贪心 | 可证明局部选择性质 | 只有直觉无证明 |
| 回溯 | 需要枚举选择树 | 无上限/剪枝导致爆炸 |
| DP | 重叠+最优子结构 | 状态缺变量导致错误合并 |

## 9. 完整案例：最小成本爬楼梯并恢复路径

### 9.1 契约

有 n 级台阶，位置 0 为起点，目标 n。每次走 1 或 2 级，落到台阶 i（1..n）支付 `cost[i-1]`。求最小总成本；成本非负 int64。相同成本优先更少步，再优先走 2 级。

输入 costs=`[4,2,7,1]`。可行路径：0→2→4，成本 2+1=3；0→1→2→4 成本 4+2+1=7。输出成本 3、路径 `[0,2,4]`。

### 9.2 状态与转移

`dp[i]` 是到达位置 i 的最小成本，`steps[i]` 是在最小成本下最少步数，`prev[i]` 保存前驱。候选来自 i-1 和 i-2：

```text
candidateCost = dp[j] + cost[i-1]
candidateSteps = steps[j] + 1
按 cost、steps、偏好比较候选
```

基例 dp[0]=0,steps[0]=0。i 从 1 到 n 递增，依赖均已计算。状态 O(n)，每状态最多 2 转移，总 O(n)，恢复路径 O(n)。

### 9.3 Go 实现

```go
package staircost

import (
    "errors"
    "math"
    "slices"
)

type answer struct { cost int64; steps int; previous int }

func Solve(costs []int64) (int64, []int, error) {
    n := len(costs)
    dp := make([]answer, n+1)
    dp[0] = answer{cost: 0, steps: 0, previous: -1}
    for i := 1; i <= n; i++ {
        if costs[i-1] < 0 { return 0, nil, errors.New("cost must be non-negative") }
        best := answer{cost: math.MaxInt64, steps: math.MaxInt, previous: -1}
        for jump := 2; jump >= 1; jump-- {
            from := i - jump
            if from < 0 { continue }
            if dp[from].cost > math.MaxInt64-costs[i-1] {
                continue
            }
            candidate := answer{
                cost: dp[from].cost + costs[i-1],
                steps: dp[from].steps + 1,
                previous: from,
            }
            if candidate.cost < best.cost ||
                (candidate.cost == best.cost && candidate.steps < best.steps) {
                best = candidate
            }
        }
        if best.previous == -1 { return 0, nil, errors.New("all path costs overflow") }
        dp[i] = best
    }
    path := make([]int, 0, n+1)
    for at := n; at >= 0; at = dp[at].previous {
        path = append(path, at)
        if at == 0 { break }
    }
    slices.Reverse(path)
    return dp[n].cost, path, nil
}
```

循环 jump 从 2 到 1，完全相同 cost/steps 时保留先看到的两级跳，实现 tie-break。比较逻辑若未来增加条件，最好抽成命名函数并测试。

### 9.4 验证

```go
func TestSolve(t *testing.T) {
    cost, path, err := Solve([]int64{4,2,7,1})
    if err != nil { t.Fatal(err) }
    if cost != 3 || !slices.Equal(path, []int{0,2,4}) {
        t.Fatalf("cost=%d path=%v", cost, path)
    }
}
```

验证路径不能只匹配预期：还应断言首 0、尾 n、每步差 1/2，并重新按落点求和等于 cost。对随机 n<=20，可用回溯枚举所有路径作为 oracle，对比 DP 最优值。

### 9.5 失败与边界

空 costs 表示起点即目标，返回成本 0、路径 `[0]`。单项 `[5]` 返回 `[0,1]` 成本 5。负成本当前拒绝；若允许负数，DAG 上 DP仍可计算，但问题契约与溢出处理要改。

输入 `[MaxInt64,1]` 可走 0→2 只付 1，不应因未选择候选的潜在溢出过早失败。实现遍历 i=1 时 dp[1]=MaxInt64，i=2 先从 0 得 1，随后发现从 dp1 加 1 溢出并跳过该候选，因此仍返回最优路径 `[0,2]`。只有一个位置的所有候选都溢出时才返回 `all path costs overflow`。这个分支说明 DP 必须逐候选定义异常语义，不能让无关坏候选阻断有效答案。

## 10. 贪心案例：区间调度

输入半开区间，先验证 start<=end，按 end 升序、start 升序、ID 排序。遍历若 start>=lastEnd 就选择。空区间是否占资源需要契约；若 `[x,x)` 不占时长，可允许多个同点。

交换论证：最早结束区间替换任一最优解首区间，不会缩小之后可用时间，因此存在包含该选择的最优解。重复应用得全局最优数量。

若目标改为最大总价值，最早结束贪心不再保证，需要加权区间调度 DP。

## 11. 回溯案例：唯一组合

候选有重复值时先排序，同一递归层跳过相同值，避免重复组合；但不能跨层全部跳过，否则会漏掉使用多个相同元素。剩余和与候选非负时可在超过目标剪枝，允许负数则该剪枝不正确。

输出规模本身可指数级，算法至少要花 O(output size)。API 设置最大解数量或提供迭代器/callback，并支持 context 取消。

## 12. 调试清单

- 双指针漏解：写出移动能排除哪些候选的证明。
- 窗口死循环：每轮 left/right 是否推进，计数删除是否归零。
- 分治栈溢出：基例和规模是否严格减小，是否可迭代。
- 贪心反例：穷举小输入与最优 oracle 对比，寻找失败。
- 回溯答案相同：记录结果时是否 clone 当前路径。
- DP 值正确路径错：prev 与 tie-break 是否同步更新。
- DP 溢出：候选加法前检查，不让无关候选阻断合法解。

## 13. 练习

1. 为 Solve 加入 `[MaxInt64,1]` 与所有候选溢出的回归测试。
2. 用穷举 oracle 对 n<=20 的随机 costs 验证 DP。
3. 实现有序两数和，返回原索引并处理重复值。
4. 给 LongestAtMostK 增加属性测试，与 O(n²) 暴力比较。
5. 实现加权区间调度 DP，明确 predecessor 的二分查找。

## 来源

- [MIT OpenCourseWare 6.006：Introduction to Algorithms](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/)（访问日期：2026-07-17）
- [MIT OpenCourseWare 6.046J：Design and Analysis of Algorithms](https://ocw.mit.edu/courses/6-046j-design-and-analysis-of-algorithms-spring-2015/)（访问日期：2026-07-17）
- [Go 标准库：slices](https://pkg.go.dev/slices)（访问日期：2026-07-17）
- [Go 语言规范：For statements](https://go.dev/ref/spec#For_statements)（访问日期：2026-07-17）
