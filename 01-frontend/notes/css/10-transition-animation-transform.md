# Transition、Animation、Transform 与关键帧

## 是什么

transition 在属性值变化间插值；animation 通过 @keyframes 驱动时间序列；transform 改变元素坐标空间而通常不重排周围布局。

## 为什么需要

用于表达状态变化和空间关系，并避免用高频 JavaScript 手工逐帧改样式。

## 关键特性与规则

只 transition 明确属性；动画定义 duration、timing、iteration、fill；transform 函数顺序影响结果；关键交互不能仅靠动画传达。

## 实际使用

```css
.button { transition:background-color .15s ease, transform .15s ease; }
.button:active { transform:scale(.98); }
@keyframes spin { to { transform:rotate(1turn); } }
```

## 常见错误与边界

transition:all 会意外动画尺寸等属性；display 通常不能传统插值；transform 会创建 containing block/stacking context；无限动画消耗资源。

## 相关补充知识

Web Animations API 可在 JS 中控制浏览器动画时间线。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_transitions/Using_CSS_transitions)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_animations/Using_CSS_animations)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/transform)

访问日期：2026-07-16。

