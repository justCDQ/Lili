# head、元数据、favicon、语言与页面标题

## 是什么与为什么需要

`head` 保存文档元数据和资源关系，不是可见主内容。字符集影响解码；`lang` 声明主要语言；`title` 命名文档并用于标签页、历史和搜索结果；favicon 标识站点；viewport 控制移动视口。

## 关键特性与规则

- `meta charset` 应尽早声明并与响应头和文件实际编码一致。
- `lang` 写在 `html` 上，局部语言变化再由子元素覆盖。
- 每个文档必须有一个非空且能区分页面的 `title`。
- `meta viewport` 应允许用户缩放，不能用它解决布局问题。
- `link` 表达 favicon、样式表、canonical 等资源关系，URL 解析受文档地址影响。

## 实际使用

```html
<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>购物车（2 件商品）｜狸力商店</title>
  <meta name="description" content="查看并修改购物车商品，继续结算。">
  <link rel="icon" href="/favicon.svg" type="image/svg+xml">
  <link rel="stylesheet" href="/assets/app.css">
</head>
<body>...</body></html>
```

字符集声明应位于文档前部。每页只有一个非空 `title`，内容具体且能区分页面。`lang` 使用有效 BCP 47 标签，局部语言变化在子元素覆盖。favicon 可提供多个尺寸/格式，路径须部署可达。

## 常见错误与边界

`title` 属性工具提示不等于页面 `<title>`。`meta keywords` 对主流搜索排序通常无实际作用。description 可能被搜索引擎改写。不要禁止缩放。资源相对路径受文档 URL 和可选 `base` 影响，`base` 会改变全部相对 URL，应谨慎。

## 补充知识

HTTP `Content-Type` 的 charset 可能影响解码；文档声明与实际编码应一致。动态单页应用导航后也应更新标题并管理焦点。

## 来源

- [WHATWG HTML：Document metadata](https://html.spec.whatwg.org/multipage/semantics.html#semantics)
- [MDN：What's in the head](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/Webpage_metadata)
- [W3C WAI：Page title](https://www.w3.org/WAI/test-evaluate/easy-checks/page-title/)

访问日期：2026-07-16。
