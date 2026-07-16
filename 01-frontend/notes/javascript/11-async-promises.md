# 同步、异步、Callback、Promise 与 async/await

## 是什么

同步代码当前栈完成后继续；异步操作稍后通知。callback 是传入的后续函数；Promise 表示最终完成/失败；async 函数总返回 Promise，await 暂停该函数而不阻塞线程。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

并行任务先创建 Promise 再 await；链必须 return；异常由 reject/throw 传播；清理放 finally。

## 实际使用

```js
async function load(){
 try { const [a,b]=await Promise.all([fetch('/a'),fetch('/b')]); return [await a.json(),await b.json()]; }
 catch(error){throw new Error('load failed',{cause:error});}
}
```

## 常见错误与边界

在循环中顺序 await 可能无意串行；未处理 rejection 难定位；await 不会自动取消底层工作。

## 相关补充知识

Promise 表示一次最终结果，不等于任务线程，也不能自行取消底层操作。并发独立任务可用组合方法，串行依赖才逐个 await；所有拒绝路径都应被观察和处理。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Extensions/Async_JS/Promises)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Using_promises)

访问日期：2026-07-16。
