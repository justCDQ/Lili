# 深色模式、RTL、高对比度与打印样式

## 是什么

这些是不同用户环境：color-scheme/prefers-color-scheme 处理主题，direction 与逻辑属性处理 RTL，forced-colors/prefers-contrast 处理高对比，print 媒体处理纸张输出。

## 为什么需要

只按默认浅色 LTR 屏幕设计会导致内容不可读或布局失效。

## 关键特性与规则

优先系统偏好并允许用户覆盖；RTL 使用真实 dir=rtl 测试；forced colors 下依赖系统颜色；打印移除交互控件并控制分页。

## 实际使用

```css
:root { color-scheme:light dark; }
@media (forced-colors:active) { .icon { forced-color-adjust:auto; } }
@media print { nav { display:none; } a::after { content:' (' attr(href) ')'; } }
```

## 常见错误与边界

filter 反色不是完整深色主题；不要用 CSS direction 修正错误语言标记；高对比模式不要强制保留品牌色而牺牲可读性。

## 相关补充知识

测试 Windows forced colors、浏览器打印预览及混合双向文本。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-color-scheme)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@media/forced-colors)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_media_queries/Printing)
- [W3C](https://www.w3.org/International/questions/qa-html-dir)

访问日期：2026-07-16。

