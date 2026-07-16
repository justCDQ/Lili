# Box Model、Normal Flow、display、溢出与滚动

## 是什么

CSS 盒由 content、padding、border、margin 构成。正常流按块/行内格式化排列。display 定义外部和内部显示类型；overflow 处理内容超出盒的行为。

## 为什么需要

尺寸、间距、换行和滚动问题都依赖盒模型与格式化上下文。

## 关键特性与规则

content-box 的声明宽度不含 padding/border；border-box 包含它们；块级竖直 margin 可能折叠；overflow:auto 只在需要时产生滚动机制。

## 实际使用

```css
*, *::before, *::after { box-sizing: border-box; }
.panel { inline-size: 20rem; padding: 1rem; overflow: auto; }
```

## 常见错误与边界

固定高度易造成文本溢出；overflow:hidden 会裁剪内容并可能隐藏焦点；display:none 从布局和可访问树移除。

## 相关补充知识

逻辑尺寸 inline-size/block-size 适配不同书写模式。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Box_model)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/CSS_layout/Introduction)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/overflow)

访问日期：2026-07-16。

