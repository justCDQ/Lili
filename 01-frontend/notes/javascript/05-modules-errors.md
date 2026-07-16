# 模块、导入导出、错误与异常处理

## 是什么

ES module 以 import/export 建立静态依赖；异常沿调用栈传播直到 catch；Error 保存名称、消息和栈信息。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

浏览器模块脚本用 type=module；相对说明符通常含扩展名；只捕获能处理的错误并保留 cause；finally 用于清理。

## 实际使用

```js
// math.js
export function divide(a,b){if(b===0) throw new RangeError('zero'); return a/b;}
// app.js
import {divide} from './math.js';
try{console.log(divide(4,0));}catch(error){console.error(error); }
```

## 常见错误与边界

catch 后静默吞错会隐藏故障；动态 import 返回 Promise；模块绑定是 live binding 不是普通复制。

## 相关补充知识

ES 模块导入是实时绑定并按 URL 标识模块实例，循环依赖可能读取到尚未初始化的绑定。错误包装应使用 `cause` 保留原始错误，捕获后不能无记录吞掉失败。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Control_flow_and_error_handling)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error)

访问日期：2026-07-16。
