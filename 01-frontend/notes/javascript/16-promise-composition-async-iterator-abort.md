# Promise 组合、Async Iterator 与 AbortController

## 是什么

`Promise.all`、`allSettled`、`any`、`race` 组合多个 Promise；异步迭代器让 `for await...of` 逐个等待值；`AbortController` 通过 `signal` 向支持取消的 API 广播终止请求。

## 为什么需要

真实请求需要控制并行、部分失败、增量结果与用户取消，单个 `await` 无法表达全部策略。

## 关键特性与规则

`all` 遇首个 rejection 快速失败；`allSettled` 等全部完成；`any` 取首个 fulfilled；`race` 取首个 settled。取消必须把同一 signal 传到每层支持者，并区分 AbortError 与真实失败。

## 实际使用

```js
const controller = new AbortController();
const jobs = urls.map(url => fetch(url, { signal: controller.signal }));
const results = await Promise.allSettled(jobs);
for await (const item of asyncSource) process(item);
controller.abort('navigation');
```

## 常见错误与边界

Promise 本身没有通用取消；abort 不会自动撤销已完成服务端写入。`race` 做超时若不取消原任务仍会占资源。无限异步序列需要退出与清理协议。

## 相关补充知识

并发数受限时使用 worker pool，不要一次创建无界 Promise；取消原因可通过 `signal.reason` 传播。

## 来源

- [MDN：Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise)
- [MDN：for await...of](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/for-await...of)
- [MDN：AbortController](https://developer.mozilla.org/en-US/docs/Web/API/AbortController)

访问日期：2026-07-16。
