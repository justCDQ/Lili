# BEM、CSS Modules、CSS-in-JS 与 Utility CSS

## 是什么

- **BEM** 用 `block__element--modifier` 命名表达组件、组成部分与变体，类名仍在全局作用域。
- **CSS Modules** 在构建时把局部类名映射为唯一名称，通过模块导入把样式与组件关联。
- **CSS-in-JS** 在 JavaScript/TypeScript 中声明样式；实现可能在运行时注入，也可能在构建时提取静态 CSS。
- **Utility CSS** 提供单一职责的工具类，由标记组合出布局、间距、颜色和状态。

## 为什么需要

大型样式表需要控制选择器作用域、组件变体、删除安全、样式复用和团队约定。不同方法把复杂度放在不同位置：BEM 依赖命名纪律，Modules 依赖构建，CSS-in-JS 可能增加运行时和水合成本，Utility CSS 增加标记密度并依赖受控配置。

## 关键特性与规则

- 方法论不能替代 Cascade、Specificity 和继承规则；即使类名局部化，元素样式、全局变量和层叠顺序仍会影响结果。
- 先定义组件边界、状态、响应式和主题需求，再选择作用域方案。
- 动态值优先使用 CSS 自定义属性、属性或有限变体，避免为任意值生成大量规则。
- Design Token 应独立于 BEM、Modules 或 CSS-in-JS，确保主题与组件实现解耦。
- 服务端渲染需要验证样式提取、加载顺序、水合一致性和无 JavaScript 首屏。
- 团队应统一目录、命名、变体、全局入口和废弃样式策略，避免同一项目并存多套无边界方案。

## 实际使用

```css
/* BEM */
.card { padding: var(--space-4); }
.card__title { font-weight: 700; }
.card--featured { border-color: var(--color-accent); }

/* Card.module.css：构建后类名局部化 */
.title { color: var(--color-text); }
```

选择步骤：

1. 列出是否需要全局主题、运行时动态值、SSR、静态提取和跨项目复用。
2. 用一个真实组件实现默认、悬停、焦点、禁用、错误和响应式状态。
3. 检查构建产物中的重复规则、加载顺序、未使用样式和缓存行为。
4. 记录新增变体的写法以及何时允许全局规则。
5. 在多组件和主题切换后再确定方案，不凭单个演示组件决定。

## 常见错误与边界

- BEM 不能阻止全局层叠，过长名称也可能反映组件边界不清。
- CSS Modules 不自动解决设计一致性、选择器复杂度或跨模块全局规则。
- 运行时 CSS-in-JS 可能产生注入、序列化、缓存和水合成本；必须按具体实现测量。
- Utility CSS 可能产生重复组合和难审查差异，应通过组件封装或受控变体减少散落规则。
- 不要只按流行度切换方案；迁移会影响 HTML、测试选择器、SSR、缓存和设计系统。

## 相关补充知识

原生 Cascade Layers 可以固定样式来源优先级，`@scope` 可限制选择器匹配范围，自定义属性可传递主题和动态值。这些能力能减少部分工具需求，但不负责组件 API、文件组织和团队治理。混合方案是允许的，例如 CSS Modules 配合 Utility 类，但必须规定各自负责的范围。

## 来源

- [getbem.com](https://getbem.com/naming/)
- [github.com](https://github.com/css-modules/css-modules)
- [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/CSS/@scope)

访问日期：2026-07-16。
