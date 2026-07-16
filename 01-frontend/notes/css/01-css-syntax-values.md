# 选择器、声明、单位、颜色、背景、边框与字体

## 是什么

CSS 规则由选择器和声明块组成；声明是属性与合法值的组合。单位决定长度参照，颜色控制前景，背景在边框以内绘制，边框参与盒尺寸，字体属性控制排版。

## 为什么需要

将结构与呈现分离，并让同一规则复用于多个元素。理解值类型可避免声明被浏览器判为无效。

## 关键特性与规则

优先 class 选择器组织组件；相对字体用 rem/em，视口尺寸慎用纯 vh；背景图只承载装饰内容；字体列表提供通用族回退。

## 实际使用

```css
h1 { color: rgb(20 40 80); font: 700 2rem/1.2 system-ui; }
.card { padding: 1rem; border: 1px solid #ccd; background: white url(bg.svg) no-repeat right top / 4rem; }
```

## 常见错误与边界

无效声明会被忽略；简写会重置未写出的长属性；百分比参照因属性而异；CSS px 不是固定物理像素。

## 相关补充知识

CSS 支持 calc/min/max/clamp 与现代颜色函数；上线前查属性兼容性。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Getting_started)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Values_and_units)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Styling_basics/Backgrounds_and_borders)

访问日期：2026-07-16。

