# 图、DFS、BFS、拓扑排序、最短路径与并查集

## 学习目标

本文从图的数据模型开始，逐项说明遍历、依赖排序、路径成本和动态连通性算法的前提、不变量与复杂度，并实现可检测未知节点、重复边和环的稳定拓扑排序。

## 1. 图模型

图 G=(V,E) 由顶点集合和边集合组成。有向边 from→to 有方向，无向边表示双向关系；加权边附成本。先定义自环、重复边、孤立点是否允许。

业务映射示例：包依赖是有向图，社交好友常建无向图，道路可能有向加权，任务前置关系应是 DAG。边权必须说明是距离、时间、金额还是概率，算法对权重条件不同。

## 2. 表示方式

邻接表为每节点保存出边，空间 O(V+E)，遍历所有邻居高效，适合稀疏图。无向边通常存两次，E 的计数要说明按逻辑边还是邻接项。

邻接矩阵使用 V×V 单元，空间 O(V²)，边存在查询 O(1)，适合顶点较少的稠密图。节点 ID 稀疏时先映射到连续索引，不能按最大 ID 分配巨大矩阵。

边列表空间 O(E)，适合批处理边、Kruskal 等，但查某节点邻居需扫描或额外索引。

Go 可用 `map[string][]string`，但 map 迭代和输入边顺序会影响遍历输出。若协议要求稳定，邻接列表和初始节点都排序。

## 3. DFS

深度优先搜索沿一条路径深入再回退。用 seen 防止环导致无限循环。邻接表下每顶点和边最多处理常数次，时间 O(V+E)，额外空间 O(V)。

递归 DFS 简洁但深链可能栈耗尽。显式栈保存待访问节点，可设置最大顶点/边并控制顺序。若希望与递归顺序一致，压栈邻居要逆序。

DFS 用于连通分量、环检测、拓扑排序的一种实现、路径存在和回溯。无向图环检测要忽略指向父节点的反向边；有向图用白/灰/黑颜色，遇到指向灰色栈内节点的边表示环。

## 4. BFS

广度优先搜索从起点逐层访问，用 FIFO 队列。非加权图中，首次发现节点的层数是最少边数距离，因为所有更短层已经处理。

```go
func Distances(graph map[string][]string, start string) map[string]int {
    distance := map[string]int{start: 0}
    queue := []string{start}
    for head := 0; head < len(queue); head++ {
        current := queue[head]
        for _, next := range graph[current] {
            if _, seen := distance[next]; seen { continue }
            distance[next] = distance[current] + 1
            queue = append(queue, next)
        }
    }
    return distance
}
```

要恢复路径，另存 `previous[next]=current`，从目标回溯到起点再反转。若目标不可达，不存在 previous；起点到自身距离 0。

BFS 不能直接求加权最短路径。权重全为 0/1 可用双端队列 0-1 BFS；一般非负权用 Dijkstra。

## 5. 拓扑排序

拓扑序是有向图顶点线性排列，使每条边 u→v 中 u 在 v 前。只有 DAG 存在，结果可能不唯一。

Kahn 算法计算每个节点入度，把入度 0 节点放队列；取出节点后删除其出边，后继入度降到 0 入队。若最终输出少于 V，剩余部分含环。

复杂度 O(V+E)，前提是建表和每条边只处理一次。重复边若保留，入度和删除次数必须同时保留；若业务把重复依赖视为输入错误，应在建图时拒绝。

需要稳定输出时用按键排序的优先队列选入度 0 节点，复杂度变为 O((V+E)log V) 的上界，不再是纯 O(V+E)。

## 6. Dijkstra

Dijkstra 求非负权图的单源最短路径。维护暂定距离和最小堆，每次取距离最小节点并松弛出边：若 `dist[u]+w < dist[v]`，更新 v。

非负权是正确性前提：已取出的最小距离不会被之后路径降低。负权会破坏该性质；有负权可用 Bellman-Ford O(VE)，有负环时最短路径未定义。

堆中常允许同一节点多个旧距离，Pop 时若与当前 dist 不同就跳过，避免实现 decrease-key。用二叉堆复杂度 O((V+E)log V)。

整数距离要防溢出，不用 `MaxInt + weight`。浮点权要处理 NaN、负零和比较误差；很多协议使用固定整数单位。

## 7. 并查集

并查集维护互不相交集合，支持 Find 找代表元、Union 合并。路径压缩让查找路径变短，按秩/大小合并把小树挂到大树。

m 次操作的摊还复杂度 O(m α(n))，α 是反 Ackermann 函数，在实际规模很小；空间 O(n)。

并查集适合无向连通分量、Kruskal、逐步增加边的连通查询。它不保存具体路径，也不支持低成本删除边；需要删除的动态连通性用更复杂结构或离线算法。

不变量：每节点 parent 最终形成根为自身的树；size/rank 只对根有效；Union 已同集合不重复增加组件数。

## 8. 完整案例：稳定任务依赖排序

### 8.1 契约

输入所有任务 nodes 和依赖 edges `[before,after]`。节点 ID 非空且唯一，边端点必须存在，重复边视为错误。输出所有合法序中按字典序最小的一个；有环返回错误。

### 8.2 实现

```go
package topo

import (
    "container/heap"
    "errors"
    "fmt"
)

type stringHeap []string
func (h stringHeap) Len() int { return len(h) }
func (h stringHeap) Less(i,j int) bool { return h[i] < h[j] }
func (h stringHeap) Swap(i,j int) { h[i],h[j] = h[j],h[i] }
func (h *stringHeap) Push(x any) { *h = append(*h, x.(string)) }
func (h *stringHeap) Pop() any {
    old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x
}

type Edge struct { Before, After string }

func Sort(nodes []string, edges []Edge) ([]string, error) {
    next := make(map[string][]string, len(nodes))
    indegree := make(map[string]int, len(nodes))
    for _, node := range nodes {
        if node == "" { return nil, errors.New("node id is empty") }
        if _, exists := next[node]; exists { return nil, fmt.Errorf("duplicate node %q", node) }
        next[node] = nil
        indegree[node] = 0
    }
    seenEdge := make(map[Edge]struct{}, len(edges))
    for _, edge := range edges {
        if _, ok := next[edge.Before]; !ok { return nil, fmt.Errorf("unknown node %q", edge.Before) }
        if _, ok := next[edge.After]; !ok { return nil, fmt.Errorf("unknown node %q", edge.After) }
        if _, duplicate := seenEdge[edge]; duplicate { return nil, fmt.Errorf("duplicate edge %+v", edge) }
        seenEdge[edge] = struct{}{}
        next[edge.Before] = append(next[edge.Before], edge.After)
        indegree[edge.After]++
    }

    ready := &stringHeap{}
    heap.Init(ready)
    for node, degree := range indegree {
        if degree == 0 { heap.Push(ready, node) }
    }
    order := make([]string, 0, len(nodes))
    for ready.Len() > 0 {
        node := heap.Pop(ready).(string)
        order = append(order, node)
        for _, successor := range next[node] {
            indegree[successor]--
            if indegree[successor] == 0 { heap.Push(ready, successor) }
        }
    }
    if len(order) != len(nodes) { return nil, errors.New("dependency cycle") }
    return order, nil
}
```

### 8.3 输入、步骤与输出

nodes=`[build,deploy,lint,test]`，edges=`lint→build,test→build,build→deploy`。初始入度 0 为 lint/test，最小堆先取 lint 再 test；build 入度归零后取 build；最后 deploy。输出 `[lint,test,build,deploy]`。

该顺序满足所有边且在每次可选集合取字典序最小。若不要求稳定最小序，可用普通队列达到 O(V+E)。当前堆版 O((V+E)log V)，空间 O(V+E)。

### 8.4 验证函数

测试不应只匹配一个顺序；通用验证建立 position map，对每条边断言 `pos[before] < pos[after]`，并确认每节点恰好一次。案例额外断言字典序最小结果。

```go
func valid(order []string, edges []Edge) bool {
    position := make(map[string]int, len(order))
    for i, node := range order { position[node] = i }
    for _, edge := range edges {
        if position[edge.Before] >= position[edge.After] { return false }
    }
    return true
}
```

### 8.5 失败分支

edges 加 `deploy→lint` 形成环，输出数量不足并返回 dependency cycle。自环 `build→build` 也失败。未知节点、重复 node/edge 在建图时返回定位错误。空图返回空切片而非 nil 是否重要由 API 契约决定。

Kahn 仅告诉存在环，不直接给环路径。需要诊断可对剩余节点运行带颜色 DFS，记录父指针恢复一个环。

仓库中的[可运行 Topological Sort 示例](../../examples/algorithms/topo/)保存了稳定顺序与环测试。

## 9. BFS 案例验证

对图 A→B,A→C,B→D,C→D，BFS(A) 得 A=0,B=1,C=1,D=2。邻接顺序可能改变 D 的 previous 是 B 或 C，但距离相同。测试最短距离，不强制某条等长路径，除非定义 tie-breaker。

起点不存在时 Distances 当前仍返回 start=0，这可能隐藏输入错误。生产 API 应先验证节点集合，返回 unknown start。

## 10. Dijkstra 与业务单位

道路时间权重必须非负并用统一单位。若边权来自实时预测，算法运行期间变化会使结果对应混合快照；先获取一致输入版本或接受近似并标记版本。

不可达节点距离保持无穷/缺失。路径恢复前防止 predecessor 环。对同距路径需要稳定选择时，在堆比较加入节点 ID，并在相等松弛时明确策略。

## 11. 调试清单

- BFS 距离太大：节点是否入队前就标 seen，避免重复覆盖。
- DFS 无限循环：是否在扩展前标记，图是否包含自环。
- topo 少节点：检查环，也检查孤立点是否未加入 nodes。
- topo 入度变负：重复边处理是否一致。
- Dijkstra 对负权错误：输入边验证必须在算法开始完成。
- 距离溢出：加法前检查，选择明确宽度与单位。
- 并查集组件数错误：重复 Union 是否误减，路径压缩是否破坏根。

## 12. 练习

1. 为 Sort 写随机 DAG 测试，用 position 验证所有边。
2. 增加环路径诊断，返回一个具体闭环。
3. 实现 BFS previous 路径，覆盖不可达与等长路径。
4. 实现非负整数 Dijkstra，与 Bellman-Ford 在随机小图上对照。
5. 实现并查集，处理 0、1、重复 Union，并用于无向图连通分量。

## 来源

- [Open Data Structures：Graphs](https://opendatastructures.org/ods-java/12_Graphs.html)（访问日期：2026-07-17）
- [Go 标准库：container/heap](https://pkg.go.dev/container/heap)（访问日期：2026-07-17）
- [MIT OpenCourseWare 6.006：Graph Search, BFS, DFS](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/)（访问日期：2026-07-17）
- [MIT OpenCourseWare 6.046J：Design and Analysis of Algorithms](https://ocw.mit.edu/courses/6-046j-design-and-analysis-of-algorithms-spring-2015/)（访问日期：2026-07-17）
