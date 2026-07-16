# Event Loop、Call Stack、Task 与 Microtask

## 是什么

调用栈执行同步帧；事件循环在栈空时取 task；每个 task 后浏览器清空 microtask 队列，再获得渲染机会。Promise reaction 和 queueMicrotask 进入 microtask。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

单线程 JS 仍可通过宿主并发 I/O；微任务会在下一 task 前清空；渲染时机由浏览器决定。

## 实际使用

```js
console.log('A');
setTimeout(()=>console.log('task'));
queueMicrotask(()=>console.log('microtask'));
console.log('B');
// A B microtask task
```

## 常见错误与边界

递归排入 microtask 会饿死任务和渲染；setTimeout(0) 不是立即执行；async/await 后续通常是 Promise job。

## 相关补充知识

每个 Task 结束后会清空 Microtask 队列，连续创建微任务可能推迟渲染和其他任务。事件循环顺序应通过 Performance 时间线和最小示例验证，不能只凭 `setTimeout(0)` 推断立即执行。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Event_loop)
- [WHATWG HTML Standard](https://html.spec.whatwg.org/multipage/webappapis.html#event-loops)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/HTML_DOM_API/Microtask_guide)

访问日期：2026-07-16。
