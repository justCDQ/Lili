# 标题层级、替代文本、焦点顺序与键盘操作

## 是什么与为什么需要

标题建立内容层级；替代文本提供非视觉等价信息；焦点是当前接收键盘输入的元素；焦点顺序决定顺序导航路径。它们使屏幕阅读器和键盘用户能理解、定位并操作页面。

## 关键规则

## 实际使用

- 页面以描述性 `h1` 开始，章节按实际层级用 `h2`、`h3`；不因字号跳级。
- 信息图片的 `alt` 表达用途/信息，装饰图 `alt=""`；复杂图在正文提供数据或说明。
- DOM 顺序应符合阅读与操作顺序；使用自然可聚焦的 `a[href]`、`button`、表单控件。
- 避免正 `tabindex`；需程序聚焦的容器可用 `tabindex="-1"`。所有鼠标功能必须可由键盘完成并有可见焦点。

```html
<a href="#main" class="skip-link">跳到主要内容</a>
<main id="main" tabindex="-1"><h1>账户设置</h1></main>
```

## 常见错误与边界

CSS 重排视觉顺序可能与 DOM/焦点顺序分离。不要移除 outline 而不提供等价焦点样式。不可给每个静态元素 `tabindex="0"`。图中文字若正文已完整重复可避免冗余 alt。

## 补充知识

键盘测试至少覆盖 Tab、Shift+Tab、Enter、Space、方向键和 Escape，具体键取决于原生控件或规范化组件模式。

## 来源

- [W3C WAI：Accessibility principles](https://www.w3.org/WAI/fundamentals/accessibility-principles/)
- [W3C：Understanding Focus Order](https://www.w3.org/WAI/WCAG22/Understanding/focus-order.html)
- [W3C WAI：Images tutorial](https://www.w3.org/WAI/tutorials/images/)

访问日期：2026-07-16。
