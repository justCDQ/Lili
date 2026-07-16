# Iterator、Generator、Proxy 与 Reflect

## 是什么

iterator 用 next() 产出序列；iterable 通过 Symbol.iterator 提供迭代器；generator 以 function* 和 yield 简化状态机；Proxy 拦截对象操作，Reflect 提供对应默认操作。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

迭代器结果含 value/done；generator 可暂停恢复；Proxy trap 必须遵守不变量；Reflect 便于正确转发。

## 实际使用

```js
function* range(n){for(let i=0;i<n;i++) yield i;}
const model=new Proxy(target,{set(obj,key,value){if(key==='age'&&value<0)return false; return Reflect.set(obj,key,value);}});
```

## 常见错误与边界

Proxy 有身份和性能成本且不能代理内部槽行为；无限迭代器配展开会不终止；generator 不是并行线程。

## 相关补充知识

Iterator 协议把遍历与集合实现分离，Generator 可暂停并生成序列。Proxy 必须遵守语言不变量且会增加调试和性能成本；Reflect 提供与内部对象操作对应的函数接口。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Iterators_and_generators)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Reflect)

访问日期：2026-07-16。
