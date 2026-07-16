# Streams、Structured Clone 与 Transferable Object

## 是什么

Streams API 增量读取、转换和写入数据，并用背压协调生产与消费。structured clone 算法深复制支持的结构。Transferable 对象可转移底层资源所有权，避免复制大缓冲区。

## 为什么需要

大响应、文件和 Worker 消息若一次完整加载或复制，会增加延迟、峰值内存与主线程工作。

## 关键特性与规则

ReadableStream 通过 reader 或异步迭代消费；TransformStream 连接处理阶段；锁定 reader 后不能被其他消费者读取。`structuredClone(value, {transfer})` 转移后原 `ArrayBuffer` 被分离。

## 实际使用

```js
const response = await fetch('/events');
const textStream = response.body.pipeThrough(new TextDecoderStream());
for await (const chunk of textStream) parseChunk(chunk);

const buffer = new ArrayBuffer(1024);
worker.postMessage({ buffer }, [buffer]);
console.log(buffer.byteLength); // 0，所有权已转移
```

## 常见错误与边界

分块边界不等于文本行或 JSON 对象边界，解析器必须缓存残片。structuredClone 不能复制函数和 DOM 节点。转移是破坏性操作，之后不能继续使用原资源。

## 相关补充知识

`Response.clone()` 基于流 tee 行为，两个消费者速率差异可能产生缓冲；处理大型数据时应测量内存。

## 来源

- [MDN：Streams API](https://developer.mozilla.org/en-US/docs/Web/API/Streams_API)
- [MDN：Structured clone algorithm](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Structured_clone_algorithm)
- [MDN：Transferable objects](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Transferable_objects)

访问日期：2026-07-16。
