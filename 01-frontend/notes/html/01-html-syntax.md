# HTML 基础语法：元素、属性、嵌套、注释、空元素与字符引用

## 是什么与为什么需要

元素由开始标签、内容和结束标签构成并表达语义；属性为元素提供配置；嵌套建立文档树。注释 `<!-- -->` 不渲染。空元素不能包含子节点且只有开始标签。字符引用在文本中表达会被当作语法的字符。

## 关键规则

## 实际使用

```html
<!doctype html>
<p class="summary">使用 <strong>正确嵌套</strong>。</p>
<img src="logo.svg" alt="产品标志">
<!-- 待补充发布日期 -->
<p>&lt;button&gt; 表示按钮元素，AT&amp;T 含与号。</p>
```

标签按后进先出闭合，不能写成 `<p><strong>文本</p></strong>`。属性值统一加引号。布尔属性出现即为真，`disabled="false"` 仍为禁用。HTML 空元素包括 `img`、`input`、`meta`、`link`、`br` 等，不写结束标签；`<div />` 不会自闭合。UTF-8 下普通 Unicode 字符可直接书写，至少在文本中转义 `<`、`&`，属性中还要处理所用引号。

## 常见错误与边界

浏览器会容错修复错误标记，但生成 DOM 可能不同于源码。注释不是秘密，会随文档下载。字符引用需以 `&` 开始并通常以 `;` 结束。HTML、SVG、MathML 的自闭合规则不同。

## 补充知识

元素允许的父子关系由内容模型定义，不是任何标签都能任意嵌套。使用 HTML 校验器发现结构错误，再以浏览器 DOM 检查实际解析结果。

## 来源

- [MDN：Basic HTML syntax](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/Basic_HTML_syntax)
- [MDN：Void element](https://developer.mozilla.org/en-US/docs/Glossary/Void_element)
- [WHATWG HTML：Syntax](https://html.spec.whatwg.org/multipage/syntax.html)

访问日期：2026-07-16。
