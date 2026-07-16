# 语义页面区域：header、nav、main、article、section、aside、footer

## 是什么与为什么需要

这些元素表达页面区域与内容关系，并可映射为辅助技术地标。`header/footer` 是所属页面或章节的头尾；`nav` 是主要导航；`main` 是页面独有主要内容；`article` 可独立分发；`section` 是有主题的章节；`aside` 与周边内容间接相关。

## 关键特性与规则

- 元素按内容关系选择，不能根据默认样式或视觉位置选择。
- 页面通常只有一个可见 `main`；重复类型地标应提供可区分名称。
- `section` 表示有主题的章节且通常有标题，纯布局容器使用 `div`。
- `article` 应能在离开当前页面上下文后独立理解或分发。
- DOM 源顺序应先满足阅读和键盘顺序，再使用 CSS 改变视觉布局。

## 实际使用

```html
<header><h1>工程周刊</h1><nav aria-label="主导航">...</nav></header>
<main>
  <article><h2>HTTP 缓存</h2><section><h3>验证缓存</h3>...</section></article>
  <aside><h2>相关阅读</h2>...</aside>
</main>
<footer>版权与联系信息</footer>
```

页面通常只有一个可见 `main`。多个 `nav` 用可访问名称区分。`section` 通常应有标题；无章节语义的纯样式容器使用 `div`。`article` 内可拥有自己的 header/footer。

## 常见错误与边界

不要把所有容器改成 `section`。`header` 不等于只能出现一次；`footer` 不保证固定在视口底部。视觉两栏不自动意味着 `main + aside`，需看内容关系。源顺序应先有逻辑，再用 CSS 布局。

## 补充知识

原生元素已有隐式角色时通常无需重复 `role`。地标过多会增加导航负担，命名应简洁且相同地标类型可区分。

## 来源

- [MDN：HTML accessibility](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Accessibility/HTML)
- [MDN：HTML elements reference](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements)
- [W3C APG：Landmark regions](https://www.w3.org/WAI/ARIA/apg/practices/landmark-regions/)

访问日期：2026-07-16。
