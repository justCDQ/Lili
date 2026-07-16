# JSON、Fetch、HTTP 方法、状态码与请求错误

## 是什么

JSON 是数据交换语法；fetch 返回 Response Promise；HTTP 方法表达动作语义，状态码表达协议结果；网络、HTTP、解析和业务错误需分层处理。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

GET/HEAD 应安全，PUT/DELETE 等幂等语义由服务端正确实现；检查 content-type；超时/取消用 AbortSignal；错误保留状态与响应信息。

## 实际使用

```js
const r=await fetch('/api/items',{headers:{Accept:'application/json'},signal});
if(!r.ok) throw new Error(`HTTP ${r.status}`);
const data=await r.json();
```

## 常见错误与边界

fetch 对 404/500 不 reject；JSON 不支持 undefined/BigInt/循环引用；重试非幂等请求可能重复写入。

## 相关补充知识

Fetch 只在网络层失败时拒绝，HTTP 4xx/5xx 仍需检查 `ok`。解析前要核对状态和内容类型；超时、取消、重试与重复写操作还需结合 AbortSignal 和幂等策略。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Methods)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)
- [www.rfc-editor.org](https://www.rfc-editor.org/rfc/rfc8259)

访问日期：2026-07-16。
