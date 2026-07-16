# GC、WeakMap、WeakSet 与常见内存泄漏

## 是什么

垃圾回收器回收从根不可达的对象。WeakMap 的键和 WeakSet 的值是弱持有，不阻止对象被回收；二者不可枚举。内存泄漏是程序仍持有已不再需要对象的引用，导致其保持可达。

## 为什么需要

长寿命单页应用中的监听器、定时器、闭包、缓存和 detached DOM 会持续累积，引起内存增长、GC 停顿和崩溃。

## 关键特性与规则

组件销毁时移除监听器、取消 timer/observer/request、终止 Worker、撤销 Blob URL。缓存设置容量或过期策略。WeakMap 适合给对象附加不影响生命周期的元数据，但不能替代所有缓存。

## 实际使用

```js
const metadata = new WeakMap();
metadata.set(element, { measuredAt: performance.now() });

const controller = new AbortController();
window.addEventListener('resize', onResize, { signal: controller.signal });
// 销毁时
controller.abort();
clearInterval(timerId);
URL.revokeObjectURL(previewUrl);
```

## 常见错误与边界

不能依赖 GC 发生时间或 finalizer 执行业务逻辑。WeakMap 只对键弱引用，value 若引用其他大对象仍会随键存活。一次 Heap Snapshot 大不能直接证明泄漏，需重复操作后比较保留路径。

## 相关补充知识

Chrome Memory 面板可用 Heap Snapshot、Allocation instrumentation 和 detached elements；先建立稳定复现，再查看 dominator 与 retainer path。

## 来源

- [MDN：Memory management](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Memory_management)
- [MDN：WeakMap](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/WeakMap)
- [Chrome：Fix memory problems](https://developer.chrome.com/docs/devtools/memory-problems/)

访问日期：2026-07-16。
