# Web Streams、结构化克隆与所有权转移

Streams API 让程序逐块读取、转换和写入数据，并通过背压协调生产速度与消费速度。结构化克隆在不同 realm 或线程之间复制支持的数据图。Transferable 对象则移动底层资源的所有权，减少大块数据复制。

## 1. 为什么需要流

一次性读取完整响应会等待全部数据到达，并需要容纳完整结果。流允许首块到达后立即处理，降低首个结果延迟和峰值内存。

适合流式处理的场景包括：

- 大文件上传和下载；
- 持续事件或日志；
- 分块文本解析；
- 压缩、解压和编码转换；
- 从网络到文件或其他目标的管道。

流不会自动让算法更快。分块管理、解码和调度有成本；小数据直接使用 `response.json()` 或 `response.text()` 通常更简单。

## 2. 三类流

### 2.1 ReadableStream

可读流表示数据来源。它具有 readable、closed 或 errored 状态，内部队列保存已产生但未消费的 chunk。

```js
function streamFromArray(values) {
  let index = 0;

  return new ReadableStream({
    pull(controller) {
      if (index === values.length) {
        controller.close();
        return;
      }
      controller.enqueue(values[index]);
      index += 1;
    },
    cancel(reason) {
      console.log("消费方取消", reason);
    },
  });
}

const stream = streamFromArray(["HTML", "CSS", "JavaScript"]);
for await (const chunk of stream) {
  console.log(chunk);
}
```

底层来源常用回调：

- `start(controller)`：构造流时初始化，可同步或返回 Promise；
- `pull(controller)`：需要更多数据时调用；
- `cancel(reason)`：消费方不再需要数据时释放来源资源。

控制器的 `enqueue(chunk)` 入队、`close()` 正常结束、`error(reason)` 使流进入错误状态。

### 2.2 WritableStream

可写流表示数据目标。底层 sink 常用 `write(chunk)`、`close()`、`abort(reason)`。

```js
const received = [];

const sink = new WritableStream({
  async write(chunk) {
    await new Promise((resolve) => setTimeout(resolve, 5));
    received.push(chunk);
  },
  close() {
    console.log("写入完成", received);
  },
  abort(reason) {
    console.error("写入中止", reason);
  },
});

const writer = sink.getWriter();
await writer.write("one");
await writer.write("two");
await writer.close();
```

`write()` 返回的 Promise 表示该写入被处理的进度。直接写入时，应观察 `writer.ready` 和 `writer.desiredSize`，避免无界排队。

### 2.3 TransformStream

转换流同时暴露 writable 输入侧和 readable 输出侧。`transform(chunk, controller)` 可产出零个、一个或多个结果；`flush(controller)` 在输入关闭时处理尾部缓存。

```js
const upperCase = new TransformStream({
  transform(chunk, controller) {
    controller.enqueue(String(chunk).toUpperCase());
  },
});

await streamFromArray(["html", "css"])
  .pipeThrough(upperCase)
  .pipeTo(new WritableStream({
    write(chunk) {
      console.log(chunk);
    },
  }));
```

## 3. 读取、锁与取消

### 3.1 使用 reader

`stream.getReader()` 取得 reader 并锁定流。同一时间只能有一个活动 reader；锁定期间不能再次取得 reader，也不能直接 `pipeTo()`。

```js
async function collect(stream) {
  const reader = stream.getReader();
  const chunks = [];

  try {
    while (true) {
      const { value, done } = await reader.read();
      if (done) return chunks;
      chunks.push(value);
    }
  } finally {
    reader.releaseLock();
  }
}
```

`read()` 返回 Promise，完成值为 `{ value, done }`。流关闭后 `done` 为 true；流出错时 Promise rejected。

只有在没有待处理读请求时才可安全释放 reader 锁。释放锁不等于取消来源；不再消费时调用 `reader.cancel(reason)`。

### 3.2 异步迭代

现代可读流支持 `for await...of`。提前 `break` 默认会取消流，可用迭代选项控制是否保留，但需要按目标环境检查支持。

```js
async function findFirst(stream, predicate) {
  for await (const chunk of stream) {
    if (predicate(chunk)) return chunk;
  }
  return undefined;
}
```

### 3.3 管道

`readable.pipeTo(writable, options)` 返回表示整条传输完成的 Promise，并在管道期间锁定两端。`pipeThrough(transform)` 返回转换流的 readable 侧，便于链式组合。

```js
await response.body
  .pipeThrough(new TextDecoderStream())
  .pipeThrough(createLineStream())
  .pipeTo(new WritableStream({
    write(line) {
      console.log(line);
    },
  }));
```

`pipeTo` 默认传播关闭、错误和中止。选项 `preventClose`、`preventAbort`、`preventCancel` 可阻止对应传播，`signal` 可取消管道。改变默认传播必须有明确资源所有权理由。

## 4. 背压

生产者快于消费者时，未处理 chunk 会在内部队列增长。背压把下游容量不足的信号沿管道向上游传播。

队列策略包含：

- `highWaterMark`：队列期望容量阈值；
- `size(chunk)`：每个 chunk 对队列的计量值；
- `controller.desiredSize`：阈值减去当前队列总大小的近似容量信号。

```js
function createCounterStream(limit) {
  let value = 0;

  return new ReadableStream(
    {
      pull(controller) {
        if (value === limit) {
          controller.close();
          return;
        }
        controller.enqueue(value++);
        console.log("剩余容量", controller.desiredSize);
      },
    },
    new CountQueuingStrategy({ highWaterMark: 2 }),
  );
}
```

对于 pull source，平台在需要数据时调用 `pull()`。无法暂停的 push source 必须自行决定丢弃、合并、限制缓冲或让底层协议降速；Stream API 不能凭空为不支持背压的来源提供无损控制。

可写侧直接写入时：

```js
async function writeAll(stream, chunks) {
  const writer = stream.getWriter();
  try {
    for (const chunk of chunks) {
      await writer.ready;
      await writer.write(chunk);
    }
    await writer.close();
  } catch (error) {
    await writer.abort(error).catch(() => {});
    throw error;
  } finally {
    writer.releaseLock();
  }
}
```

## 5. 文本与分块边界

网络 chunk 边界不等于字符、文本行或 JSON 对象边界。UTF-8 多字节字符可能跨 chunk，必须使用状态化解码器，如 `TextDecoderStream` 或 `TextDecoder.decode(bytes, { stream: true })`。

下面的转换器缓存不完整行：

```js
function createLineStream() {
  let remainder = "";

  return new TransformStream({
    transform(chunk, controller) {
      const parts = (remainder + chunk).split(/\r?\n/);
      remainder = parts.pop() ?? "";
      for (const line of parts) controller.enqueue(line);
    },
    flush(controller) {
      if (remainder !== "") controller.enqueue(remainder);
    },
  });
}
```

不能对每个字节 chunk 单独 `new TextDecoder().decode(chunk)` 后拼接，这可能破坏跨块字符。逐行 JSON 还应对每行单独捕获解析错误并记录行号。

## 6. tee 与多消费者

`readable.tee()` 返回两个分支，来源 chunk 被送往两边。两个消费者速率不同会导致快分支继续推进，而慢分支累积数据；大流可能显著增加内存。

```js
const [forDisplay, forCache] = response.body.tee();

await Promise.all([
  consumeForDisplay(forDisplay),
  saveToCache(forCache),
]);
```

`Response.clone()` 的 body 也建立分支语义。不要把 clone 当成无成本深复制；需要测量慢分支和大响应的缓冲。

## 7. 字节流与 BYOB

可读字节流使用 `{ type: "bytes" }`。BYOB reader 允许消费方提供缓冲区，减少中间分配和复制，适用于二进制协议和高吞吐路径。

```js
async function readBytes(byteStream) {
  const reader = byteStream.getReader({ mode: "byob" });
  try {
    let buffer = new Uint8Array(4096);
    while (true) {
      const { value, done } = await reader.read(buffer);
      if (done) break;
      processBytes(value);
      buffer = new Uint8Array(4096);
    }
  } finally {
    reader.releaseLock();
  }
}
```

BYOB 涉及缓冲区分离和视图生命周期，不能假设每次读后原视图仍可复用。应按 API 返回的 `value` 视图处理，并在目标环境验证行为。

## 8. 结构化克隆

`structuredClone(value, options)` 使用结构化克隆算法复制数据图。它支持循环引用，并保留同一图中重复引用之间的关系。

```js
const original = {
  createdAt: new Date("2026-07-17T00:00:00Z"),
  tags: new Set(["stream", "worker"]),
};
original.self = original;

const copy = structuredClone(original);

console.log(copy !== original); // true
console.log(copy.self === copy); // true
console.log(copy.createdAt instanceof Date); // true
console.log(copy.tags instanceof Set); // true
```

结构化克隆不是 JSON 往返：它支持 `Map`、`Set`、`Date`、`RegExp`、ArrayBuffer、TypedArray 等多种类型和循环引用。具体 Web API 类型是否可序列化由相应规范定义。

### 8.1 不支持与丢失的语义

函数和 DOM 节点不能结构化克隆，通常抛出 `DataCloneError`。类实例的自定义原型、访问器和完整属性描述符不会作为普通对象行为完整复制；私有字段也不会以公共数据方式出现。

```js
try {
  structuredClone({ render() {} });
} catch (error) {
  console.log(error.name); // DataCloneError
}
```

跨线程消息应使用显式数据传输对象：只包含协议所需字段、版本、类型标签和可克隆数据。不要直接发送富业务实例并期待方法保留。

## 9. Transferable 与所有权移动

Transferable 对象的底层资源可以移动到新上下文。以 `ArrayBuffer` 为例，转移后发送方 buffer 被分离，`byteLength` 变为 0；接收方获得资源。

```js
const bytes = new Uint8Array([10, 20, 30]);

const moved = structuredClone(
  { payload: bytes },
  { transfer: [bytes.buffer] },
);

console.log(bytes.byteLength); // 0
console.log([...moved.payload]); // [10, 20, 30]
```

TypedArray 视图本身通常不是 transfer list 中的对象，应转移其 `buffer`。要接收资源，transfer list 中的资源还必须从所发送的数据图可达；只列在 transfer list 却不放进消息可能造成原资源分离但接收方无法访问。

Worker 传输：

```js
const buffer = new ArrayBuffer(1024 * 1024);
const view = new Uint8Array(buffer);
view[0] = 42;

worker.postMessage(
  { type: "PROCESS_BUFFER", buffer },
  [buffer],
);

console.log(buffer.byteLength); // 0
```

转移适合明确交接所有权的大型缓冲区、`MessagePort` 或支持转移的流。调用后继续使用原资源是错误。转移常可避免内存块复制，但规范并不把所有 transferable 的底层实现都保证为字面零复制。

## 10. 完整案例：流式 NDJSON 解析

案例从 `Response.body` 增量解码 UTF-8，按行切分，解析每条 JSON，并通过 signal 取消整条管道。

```js
function createLineSplitter() {
  let buffer = "";
  let lineNumber = 0;

  return new TransformStream({
    transform(chunk, controller) {
      buffer += chunk;
      const lines = buffer.split(/\r?\n/);
      buffer = lines.pop() ?? "";

      for (const line of lines) {
        lineNumber += 1;
        if (line.trim() !== "") controller.enqueue({ line, lineNumber });
      }
    },
    flush(controller) {
      if (buffer.trim() !== "") {
        lineNumber += 1;
        controller.enqueue({ line: buffer, lineNumber });
      }
    },
  });
}

function createJSONParser() {
  return new TransformStream({
    transform(entry, controller) {
      try {
        controller.enqueue({
          ok: true,
          lineNumber: entry.lineNumber,
          value: JSON.parse(entry.line),
        });
      } catch (error) {
        controller.enqueue({
          ok: false,
          lineNumber: entry.lineNumber,
          error,
          source: entry.line,
        });
      }
    },
  });
}

async function readNDJSON(url, { signal, onItem, onInvalid }) {
  const response = await fetch(url, { signal });
  if (!response.ok) throw new Error(`HTTP ${response.status}`);
  if (!response.body) throw new Error("响应没有可读 body");

  const sink = new WritableStream({
    async write(result) {
      if (result.ok) await onItem(result.value, result.lineNumber);
      else await onInvalid(result);
    },
  });

  await response.body
    .pipeThrough(new TextDecoderStream())
    .pipeThrough(createLineSplitter())
    .pipeThrough(createJSONParser())
    .pipeTo(sink, { signal });
}

const controller = new AbortController();

readNDJSON("/api/notes.ndjson", {
  signal: controller.signal,
  async onItem(note, lineNumber) {
    console.log("有效", lineNumber, note.id);
  },
  async onInvalid(result) {
    console.error("无效行", result.lineNumber, result.error.message);
  },
}).catch((error) => {
  if (controller.signal.aborted) {
    console.log("读取取消", controller.signal.reason);
    return;
  }
  console.error("流失败", error);
});
```

验证路径：

1. 一条 JSON 被任意拆成多个字节 chunk，仍应正确解码和拼行；
2. 多条记录落在同一个 chunk，逐条输出；
3. 中间出现无效 JSON，只记录该行，后续有效行继续；
4. `onItem` 返回慢 Promise，管道通过写入等待向上游施加背压；
5. 调用 `controller.abort()`，fetch 和管道停止并进入取消分支；
6. HTTP 非成功状态和无 body 分别产生明确错误。

若业务要求任何无效行都终止，应在 JSON 转换器中 `controller.error()` 或直接抛错，而不是产出失败结果。两种策略必须由数据契约决定。

## 11. 常见错误与调试清单

### 11.1 常见错误

1. 把网络 chunk 当成完整字符串、行或 JSON。
2. 每块单独解码 UTF-8，破坏跨块字符。
3. 取得 reader 后忘记释放锁或取消来源。
4. 不等待 `writer.ready` 或写入 Promise，造成无界缓冲。
5. 认为所有 push source 都能响应背压。
6. 使用 `tee()` 后让一个分支长期不消费。
7. 忘记处理 `pipeTo()` 返回的 rejected Promise。
8. 认为 `structuredClone` 保留类方法、访问器和描述符。
9. 尝试克隆函数或 DOM 节点。
10. 转移 buffer 后继续使用原 TypedArray。
11. 把资源列入 transfer list，却没有放入消息数据图。
12. 把“可转移”无条件理解为所有实现都零复制。

### 11.2 调试清单

- 记录 chunk 类型、字节数和累计队列，不记录敏感原文；
- 用单字节、小 chunk 和随机边界测试解析器；
- 检查 `stream.locked` 与 reader/writer 生命周期；
- 检查正常关闭、来源错误、目标错误、取消的传播；
- 模拟慢 sink 观察背压和内存；
- 观察 `desiredSize`，但不把它当精确业务计数；
- 检查页面卸载或操作取消时底层资源是否关闭；
- 克隆前列出数据类型，明确不可克隆字段；
- 转移后断言原 buffer 已分离；
- 用浏览器内存和性能工具测量，而不是假设流一定省内存。

## 12. 练习

### 练习一：CSV 分块解析

实现能处理跨 chunk 换行和引号字段的 CSV 转换流。至少验证引号内换行、转义引号、最后一行无换行和无效格式。

### 练习二：背压

创建快速 pull source 与每项延迟 50ms 的 sink。记录 `desiredSize`，比较等待写入与不受控写入的内存和完成顺序。

### 练习三：管道错误

让 transform 在第五项抛错，验证来源取消和目标 abort。再分别使用 `preventCancel`、`preventAbort`，说明资源由谁清理。

### 练习四：Worker 转移

把 16MB `ArrayBuffer` 发送给 Worker 计算校验和。分别测试结构化克隆与转移，验证结果、原 buffer 状态与耗时。

### 练习五：消息协议

设计主线程与 Worker 的消息 DTO，包含 `type`、`requestId`、版本、payload 和错误结果。禁止发送函数与类实例，并为未知消息类型设计失败分支。

## 13. 补充知识

- Fetch 的 `Response.body` 是可读字节流；body 通常只能消费一次，读取后 `bodyUsed` 会反映状态。
- `CompressionStream`、`DecompressionStream`、`TextEncoderStream` 和 `TextDecoderStream` 可直接组成管道，但需检查目标环境支持。
- 流可以作为 transferable 在支持的环境间转移，转移后原上下文不再拥有该流。
- 结构化克隆维护已访问引用映射，因此能复制循环数据图；它不是通用持久化格式。
- IndexedDB、`postMessage()` 和 `structuredClone()` 都使用结构化序列化体系，但具体上下文允许的对象类型可能不同。

## 来源

- [WHATWG Streams Standard](https://streams.spec.whatwg.org/)（访问日期：2026-07-17）
- [HTML Standard：Safe passing of structured data](https://html.spec.whatwg.org/multipage/structured-data.html)（访问日期：2026-07-17）
- [HTML Standard：Web messaging](https://html.spec.whatwg.org/multipage/web-messaging.html)（访问日期：2026-07-17）
- [MDN：The structured clone algorithm](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Structured_clone_algorithm)（访问日期：2026-07-17）
- [MDN：Transferable objects](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Transferable_objects)（访问日期：2026-07-17）
