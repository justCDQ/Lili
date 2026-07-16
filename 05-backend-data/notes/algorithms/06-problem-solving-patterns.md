# 双指针、滑动窗口、分治、贪心、回溯与动态规划

## 是什么

- 双指针：两个索引协调扫描，常把 O(n²) 降至 O(n)，要求排序或可维护关系。
- 滑动窗口：维护连续区间状态，若左右边界单调前进通常 O(n)。
- 分治：拆成独立子问题再合并，典型递归式决定复杂度。
- 贪心：每步局部选择，只有能证明交换性质/最优子结构才保证全局最优。
- 回溯：枚举选择树并剪枝，最坏通常指数级。
- 动态规划：重叠子问题与最优子结构，保存状态避免重复计算。

```js
function longestAtMostKDistinct(xs, k) {
  let left = 0, best = 0; const count = new Map();
  for (let right = 0; right < xs.length; right++) {
    count.set(xs[right], (count.get(xs[right]) ?? 0) + 1);
    while (count.size > k) {
      const x = xs[left++], n = count.get(x) - 1;
      n ? count.set(x, n) : count.delete(x);
    }
    best = Math.max(best, right - left + 1);
  }
  return best;
}
```

该窗口时间 O(n)、空间 O(k)。

## 常见错误与边界

先定义状态、不变量和转移，再写代码。滑动窗口只适合连续范围且能随边界更新；负数会破坏某些“和单调”窗口。DP 需估算状态数×转移成本，并可用滚动数组降空间。回溯必须设置输入上限和剪枝。贪心若无证明只能是启发式。

## 为什么需要

数据结构和算法决定操作成本随数据规模增长的方式。明确复杂度、数据分布和操作频率，才能在延迟、内存、实现复杂度与稳定性之间选择，而不是根据题型名称套用结构。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 实际怎么使用

先写输入规模 n、关键操作频率、内存上限和正确性不变量，再运行本文实现。为 n=0、1、重复值、最坏分布和大规模随机数据编写测试；用 JavaScript `performance.now()` 或 Go Benchmark 测量，并核对结果是否符合给出的时间与空间复杂度。

## 补充知识

Big-O 忽略常数、缓存局部性、分配和 I/O，不能直接预测实际延迟。生产选择还需比较标准库实现、最坏输入、稳定性、并发访问和数据是否超过内存。

## 来源

- [MIT OpenCourseWare：Introduction to Algorithms](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/)（访问日期：2026-07-16）
- [Cormen et al. Introduction to Algorithms](https://mitpress.mit.edu/9780262046305/introduction-to-algorithms/)（访问日期：2026-07-16）
- [MDN：Map](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map)（访问日期：2026-07-16）
