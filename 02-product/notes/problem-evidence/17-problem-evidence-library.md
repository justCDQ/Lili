# 问题证据库：来源、日期与可追溯性

## 是什么

问题证据库把问题陈述、原始证据、分析判断和决策连接起来，使结论可复核、可更新。它不是功能愿望池，也不是复制所有原文的资料仓库。

## 为什么需要

公开网页、产品版本和判断都会变化。证据库保存来源、时间与推断链，避免重复研究，并能在新证据出现时定位受影响的结论和决策。

## 关键特性或判断规则

- 每个问题有稳定 ID，原始证据与分析结论分开保存。
- 来源包含原链接、发生/发布日期、访问日期、版本和地区。
- 同义问题可合并，但原始记录、反证和状态历史不能覆盖。
- 只保留最少必要摘录和个人数据，并设置定期复核日期。

## 数据结构

每条至少包含：`问题 ID｜角色/场景｜问题陈述｜事实｜来源链接｜来源类型｜发生/发布日期｜访问日期｜产品版本/地区｜频率｜严重度｜范围｜可信度｜反证｜当前做法｜替代方案｜待验证假设｜状态｜最后复核日`。

## 实际怎么使用

新证据先保存原链接和最小必要摘录，再关联已有问题；同义问题合并但保留原记录；事实与推断分字段；每月检查失效链接和状态；版本变化后重新验证；重要决策引用问题 ID。个人信息只保存分析所必需内容。

## 常见错误与边界

只存截图导致无法追溯；没有访问日期；复制整篇受版权保护内容；将来源删除或页面更新误认为原结论仍有效；用“已解决”覆盖历史证据。无法确认许可时保存链接和自己的摘要。

## 补充知识

W3C PROV 将实体、活动和责任主体之间的来源关系标准化。个人仓库无需实现完整标准，但至少应能从结论回到证据，再回到产生该证据的时间和条件。

## 来源

- [W3C PROV Overview](https://www.w3.org/TR/prov-overview/)（访问日期：2026-07-16）
- [GitHub Docs：About issue and pull request templates](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates)（访问日期：2026-07-16）
- [UK ICO：Data minimisation](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/data-protection-principles/a-guide-to-the-data-protection-principles/)（访问日期：2026-07-16）
