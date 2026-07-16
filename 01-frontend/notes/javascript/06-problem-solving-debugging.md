# 伪代码、断点与日志调试

## 是什么

伪代码先明确输入、输出、步骤和边界；断点暂停运行检查作用域和调用栈；结构化日志记录关键状态与事件。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

先构造最小复现和预期/实际差异；用条件断点、watch、step over/into；日志包含上下文但不含秘密。

## 实际使用

```js
// 输入: numbers；输出: 最大值或 null
function maxOrNull(numbers){
 debugger;
 if(numbers.length===0) return null;
 return Math.max(...numbers);
}
```

## 常见错误与边界

大量 console.log 改变时序且难关联；只修症状不验证根因；生产日志可能泄漏个人数据。

## 相关补充知识

调试先固定输入、预期和最小复现，再观察控制流与状态，不用连续添加日志猜测。条件断点、调用栈、作用域和 Network 请求可分别定位计算、时序与接口问题。

## 来源

- [Chrome Developers](https://developer.chrome.com/docs/devtools/javascript/)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Scripting/Debugging_JavaScript)
- [Visual Studio Code Documentation](https://code.visualstudio.com/docs/editor/debugging)

访问日期：2026-07-16。
