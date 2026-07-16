# SEO、Open Graph 与结构化内容基础

## 是什么与为什么需要

SEO 让搜索引擎能抓取、理解并呈现页面；Open Graph 用 `meta property` 描述社交分享对象；结构化数据以机器可读词汇声明实体及属性。它们改善发现与展示，不保证排名或富结果。

## 关键特性与规则

- 可抓取的真实正文、语义结构、有效链接和正确 HTTP 状态是基础。
- canonical、robots 和 sitemap 各自解决规范 URL、抓取提示和 URL 发现，不能互相替代。
- Open Graph 面向分享预览，应使用绝对 URL 和可访问图片。
- 结构化数据必须与可见内容一致，并符合目标搜索平台支持的类型和必需字段。
- SEO 配置只能提供信号，不能保证索引、排名、摘要或富结果展示。

## 实际使用

```html
<title>HTML 表单校验指南｜狸力</title>
<meta name="description" content="解释原生约束校验、错误反馈和服务端验证。">
<link rel="canonical" href="https://example.com/html/forms">
<meta property="og:title" content="HTML 表单校验指南">
<meta property="og:type" content="article">
<meta property="og:url" content="https://example.com/html/forms">
<meta property="og:image" content="https://example.com/images/forms-cover.png">
<script type="application/ld+json">{"@context":"https://schema.org","@type":"TechArticle","headline":"HTML 表单校验指南"}</script>
```

使用描述性标题、语义标题层级、可抓取链接和真实正文；返回正确状态码并提供 sitemap/robots 配置。结构化数据必须与页面可见内容一致、符合所用搜索平台支持类型，并用验证工具检查。

## 常见错误与边界

`robots.txt` 不是访问控制；敏感资源需鉴权。canonical 是提示。关键词堆砌、隐藏文本和虚假结构化数据可能被忽略或处罚。OG 与标准 meta/结构化数据职责不同，应分别维护绝对 URL 和合适图片。

## 补充知识

JavaScript 页面仍需保证爬虫最终能取得主要内容。站点变更后用搜索平台检查抓取、索引、结构化数据错误和实际规范 URL。

## 来源

- [Google Search：SEO starter guide](https://developers.google.com/search/docs/fundamentals/seo-starter-guide)
- [Google Search：Structured data](https://developers.google.com/search/docs/appearance/structured-data/intro-structured-data)
- [Open Graph protocol](https://ogp.me/)

访问日期：2026-07-16。
