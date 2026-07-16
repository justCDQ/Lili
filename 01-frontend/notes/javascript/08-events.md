# Event、冒泡、捕获、委托、默认行为与自定义事件

## 是什么

事件沿捕获阶段到目标再冒泡；委托在共同祖先处理后代事件；默认行为由浏览器在未取消时执行；CustomEvent 可携带应用数据。

## 为什么需要

这些能力用于建立可预测的程序状态、控制流和浏览器交互，也是框架与工程工具的运行基础。

## 关键特性与规则

监听器默认冒泡阶段；preventDefault 只取消可取消默认行为；stopPropagation 阻止传播但不等于取消默认；委托需验证 closest 仍在容器内。

## 实际使用

```js
list.addEventListener('click',e=>{
 const button=e.target.closest('button[data-id]');
 if(!button||!list.contains(button)) return;
 remove(button.dataset.id);
});
```

## 常见错误与边界

用箭头函数临时注册后无法按同一引用移除；被动监听器不能 preventDefault；自定义事件不应假装可信用户事件。

## 相关补充知识

事件传播包含捕获、目标和冒泡阶段；委托依赖冒泡并应使用 `closest` 后验证容器归属。`preventDefault` 只取消默认行为，`stopPropagation` 会影响其他监听器，不能作为常规控制流。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Scripting/Event_bubbling)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Event)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/CustomEvent)

访问日期：2026-07-16。
