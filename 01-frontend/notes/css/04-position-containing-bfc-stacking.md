# Position、Containing Block、BFC 与 Stacking Context

## 是什么

position 决定定位方案；绝对偏移相对 containing block。BFC 隔离块布局、浮动和 margin 行为。stacking context 把后代作为整体参与父级堆叠。

## 为什么需要

定位和 z-index 异常通常来自参照块或新堆叠上下文，而非数值不够大。

## 关键特性与规则

absolute 脱离正常流；fixed 通常相对视口但 transform 祖先可改变参照；sticky 受最近滚动容器和 inset 限制；多种属性会创建 stacking context。

## 实际使用

```css
.shell { position: relative; }
.badge { position: absolute; inset-block-start: 0; inset-inline-end: 0; z-index: 1; }
.flow-root { display: flow-root; }
```

## 常见错误与边界

z-index 不能跨越祖先 stacking context；sticky 没有可滚动空间不会移动；用 overflow 创建 BFC 可能意外裁剪。

## 相关补充知识

用 DevTools 检查 containing block、scroll container 和 stacking context。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/position)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Display/Containing_block)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Display/Block_formatting_context)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Positioned_layout/Stacking_context)

访问日期：2026-07-16。

