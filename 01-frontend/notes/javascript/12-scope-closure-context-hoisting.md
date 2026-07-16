# Scope、Closure、Execution Context 与 Hoisting

## 是什么

作用域控制标识符可见性；闭包是函数及其词法环境；执行上下文保存当前代码所需绑定；hoisting 描述声明在执行前被实例化的可观察效果。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

let/const 有暂时性死区；函数声明可在声明前调用；闭包捕获绑定而非冻结值；每次函数调用有独立执行上下文。

## 实际使用

```js
function makeCounter(){let n=0; return ()=>++n;}
const next=makeCounter();
console.log(next(),next());
```

## 常见错误与边界

把 var 当块级会泄漏；循环闭包配 var 共享绑定；闭包长期引用大对象会延长生命周期。

## 相关补充知识

闭包保留词法环境，可用于封装状态，也可能长期持有大型对象。`let/const` 绑定已创建但在初始化前处于暂时性死区；函数声明、`var` 和类的提升行为不同。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Glossary/Scope)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Closures)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Glossary/Hoisting)

访问日期：2026-07-16。
