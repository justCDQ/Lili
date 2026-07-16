# label、fieldset、legend、错误提示与帮助文本

## 是什么与为什么需要

`label` 为控件提供可访问名称并扩大点击区域；`fieldset/legend` 为相关控件建立组名；帮助文本说明格式，错误提示说明问题与修正方法。这些关系必须能被程序确定，不能只靠位置或颜色。

## 实际使用与规则

```html
<fieldset><legend>通知方式</legend>
 <label><input type="radio" name="notify" value="email"> 邮件</label>
 <label><input type="radio" name="notify" value="sms"> 短信</label>
</fieldset>
<label for="username">用户名</label>
<input id="username" name="username" aria-describedby="username-help username-error" aria-invalid="true">
<p id="username-help">4–20 个字母或数字。</p>
<p id="username-error">用户名含空格，请删除空格。</p>
```

显式 label 的 `for` 必须精确匹配唯一 `id`。错误出现后设置 `aria-invalid="true"`，用 `aria-describedby` 关联说明；提交后提供错误摘要、指向字段的链接并把焦点移到合适位置。

## 常见错误与边界

不要仅用 placeholder、`title` 或图标代替可见 label。不要把整个复杂区域包进 label。`aria-describedby` 不控制视觉显示。动态错误可用 live region 通知，但避免每次输入都打断用户。服务端错误返回时保留非敏感输入。

## 补充知识

必填状态应在 label/legend 文本中可理解；星号需解释。只收集完成任务所需字段可降低填写负担。

## 来源

- [W3C WAI：Labeling controls](https://www.w3.org/WAI/tutorials/forms/labels/)
- [W3C WAI：Grouping controls](https://www.w3.org/WAI/tutorials/forms/grouping/)
- [W3C WAI：User notifications](https://www.w3.org/WAI/tutorials/forms/notifications/)

访问日期：2026-07-16。
