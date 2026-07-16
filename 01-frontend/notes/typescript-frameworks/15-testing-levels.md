# Unit、Component 与 E2E 测试

## 是什么

单元测试验证小逻辑单元；组件测试验证渲染与交互边界；E2E 在真实浏览器验证用户流程。层级越高覆盖集成越多，速度和诊断成本通常也更高。

## 为什么需要

单元、组件和端到端测试覆盖不同边界：纯逻辑、可交互 UI 以及真实集成流程。合理分层能在反馈速度、故障定位和环境真实性之间取得平衡。

## 关键特性与规则

测试可观察行为；隔离数据；关键流程少量 E2E，纯逻辑单测，复杂组件做交互测试；失败保存 trace。

## 实际使用

```tsx
// Playwright
test('login',async({page})=>{await page.goto('/login');await page.getByLabel('邮箱').fill('a@b.com');await page.getByRole('button',{name:'登录'}).click();await expect(page).toHaveURL('/');});
```

## 常见错误与边界

按实现细节断言导致脆弱；过度 mock 失去集成价值；覆盖率高不等于需求正确。

## 相关补充知识

测试应通过角色、标签和可见文本操作界面，避免绑定内部 DOM。Mock 适合隔离边界但可能偏离真实协议；关键流程仍需在接近生产的浏览器、网络和服务环境回归。

## 来源

- [web.dev](https://web.dev/learn/testing)
- [Playwright Documentation](https://playwright.dev/docs/intro)
- [Playwright Documentation](https://playwright.dev/docs/best-practices)

访问日期：2026-07-16。
