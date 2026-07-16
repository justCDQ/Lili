# 表单控件、按钮、原生校验、自动填充与提交

## 是什么与为什么需要

`form` 组织可提交控件。`input`、`select`、`textarea` 收集值，`button` 触发动作。约束属性提供客户端校验，`autocomplete` 令浏览器识别字段，提交把“成功控件”的 `name=value` 按方法和编码发送。

## 关键特性与规则

- 控件需要稳定 `name` 才能进入表单数据，并用可见 `label` 提供名称。
- `type`、`required`、`min/max/step` 和 `pattern` 描述浏览器可执行的基础约束。
- `autocomplete` 使用标准 token 描述字段用途，而不是把关闭自动填充当安全措施。
- 提交前的客户端校验改善反馈，服务端仍必须重新验证、授权和防重放。
- 提交失败应保留有效输入，并把错误关联到字段和错误摘要。

## 实际使用

```html
<form action="/signup" method="post">
 <label for="email">邮箱</label>
 <input id="email" name="email" type="email" autocomplete="email" required>
 <label for="password">密码</label>
 <input id="password" name="password" type="password" autocomplete="new-password" minlength="12" required>
 <button type="submit">注册</button>
 <button type="button" id="preview">预览</button>
</form>
```

控件必须有 `name` 才进入提交数据；checkbox 未选中通常不提交。表单内 `button` 默认是 submit，因此非提交按钮显式写 `type="button"`。使用准确 `type`、`required`、`min/max/step/pattern`；复杂跨字段规则再调用 Constraint Validation API。

## 常见错误与边界

客户端校验可被绕过，服务端必须重新验证、授权和防伪造。`placeholder` 不是 label。`autocomplete="off"` 可能被浏览器忽略，登录字段应使用标准 token。`form.submit()` 不触发约束校验和 submit 事件，通常用 `requestSubmit()`。

## 补充知识

GET 把表单数据放查询串，适合安全、幂等检索；敏感或改变状态通常 POST，但 HTTPS 才保护传输。文件上传需 `enctype="multipart/form-data"`。

## 来源

- [MDN：Forms and buttons](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Structuring_content/HTML_forms)
- [MDN：Constraint validation](https://developer.mozilla.org/en-US/docs/Web/HTML/Guides/Constraint_validation)
- [MDN：form](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form)

访问日期：2026-07-16。
