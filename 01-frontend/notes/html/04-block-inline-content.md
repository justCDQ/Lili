# 块级与行内内容

## 是什么与为什么需要

“块级/行内元素”是历史简化说法。现代规范把 HTML 语义分类与 CSS 布局分开：CSS `display` 的外部显示类型决定盒子在父级中以 block 或 inline 参与布局，内部显示类型决定子内容布局。理解它能解释换行、尺寸和正常流。

## 关键特性

## 实际使用

默认情况下 `p`、`div` 常生成块盒并占据可用行向空间；`a`、`span` 常生成行内盒并随文本断行。CSS 可改变：

```css
a.card { display: block; }
.badge { display: inline-block; width: 6rem; }
.actions { display: inline-flex; gap: .5rem; }
```

`inline` 盒的 `width/height` 通常不按块盒方式生效；`inline-block` 在外部随行排布，内部形成独立盒。改变 display 不改变 HTML 语义，`span {display:block}` 仍不是段落或区域。

## 常见错误与边界

不要依据“是否换行”选择 HTML 元素。块元素可包含什么由 HTML 内容模型决定，不由 CSS display 决定。行内格式化中的空白、基线和可断行位置会影响间隙；图片底部空隙常来自基线对齐。

## 补充知识

书写模式会改变 block/inline 物理方向，因此现代 CSS 使用 `inline-size`、`block-size` 等逻辑属性。`display: contents` 会移除自身盒但可能影响可访问性实现，使用前测试。

## 来源

- [MDN：Inline-level content](https://developer.mozilla.org/en-US/docs/Glossary/Inline-level_content)
- [MDN：Block and inline layout](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Display/Block_and_inline_layout)
- [MDN：display](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/display)

访问日期：2026-07-16。
