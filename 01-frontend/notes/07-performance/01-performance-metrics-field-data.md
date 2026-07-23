---
title: 性能指标、真实用户数据与性能预算
stage: intermediate
direction: frontend
tags:
  - frontend
  - 07-performance
---

# 性能指标、真实用户数据与性能预算

性能指标把“页面是否及时呈现、是否及时响应、是否稳定、是否耗尽资源”转成可观察数据。正确使用指标的前提是分清测量对象：TTFB 主要描述请求开始阶段；FCP 和 LCP 描述内容何时出现；INP 描述交互到下一次绘制；CLS 描述非预期位移；Long Task、FPS 和 Heap 是定位运行时瓶颈的诊断信号。它们不能由一次本地刷新或一个平均值替代。

## 前置知识与边界

- [从 DNS 到 HTTP：浏览器请求的连接路径与诊断](../05-browser-runtime/01-dns-tcp-tls-http.md)
- [Style、Layout、Paint 与 Composite](../05-browser-runtime/07-rendering-pipeline.md)
- [Long Task 与 Layout Thrashing](../05-browser-runtime/09-long-task-layout-thrashing.md)
- [环境变量、配置、构建、部署与错误监控](../04-typescript-frameworks/16-config-deploy-monitoring.md)

本文讨论浏览器侧性能证据。它不能证明服务端授权、库存、金额或审计正确；这些不变量必须继续在服务端验证。性能采集也不是无边界日志系统：URL query、表单、认证信息、完整 IP、稳定跨站标识和用户内容都要按数据分类最小化处理。

## 1. 两种数据：实验室数据与真实用户数据

实验室数据（lab data）在固定浏览器、CPU、网络和脚本条件下运行。它适合复现某次回归、比较两个提交、在 CI 设门禁。它的弱点是环境单一：模拟网络、干净缓存和自动化操作不能代表所有设备、地区、登录态或第三方响应。

真实用户数据（RUM，Real User Monitoring）由实际访问产生。它能反映真实设备性能、缓存命中、网络抖动和真实交互，但有采样偏差、隐私义务、灰度版本混杂和用户行为差异。RUM 用于确认影响范围和持续回归；它不应在没有样本量、版本覆盖率和路由分组的情况下被拿来宣布因果关系。

```mermaid
flowchart LR
    A["提交或发布候选"] --> B["实验室 trace 与预算"]
    B --> C["小范围发布"]
    C --> D["RUM：路由、设备、release"]
    D --> E["分位数、错误与业务结果"]
    E --> F["修复、回滚或扩大发布"]
    F --> B
```

实验室与 RUM 必须用相同的页面版本、路由和场景名称连接。例如实验室测“未登录移动端首页的 LCP”，线上也按该路由、设备类别和 release 查看；把登录态后台页与公开首页混合，任何分位数都失去解释力。

## 2. 指标地图：每个数字回答不同问题

| 指标 | 用户可观察的问题 | 最常见的误读 | 首先查看的工具 |
| --- | --- | --- | --- |
| TTFB | 请求为什么迟迟没有开始返回内容 | 把它当成完整页面加载时间 | Network、服务端/CDN trace |
| FCP | 用户何时第一次看到内容 | 把空白背景或骨架动画当内容绘制 | Performance、Lighthouse、RUM |
| LCP | 首屏主要内容何时完成绘制 | 只压缩图片，不看发现时机和服务端响应 | LCP 归因、Network、Performance |
| INP | 点击、键盘输入后何时看见更新 | 只测 click handler 同步时间 | Event Timing、Performance trace |
| CLS | 为什么页面在没有用户触发时跳动 | 认为所有移动都计入 CLS | Layout Shift entries、DevTools |
| Long Task | 哪段主线程工作连续占用太久 | 认为没有 Long Task 就没有输入延迟 | Performance trace、Long Tasks API |
| FPS | 连续动画/滚动是否丢帧 | 用单次 FPS 代替交互响应 | Frames track、真实设备录制 |
| JS Heap | 内存是否持续增长或峰值过高 | 把一次大 heap 直接判定为泄漏 | Heap snapshot、Allocation timeline |

### 2.1 TTFB：从开始请求到收到响应首字节

TTFB（Time to First Byte）是在用户代理开始获取资源后，到收到响应第一个字节的时间。它包含连接复用或建连、请求发送、CDN/服务端处理和响应开始传输等阶段；不同工具的时间起点可能略有差异，因此比较时必须使用同一采集实现。

TTFB 高通常提示：缓存未命中、边缘节点距离远、服务端渲染/数据库慢、重定向过多、连接没有复用，或 TLS/DNS 建连代价高。它不说明 JavaScript 是否阻塞，也不说明首屏图片何时出现。若 HTML 的 TTFB 是 200ms 而 LCP 是 4s，瓶颈常在后续资源发现、下载、解码或渲染；若 HTML 的 TTFB 是 2.5s，先调查服务端与 CDN 通路。

在 Network 面板中选择文档请求，分别查看排队、DNS、连接、请求发送、等待响应与内容下载。诊断不能只看一个总时间：同样 1.5s 的 TTFB，可能是首次跨洲 TLS 连接，也可能是每次都执行慢查询。把 CDN cache status、服务端 request ID、区域、协议和 release 关联，才能把浏览器症状转给正确所有者。

### 2.2 FCP：第一块内容真正画到屏幕的时刻

FCP（First Contentful Paint）在浏览器首次绘制文本、图片（包括背景图）、SVG 或 canvas 内容时记录。仅改变非内容背景色不构成 FCP。它回答“用户是否还面对空白页”，不保证用户已经看到主要任务或可以交互。

影响 FCP 的典型因素是 HTML 是否及时到达、渲染阻塞 CSS、同步脚本、字体策略和首批可见文本。把所有 CSS 拆成多个晚到的异步文件可能降低阻塞，却造成首屏无样式或闪烁；删除关键 CSS 也不是普适优化。应优先让首屏必要样式和文本可被早发现，而把不影响首屏的脚本/样式推迟到合适边界。

FCP 适合检测空白屏回归。对于一个纯图片 landing page，FCP 可能由小图标触发而主视觉仍很晚，因此不能用“FCP 变好”替代 LCP 验收。对于已登录 SPA，路由切换的用户感受也未必由导航时的 FCP 代表，要另行观察交互和页面骨架到真实数据的状态。

### 2.3 LCP：最大可见内容何时完成渲染

LCP（Largest Contentful Paint）报告页面生命周期中最大的候选内容元素完成渲染的时间。候选通常是 `<img>`、视频 poster、带背景图的元素或文本块；随着页面加载和布局变化，候选可以更换。用户输入后，后续 LCP 变化不再是页面加载体验的同一类信号，因此诊断应查看最终候选及其资源链。

优化 LCP 先分解时间：

1. **TTFB**：HTML 或资源请求何时能开始。
2. **资源发现延迟**：主视觉是否由 CSS/JS 很晚创建，浏览器何时知道 URL。
3. **资源下载延迟和时长**：请求优先级、竞争、压缩、尺寸、CDN 与缓存是否合理。
4. **元素渲染延迟**：下载后是否被主线程、字体、解码、hydration 或样式阻塞。

例如首屏图片放在客户端组件中，等 JavaScript hydration 后才创建 `<img>`，即使图片文件很小也会晚发现。更合适的做法是将它作为可在初始 HTML 中发现的语义图片，提供正确尺寸与 `srcset/sizes`，并只在确有证据时使用 preload 或优先级提示。对位于首屏外的图片滥用 preload 会抢占关键资源，反而损害 LCP。

```html
<!-- 仅对实际首屏 LCP 候选使用；width/height 用于预留布局空间。 -->
<img
  src="/hero-1280.avif"
  srcset="/hero-640.avif 640w, /hero-1280.avif 1280w"
  sizes="(max-width: 700px) 100vw, 1200px"
  width="1280"
  height="720"
  fetchpriority="high"
  alt="团队项目看板"
>
```

验证时在 Performance 或 RUM 的 LCP 归因中确认最终元素确实是这张图片，而不是标题或广告位；同时检查移动端是否选择了合适候选。资源体积变小但 LCP 未改善，常见原因是发现时间或渲染延迟仍未改变。

### 2.4 INP：从一次交互到下一帧可见反馈

INP（Interaction to Next Paint）衡量一次访问中交互响应性的代表性高延迟值。交互包括点击、触摸和键盘输入；它从输入发生到浏览器绘制下一次视觉反馈的时间。诊断时把一次交互拆成输入延迟、处理时长和呈现延迟：事件在主线程队列中等待、handler/Promise 回调运行太久、或计算结束后样式/布局/绘制迟迟无法完成，都会拉长 INP。

INP 不是“按钮函数执行时间”。下面的 handler 自身可能只运行 5ms，但如果前面有长脚本，输入先排队 200ms；或者 handler 更新大量 DOM，下一帧 layout 和 paint 又花 150ms，用户仍然会感到迟缓。

```ts
searchInput.addEventListener('input', (event) => {
  const value = (event.target as HTMLInputElement).value;
  // 反例：每次输入同步过滤十万项并重建所有行。
  // const rows = allItems.filter((item) => item.name.includes(value));
  // list.replaceChildren(...rows.map(renderRow));

  scheduleSearch(value); // 由后续调度、窗口化和取消逻辑处理
});
```

修复顺序是先 trace 再选方案：缩小被处理的数据；对列表窗口化；把可序列化的重计算移给 Worker；把非关键视觉工作推到后续帧；对旧请求用 `AbortController` 取消。仅加 debounce 会减少执行次数，却可能增加用户等待，并且不解决点击、拖拽或提交按钮的主线程阻塞。所有调度方案都要测试键盘连续输入、清空、请求乱序和页面离开。

### 2.5 CLS：没有预期时发生的布局位移

CLS（Cumulative Layout Shift）累计非预期布局位移的分数。一次位移分数与受影响区域和移动距离有关；分数按会话窗口组合，并不是“页面里元素动过多少次”的简单计数。用户刚点击导致的合理展开、滚动或拖拽不会按同样方式算作意外位移，但依赖“用户刚操作过”的例外来掩盖迟到内容不是解决方案。

最常见来源是图片/广告/iframe 未预留尺寸、晚到的字体导致换行变化、异步通知插入文首、Cookie 横幅压下已有内容、服务端与客户端渲染结构不同。优先使用 `width`/`height` 或 `aspect-ratio` 为媒体和嵌入内容保留空间；把不可预估的动态内容放在不会推动主内容的位置，或明确占位；字体使用合适的 fallback 度量和加载策略。

```css
.video-frame {
  aspect-ratio: 16 / 9;
  background: #eef1f5;
}

.video-frame > iframe {
  width: 100%;
  height: 100%;
  border: 0;
}
```

DevTools 的 Performance 面板会标出 layout shift，RUM 中可记录 `LayoutShift` entry 的 value 和相关元素选择器摘要。不要上传完整 DOM 文本或用户内容。若 CLS 只在一种语言出现，检查翻译长度和字体回退，而不是仅为默认语言压缩字距。

### 2.6 Long Task：主线程超过 50ms 的连续任务

Long Tasks API 将主线程中超过 50ms 的连续任务标为 long task。它是定位信号：任务结束前浏览器难以运行输入处理、定时器和渲染工作，因此长任务会造成交互排队和掉帧。它不是完整性能评分，也无法捕捉多个 45ms 连续任务、GPU 卡顿或网络等待。

在 Performance trace 中展开长任务，区分脚本执行、样式计算、布局、垃圾回收与第三方代码。`setTimeout` 或 `await Promise.resolve()` 不保证给渲染让步，因为 microtask 会在浏览器有机会绘制前继续清空；长计算可分块并在帧/任务边界检查取消，也可转移给 Worker。拆分必须保持状态一致性，不能让旧搜索结果在新输入之后写回页面。

```js
const observer = new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    console.info('long_task', { duration: entry.duration, startTime: entry.startTime });
  }
});

observer.observe({ type: 'longtask', buffered: true });
```

该观察器只能说明浏览器观测到的连续阻塞，不会自动告诉你业务函数名。生产中将它与 route、release、设备类别和采样 ID 关联，开发中再用 source map 和 trace 定位代码；不要向遥测传递函数参数或用户输入。

### 2.7 FPS：连续视觉更新是否跟得上显示节奏

FPS（frames per second）适合动画、滚动、Canvas、拖拽和实时图表。显示器刷新率可能是 60Hz、90Hz 或更高，因此“60 FPS”不是通用目标；对 60Hz 屏幕，单帧约 16.7ms 内完成更容易连续呈现，而高刷新设备预算更紧。静态表单加载后测得 60 FPS 并不说明输入响应好，INP/Long Task 才更合适。

诊断 FPS 应录制真实交互：持续滚动表格、拖动节点、缩放图表。查看 Frames track 是否出现长帧，再回到主线程、合成、光栅或 GPU 轨道寻找原因。避免每个 pointermove 都触发同步 layout；收集最新坐标，在 `requestAnimationFrame` 中每帧最多提交一次视觉更新。对非必要动画尊重 `prefers-reduced-motion`，并让 CPU 较弱设备可关闭高频效果。

### 2.8 JS Heap：内存占用、峰值和泄漏不是同一结论

JS heap 是 JavaScript 可达对象占用的托管内存。一次大 heap 可能是合理缓存、解析大文档或初始化编辑器；泄漏是完成任务、释放引用并触发合理 GC 后，基线仍随重复操作持续上升。浏览器不承诺何时 GC，因此不能用“点一次按钮 heap 没马上下降”判定泄漏。

验证方式是固定场景重复执行，例如打开/关闭编辑器十次、添加/删除十次图层、进入/离开页面十次；在相似稳定点采集 heap snapshot，比较 retaining path。常见根因是全局事件监听器未移除、timer/observer/Worker 未终止、闭包持有大数组、detached DOM 仍被缓存、无限 Map 缓存或 Blob URL 未 `revokeObjectURL()`。

内存指标常需要实验室或受控诊断采集；某些 memory API 的精度和支持范围有限。线上优先观测崩溃、页面冻结、长会话的受控样本与资源错误，避免把对象内容、页面文本或用户数据上传为“内存诊断”。

## 3. 分位数：为什么 p75 常用且不能脱离样本量

p75（75th percentile）表示排序后约 75% 的样本不高于该值，约 25% 更差。它比平均值更不容易被少数极快样本掩盖，同时不像 p95 那样对低样本高度不稳定。Core Web Vitals 的页面体验评估常采用 75th percentile；团队内部阈值仍应按任务、设备、人群和风险定义。

假设同一路由同一 release 的 LCP 样本（毫秒）为：`[1200, 1400, 1600, 1800, 2200, 2400, 2800, 4100]`。p50 位于中间附近，能描述典型访问；p75 接近 2800ms，揭示四分之一用户已进入较慢区间。若把桌面光纤用户与低端移动网络混合，p75 不能告诉你谁受影响；至少按设备类别、路由、release、登录/缓存场景和必要的地区维度切分。

分位数不是小样本魔法。只有十次访问时，单个异常就会显著改变 p75；发布灰度只覆盖 5% 用户时，新旧 release 的样本不可直接比较。报告必须同时写样本数、窗口、覆盖率和缺失/采样规则。出现改善时还要检查是否只是流量来源、CDN 命中或设备构成改变。

## 4. 性能预算：把“应该快”写成可执行约束

预算是发布前或发布后触发行动的阈值，不是单个漂亮分数。预算至少指明对象、环境、统计方法、阈值、例外流程和负责人。资源预算与体验预算分别约束不同风险：前者更适合 CI，后者更适合 RUM 与合成测试。

| 预算类别 | 示例 | 检查位置 | 超标后的动作 |
| --- | --- | --- | --- |
| 初始 JavaScript | 移动首页 gzip 后首次加载 ≤ 170KB | 构建分析/CI | 阻断或附带明确豁免与删除日期 |
| LCP | 移动首页 RUM p75 ≤ 2.5s | 发布后分组仪表板 | 暂停扩量、定位资源链、必要时回滚 |
| INP | 搜索页 RUM p75 ≤ 200ms | RUM 与交互 trace | 检查长任务、渲染和数据规模 |
| CLS | 结算页 RUM p75 ≤ 0.1 | RUM/layout shift 调试 | 修复预留空间或异步插入 |
| 内存 | 编辑器连续操作后 heap 基线不单调上升 | 受控浏览器实验 | 排查 retaining path 和缓存上限 |

阈值不是脱离业务的法律。内部后台可能允许较大首包但要求复杂表格滚动稳定；营销页可能对 LCP 更严格；支付页还要优先保证错误恢复和数据正确。豁免必须可追踪：为什么超标、影响人群、临时开关、修复责任人和截止日期。没有截止日期的“临时”例外会成为永久退化。

## 5. 案例一：用 LCP 资源链修复移动端首页回归

### 现象与证据

发布 `web@2026.07.23.1` 后，移动端首页 RUM p75 LCP 从 2.2s 升到 3.4s，桌面端无明显变化。先按 `route=/`、`device=mobile`、release、网络类别切分，确认新 release 有足够样本且异常不只集中在一个地区。LCP 元素归因显示 86% 是首屏主视觉图片。

Performance trace 显示：HTML 在 420ms 收到首字节；分析脚本和应用 bundle 先下载；应用 hydration 后才由组件插入主视觉 `<img>`；图片请求到 1.7s 才开始。此时压缩图片本身只能减少后半段下载，根因是资源发现迟。

### 变更

将主视觉改为服务端/初始 HTML 可发现的 `<img>`，给出尺寸和 responsive 候选；把非首屏分析 SDK 推迟到同意和空闲边界。没有给所有图片加 preload，仅对最终 LCP 元素在实验分支验证优先级。保留 feature flag，使新图片路径可按 release 回退。

### 验证

1. 实验室用相同移动网络/CPU 预设录制 before/after，确认图片请求早于非关键脚本。
2. 检查 `srcset` 在窄屏选择较小资源，并确认图片尺寸不会造成 CLS。
3. 在 canary 中看移动首页新 release 的 LCP p75、TTFB、请求失败率和分析事件量；同时比较桌面，防止资源竞争伤害其他人群。
4. 若 LCP 没有改善，检查最终候选是否已经变成标题或广告，而不是继续压缩原图片。

### 失败分支

若 preload 让字体或关键 CSS 更晚，FCP/INP 可能变坏；应撤回该提示并优先减少竞争。若 CDN miss 使 TTFB 上升，图片发现优化不能掩盖服务端回归，应分别追踪文档与图片请求的 cache status。结论是“在移动首页、新 release、该候选元素下资源发现提前改善 LCP”，不是“preload 永远让网站更快”。

## 6. 案例二：从 420ms INP 找到搜索筛选的主线程问题

### 现象与证据

RUM 报告搜索页 p75 INP 为 420ms，发生最多的是键盘输入。Performance trace 显示输入事件前有约 110ms 脚本队列等待；handler 对 80,000 条记录执行同步过滤和排序，耗时 180ms；随后创建数千个行节点和 layout，下一次绘制又耗时约 130ms。这里的主要问题不是网络，也不是单纯 debounce。

### 变更

索引预处理在数据加载后移到 Worker；输入每次只发送查询词和版本号。主线程只维护当前请求序号，并将结果交给虚拟列表渲染可见行。输入新字符时通过 `AbortController` 或 worker 消息协议废弃旧任务；视觉更新在帧边界合并，避免一次输入触发多次 layout。

```ts
let latestRequest = 0;

function requestSearch(query: string) {
  const request = ++latestRequest;
  worker.postMessage({ type: 'search', request, query });
}

worker.addEventListener('message', (event) => {
  const { request, rows } = event.data as { request: number; rows: unknown[] };
  if (request !== latestRequest) return; // 旧输入的结果不得覆盖新输入
  renderVisibleRows(rows);
});
```

### 验证与失败分支

测试连续输入、退格、清空、切换筛选、Worker 错误和离开页面。以真实设备的 trace 确认 handler、渲染和 long task 均下降，再在 RUM 观察按设备切分后的 INP p75。若 Worker 启动和数据复制成本使小数据集更慢，可以设置数据量阈值或让轻量路径留在主线程；不要仅因 Worker 是“高级方案”而强制使用。

若虚拟化后键盘焦点滚出可见区域，说明性能优化破坏了可操作性。修复应包括焦点保留、`aria-rowcount`/行语义的合适实现、加载状态和可恢复错误，而不是删除键盘支持来换取更小 DOM。

## 7. RUM 事件设计、采样与隐私

RUM 数据应有稳定 schema。最小事件可包含指标名、值、rating、页面路径模板、release、设备类别、连接粗分类、匿名会话短 ID、采集时间和采样标记。路径应去掉用户 ID、订单号、搜索词和敏感 query；错误、请求和业务事件使用受控枚举与服务端关联 ID，而不是把任意对象 `JSON.stringify` 后上报。

```ts
type VitalEvent = {
  name: 'LCP' | 'INP' | 'CLS' | 'FCP' | 'TTFB';
  value: number;
  rating: 'good' | 'needs-improvement' | 'poor';
  route: '/home' | '/search' | '/checkout';
  release: string;
  deviceClass: 'mobile' | 'desktop' | 'other';
};

export function sendVital(event: VitalEvent) {
  const payload = JSON.stringify(event);
  if (navigator.sendBeacon) return navigator.sendBeacon('/rum/vitals', payload);
  return fetch('/rum/vitals', { method: 'POST', body: payload, keepalive: true }).then(() => undefined);
}
```

传输函数并不保证数据合规。服务端还需验证 schema、限制大小与速率、拒绝未知字段、执行访问控制、定义保留期限和删除流程。用户同意、地区法规、儿童数据、企业数据隔离等要求由产品/法务/安全规则决定，不能由前端偷偷采集来规避。

采样按用途区分：核心指标可做有代表性的概率采样；错误和严重回归可提高采样但仍限流；调试 session 需要明确授权和更短保留。采样率、SDK 版本和 release 也要写入数据，否则一次采集策略变化会伪造趋势变化。

## 8. 调试顺序与反例

当用户报告“慢”时，先确定是加载、交互、位移、动画还是内存问题，再选择指标。以下反例常导致错误优化：

| 观察到的现象 | 不可靠结论 | 更可靠的下一步 |
| --- | --- | --- |
| Lighthouse 变高 | 线上所有用户变快 | 查看同 release、同路由 RUM 分位数与样本 |
| 平均 LCP 下降 | 尾部体验已修复 | 看 p75/p95、设备和网络分组 |
| handler 很短 | INP 不会差 | 检查输入排队与下一帧渲染 |
| 没有 single Long Task | 主线程没有问题 | 检查连续短任务、layout、GC 和第三方脚本 |
| JS heap 较高 | 一定内存泄漏 | 重复固定场景、比较稳定点 retaining path |
| 给所有图片 preload | LCP 必然下降 | 验证最终 LCP 候选与资源竞争 |
| 采集更多字段 | 更容易定位 | 先确认字段用途、隐私、成本和保留期限 |

可重复调试记录应写明：用户任务、路由、release、设备/网络、样本数、指标与分位数、trace/瀑布图、假设、唯一变更、失败分支、回滚开关和发布后观察窗口。一次优化若只提升桌面缓存命中而伤害移动首访，记录必须保留这个反例。

## 9. 综合练习：建立一个可回滚的性能门禁

选择一个公开路由和一个高频交互，完成以下工作：

1. 为路由定义移动和桌面的 FCP/LCP/CLS 预算，为交互定义 INP 预算，为构建定义资源预算。
2. 在实验室固定设备、网络、登录态和缓存条件，保存 before trace 与构建分析。
3. 通过 RUM 发送最小指标 schema，按 release、route、deviceClass 和样本数展示 p50/p75/p95。
4. 制造一个 LCP 资源发现回归和一个输入主线程回归，分别用 trace 找到根因。
5. 用 feature flag 或可回滚制品发布修复，制定“何种 release 覆盖率、观察多久、何种阈值回滚”的规则。
6. 审查事件字段：不得包含 token、表单内容、完整 query、用户生成内容或不必要的稳定身份。

验收：能说明每个预算对应哪个用户任务；能区分 TTFB、FCP、LCP、INP、CLS 的责任边界；出现一次线上回归时能定位到 route、release、设备人群和资源/主线程证据；关闭采集或撤回 release 后系统仍保持核心功能可用。

## 来源

- [web.dev：Web Vitals](https://web.dev/articles/vitals)（访问日期：2026-07-23）
- [web.dev：Optimize LCP](https://web.dev/articles/optimize-lcp)（访问日期：2026-07-23）
- [web.dev：Optimize INP](https://web.dev/articles/optimize-inp)（访问日期：2026-07-23）
- [MDN：Performance API](https://developer.mozilla.org/en-US/docs/Web/API/Performance_API)（访问日期：2026-07-23）
- [W3C：Long Tasks API](https://w3c.github.io/longtasks/)（访问日期：2026-07-23）
