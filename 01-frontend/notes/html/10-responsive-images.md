# 响应式图片、picture、srcset 与懒加载

## 是什么与为什么需要

`srcset` 提供候选资源，`sizes`描述图片槽位宽度，浏览器结合视口与像素密度选择候选；`picture/source` 用于艺术方向或格式切换；`loading="lazy"` 延迟非关键图加载。目标是降低传输、适配密度并控制裁切。

## 关键特性与规则

- `w` 描述符声明资源固有宽度并与 `sizes` 配合；`x` 描述符用于固定槽位的像素密度候选。
- `picture` 按 source 顺序匹配媒体条件和类型，最终 `img` 提供语义与回退。
- 浏览器根据当前环境选择候选，开发者不能假设某个候选一定被下载。
- 图片显式 `width` 和 `height` 可建立宽高比，减少加载期间布局偏移。
- 首屏关键图通常不应懒加载，非关键图再使用 `loading="lazy"`。

## 实际使用

```html
<picture>
 <source media="(max-width: 40rem)" srcset="hero-crop.avif" type="image/avif">
 <source srcset="hero.avif" type="image/avif">
 <img src="hero.jpg" srcset="hero-640.jpg 640w, hero-1280.jpg 1280w" sizes="(max-width: 40rem) 100vw, 60vw" width="1280" height="720" alt="团队在评审界面">
</picture>
<img src="related.jpg" loading="lazy" width="640" height="360" alt="相关项目截图">
```

`w` 描述符必须是资源固有宽度，搭配 `sizes`；固定 CSS 尺寸且仅适配密度时用 `1x/2x`。picture 内最终必须有 img 作为语义载体与回退。首屏关键图通常不懒加载；始终提供 width/height 减少布局偏移。

## 常见错误与边界

浏览器选择是提示驱动，不保证每次换宽度重新下载。不要在同一候选混用 `w` 与 `x`。CSS/JS 替换可能在预加载扫描后造成重复下载。`alt` 描述内容，与响应式资源无关。

## 补充知识

可在控制台读取 `img.currentSrc`，并在 Network 禁用缓存后核对实际候选；现代格式切换应保留浏览器可用的回退资源。

## 来源

- [MDN：Responsive images](https://developer.mozilla.org/en-US/docs/Web/HTML/Guides/Responsive_images)
- [MDN：img](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img)
- [MDN：picture](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture)

访问日期：2026-07-16。
