# Cascade、Specificity、Inheritance 与默认样式

## 是什么

层叠按来源、重要性、层、优先级、作用域接近度和源码顺序决定胜出声明；继承把部分属性的计算值传给子元素；浏览器默认样式来自用户代理样式表。

## 为什么需要

当多个规则命中同一元素时，必须能解释最终值来源，才能稳定覆盖而不是堆叠 !important。

## 关键特性与规则

ID 比 class 优先级高，class/属性/伪类高于类型；:where() 优先级为零；继承属性可用 inherit，重置可用 initial/unset/revert。

## 实际使用

```css
@layer reset, base, components;
@layer base { body { color: #222; } }
@layer components { .button.primary { color: white; } }
```

## 常见错误与边界

!important 会改变层叠优先顺序且难覆盖；源码靠后只在前序条件相同才胜出；继承与子选择器不是一回事。

## 相关补充知识

DevTools Computed 面板可查看被覆盖声明和继承链。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Handling_conflicts)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Cascade/Specificity)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/inheritance)

访问日期：2026-07-16。

