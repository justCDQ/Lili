# 表单状态、校验、URL 与 Web Storage

## 是什么

表单状态包括当前值、有效性和提交状态；URL API 解析/生成地址；localStorage 跨会话，sessionStorage 限标签页会话，值均为字符串。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

URL 适合可分享筛选状态；存储前版本化和校验解析结果；storage 同源隔离；敏感状态放安全服务端会话。

## 实际使用

```js
const data=new FormData(form);
if(!form.reportValidity()) return;
const url=new URL(location.href); url.searchParams.set('q',data.get('q'));
history.replaceState(null,'',url);
localStorage.setItem('prefs',JSON.stringify(prefs));
```

## 常见错误与边界

Web Storage 同步阻塞且容量有限；用户可修改/清除；XSS 可读取 localStorage，勿存认证秘密。

## 相关补充知识

URL 适合可分享和可返回的非敏感状态；Web Storage 同步且可被同源脚本读取，不适合 Secret。表单提交应以 FormData 和原生校验为基础，并在服务端重新验证。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/FormData)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/URL)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)

访问日期：2026-07-16。
