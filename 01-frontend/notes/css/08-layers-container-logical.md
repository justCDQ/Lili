# Cascade Layers、Container Queries 与 Logical Properties

## 是什么

层叠层显式控制规则组优先顺序；容器查询根据祖先容器条件应用样式；逻辑属性用 inline/block 方向替代物理上下左右。

## 为什么需要

它们分别解决可预测覆盖、可复用组件响应和多书写模式适配。

## 关键特性与规则

未分层普通规则优先于分层普通规则；容器需声明 containment 类型；查询不能用被查询元素自身作为容器；逻辑属性随 writing-mode/direction 映射。

## 实际使用

```css
@layer reset, base, components, overrides;
.card-wrap { container-type:inline-size; }
@container (width > 30rem) { .card { grid-template-columns:8rem 1fr; } }
.card { padding-inline:1rem; margin-block:.5rem; }
```

## 常见错误与边界

混用物理与逻辑属性可能互相覆盖；容器没有可查询尺寸时规则不生效；层顺序应集中声明。

## 相关补充知识

样式查询支持范围与兼容性需逐项查证。

## 来源

- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@layer)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Containment/Container_queries)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_logical_properties_and_values)

访问日期：2026-07-16。

