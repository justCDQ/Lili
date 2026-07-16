# 标题、段落、列表、链接、图片、音视频与表格

## 是什么与为什么需要

这些元素把内容关系写入文档：`h1`–`h6` 标题，`p` 段落，`ul/ol/li` 列表，`a` 链接，`img` 图片，`audio/video` 媒体，`table` 表格数据。正确元素可被浏览器、搜索引擎和辅助技术识别。

## 实际使用与规则

```html
<h1>订单</h1><p>最近 30 天订单。</p>
<ol><li>确认商品</li><li>支付</li></ol>
<a href="/help/payment">支付帮助</a>
<img src="chart.png" alt="五月销量较四月增长 12%" width="800" height="450">
<video controls width="640"><source src="demo.webm" type="video/webm"><track kind="captions" src="zh.vtt" srclang="zh" label="中文"></video>
<table><caption>月度销量</caption><thead><tr><th scope="col">月份</th><th scope="col">销量</th></tr></thead><tbody><tr><th scope="row">五月</th><td>112</td></tr></tbody></table>
```

链接文字应脱离上下文仍说明目标；外部新窗口若使用 `target="_blank"` 应明确告知。内容图片提供有意义 `alt`，装饰图用 `alt=""`。媒体提供 `controls` 和字幕/文本替代。表格仅表达二维数据，使用 `caption`、`th`、`scope` 建立关系，不用于布局。

## 常见错误与边界

标题级别不为字号服务；用 CSS 调整外观。列表项必须在列表容器中。图片 `alt` 不重复“图片”。自动播放常受限且妨碍用户。复杂表格可能需 `headers/id`；移动端还需处理溢出。

## 补充知识

相对链接以当前文档 URL 为基准。为图片设置固有 `width/height` 可提前保留宽高比空间。音视频编码支持存在浏览器差异，可提供多个 `source` 和下载链接。

## 来源

- [MDN：HTML elements reference](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements)
- [MDN：Creating links](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/Creating_links)
- [MDN：HTML video and audio](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/HTML_video_and_audio)
- [W3C WAI：Tables tutorial](https://www.w3.org/WAI/tutorials/tables/)

访问日期：2026-07-16。
