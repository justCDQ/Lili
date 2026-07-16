# ARIA 的使用边界：优先正确的原生元素

## 是什么与为什么需要

ARIA 通过 role、state、property 补充可访问性语义，主要用于 HTML 没有等价原生语义的复杂组件。它不会自动添加键盘行为、焦点管理、样式、校验或业务逻辑。

## 关键规则

## 实际使用

优先 `<button>` 而非 `<div role="button">`，优先 `<nav>` 而非 `<div role="navigation">`。原生语义不足时再补充：

```html
<button aria-expanded="false" aria-controls="filters">筛选条件</button>
<section id="filters" hidden>...</section>
<nav aria-label="页脚导航">...</nav>
```

ARIA 不得覆盖成与原生行为冲突的角色。状态必须与实际 UI 同步。每个交互组件需有可访问名称、正确键盘模式、焦点处理，并在真实浏览器与辅助技术组合中测试。

## 常见错误与边界

`role="button"` 不会响应 Space/Enter。`aria-hidden="true"` 不应放在仍可聚焦元素或其祖先上。名称、描述和内容不是一回事，不要滥用 `aria-label` 覆盖有用可见文本。“无 ARIA”优于错误 ARIA，但缺失必要语义也需修复。

## 补充知识

APG 是实现模式指南而非规范或 UI 设计系统；规范符合不代表所有辅助技术实现一致。自动化检查只能覆盖部分问题，必须键盘和屏幕阅读器人工验证。

## 来源

- [W3C APG：Read Me First](https://www.w3.org/WAI/ARIA/apg/practices/read-me-first/)
- [W3C APG：Structural roles](https://www.w3.org/WAI/ARIA/apg/practices/structural-roles/)
- [W3C：ARIA in HTML](https://www.w3.org/TR/html-aria/)

访问日期：2026-07-16。
