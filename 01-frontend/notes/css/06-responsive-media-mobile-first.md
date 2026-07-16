# 响应式设计、媒体查询、移动优先与断点

## 是什么

响应式设计让布局适应可用空间、输入和用户偏好。媒体查询按环境条件应用规则；移动优先从窄屏基础样式开始，再用 min-width 增强。

## 为什么需要

设备尺寸与窗口变化不可枚举，应按内容失效点设计流式布局。

## 关键特性与规则

设置 viewport；用相对单位、min/max/clamp、可换行布局；断点根据内容何时拥挤确定，不按具体机型；同时测试缩放、横竖屏和大字体。

## 实际使用

```css
.page { padding:1rem; }
@media (width >= 48rem) { .page { display:grid; grid-template-columns:16rem 1fr; } }
```

## 常见错误与边界

隐藏内容不等于响应式；固定像素宽度易横向滚动；hover 查询不能替代键盘支持；媒体查询也可检测 reduced-motion、contrast、color-scheme。

## 相关补充知识

容器查询适合组件根据容器而非视口响应。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/CSS_layout/Responsive_Design)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Media_queries)
- [web.dev](https://web.dev/articles/responsive-web-design-basics)

访问日期：2026-07-16。

