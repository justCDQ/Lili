---
type: ai-note
stage: beginner
topic: secrets-permissions-cost
verified: 2026-07-16
tags: [ai, security, secrets, least-privilege, cost]
---

# Secret、最小权限与费用上限

## 是什么

Secret 是能够授予系统访问能力的敏感值，包括 API Key、访问 Token、密码、私钥和连接串。最小权限要求身份只获得完成当前任务所需的最少资源、操作和时间范围。费用上限通过账户预算、项目配额、速率限制、请求限制和应用内计量控制模型调用的财务风险。

## 为什么需要

模型 API 通常按输入、输出或工具使用计费。Secret 泄露会导致未授权访问、数据暴露和费用损失；Agent 或循环代码即使没有攻击者也可能因失控重试和长上下文快速增加费用。

## 关键特性

### Secret 生命周期

Secret 需要创建、分发、存储、使用、轮换、吊销和审计。长期静态 Key 风险高于短期凭据。Secret 不应出现在源代码、Git、URL、客户端包、截图、日志和模型上下文中。

### 最小权限

开发、测试、CI 和生产使用不同身份。只读与写入权限分开；不同项目和租户分开；能够限制模型、区域或预算时应显式限制。高风险操作需要服务端再次授权，不能相信模型生成的参数。

### 费用边界

控制分为多层：供应商账户预算与告警、项目级配额、用户/租户速率限制、单请求最大输入输出、Agent 最大步骤和总超时、队列并发、异常检测与熔断。告警不是硬限制，必须确认供应商预算功能是否会自动停止请求。

## 实际怎么使用

上线前检查：

```text
[ ] Secret 由服务端 Secret Manager 或受控运行环境提供
[ ] 本地 Secret 文件已忽略，仓库启用 Secret 扫描
[ ] 开发、CI、生产使用不同凭据
[ ] Key 可立即吊销和轮换
[ ] 日志和错误报告会脱敏
[ ] 设置账户预算、项目告警和应用硬上限
[ ] 每请求限制输入、输出、重试、步骤和总时间
[ ] Usage 按用户、租户、模型、功能和版本记录
[ ] 费用异常能够停止队列或关闭功能
```

发生泄露时先吊销/轮换，不要先花时间清理 Git 历史。随后确认使用日志、影响范围和下游凭据，再清理历史并补检测规则。

## 常见错误与边界

- 把 Secret 放入 `.env` 后认为绝对安全；`.env` 只是本地配置方式，仍可能被误提交、备份或读取。
- 在浏览器调用模型供应商并暴露长期 Key。
- 只设费用告警，不设应用级最大 Token、步骤、并发和速率。
- 多租户共享高权限凭据且不记录租户 Usage，无法归责和隔离。
- 把 Tool 输出和外部文档当可信指令，导致模型诱导高成本或高权限调用。

## 补充知识

费用治理也影响可靠性：限制过严会造成正常请求失败，因此应区分硬上限、软告警和降级策略。可以在达到阈值后切换低成本模型、减少上下文、禁用非必要 Tool 或转人工处理，但必须记录质量影响。

## 来源

- [OWASP Secrets Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)（访问日期：2026-07-16）
- [OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)（访问日期：2026-07-16）
- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)（访问日期：2026-07-16）

