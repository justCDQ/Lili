---
title: URL State：把可导航、可分享的界面状态交给地址
stage: intermediate
direction: frontend
tags:
  - architecture
  - url
  - routing
---

# URL State：把可导航、可分享的界面状态交给地址

URL State 是能够通过地址恢复导航结果的状态，包括路径参数、查询参数和片段。它由浏览器历史记录承载，适合筛选、分页、选中资源和可分享视图，不适合秘密、瞬时动画或庞大草稿。

## 前置知识与能力边界

- [单一职责与组合](01-single-responsibility-composition.md)
- [Controlled 与 Uncontrolled](02-controlled-uncontrolled.md)
- React State、Context、Effect 与 TypeScript 判别联合
- 浏览器事件、HTTP 和可访问性基础

本文处理 URL 解析、规范化、History、路由同步与 SSR；HTTP URL 语法的全部标准细节不逐项展开。

## 1. 定义、所有权与数据流

URL State 是能够通过地址恢复导航结果的状态，包括路径参数、查询参数和片段。它由浏览器历史记录承载，适合筛选、分页、选中资源和可分享视图，不适合秘密、瞬时动画或庞大草稿。 URL 是地址栏和 History 的权威输入；组件不再保存一份长期同步的查询副本。用户事件通过导航更新地址，下一次渲染重新解析地址。

```mermaid
flowchart LR
    A0["地址字符串"] --> A1
    A1["路由解析"] --> A2
    A2["schema 校验"] --> A3
    A3["领域查询"] --> A4
    A4["渲染结果"] --> A5
    A5["历史导航"]
```

地址先经过 URL 语法解析，再经过应用 schema 转换为领域查询。请求 key 和界面都使用同一个解析结果；浏览器前进、后退或直接打开链接时，不需要额外同步 Effect。

## 2. 关键机制

### 2.1 pathname

表达资源层级和稳定身份，路径变化通常代表主要导航。

若边界缺失，把临时筛选塞进路径造成路由爆炸。

验证：直接打开每个路径并检查同一资源。

### 2.2 search params

表达可选筛选、排序、页码；同名键可重复且顺序可保留。

若边界缺失，把 URLSearchParams 当普通对象会丢重复值。

验证：测试 getAll、空值和编码字符。

### 2.3 hash

片段不随 HTTP 请求发送，常用于文档内定位。

若边界缺失，把授权或服务端查询依赖 hash。

验证：服务端日志确认无法读取片段。

### 2.4 解析与校验

字符串先按 schema 转成 number、enum、array，再进入业务层。

若边界缺失，Number('') 得 0 导致缺失页码被误读。

验证：fixture 覆盖缺失、重复、非法和超范围值。

### 2.5 规范化

无效值回退，默认值可从 URL 删除，键顺序稳定以减少重复地址。

若边界缺失，每次 render 重排查询导致循环导航。

验证：规范化函数满足 normalize(normalize(x)) 不变。

### 2.6 push/replace

用户产生的新导航用 push；修正默认值或高频输入常用 replace。

若边界缺失，每个键击 push 让返回键需走几十步。

验证：检查 history.length 与返回行为。

### 2.7 双向同步

URL 是权威源时 UI 从路由快照派生，事件只写 URL。

若边界缺失，同时维护 useState 形成竞态和回弹。

验证：删除本地副本后前进后退仍正确。

### 2.8 SSR

服务端和客户端必须使用相同解析与默认规则。

若边界缺失，服务端默认 page=1、客户端 page=0 产生 hydration 差异。

验证：对同一 URL 比较两端解析 JSON。

### 2.9 编码

URLSearchParams 负责 percent-encoding；空格序列化常显示为 +。

若边界缺失，手工拼接 & 和非 ASCII 破坏地址。

验证：往返测试中文、加号、空字符串。

### 2.10 长度与隐私

地址会进入历史、日志、Referer 和截图，必须保持有限且无秘密。

若边界缺失，把 token 和个人数据放入查询参数。

验证：安全扫描 URL 与日志字段。

## 3. URL 的三种状态载体

### 3.1 路径参数标识资源

`/workspaces/acme/orders/42` 表达资源层级和身份。路径通常对应页面主内容，并参与服务端路由、权限检查和缓存。把价格范围等可选筛选做成层层路径会产生大量等价地址，且删除某个筛选不自然。

### 3.2 查询参数表达可选视图

`URLSearchParams` 是有序键值序列，不是普通对象。同名键可以重复，`get()` 只返回第一个，`getAll()` 返回全部。键存在但值为空与键缺失不同：`?q=` 的 `get("q")` 是空字符串，缺失键返回 `null`。

`URLSearchParams.toString()` 不包含开头问号。它使用表单 URL 编码规则，空格序列化为 `+`，字面加号会 percent-encode。不要用字符串拼接处理中文、`&`、`+` 或重复参数。

### 3.3 fragment 只在客户端

`#section` 不随 HTTP 请求发送。它适合文档定位和同页客户端视图，服务端不能根据 fragment 返回不同数据。访问令牌放进 fragment 仍会暴露给页面脚本、浏览器历史和截图，因此也不是安全存储。

## 4. 解析、规范化与序列化

解析函数接收字符串，只输出已校验的领域值。页码需要同时满足 finite、integer 和大于零；排序只接受白名单；重复 category 要去空、去重并决定是否排序。

规范化函数产生唯一表示：

- 缺失页码等于第 1 页，序列化时省略 `page=1`；
- 默认排序 `relevance` 不写入地址；
- 多选类别去重并按稳定顺序输出；
- 未知键若不承诺透传则删除；
- 非法值回退后使用 replace 修正，避免给 History 增加坏地址。

应测试幂等性质：`serialize(parse(serialize(parse(input))))` 与第一次规范结果相同。该性质可发现键顺序不稳、默认值往返变化和空值处理不一致。

## 5. History：push 与 replace 的用户语义

`pushState` 创建新历史条目，适合用户明确提交一次新导航，例如选择类别或打开详情。`replaceState` 替换当前条目，适合修正非法地址、写入默认值或更新尚未提交的高频搜索输入。

两者都不会自动触发 `popstate`。用户前进或后退激活历史条目时才收到相应导航通知。路由库会封装这些细节，但应用仍需决定一次操作应当“可返回”还是只修正当前地址。

搜索框每个键击都 push 会让返回键逐字回退。可选策略是输入留在局部草稿，按 Enter 后 push；或防抖后 replace，用户选择建议项时再 push。

## 6. 案例一：可分享的商品搜索

地址为 `/products?q=机械键盘&category=input&category=office&sort=price-asc&page=3`。解析输出直接生成查询 key：

1. pathname 确认当前资源是商品集合。
2. `q` 保留用户输入但在发请求前按产品规则 trim；不要擅自改变大小写语义。
3. `category` 用 `getAll` 读取、去空、去重并验证允许值。
4. `sort` 解析为联合类型，未知值回退 relevance。
5. `page` 只接受正整数；修改 query、category 或 sort 时写回 page=1。
6. 规范结果同时传给查询缓存和筛选控件，没有第二份 filter state。

验证依次执行：复制地址到新标签页、刷新、前进后退、直接输入非法 `page=-1`、加入未知 sort、传两个相同类别。每次网络请求的 key 与页面控件必须一致。

失败分支是快速切换 category 时旧请求晚返回。URL 负责当前参数，查询层还需要按 key 隔离或取消旧请求；URL State 本身不解决网络竞态。

## 7. 案例二：列表详情抽屉

列表保留 `/orders?status=pending&item=42`。其中 `item` 表示可导航选择，打开抽屉使用 push，关闭使用 history back 或删除参数，取决于条目是怎样进入的。

直接打开带 item 的链接时先渲染列表骨架与详情加载态。订单不存在返回 404，用户可以关闭抽屉并保留原筛选；没有权限返回 403，不能仅把参数静默删除造成“什么也没发生”。

从列表点击 42，再点击 43，返回键应回到 42；若产品要求返回直接关闭抽屉，则第二次选择应 replace。两种行为都可实现，但必须由交互语义决定并写进自动化测试。

抽屉关闭后焦点恢复到仍存在的触发行。直接链接没有触发行时，焦点移到列表主标题。URL 决定打开对象，焦点仍属于 Interaction State。

## 8. TypeScript 核心实现

下面代码只实现本主题的核心契约；网络、DOM 或存储副作用留在调用边界。

```tsx
type Sort = "relevance" | "price-asc" | "price-desc";
type SearchState = { query: string; categories: string[]; page: number; sort: Sort };

export function parseSearch(input: string): SearchState {
  const params = new URLSearchParams(input);
  const rawPage = Number(params.get("page") ?? "1");
  const rawSort = params.get("sort");
  return {
    query: params.get("q") ?? "",
    categories: params.getAll("category").filter(Boolean),
    page: Number.isInteger(rawPage) && rawPage > 0 ? rawPage : 1,
    sort: rawSort === "price-asc" || rawSort === "price-desc" ? rawSort : "relevance",
  };
}

export function serializeSearch(state: SearchState): string {
  const p = new URLSearchParams();
  if (state.query) p.set("q", state.query);
  for (const category of [...new Set(state.categories)].sort()) p.append("category", category);
  if (state.page !== 1) p.set("page", String(state.page));
  if (state.sort !== "relevance") p.set("sort", state.sort);
  return p.toString();
}
```

类型检查用于排除结构错误，运行时仍需校验外部输入、测试时序并执行安全约束。

## 9. 方案选择

| 方案 | 适用条件 | 成本与限制 |
|---|---|---|
| 组件本地 State | 弹层动画、hover 等不可导航状态 | 刷新与分享不保留 |
| URL 查询 | 筛选、分页、可选视图 | 需要解析、长度与隐私控制 |
| 路径参数 | 稳定资源身份和主要层级 | 结构变化具有导航语义 |

选择应以所有权、生命周期、订阅范围和失败成本为依据。引入库不能替代这些判断；库只提供实现机制。

## 10. 调试与失败注入

| 现象 | 检查 | 修正 |
|---|---|---|
| 返回键需按很多次 | 输入是否每次 push | 高频修改改用 replace 或提交后 push |
| 刷新丢筛选 | 是否存在本地副本 | URL 成为单一源 |
| 重复请求 | 序列化顺序是否稳定 | 先规范化并固定排序 |
| 中文参数损坏 | 是否手工拼接 | 使用 URL/URLSearchParams |
| SSR hydration 警告 | 两端默认是否一致 | 共享纯解析函数 |
| page 变成 NaN | 是否只做类型断言 | 运行时校验并回退 |
| 秘密出现在日志 | 是否编码 token | 移到受保护会话或请求体 |
| 后退后 UI 不更新 | 是否监听路由快照 | 从 location 派生而非初始化一次 |

调试顺序是：确认输入事实，再检查所有者和转换，随后检查订阅与渲染，最后检查异步资源。跳过前序证据直接增加 Effect，通常会制造第二个状态源。

## 11. 性能、安全与运维边界

- URL schema 需有兼容策略，旧书签应迁移或给出明确错误。
- 查询参数白名单可阻止未知键污染缓存 key。
- 不要在 URL 放访问令牌、身份证号或未脱敏搜索内容。
- 服务端与客户端复用解析器并锁定默认值。
- 规范地址可用 replace，避免重复历史条目。
- 统计超长 URL、解析失败和空结果率。
- 分页与排序变化应取消旧请求。
- 无障碍焦点在导航后移到主标题或明确目标。

生产验证至少记录一次正常路径和一次故障路径；对“URL State”的结论必须能关联到日志、Profile、网络记录或自动化测试。

## 12. 与其他架构模块集成

- 与 Server State 结合：规范化 URL 生成查询 key。
- 与 Persistent State 结合：偏好可作默认值，但显式 URL 优先。
- 与表单结合：高频输入先本地编辑，提交后写 URL。
- 与路由错误边界结合：无效资源有可恢复出口。

集成时先画出事实所有者，跨边界只传递稳定契约。不要为了减少一层调用而复制同一事实。

## 13. 综合练习

实现一个可分享的商品搜索页和订单详情抽屉。解析器与序列化器不得依赖 React，路由层只负责读取 location 和执行导航。

### 验收标准

- [ ] 覆盖缺失、空值、重复键、非法数字、未知枚举、中文、空格与加号。
- [ ] 证明规范化幂等，默认参数不会反复写入 History。
- [ ] 筛选变化重置页码，前进后退时控件和请求 key 同步。
- [ ] 高频输入不会产生逐字符历史条目。
- [ ] 详情 404 与 403 都有可恢复出口和正确焦点。
- [ ] URL、日志和 Referer 检查中没有 token 或敏感个人数据。
- [ ] 服务端与客户端对同一 URL 产生完全相同的解析 JSON。

## 来源

- [WHATWG URL Standard](https://url.spec.whatwg.org/)（访问日期：2026-07-18）
- [MDN：URLSearchParams](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams)（访问日期：2026-07-18）
- [MDN：History API](https://developer.mozilla.org/en-US/docs/Web/API/History_API)（访问日期：2026-07-18）
- [React Router：URL Values](https://reactrouter.com/start/framework/url-values)（访问日期：2026-07-18）
