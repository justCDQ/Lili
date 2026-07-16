# prefers-reduced-motion 与动效性能

## 是什么

prefers-reduced-motion 反映用户减少非必要运动的系统偏好。动效性能取决于每帧工作量；transform/opacity 常可避免布局和绘制，但合成仍有成本。

## 为什么需要

大幅移动、缩放和视差可能造成不适；低效动画会阻塞交互和掉帧。

## 关键特性与规则

为重要状态提供无运动替代而非简单删除信息；用 Performance 面板验证 style/layout/paint/composite；避免动画几何属性和大面积滤镜。

## 实际使用

```css
@media (prefers-reduced-motion:reduce) { *,*::before,*::after { scroll-behavior:auto; animation-duration:.01ms; animation-iteration-count:1; } }
```

## 常见错误与边界

will-change 不是通用加速开关，会增加内存和图层；短 duration 仍可能不适；CSS 动画也可阻塞主线程相关工作。

## 相关补充知识

优先让用户设置生效，并提供应用内持久偏好。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion)
- [web.dev](https://web.dev/articles/animations-guide)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/Performance/Guides/CSS_JavaScript_animation_performance)

访问日期：2026-07-16。

