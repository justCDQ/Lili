# Custom Properties、Design Token 与主题切换

## 是什么

自定义属性是参与层叠和继承的 CSS 值容器；Design Token 是跨设计与代码命名的决策值；主题通过在作用域重定义 token 改变呈现。

## 为什么需要

集中语义值可减少重复并支持主题、品牌和状态一致性。

## 关键特性与规则

token 名称优先表达语义而非具体颜色；var() 可提供回退；自定义属性区分大小写且通常继承；主题状态放在根或局部容器。

## 实际使用

```css
:root { --color-surface:#fff; --color-text:#17202a; --space-2:.5rem; }
[data-theme=dark] { --color-surface:#111; --color-text:#eee; }
.card { color:var(--color-text); background:var(--color-surface); }
```

## 常见错误与边界

自定义属性保存的是 token 流，不自动类型检查；循环引用使值无效；仅换颜色仍需验证对比度。

## 相关补充知识

@property 可声明语法、初始值和继承行为。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Cascading_variables/Using_custom_properties)
- [W3C](https://www.w3.org/community/design-tokens/)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@property)

访问日期：2026-07-16。

