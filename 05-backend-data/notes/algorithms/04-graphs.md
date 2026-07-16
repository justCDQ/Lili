# 图、DFS、BFS、拓扑排序、最短路径与并查集

## 是什么

图由顶点和边组成，可有向/无向、加权/非加权。邻接表空间 O(V+E)。DFS/BFS 时间 O(V+E)：DFS 适合连通、环、回溯；BFS 求非加权最短边数。Kahn 拓扑排序 O(V+E)，仅适用于有向无环图。Dijkstra 用堆时 O((V+E)log V)，要求非负权。并查集配合路径压缩与按秩合并，m 次操作接近 O(m α(n))。

```js
function bfs(graph, start) {
  const q = [start], seen = new Set([start]), dist = new Map([[start, 0]]);
  for (let i = 0; i < q.length; i++) {
    for (const v of graph.get(q[i]) ?? []) if (!seen.has(v)) {
      seen.add(v); dist.set(v, dist.get(q[i]) + 1); q.push(v);
    }
  }
  return dist;
}
```

## 关键特性或规则

任务依赖调度用有向图和拓扑排序；层级最短路径用 BFS；非负成本路由用 Dijkstra；动态合并连通分量用并查集。先明确边方向、权重语义、重复边和孤立点。

## 常见错误与边界

拓扑结果可能不唯一；结果数量少于 V 表示有环。Dijkstra 不能处理负权，存在负权可选 Bellman-Ford O(VE)。递归 DFS 可能栈溢出，大图用显式栈。节点 ID 若稀疏，不要直接分配巨大数组。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 实际怎么使用

任务依赖调度器先统计每个节点的入度，再把入度为 0 的节点放入队列。每取出一个节点，就降低其后继入度；降到 0 时入队。最终输出数量少于节点数表示存在环。

```js
function topo(nodes, edges) {
  const next = new Map(nodes.map(x => [x, []]));
  const indegree = new Map(nodes.map(x => [x, 0]));
  for (const [from, to] of edges) {
    if (!next.has(from) || !next.has(to)) throw new Error("unknown node");
    next.get(from).push(to);
    indegree.set(to, indegree.get(to) + 1);
  }
  const queue = nodes.filter(x => indegree.get(x) === 0);
  const order = [];
  for (let i = 0; i < queue.length; i++) {
    const node = queue[i];
    order.push(node);
    for (const to of next.get(node)) {
      indegree.set(to, indegree.get(to) - 1);
      if (indegree.get(to) === 0) queue.push(to);
    }
  }
  if (order.length !== nodes.length) throw new Error("dependency cycle");
  return order;
}
```

测试空图、孤立节点、重复边、多个合法顺序、环和未知节点。BFS 验证起点距离为 0、每条树边距离递增 1；Dijkstra 验证所有权重非负；并查集验证合并的自反、对称和传递结果。所有邻接表算法先按 O(V+E) 预算空间。

## 补充知识

稠密图中 E 接近 V²，邻接矩阵的 O(V²) 空间可能换取 O(1) 边查询；稀疏图通常选邻接表。A* 在有可采纳启发函数时减少最短路径搜索范围，但最坏复杂度仍可能很高。并查集适合只增加连接的离线/增量连通性，不支持低成本删除边。

## 来源

- [Open Data Structures：Graphs](https://opendatastructures.org/ods-java/12_Graphs.html)（访问日期：2026-07-16）
- [Go：container/heap](https://pkg.go.dev/container/heap)（访问日期：2026-07-16）
- [Cormen et al. Algorithms resources](https://mitpress.mit.edu/9780262046305/introduction-to-algorithms/)（访问日期：2026-07-16）
