# DOM 树、查询、创建、更新、删除与样式

## 是什么

DOM 是浏览器把文档表示为节点树的 API；查询取得节点，创建和更新改变树，class/style 改变呈现。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

优先 textContent 放不可信文本；批量变更可用 DocumentFragment；保留稳定节点引用；样式状态优先切 class。

## 实际使用

```js
const list=document.querySelector('#items');
const li=document.createElement('li');
li.textContent=user.name;
li.classList.add('item');
list.append(li);
// li.remove();
```

## 常见错误与边界

innerHTML 接受未净化输入会 XSS；重复查询和交错读写可能低效；DOM 节点移动不是复制。

## 相关补充知识

批量 DOM 读写要避免在循环中交替触发布局计算。用户文本用 `textContent`，不能直接拼入 `innerHTML`；动态移除节点时同步清理观察器、事件和外部资源。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Introduction)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelector)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Document/createElement)

访问日期：2026-07-16。
