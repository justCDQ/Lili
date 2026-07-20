# Promise 组合、异步迭代与取消协议

单个 Promise 表示一个异步结果；组合器表达多个结果的完成策略；异步迭代器表达随时间逐项到达的数据；`AbortSignal` 表达调用方不再需要结果。真实异步流程通常需要同时处理并发、失败、背压和取消。

## 1. Promise 的状态与命运

Promise 初始为 pending，之后只能变为 fulfilled 或 rejected，且状态一旦确定不能改变。fulfilled 与 rejected 统称 settled。

“resolved”不一定等于 fulfilled：Promise 被另一个 pending Promise 或 thenable 接管后已经 resolved，但仍可能处于 pending，最终也可能 rejected。

```js
const inner = new Promise((resolve) => {
  setTimeout(() => resolve("完成"), 20);
});

const outer = new Promise((resolve) => {
  resolve(inner);
});

outer.then(console.log); // 等待 inner 后输出“完成”
```

构造器 executor 会同步执行；`then`、`catch`、`finally` 注册的 reaction 始终异步运行。

```js
console.log("A");

const promise = new Promise((resolve) => {
  console.log("B");
  resolve("D");
});

promise.then(console.log);
console.log("C");
// A, B, C, D
```

### 1.1 链式返回规则

`then()` 返回一个新 Promise。回调的结果决定新 Promise：

- 返回普通值：新 Promise fulfilled；
- 返回 Promise 或 thenable：采用其最终状态；
- 抛出错误：新 Promise rejected；
- 没有 `return`：以 `undefined` fulfilled。

```js
function loadProfile() {
  return Promise.resolve({ id: "u-1" })
    .then((profile) => ({ ...profile, active: true }))
    .then((profile) => {
      if (!profile.active) throw new Error("账号不可用");
      return profile;
    });
}

loadProfile().then(console.log, console.error);
```

忘记返回嵌套 Promise，会使外层链提前完成，后续错误也无法被该链捕获。

```js
function incorrect() {
  return Promise.resolve().then(() => {
    fetch("/api/profile"); // 没有返回
  });
}

function correct() {
  return Promise.resolve().then(() => fetch("/api/profile"));
}
```

### 1.2 catch 与 finally

`catch(onRejected)` 等价于 `then(undefined, onRejected)`。catch 返回替代值会恢复为 fulfilled；重新抛出才继续拒绝。

```js
async function readOptionalSettings(load) {
  try {
    return await load();
  } catch (error) {
    if (error.code === "NOT_FOUND") return {};
    throw error;
  }
}
```

`finally(callback)` 用于不依赖结果的清理。正常返回不会替换原值；若 finally 抛错或返回 rejected Promise，则以新错误替换原结果。

## 2. 四个 Promise 组合器

四个组合器都接收 iterable，并把其中的普通值通过当前 Promise 构造器的 resolve 规则处理。

| API | fulfilled 条件 | rejected 条件 | 结果 |
| --- | --- | --- | --- |
| `Promise.all` | 全部 fulfilled | 任一个 rejected | 按输入顺序的值数组 |
| `Promise.allSettled` | 全部 settled | 通常不因成员失败而 rejected | 状态结果数组 |
| `Promise.any` | 任一个 fulfilled | 全部 rejected | 首个成功值或 `AggregateError` |
| `Promise.race` | 首个成员 settled | 首个成员 rejected | 首个 settled 的值或原因 |

### 2.1 Promise.all

适合所有结果都必需的独立操作。结果数组保持输入顺序，不按完成顺序排列。

```js
async function loadDashboard(api) {
  const [profile, projects, notices] = await Promise.all([
    api.getProfile(),
    api.getProjects(),
    api.getNotices(),
  ]);

  return { profile, projects, notices };
}
```

某一项拒绝时，组合 Promise 立即拒绝，但其他操作不会自动取消。若支持取消，应共享 signal 并在失败策略中明确中止。

空 iterable 的 `Promise.all([])` 以空数组 fulfilled。

### 2.2 Promise.allSettled

适合需要汇总每项成功与失败的批处理。

```js
const results = await Promise.allSettled([
  Promise.resolve("HTML"),
  Promise.reject(new Error("CSS 加载失败")),
]);

for (const result of results) {
  if (result.status === "fulfilled") {
    console.log(result.value);
  } else {
    console.error(result.reason);
  }
}
```

每项结果是 `{ status: "fulfilled", value }` 或 `{ status: "rejected", reason }`。它不会把业务失败变成成功；调用方仍需检查每个状态。

### 2.3 Promise.any

适合多个等价来源中任一成功即可的场景。失败会被忽略，直到全部失败；全部失败时以 `AggregateError` 拒绝，其 `errors` 按输入顺序保存原因。

```js
async function loadFromMirrors(loaders) {
  try {
    return await Promise.any(loaders.map((load) => load()));
  } catch (error) {
    if (error instanceof AggregateError) {
      console.error(error.errors);
    }
    throw error;
  }
}
```

空 iterable 的 `Promise.any([])` 会以 `AggregateError` 拒绝。

### 2.4 Promise.race

`race` 采用首个 settled 成员，不要求成功。常用于竞争事件，但不负责取消输家。

```js
function delay(ms, value) {
  return new Promise((resolve) => setTimeout(resolve, ms, value));
}

console.log(await Promise.race([
  delay(20, "first"),
  delay(40, "second"),
]));
```

空 iterable 的 `Promise.race([])` 会永久保持 pending。用 `race` 实现超时但不取消原操作，会留下仍占用网络、计时器或连接的工作。

## 3. 并发与顺序

连续 `await` 会在前一步完成后才启动下一步。若任务相互独立，应先启动再组合等待。

```js
async function sequential(api) {
  const profile = await api.getProfile();
  const projects = await api.getProjects();
  return { profile, projects };
}

async function concurrent(api) {
  const profilePromise = api.getProfile();
  const projectsPromise = api.getProjects();
  const [profile, projects] = await Promise.all([
    profilePromise,
    projectsPromise,
  ]);
  return { profile, projects };
}
```

并发不是越多越好。无界创建 Promise 可能耗尽连接、触发限流或占用大量内存。限制并发数时可用 worker pool。

```js
async function mapWithConcurrency(items, limit, mapper) {
  if (!Number.isInteger(limit) || limit <= 0) {
    throw new RangeError("limit 必须是正整数");
  }

  const results = new Array(items.length);
  let nextIndex = 0;

  async function worker() {
    while (true) {
      const index = nextIndex;
      nextIndex += 1;
      if (index >= items.length) return;
      results[index] = await mapper(items[index], index);
    }
  }

  const workerCount = Math.min(limit, items.length);
  await Promise.all(Array.from({ length: workerCount }, () => worker()));
  return results;
}
```

此实现快速失败，但已启动的 worker 不自动停止；需要取消时把 signal 传给 mapper 并检查。

## 4. async/await 的错误边界

async 函数总是返回 Promise。函数内 `return value` 产生 fulfilled Promise，抛错产生 rejected Promise。

```js
async function parseResponse(response) {
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  return response.json();
}
```

`try/catch` 只能捕获其动态执行范围内实际等待的拒绝。启动 Promise 后既不 `await` 也不返回，错误会脱离该边界。

```js
async function saveAll(save) {
  try {
    await save("profile");
    await save("settings");
  } catch (error) {
    throw new Error("保存失败", { cause: error });
  }
}
```

不要用一个巨大 catch 把校验、请求、解析和渲染错误都转换成同一消息。只捕获能处理或需要增加上下文的层级，并保留 `cause`。

## 5. 异步迭代协议

异步 iterable 实现 `Symbol.asyncIterator`，返回异步迭代器。其 `next()` 返回 Promise，该 Promise fulfilled 后得到 `{ value, done }`。

```js
const asyncRange = {
  from: 1,
  to: 3,

  [Symbol.asyncIterator]() {
    let current = this.from;
    const to = this.to;

    return {
      async next() {
        await new Promise((resolve) => setTimeout(resolve, 10));
        if (current > to) return { value: undefined, done: true };
        return { value: current++, done: false };
      },
    };
  },
};

for await (const value of asyncRange) {
  console.log(value);
}
```

`for await...of` 可以消费异步 iterable，也可包装同步 iterable；每轮会等待迭代结果及其中的异步值。它只能出现在允许 `await` 的上下文，如 async 函数或模块顶层。

### 5.1 异步生成器

`async function*` 简化异步 iterable。函数体可 `await`，`yield` 逐项产出。

```js
async function* paginate(fetchPage) {
  let cursor = null;

  do {
    const page = await fetchPage(cursor);
    for (const item of page.items) yield item;
    cursor = page.nextCursor;
  } while (cursor !== null);
}
```

消费者每次请求下一项后，生产者才继续推进，天然形成一对一背压。若生产者内部主动缓冲多页，缓冲策略仍需自行控制。

### 5.2 提前退出与清理

`break` 或异常离开 `for await...of` 时，消费方会调用异步迭代器的 `return()` 并等待清理结果。异步生成器可用 `try/finally` 释放资源。

```js
async function* readMessages(openConnection) {
  const connection = await openConnection();
  try {
    while (true) {
      const message = await connection.read();
      if (message === null) return;
      yield message;
    }
  } finally {
    await connection.close();
  }
}
```

## 6. AbortController 与 AbortSignal

`AbortController` 是发出取消请求的一方；`controller.signal` 是只读信号，传给支持取消的操作；`controller.abort(reason)` 使信号永久进入 aborted 状态。

```js
const controller = new AbortController();

fetch("/api/notes", { signal: controller.signal })
  .catch((error) => {
    if (controller.signal.aborted) {
      console.log("取消原因", controller.signal.reason);
      return;
    }
    throw error;
  });

controller.abort(new DOMException("页面离开", "AbortError"));
```

信号是一次性的，取消后不能复位。新的操作需要新的 controller。取消表示调用方不再等待，不保证撤销服务器已执行的写入。

### 6.1 让自定义 API 支持取消

自定义异步操作应：

1. 入口调用 `signal?.throwIfAborted()`，处理已取消信号；
2. 监听一次 `abort`；
3. 取消底层资源并以 `signal.reason` 拒绝；
4. 正常完成和失败时移除监听；
5. 避免重复 settle 和监听器泄漏。

```js
function abortableDelay(ms, { signal } = {}) {
  return new Promise((resolve, reject) => {
    signal?.throwIfAborted();

    const timerId = setTimeout(() => {
      signal?.removeEventListener("abort", onAbort);
      resolve();
    }, ms);

    function onAbort() {
      clearTimeout(timerId);
      reject(signal.reason);
    }

    signal?.addEventListener("abort", onAbort, { once: true });
  });
}
```

### 6.2 超时与组合信号

当前平台可用 `AbortSignal.timeout(ms)` 创建超时信号，用 `AbortSignal.any(signals)` 创建任一输入取消即取消的组合信号。使用前应按目标运行环境检查兼容性。

```js
async function loadWithPolicy(url, userSignal) {
  const timeoutSignal = AbortSignal.timeout(5000);
  const signal = userSignal
    ? AbortSignal.any([userSignal, timeoutSignal])
    : timeoutSignal;

  const response = await fetch(url, { signal });
  if (!response.ok) throw new Error(`HTTP ${response.status}`);
  return response.json();
}
```

组合取消后应检查 `signal.reason`。超时和用户取消可有不同原因，但通过 `AbortSignal.any()` 后不能只凭某个原始 signal 的时间状态推断实际先后，应保留所需上下文。

## 7. 完整案例：有限并发、增量结果与整体取消

```js
function createAbortError(message) {
  return new DOMException(message, "AbortError");
}

async function* loadNotes({ ids, load, concurrency = 3, signal }) {
  if (!Array.isArray(ids)) throw new TypeError("ids 必须是数组");
  if (typeof load !== "function") throw new TypeError("load 必须是函数");
  if (!Number.isInteger(concurrency) || concurrency <= 0) {
    throw new RangeError("concurrency 必须是正整数");
  }

  signal?.throwIfAborted();
  const executing = new Set();
  let nextIndex = 0;

  function start(index) {
    const promise = Promise.resolve()
      .then(() => load(ids[index], { signal }))
      .then(
        (value) => ({ index, status: "fulfilled", value }),
        (reason) => ({ index, status: "rejected", reason }),
      );

    const tracked = promise.finally(() => executing.delete(tracked));
    executing.add(tracked);
  }

  try {
    while (nextIndex < ids.length || executing.size > 0) {
      signal?.throwIfAborted();

      while (nextIndex < ids.length && executing.size < concurrency) {
        start(nextIndex);
        nextIndex += 1;
      }

      if (executing.size > 0) {
        yield await Promise.race(executing);
      }
    }
  } finally {
    // 底层 load 必须实际使用同一个 signal，才能停止已启动工作。
  }
}

async function collectNotes(options) {
  const values = [];
  const failures = [];

  for await (const result of loadNotes(options)) {
    if (result.status === "fulfilled") {
      values.push(result.value);
    } else if (options.signal?.aborted) {
      throw options.signal.reason;
    } else {
      failures.push({ id: options.ids[result.index], reason: result.reason });
    }
  }

  return { values, failures };
}

const controller = new AbortController();
const fakeLoad = async (id, { signal }) => {
  await abortableDelay(10, { signal });
  if (id === "bad") throw new Error("模拟加载失败");
  return { id, title: `笔记 ${id}` };
};

collectNotes({
  ids: ["js-13", "bad", "js-15", "js-16"],
  load: fakeLoad,
  concurrency: 2,
  signal: controller.signal,
})
  .then(({ values, failures }) => {
    console.log(values.map((value) => value.id));
    console.log(failures.map((failure) => failure.id));
  })
  .catch((error) => {
    if (error.name === "AbortError") console.log("加载已取消");
    else console.error(error);
  });
```

运行预期：同时最多两个 `load`；完成结果按实际完成顺序产出；`bad` 被记录为局部失败；取消时底层延迟及时拒绝。若业务要求保持输入顺序，可在收集后按 `index` 排序，代价是可能延迟展示早完成的后项。

案例中的 `Promise.race` 只选择下一条完成结果；它不会取消其他成员。真正取消来自贯穿每层的同一 signal。

## 8. 常见错误与调试清单

### 8.1 常见错误

1. 忘记从 `then` 返回 Promise，导致链提前完成。
2. 用 `forEach(async () => {})` 后期待外层等待回调。
3. 对本应并发的独立请求逐个 `await`。
4. 一次启动无界请求，触发限流和内存压力。
5. 认为 `Promise.all` 首次失败会取消其他成员。
6. 用 `race` 超时但不取消输掉的操作。
7. 忽略 `allSettled` 中的 rejected 项。
8. 把取消与网络失败统一吞掉。
9. 复用已经 aborted 的 signal。
10. 自定义 API 注册 abort 监听后不清理。
11. 提前退出异步迭代却未实现资源清理。
12. 认为 abort 能撤销服务端已经提交的副作用。

### 8.2 调试清单

- 标记每个操作的启动、settled 与消费时间；
- 检查 Promise 是否被 `return` 或 `await`；
- 确认组合器语义与产品失败策略一致；
- 检查结果顺序是输入顺序还是完成顺序；
- 记录当前并发数并验证上限；
- 在入口检查 `signal.aborted` 或调用 `throwIfAborted()`；
- 检查 signal 是否传到所有底层 API；
- 正常、失败、超时、用户取消分别测试；
- 提前 `break` 验证异步生成器 `finally`；
- 检查未处理 Promise rejection 和残留监听器。

## 9. 练习

### 练习一：组合器选择

分别为“全部配置必需”“批量上传允许部分失败”“多个镜像任一成功”“用户点击与超时竞争”选择组合器，并实现空输入和全部失败测试。

### 练习二：并发池

扩展 `mapWithConcurrency`，支持 signal、保持结果顺序、快速失败和收集全部失败两种模式。统计最大实际并发数验证限制。

### 练习三：异步分页

实现带游标的异步生成器，逐项产出页面内容。消费者读取五项后 `break`，验证连接或请求资源得到清理。

### 练习四：超时

实现一个既接受用户 signal 又有 2 秒超时的请求。区分用户取消、超时、HTTP 错误和 JSON 解析错误。

### 练习五：取消副作用

设计创建订单流程，说明客户端取消后为何仍需幂等键、服务端状态查询或补偿操作，而不能只依赖 abort。

## 10. 补充知识

- Promise 表示结果而不是可重复事件流；多个随时间产生的值应考虑异步 iterable 或流。
- Promise executor 中抛错会拒绝 Promise，但 executor 之外独立异步回调抛错不会自动关联。
- `for await...of` 顺序等待每项；对独立任务需要并发时，应显式使用池或组合器。
- `AbortSignal` 继承 `EventTarget`，取消事件只触发一次，`reason` 保存取消原因。
- 取消协议是协作式的：只有观察并响应 signal 的层才会停止。

## 来源

- [ECMAScript 2026：Control Abstraction Objects](https://tc39.es/ecma262/2026/multipage/control-abstraction-objects.html)（访问日期：2026-07-17）
- [ECMAScript 2026：Async Function Definitions](https://tc39.es/ecma262/2026/multipage/ecmascript-language-functions-and-classes.html#sec-async-function-definitions)（访问日期：2026-07-17）
- [DOM Standard：Aborting ongoing activities](https://dom.spec.whatwg.org/#aborting-ongoing-activities)（访问日期：2026-07-17）
- [MDN：Promise concurrency](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise#promise_concurrency)（访问日期：2026-07-17）
- [MDN：for await...of](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/for-await...of)（访问日期：2026-07-17）
