---
stage: intermediate
---

# 事件、定时器、Worker 与 Blob URL：前端资源泄漏的生命周期治理

资源泄漏指组件、路由或任务结束后，注册过的外部资源仍在运行或占用内存。它不只表现为 heap 增长：重复监听会让逻辑执行多次，timer 会持续唤醒 CPU，worker 会保留线程和独立 heap，Blob URL 会保留二进制数据，WebSocket 会占连接，observer 会持续产生回调。

## 1. 统一模型

所有资源都可抽象为：

```text
acquire → use → release
```

获取动作必须同时产生释放动作：

| acquire | release |
|---|---|
| `addEventListener` | `removeEventListener` / abort |
| `setInterval` | `clearInterval` |
| `requestAnimationFrame` | `cancelAnimationFrame` |
| `new Worker` | `terminate` |
| `URL.createObjectURL` | `URL.revokeObjectURL` |
| `new WebSocket` | `close` |
| `observe` | `disconnect` / `unobserve` |
| `subscribe` | unsubscribe |
| `new BroadcastChannel` | `close` |
| `getUserMedia` | `track.stop()` |
| `AudioContext` | `close()` |

GC 只能回收不可达对象，不能替代协议级释放。即使某个 wrapper 能被 GC，摄像头、socket、worker 或 URL registry 的生命周期也可能继续。

## 2. 所有权

资源只能有一个明确 owner：

- 组件 owner：组件卸载释放；
- 路由 owner：离开路由释放；
- 请求 owner：请求完成/取消释放；
- 会话 owner：登出释放；
- 应用 owner：页面关闭释放；
- 共享服务 owner：引用计数归零或应用关闭释放。

共享资源不能让每个消费者直接关闭。用 lease：

```js
const socketHub = createSocketHub();
const release = socketHub.acquire(subscriber);
// 组件结束
release();
```

hub 在第一个订阅时连接、最后一个退订时延迟关闭，并处理快速重连。引用计数必须与 subscriber token 绑定，避免重复 release 变负数。

## 3. EventTarget

### callback identity

错误：

```js
window.addEventListener("resize", () => layout());
window.removeEventListener("resize", () => layout());
```

两个函数不是同一对象。正确保存 callback，或使用 AbortSignal：

```js
const controller = new AbortController();
window.addEventListener("resize", layout, {
  signal: controller.signal,
  passive: true,
});
controller.abort();
```

移除匹配主要取决于 type、callback 和 capture。passive/once 不作为匹配身份，但代码应保持选项一致以便维护。

### `once`

`{ once: true }` 在事件首次触发后自动移除，适合必然发生的单次事件。若事件永不发生，监听仍存在；组件提前结束仍应 abort。

### 委托

在稳定父节点委托 click 可减少监听数量，但 handler 必须验证 `event.target`、`closest()` 结果和容器边界：

```js
function onClick(event) {
  const button = event.target.closest("[data-action]");
  if (!button || !list.contains(button)) return;
  dispatch(button.dataset.action);
}
```

委托减少注册，不自动解决父节点生命周期和闭包捕获。

## 4. 事件总线

自定义 bus 的泄漏比 DOM listener 更难被工具直接展示：

```js
function createBus() {
  const listeners = new Map();
  return {
    on(type, fn) {
      const set = listeners.get(type) ?? new Set();
      set.add(fn);
      listeners.set(type, set);
      return () => {
        set.delete(fn);
        if (set.size === 0) listeners.delete(type);
      };
    },
    emit(type, value) {
      for (const fn of [...(listeners.get(type) ?? [])]) fn(value);
    },
  };
}
```

返回 unsubscribe token；emit 复制集合以允许回调中退订；异常策略明确；开发环境暴露 listener count。不要以字符串 owner 进行模糊批量删除。

## 5. Timer

### interval 漂移与重入

`setInterval` 不保证精确周期。回调比周期长时任务会积压或连续执行；网络轮询还可能并发：

```js
async function poll(signal) {
  while (!signal.aborted) {
    await refresh({ signal });
    await delay(5000, signal);
  }
}
```

串行 loop 比 interval 更容易控制并发、错误、退避和取消。`delay` 必须在 abort 时 clearTimeout 并 reject。

### 页面可见性

后台 timer 被节流。库存、倒计时和 token 到期应基于绝对时间：

```js
const remaining = Math.max(0, expiresAt - Date.now());
```

不要每秒 `remaining -= 1000`。恢复时重新计算；关键服务端状态重新验证。

## 6. rAF

持续动画 loop：

```js
let frameId = 0;
let running = true;

function frame(now) {
  if (!running) return;
  update(now);
  frameId = requestAnimationFrame(frame);
}

frameId = requestAnimationFrame(frame);

function dispose() {
  running = false;
  cancelAnimationFrame(frameId);
}
```

仅 cancel 当前 id 仍可能遇到 callback 已开始并请求下一帧；`running` guard 提供幂等停止。页面隐藏时按业务暂停，不要让恢复时巨大 delta 穿透。

## 7. Worker

Worker 有独立全局、事件循环和 heap：

```js
const worker = new Worker(
  new URL("./search.worker.js", import.meta.url),
  { type: "module" },
);
```

`worker.terminate()` 立即请求停止，不等待 finally。需要持久化或优雅关闭时先发 shutdown 消息，等待 ack，并设置 timeout 后 terminate。

协议需包含：

- request id/version；
- cancel；
- progress；
- success/error 可序列化结构；
- shutdown/ack；
- worker crash/messageerror；
- transferable 所有权；
- 最大队列与 backpressure。

页面卸载前不一定有时间完成异步 ack，重要数据不能依赖此时保存。

## 8. SharedWorker 与 Service Worker

SharedWorker 可被多个页面连接，单页面关闭不代表 worker 终止。每个 MessagePort 应 `start()`/`close()`，共享服务维护 client heartbeat/lease。

Service Worker 生命周期由浏览器管理，页面不能把它当常驻进程。事件必须通过 `event.waitUntil()` 延长到异步工作完成；缓存和 IndexedDB 需要版本、容量和清理。泄漏更多表现为持久化数据无限增长，而不是单页面 heap。

## 9. Blob URL

```js
const url = URL.createObjectURL(blob);
download.href = url;
```

URL 映射保留 Blob。释放时机：

- 图片/视频：确认资源不再需要；
- 下载：点击后异步 revoke，避免下载尚未读取；
- worker script：Worker 已完成加载后可 revoke，但需跨浏览器测试；
- 预览：替换旧文件或关闭时 revoke；
- 多消费者：最后一个消费者结束后 revoke。

Data URL 不需 revoke，但字符串膨胀、复制和日志风险更高。Blob URL 不是跨会话持久标识。

## 10. Observer

### ResizeObserver

对每个组件创建 observer 可行，但卸载必须 disconnect。callback 中改尺寸可能触发 resize loop；把视觉写批到 rAF，并避免读写反馈。

### IntersectionObserver

适合可见性/懒加载，不保证像素精确或立即回调。图片加载完成后可 `unobserve(target)`；列表销毁 `disconnect()`。

### MutationObserver

观察 `document.body` + `subtree:true` 范围很大。限制 target、attributeFilter 和观察期；处理完 disconnect。mutation records 可能临时引用已删除节点。

### PerformanceObserver

长期性能监控可以保持，但 owner 通常是应用级服务；SDK 关闭或测试结束需 disconnect，entry buffer 和自建数组必须有上限。

## 11. 媒体与硬件

```js
const stream = await navigator.mediaDevices.getUserMedia({ video: true });
video.srcObject = stream;

function dispose() {
  for (const track of stream.getTracks()) track.stop();
  video.srcObject = null;
}
```

删除 video 不等于关闭摄像头。屏幕分享、麦克风同理。监听 track ended 处理用户从浏览器 UI 停止。

Web Audio：

- `AudioBufferSourceNode` 一次性；
- oscillator/source 要 stop；
- 断开不再需要的 node；
- `AudioContext.close()` 释放系统音频资源；
- resume 常需用户手势。

WebGL/WebGPU 资源有 buffer/texture/pipeline；释放 JavaScript 引用不保证立即归还 GPU。显式 destroy 可用时使用，并处理 context/device lost。

## 12. 网络连接

### WebSocket

组件 owner 通常不应每次渲染建 socket。连接服务处理：

- open/message/error/close；
- 指数退避与 jitter；
- 在线/可见性；
- auth 轮换；
- outbound queue 上限；
- unsubscribe；
- 正常 close code/reason；
- server heartbeat。

反复 mount 创建多个连接会造成重复消息，即使内存增长不明显。

### SSE

`EventSource.close()` 停止自动重连。服务端支持 last-event-id；客户端切路由若仍需后台同步，可由会话服务持有，不由页面组件直接持有。

### Fetch

AbortController 取消 fetch 和 body 消费：

```js
const controller = new AbortController();
const response = await fetch(url, { signal: controller.signal });
```

Promise 本身没有通用 cancel；把 signal 贯穿 API。超时用 `AbortSignal.timeout()` 或组合 signal，检查环境支持。取消后仍要处理 AbortError，不能当业务失败弹窗。

## 13. Disposer Stack

```ts
type Disposer = () => void | Promise<void>;

function createResourceScope() {
  const stack: Disposer[] = [];
  let closed = false;

  return {
    use<T>(resource: T, dispose: (resource: T) => void | Promise<void>) {
      if (closed) throw new Error("scope closed");
      stack.push(() => dispose(resource));
      return resource;
    },
    defer(dispose: Disposer) {
      if (closed) throw new Error("scope closed");
      stack.push(dispose);
    },
    async close() {
      if (closed) return;
      closed = true;
      const errors: unknown[] = [];
      for (const dispose of stack.reverse()) {
        try {
          await dispose();
        } catch (error) {
          errors.push(error);
        }
      }
      if (errors.length) throw new AggregateError(errors);
    },
  };
}
```

逆序释放依赖；close 幂等；单项失败不阻断后续。生产还应给异步 close 超时，并在 diagnostics 暴露未释放资源类型。

## 14. React 生命周期应用

```tsx
useEffect(() => {
  const controller = new AbortController();
  const worker = new Worker(workerUrl, { type: "module" });
  const url = URL.createObjectURL(file);

  window.addEventListener("online", refresh, {
    signal: controller.signal,
  });

  return () => {
    controller.abort();
    worker.terminate();
    URL.revokeObjectURL(url);
  };
}, [file, workerUrl]);
```

每次依赖变化都会先清理旧资源再创建新资源。若 worker 与 file 无关，应拆 effect，避免不必要重启。不要把 cleanup 写在事件 handler 的不可达分支。

## 15. 案例一：实时行情页

### 症状

每次切换股票，组件新建 WebSocket 和 interval；旧连接仍推送，价格跳回旧股票，10 分钟后有 30 个连接。

### 设计

会话级 SocketService 只有一个连接；组件订阅 symbol，返回 token；最后一个订阅取消服务端频道。心跳由 service 管理，后台降频但服务端序列号校验。队列上限防断网时无限缓存。

验证快速切换 100 次，活动 socket=1，symbol subscription=1；断网重连不重复；登出立即关闭；旧序列消息不覆盖新状态。

## 16. 案例二：文件预览器

用户快速选择 50 个 100 MB 文件：

1. 旧 object URL 必须 revoke；
2. 旧解码 request 取消或版本丢弃；
3. ImageBitmap close；
4. PDF worker terminate；
5. canvas 尺寸归零/释放引用；
6. 缩略图 LRU 按字节上限；
7. error/unsupported 分支也清理；
8. 路由离开执行同一 dispose。

用进程内存、worker 数、Blob URL registry 间接行为和 GC 后 heap 共同验证。

## 17. 案例三：地图组件

地图 SDK 创建 canvas、WebGL、ResizeObserver、document pointer listener、marker DOM 和内部 worker。只 `container.remove()` 会残留资源。

wrapper contract：

```ts
interface MapAdapter {
  setData(data: GeoJSON): void;
  resize(): void;
  destroy(): Promise<void>;
}
```

destroy 调 SDK remove、关闭自建 worker、清空 marker/portal、删除 bus token，并可重复调用。E2E 打开关闭地图 50 次，检查 canvas、listener、worker、GPU context 和 heap。

## 18. 案例四：协作编辑器

资源包括 WebSocket、presence interval、document listener、worker parser、IndexedDB transaction、BroadcastChannel、selection observer。offline 并不代表释放：编辑会话仍需本地持久化；关闭文档才释放文档级资源，登出再释放账户级资源。

建立层级 scope：

```text
application scope
└── account scope
    ├── sync service
    └── document scope
        ├── socket subscription
        ├── parser worker
        ├── presence
        └── UI scope
```

父 scope 关闭必须级联子 scope；文档切换不能关闭全局网络服务。

## 19. 故障注入

资源清理测试不能只走成功路径：

- acquire 一半抛错；
- callback 正在执行时 dispose；
- dispose 调两次；
- worker 不回应 shutdown；
- socket close 前网络断开；
- Blob 解码失败；
- observer callback 抛错；
- route 快速 mount/unmount；
- 页面进入 BFCache；
- HMR 重载模块；
-浏览器拒绝权限；
-后台 timer 长时间节流。

构造函数应在部分初始化失败时回滚已获得资源。可以先创建 scope，每成功一步立刻注册 disposer，再继续下一步。

## 20. 诊断与监控

开发 diagnostics：

```js
resources.counts()
// { listeners: 12, timers: 3, workers: 1, sockets: 1, blobUrls: 4 }
```

计数必须由封装 acquire/release 更新，不能扫描浏览器内部。生产只采样汇总：

- route 离开后的活动资源数；
- 每会话 socket/worker；
- 重连次数；
- timer callback 延迟；
- Blob/图片缓存字节；
- disposer 失败；
- 长会话 heap/进程内存趋势。

不上传 URL payload、DOM 或用户文件内容。

## 21. 常见错误

1. 相信 GC 会关闭外部资源；
2. removeEventListener 使用新函数；
3. 认为 once 一定会触发；
4. interval 内并发请求；
5. 只 cancel rAF id 不设运行 guard；
6. Worker 不 terminate；
7. 创建 Blob URL 不 revoke；
8. 删除 video 却不 stop tracks；
9. observer 只 unobserve 一部分；
10. 共享 socket 被任一组件关闭；
11. cleanup 抛错导致后续未执行；
12. 只测成功初始化；
13. 资源计数无 owner；
14. 把后台节流当精确计时。

## 22. 综合练习

实现实时文件协作页，包含 WebSocket、轮询 fallback、parser worker、Blob 预览、摄像头头像、observer 和动画。

验收：

1. 每次 acquire 同时注册 release；
2. 分 application/account/document/UI 四级 scope；
3. 快速切换文档 100 次资源数平台化；
4. worker 优雅关闭超时后 terminate；
5. object URL、media track、AudioContext 正确关闭；
6. socket 共享且订阅 token 幂等；
7. timer 基于绝对时间；
8. 初始化中途失败完整回滚；
9. disposer 错误聚合；
10. 输出 heap、worker、socket、listener 和进程内存对比。

## 来源

- [DOM Standard：EventTarget](https://dom.spec.whatwg.org/#interface-eventtarget)（访问日期：2026-07-17）
- [HTML Standard：Timers](https://html.spec.whatwg.org/multipage/timers-and-user-prompts.html#timers)（访问日期：2026-07-17）
- [HTML Standard：Web Workers](https://html.spec.whatwg.org/multipage/workers.html)（访问日期：2026-07-17）
- [File API：Blob URL](https://w3c.github.io/FileAPI/#url)（访问日期：2026-07-17）
- [Media Capture and Streams](https://w3c.github.io/mediacapture-main/)（访问日期：2026-07-17）
