# Dialog、Popover 与 Details 原生交互元素

## 是什么与为什么需要

`dialog` 表示模态或非模态对话框；Popover API 把元素置于 top layer 并提供声明式显示/关闭；`details/summary` 提供可展开内容。原生能力能减少自制组件所需的焦点、Escape、层级与状态代码。

## 关键特性与规则

- `showModal()` 创建真正模态的 dialog；只设置 `open` 不会阻止背景交互。
- 模态打开时焦点进入对话框，关闭后应回到触发点或合理后续位置。
- popover 始终是非模态；`auto` 支持轻关闭，`manual` 由应用显式管理。
- `summary` 是 details 的可操作标签，文案应说明将展开的内容。
- 原生元素仍需测试标题、关闭路径、焦点顺序、窄屏和辅助技术支持。

## 实际使用

```html
<button commandfor="confirm" command="show-modal">删除</button>
<dialog id="confirm"><p>确认删除？</p><form method="dialog"><button value="cancel">取消</button></form></dialog>
<button popovertarget="actions">操作</button><div id="actions" popover>...</div>
<details><summary>系统要求</summary><p>Node.js 22+</p></details>
```

模态对话框用 `showModal()`，关闭用 `close()`/表单 `method="dialog"`；必须有清楚标题和关闭途径。popover 默认 `auto` 可点击外部或 Escape 轻关闭，`manual` 需显式关闭。`summary` 应描述展开内容。

## 常见错误与边界

Popover 永远非模态，需要阻止背景交互时使用 dialog。不要给 dialog 自身设置 `tabindex`。只写 dialog 的 `open` 属性不会获得 `showModal()` 的模态行为。新属性兼容性变化快，上线前查支持并测试键盘/屏幕阅读器。

## 补充知识

这些组件进入浏览器 top layer 后不受普通 `z-index` 层级限制；可使用 `::backdrop` 设置背景层样式。

## 来源

- [MDN：dialog](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog)
- [MDN：Popover API](https://developer.mozilla.org/en-US/docs/Web/API/Popover_API)
- [MDN：details](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details)

访问日期：2026-07-16。
