# 条件、循环、函数、参数、返回值与递归

## 是什么

条件选择执行路径，循环重复处理，函数封装可调用逻辑；参数是输入绑定，return 结束调用并给出结果；递归是函数直接或间接调用自身。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

优先提前返回减少嵌套；for...of 遍历可迭代值，for...in 遍历可枚举键；递归必须有终止条件。

## 实际使用

```js
function factorial(n){
 if(!Number.isInteger(n)||n<0) throw new RangeError();
 if(n<=1) return 1;
 return n*factorial(n-1);
}
for(const n of [0,1,5]) console.log(factorial(n));
```

## 常见错误与边界

深递归可能栈溢出；遗漏 return 得 undefined；break/continue 只影响对应循环。

## 相关补充知识

递归深度受调用栈限制，处理不受控深度的树或图时可改用显式栈。函数应明确返回路径和副作用，循环则要验证终止条件、空输入与边界索引。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Control_flow_and_error_handling)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Loops_and_iteration)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Functions)

访问日期：2026-07-16。
