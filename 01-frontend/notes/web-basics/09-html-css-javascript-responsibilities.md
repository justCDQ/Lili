# HTML、CSS 与 JavaScript 的职责

## 是什么与为什么需要

HTML 表达内容结构和语义，CSS 控制呈现与布局，JavaScript 实现计算、状态变化和行为。明确职责能保留无脚本内容、复用样式，并让辅助技术理解原生语义。

## 关键规则

- HTML 用适合含义的元素，不按默认外观选标签。
- CSS 通过选择器匹配元素，层叠后形成计算样式；不应承载关键文本内容。
- JavaScript 操作 DOM 和 Web API；脚本失败时基础信息与导航应尽量可用。
- 浏览器解析 HTML 建 DOM，解析 CSS 建样式规则，执行脚本后计算布局与绘制；三者会相互影响但职责不相同。

## 实际使用

```html
<button class="save" type="button">保存</button>
<link rel="stylesheet" href="app.css">
<script src="app.js" defer></script>
```

```css
.save { color: white; background: #1463ff; }
```

```js
document.querySelector('.save').addEventListener('click', saveDraft);
```

## 常见错误与边界

不要用 `div` 模拟按钮；会丢失键盘和表单行为。不要用 JS 完成纯样式状态。CSS 也能提供有限交互，HTML 也有原生交互元素，职责边界以语义、可维护性和降级为准，而不是绝对禁止重叠。

## 补充知识

渐进增强先建立可工作的 HTML，再添加呈现和增强行为。服务端也可生成 HTML、CSS 或 JS；“生成位置”不改变文件在浏览器中的职责。

## 来源

- [MDN：The web standards model](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Web_standards/The_web_standards_model)
- [MDN：How browsers load websites](https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Web_standards/How_browsers_load_websites)

访问日期：2026-07-16。
