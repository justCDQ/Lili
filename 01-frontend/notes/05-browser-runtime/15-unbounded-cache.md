---
stage: intermediate
---

# 无界缓存：容量、淘汰、失效与并发正确性

缓存用空间和一致性复杂度换取延迟、吞吐或离线能力。无界缓存把每个新 key 永久保存，长会话、高基数输入或攻击流量会让内存与持久化空间持续增长。完整缓存必须定义 key、value、容量单位、淘汰、过期、失效、并发、错误、可观测性和数据安全。

## 1. 缓存不是 Map

```js
const cache = new Map();

export async function getUser(id) {
  if (cache.has(id)) return cache.get(id);
  const user = await fetchUser(id);
  cache.set(id, user);
  return user;
}
```

问题：

- key 数无限；
- 用户更新后数据过期；
- 并发 miss 重复请求；
- 失败是否缓存不明确；
- 登出后敏感数据仍在；
- 多租户 key 可能冲突；
- value 大小差异巨大；
- 服务端/客户端版本变化；
- 没有命中率和淘汰指标。

缓存设计先写契约，再选 Map、LRU、Cache API、IndexedDB 或服务端缓存。

## 2. 容量单位

常见上限：

- entry count；
- 估算字节；
- 解码后图像字节；
- DOM 节点数；
- 持久化配额；
- 每租户/用户配额；
- 总量 + 单项上限；
- 活跃时间窗口。

按条目限制对大小均匀的配置对象足够；图片、响应 body、编辑器文档必须按字节。JavaScript 对象精确占用难以同步计算，可用序列化大小、业务字段或资源固有尺寸近似，并保留安全裕度。

```js
function imageBytes(bitmap) {
  return bitmap.width * bitmap.height * 4;
}
```

GPU、mipmap 和副本可能更高，这只是应用预算。

## 3. LRU

Least Recently Used 淘汰最久未访问条目。Map 保持插入顺序，可实现小型 LRU：

```js
class LRUCache {
  #entries = new Map();

  constructor(maxEntries) {
    if (!Number.isInteger(maxEntries) || maxEntries <= 0) {
      throw new TypeError("maxEntries must be positive");
    }
    this.maxEntries = maxEntries;
  }

  get(key) {
    if (!this.#entries.has(key)) return undefined;
    const value = this.#entries.get(key);
    this.#entries.delete(key);
    this.#entries.set(key, value);
    return value;
  }

  set(key, value) {
    this.#entries.delete(key);
    this.#entries.set(key, value);
    while (this.#entries.size > this.maxEntries) {
      const oldest = this.#entries.keys().next().value;
      this.#entries.delete(oldest);
    }
  }

  delete(key) {
    return this.#entries.delete(key);
  }

  clear() {
    this.#entries.clear();
  }
}
```

`get()` 返回 undefined 无法区分 miss 与缓存值 undefined，可提供 `has()`/Result。高吞吐或大缓存用成熟结构；Map delete+set 会改变迭代顺序，且每次命中都是写操作。

## 4. 按字节 LRU

```js
class ByteLRU {
  #entries = new Map();
  #bytes = 0;

  constructor(maxBytes, sizeOf, dispose = () => {}) {
    this.maxBytes = maxBytes;
    this.sizeOf = sizeOf;
    this.dispose = dispose;
  }

  set(key, value) {
    this.delete(key);
    const bytes = this.sizeOf(value);
    if (bytes > this.maxBytes) {
      this.dispose(value, "oversize");
      return false;
    }
    this.#entries.set(key, { value, bytes });
    this.#bytes += bytes;
    this.#evict();
    return true;
  }

  delete(key) {
    const entry = this.#entries.get(key);
    if (!entry) return false;
    this.#entries.delete(key);
    this.#bytes -= entry.bytes;
    this.dispose(entry.value, "delete");
    return true;
  }

  #evict() {
    while (this.#bytes > this.maxBytes) {
      const key = this.#entries.keys().next().value;
      const entry = this.#entries.get(key);
      this.#entries.delete(key);
      this.#bytes -= entry.bytes;
      this.dispose(entry.value, "evict");
    }
  }
}
```

真实版还需 get 更新 recency、replace 不重复 dispose、clear、统计和异常隔离。ImageBitmap/WebGPU texture 淘汰时调用 `close()`/`destroy()`；普通对象只删除引用。

## 5. TTL 与 freshness

TTL 规定条目在某时间后过期：

```js
cache.set(key, {
  value,
  expiresAt: Date.now() + 60_000,
});
```

两种清理：

- lazy：get 时发现过期再删除；
- active：定时/堆结构清理。

只 lazy 会让永不再读取的过期 key 继续占内存；只为每条创建 timer 会引入大量 timer。可用单一周期清理、最小堆或容量淘汰兜底。

TTL 不等于业务正确性。权限、余额、库存和 feature flag 可能需要事件失效、版本号或每次验证。使用单调时间衡量进程内相对 TTL 可避免系统时钟跳变；跨会话持久过期仍用 wall clock 并容忍偏差。

## 6. Stale-While-Revalidate

SWR 允许先返回陈旧值，再后台刷新：

```text
fresh           → 直接返回
stale but usable → 返回旧值 + 单例刷新
expired hard     → 等待新值/失败
```

必须定义：

- fresh TTL；
- stale window；
- hard expiration；
- 刷新失败是否继续用旧值；
- UI 是否显示更新时间；
- 何时通知订阅者；
- 旧请求是否可覆盖新版本。

认证、价格确认等不应不加提示返回陈旧结果。内容列表、头像、非关键配置更适合。

## 7. 请求去重

把 in-flight Promise 作为短期缓存：

```js
const inFlight = new Map();

async function singleFlight(key, load) {
  if (inFlight.has(key)) return inFlight.get(key);
  const promise = Promise.resolve().then(load);
  inFlight.set(key, promise);
  try {
    return await promise;
  } finally {
    if (inFlight.get(key) === promise) inFlight.delete(key);
  }
}
```

边界：

- 失败后必须删除，除非明确负缓存；
- 各消费者 cancel 不应随意取消共享请求；
- key 包含所有影响响应的参数与身份；
- finally 的 identity check 防旧 Promise 删除新请求；
- 长期 pending 需要 timeout；
- 返回可变对象时消费者可能互相污染。

可给每个消费者独立 AbortSignal，最后一个消费者离开才 abort 底层请求。

## 8. Cache Stampede

热点 key 同时过期，许多请求一起回源。浏览器单页面也可能由多个组件触发。措施：

- single-flight；
- TTL 加随机 jitter；
- 提前异步刷新；
- stale-if-error；
- 服务端合并请求；
- 并发上限；
- retry backoff；
- 失败短时负缓存。

负缓存只保存可安全重试的失败分类。例如 404 可短时缓存，401 必须走认证刷新，网络错误不宜长缓存。错误对象不能带巨大 response/DOM。

## 9. Key 设计

```js
const key = JSON.stringify({ endpoint, params, locale, tenantId });
```

JSON stringify 的属性顺序、undefined、Date、BigInt 和对象 identity 都需处理。更安全使用规范化 key builder：

```js
function userKey({ tenantId, userId, fields, locale }) {
  const normalizedFields = [...fields].sort().join(",");
  return `user:v3:${tenantId}:${userId}:${locale}:${normalizedFields}`;
}
```

包含 schema/version 能在数据结构变更时自然 miss。不要把 access token 本身写入持久 key 或日志；租户/用户边界必须存在，登出清除。

## 10. Invalidation

“缓存失效是难题”不意味着只能缩短 TTL。策略：

### 写后更新

mutation 成功后用服务端返回值替换对应 cache。响应必须是权威新版本。

### 写后删除

删除相关 key，下次读取回源。简单但短期增加延迟。

### 事件失效

服务端推送 entity/version，客户端 invalidate。处理断线、乱序、重复和补偿同步。

### Tag

条目记录 tag，如 `project:42`，mutation 使相关 tag 失效。tag index 本身也必须清理和有界。

### Version

key 或 value 含 revision；旧请求完成时只在版本仍匹配时 commit。

## 11. WeakMap

WeakMap key 必须是对象或非注册 symbol，且不阻止 key 被 GC，适合：

- DOM → metadata；
- object → memoized derived data；
- instance → private state。

它不可枚举、无法读 size、GC 时机不可预测，不适合需要主动淘汰、命中统计、持久化或按 ID 查询的业务缓存。

```js
const metadata = new WeakMap();
metadata.set(element, { measuredWidth: 240 });
```

若 value 强引用 key 的外部 owner，整体生命周期仍可能不符合预期；WeakMap 不是泄漏免疫。

## 12. Memoization

函数 memoization 常把参数组合永久保存：

```js
const memo = new Map();
function parse(source) {
  if (!memo.has(source)) memo.set(source, expensiveParse(source));
  return memo.get(source);
}
```

编辑器中 `source` 每次键入都不同，缓存命中低却保存全文和 AST。方案：

- 只缓存最近 N 个版本；
- 以 documentId+revision 并淘汰；
- 增量解析；
- 只在单次 render 生命周期 memo；
- 统计 hit rate 后决定是否保留。

React `useMemo` 是性能提示和当前组件实例缓存，不是语义保证，也不能替代跨组件数据缓存。

## 13. 浏览器 HTTP Cache

HTTP cache 由响应头控制，浏览器负责容量和淘汰。应用内 Map 缓存 fetch response 往往重复占用：

- HTTP cache 已保存 body；
- data client 又保存 parsed object；
-组件 state 再复制；
-持久化层再保存。

定义每层目的。HTTP cache 处理传输复用；应用 cache 处理结构化数据、订阅和 mutation；IndexedDB 支持离线。避免三层都无限。

`Cache-Control: no-cache` 允许存储但使用前需验证；`no-store` 才禁止存储。敏感个性化响应正确设置 `private`、Vary 和认证边界。

## 14. Cache API

Service Worker Cache API 按 Request/Response 保存，不自动执行 HTTP freshness/容量策略：

```js
const cache = await caches.open("assets-v4");
await cache.put(request, response);
```

升级时删除旧命名 cache：

```js
const keep = new Set(["assets-v4"]);
for (const name of await caches.keys()) {
  if (!keep.has(name)) await caches.delete(name);
}
```

还需单项/总量策略、opaque response 成本、配额失败、用户清理、版本回滚和多 tab 协调。cache.put 的 response body 需可消费的 clone。

## 15. IndexedDB 缓存

IndexedDB 适合大结构化离线数据，但：

- schema migration；
- transaction 失败/中止；
- quota exceeded；
- eviction；
- 多 tab versionchange；
- 索引体积；
- 清理游标；
- 数据加密/登出；
- 服务端版本冲突。

按 `lastAccessedAt` 建索引可批量淘汰，但每次 read 更新又增加写放大。可采样更新或按 bucket 记录。持久化前先问离线价值，而不是把内存 cache 全量镜像。

## 16. 案例一：头像缓存

### 症状

聊天页面按 userId 保存 ImageBitmap，无上限；长会话看过 20k 用户，内存数百 MiB。

### 设计

- key：tenant/user/avatarVersion/size；
- value：ImageBitmap + decoded bytes；
- 最大 64 MiB、单项 4 MiB；
- byte LRU；
- 淘汰 `bitmap.close()`；
- HTTP cache 保存原图；
- 当前可见头像 pin，离屏可淘汰；
- avatarVersion 变化自然 miss；
- 登出 clear；
- 失败短时缓存，不保存巨大 Error。

验证快速滚动、DPR 切换、头像更新、断网、内存压力。

## 17. 案例二：搜索结果

每个 query 缓存结果会受到高基数输入：

```text
a, ab, abc, abcd, ...
```

设计：

- query normalize；
- 少于 2 字符不缓存；
- LRU 50；
- fresh 30s、stale 5min；
- single-flight；
- requestId 防乱序；
- 结果仅保存 ID，实体归一化缓存；
- 权限/租户进入 key；
- 统计 hit/miss/eviction；
- 输入法 composition 结束后查询。

若 hit rate <5%，缓存可能不值得其内存和一致性成本。

## 18. 案例三：离线项目

项目文档存 IndexedDB，不只是缓存，而是本地副本。不能用 LRU 静默删除未同步编辑。分类：

| 数据 | 策略 |
|---|---|
| 已同步可重取附件 | byte LRU |
| 未同步操作日志 | 不淘汰，配额预警 |
| 项目列表 | TTL + 刷新 |
| 缩略图 | LRU |
| 密钥/会话 | 登出清除 |

配额接近上限时先删可重取缓存，保留用户创作，阻止继续导入前明确提示。

## 19. 案例四：权限缓存

前端权限缓存只用于 UI 体验，后端仍必须授权。key 包含 tenant、actor、resource、action 和 policyVersion。角色变化通过事件失效；重新聚焦或关键操作前刷新。陈旧 allow 不能让后端放行，陈旧 deny 可能误隐藏功能，所以 UI 提供刷新/错误状态。

不把权限列表永久写 localStorage；登出和切换组织必须原子清除。

## 20. 可观测性

每个 cache 至少记录：

- entries/bytes；
- hit/miss/hit rate；
- stale hit；
- load latency；
- in-flight count；
- eviction reason；
- oversize reject；
- negative hit；
- refresh error；
- oldest age；
- per-tenant usage；
- quota error。

高基数 key 不进入 metric label。调试日志采样 hash/key type，不记录 token、query 原文或个人数据。

## 21. 内存压力与自适应

Web 没有通用可靠的“系统内存不足”事件供所有页面使用。应用应采用固定安全预算和设备能力档位，不能等 OOM 才清理。可根据：

- 设备类别与实测；
- 页面可见性；
- 当前活跃资源；
- cache hit/eviction；
- 长会话内存趋势；
- 浏览器存储估算。

自适应缩小缓存会降低命中率，要防反复扩缩振荡。关键数据不能在压力下无提示丢失。

## 22. 常见错误

1. 用 Map 即完成缓存；
2. 只限条目不管字节；
3. TTL 到期但不主动/容量清理；
4. 把所有失败永久缓存；
5. 并发 miss 重复回源；
6. key 缺租户、语言或版本；
7. 旧请求覆盖新值；
8. WeakMap 当 LRU；
9. localStorage 无限保存；
10. Cache API 当自动 HTTP cache；
11. 淘汰 ImageBitmap 不 close；
12. 登出不清用户数据；
13. 只看命中率不看延迟/内存；
14. 缓存不可重取的用户编辑；
15. 高基数 key 写监控标签。

## 23. 综合练习

实现一个 data client，支持内存 byte LRU、TTL/SWR、single-flight、版本失效和 IndexedDB 离线层。

验收：

1. key 包含 tenant、locale、schema version；
2. 内存同时限制条目与字节；
3. 过大 value 拒绝缓存；
4. 淘汰资源执行 disposer；
5. 并发 100 次只回源一次；
6. 请求失败后 in-flight 删除；
7. 旧响应不覆盖新 revision；
8. 过期项即使不再读取也会清理；
9. 登出清敏感数据；
10. 配额不足优先删可重取项；
11. 输出 hit/miss/bytes/eviction；
12. 长会话 key 数与字节平台化。

## 来源

- [HTTP Caching](https://httpwg.org/specs/rfc9111.html)（访问日期：2026-07-17）
- [Service Workers：Cache](https://w3c.github.io/ServiceWorker/#cache-interface)（访问日期：2026-07-17）
- [Indexed Database API 3.0](https://w3c.github.io/IndexedDB/)（访问日期：2026-07-17）
- [Storage Standard](https://storage.spec.whatwg.org/)（访问日期：2026-07-17）
- [ECMAScript：WeakMap Objects](https://tc39.es/ecma262/multipage/keyed-collections.html#sec-weakmap-objects)（访问日期：2026-07-17）
