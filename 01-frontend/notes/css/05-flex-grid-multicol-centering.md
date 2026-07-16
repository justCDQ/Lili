# Flexbox、Grid、多列布局与居中

## 是什么

Flexbox 是一维布局，Grid 是二维轨道布局，多列布局让内容流入多栏。居中方法由轴、已知尺寸和布局模式决定。

## 为什么需要

选择匹配问题维度的布局系统，可减少定位和负 margin。

## 关键特性与规则

flex 主轴由 flex-direction 决定；flex-basis 参与空间分配；Grid 项由线和区域定位；gap 优先于子项 margin；多栏适合连续文本而非精确组件网格。

## 实际使用

```css
.toolbar { display:flex; align-items:center; justify-content:space-between; gap:.75rem; }
.gallery { display:grid; grid-template-columns:repeat(auto-fit,minmax(14rem,1fr)); gap:1rem; }
```

## 常见错误与边界

flex 子项默认 min-width:auto 可能不收缩；Grid 1fr 仍受最小内容限制，可用 minmax(0,1fr)；视觉重排不要破坏 DOM 顺序。

## 相关补充知识

单项水平垂直居中可用 display:grid; place-items:center。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/CSS_layout/Flexbox)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/CSS_layout/Grids)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Multicol_layout)

访问日期：2026-07-16。

