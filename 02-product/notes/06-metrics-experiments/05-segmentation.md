---
title: 按用户群、版本与来源分析
stage: intermediate
direction: product
tags:
  - product
  - metrics
  - segmentation
  - cohort
---

# 按用户群、版本与来源分析：让总体指标暴露组成差异

分群把分析单位按预先定义的维度划为可比较子集。它用于定位差异、检验外推边界和选择行动，不会自动解释原因。Cohort 是共享某个起始条件和时间的群体，特别适合观察生命周期。

## 一、分群之前固定分析单位

用户、设备、会话、工作区和订单不能混用。若崩溃率分母是会话，版本维度取本次会话版本；若激活率按工作区，来源要定义为工作区首次合格来源，而不是任意成员最后点击。

一个单位必须明确是否能进入多个群。版本按事件时版本可跨群；首次获客来源通常固定到一个群；功能使用标签可能多值。比较前说明互斥、穷尽与 unknown。

## 二、常用维度的时间语义

### 1. 用户/账户属性

计划、规模、地区、行业会变化。使用事件时值、期末值还是首次值会回答不同问题。分析升级前体验时不能用升级后的计划回填历史。

### 2. 版本

客户端版本来自事件；服务端版本来自部署或请求 trace；实验版本来自 assignment。`latest` 是变化标签，不可用于历史复现。

### 3. 来源

first-touch、last-touch 和活动归因不同。来源参数会丢失、被隐私机制限制或被自然流量覆盖。产品行为分析应避免把营销归因算法当客观身份。

### 4. Cohort

Cohort inclusion 定义起始条件，granularity 定义日/周/月边界，return criterion 定义后续行为。Google Analytics 的标准、rolling、cumulative 计算含义不同，不能只写“留存”。

## 三、成员规则合同

| 字段 | 示例 |
|---|---|
| unit | workspace_id |
| inclusion | 首次创建有效项目 |
| dimension | acquisition_week |
| effective_at | 首次有效项目服务端时间 |
| window | 第 1–4 个完整周 |
| membership | 每个工作区一个 acquisition cohort |
| unknown | 无法恢复首次事件单独显示 |
| privacy | 小于阈值不展示 |

## 四、总体均值的组成效应

总体转化率是各群转化率的加权平均：`R = Σ n_g r_g / Σ n_g`。即使每群 r_g 不变，群权重 n_g 改变也会让总体变化。

例：旧版 90% 流量、崩溃率 1%；新版 10%、崩溃率 5%，总体 1.4%。发布扩大到新版 60%，两版率不变，总体变 3.4%。这不是每个版本突然变差，而是组成变化。

反之总体改善可能掩盖每群恶化，即 Simpson's paradox。必须同时看群内变化和权重变化。

## 五、案例一：移动端崩溃

### 指标

分析单位是冷启动会话；崩溃结果是启动后 60 秒内权威 crash report；版本取本次 binary version；设备型号和 OS 来自会话快照。

### 数据

| 版本 | OS | 会话 | 崩溃 | 率 |
|---|---|---:|---:|---:|
| 8.1 | iOS 18 | 80,000 | 800 | 1.0% |
| 8.1 | iOS 19 | 20,000 | 400 | 2.0% |
| 8.2 | iOS 18 | 20,000 | 300 | 1.5% |
| 8.2 | iOS 19 | 80,000 | 2,400 | 3.0% |

8.2 总率 2.7%，8.1 为 1.2%。两个 OS 内 8.2 都更差，方向清楚；同时 8.2 的 iOS 19 权重更高。发布决策需估算标准化到同一 OS 分布后的差异。

### 标准化

使用共同权重 50%/50%：8.1 标准化率 1.5%，8.2 为 2.25%，差 0.75 点。不能用标准化替代实际影响，实际崩溃量仍取决于真实流量。

### SQL

```sql
SELECT app_version, os_major,
       COUNT(*) AS sessions,
       COUNT(*) FILTER (WHERE crashed) AS crashes,
       AVG(CASE WHEN crashed THEN 1.0 ELSE 0.0 END) AS crash_rate
FROM startup_sessions
WHERE started_at >= :start AND started_at < :end
GROUP BY app_version, os_major;
```

### 数据质量

崩溃会阻止正常事件上传，客户端埋点可能系统性漏掉崩溃会话。应使用独立 crash reporter，并对启动 session 与报告关联率做监控。

## 六、案例二：新手激活

激活定义为工作区创建后 7 天内完成项目发布且有第二名成员参与。按 acquisition week 建 cohort，来源采用工作区创建前最后一个合格 campaign。

### 结果

总体激活从 40% 降到 36%。分群显示 organic 42%→43%，partner 55%→54%，paid 25%→25%；下降来自 paid 权重从 20% 升到 45%。

这支持“流量结构变化”诊断，不证明 paid 用户质量差的原因。可能是国家、设备、承诺或欺诈不同。下一步在 paid 内按活动和地区检查，并核对归因规则是否变化。

### Cohort 成熟

本周工作区尚未完整经历 7 天，不能与上月成熟 cohort 比。看板显示 cohort size、mature units 和 activation count，成熟率不足时不计算最终激活。

### 多成员身份

工作区来源不能由第一个完成动作的成员来源覆盖。用户级 campaign 与工作区 acquisition 是不同实体属性。合并身份时保留 lineage。

## 七、版本发布分群

按版本比较需处理自选择：主动升级用户可能更活跃。阶段发布把设备兼容、地区和随机 bucket 影响带入版本。仅观察版本差异不能证明新版造成结果。

如果 rollout bucket 随机且遵守 assignment，可按意向处理估计；若用户可切换版本，记录 assignment 与 exposure 分开。

## 八、来源归因边界

UTM 缺失不等于 organic；跨设备、广告拦截、cookie 限制会产生 unknown。不要把 unknown 按比例偷偷分摊。

first-touch 适合获客 cohort，last-touch 适合最近触点分析，多触点模型依赖额外假设。产品团队应保存原始触点和模型版本。

## 九、小样本、阈值和多重比较

分群越细，样本越小、区间越宽、偶然极值越多。每个率显示分子、分母和区间；低于隐私/稳定阈值合并或隐藏。

查看几十个维度后选择最异常群属于多重探索。把发现标为探索性，在新时间窗口或实验中验证。

## 十、维度数据质量

监控 null/unknown、基数、分布漂移、迟到更新和枚举非法。版本 `8.2`, `8.2.0`, `v8.2` 需规范化但保留 raw。

Slowly changing dimension 要保留有效期：`valid_from <= event_time < valid_to`。只用当前账户规模 join 历史事件会产生未来信息。

## 十一、分群决策树

1. 总体变化是否来自数据/口径版本。
2. 哪些预先重要维度方向一致。
3. 群内率变化还是群权重变化。
4. 差异是否有足够样本和实际意义。
5. 群定义是否在结果前确定。
6. 哪个实验或定性证据能区分解释。

## 十二、留存 cohort 的计算

Standard 第k期统计该期返回者；rolling 要求此前每期都返回；cumulative 统计截至k期至少一次返回。100人中第2周40人活跃、连续两周30人、任一周70人，则三者为40%、30%、70%。

粒度按属性时区的日/周/月边界，不等于滚动7/30天。跨时区产品需固定业务时区或按用户时区并说明。

## 十三、直接标准化

共同权重 `w_g` 下 `R_std=Σw_g r_g`。标准化回答若群结构相同率会怎样，实际事件量仍按真实权重。权重用合并总体或目标总体，需预先说明。

iOS18/19各50%时，8.1标准化1.5%，8.2为2.25%。实际8.2为2.7%，其中部分来自iOS19权重更高。

## 十四、事件时维度

```sql
SELECT e.app_version,d.plan_tier,COUNT(*) units,AVG(e.activated::int) rate
FROM activation_events e JOIN workspace_plan_history d
 ON e.workspace_id=d.workspace_id
AND e.eligible_at>=d.valid_from
AND e.eligible_at<COALESCE(d.valid_to,'infinity')
GROUP BY e.app_version,d.plan_tier;
```

维表有效期不得重叠。join后0行进入unknown质量队列，多行是错误，不能用DISTINCT隐藏。

## 十五、小群隐私

小企业、罕见设备和人口属性可能可识别。设展示阈值、访问控制和导出审计；suppressed 后可见群之和不等于总体，界面说明。

群差异不直接支持差别待遇。行动针对可改变障碍，涉及公平性时审查测量偏差与潜在伤害。

## 十六、常见错误

| 错误 | 修正 |
|---|---|
| 用当前属性回填历史 | 使用事件时维表 |
| 只报百分比 | 同报分子分母区间 |
| 无穷切片找显著 | 预注册或新样本验证 |
| 隐藏 unknown | 单列并监控 |
| 版本差异当因果 | 检查 rollout/自选择 |
| 未成熟 cohort | 明确成熟窗口 |

## 十七、综合练习

分析一次激活率下降，至少按 acquisition cohort、来源、版本、地区和账户规模分解。

### 验收标准

- [ ] 单位、成员规则与属性时间语义明确。
- [ ] 区分群内率变化与群权重变化。
- [ ] SQL 使用事件时维度。
- [ ] unknown 与小样本不被隐藏。
- [ ] cohort 成熟度可见。
- [ ] 探索发现使用新窗口或实验验证。

## 十八、来源与版本的交叉分解

单独看来源或版本可能遗漏交互。假设 paid 来源主要进入 8.2，新版恰有激活 bug，总体会看似 paid 质量差。建立来源×版本表：

| 来源 | 8.1 激活 | 8.2 激活 | 8.2 流量占比 |
|---|---:|---:|---:|
| organic | 43% | 35% | 20% |
| partner | 54% | 45% | 30% |
| paid | 41% | 24% | 85% |

三个来源内 8.2 都下降，paid 因 8.2 暴露高而总体最低。行动优先修版本，而不是暂停渠道。交叉表样本不足时合并窗口或用模型估计，但要报告区间。

## 十九、身份合并的敏感性

匿名用户登录后，历史事件是否回填到 user_id 会改变 cohort 大小。分别计算“仅登录后”和“允许确定性回填”两版，比较差异。使用模糊设备匹配会引入错误合并，不作为权威用户指标。

工作区多成员场景中，一个成员可属于多个工作区。用户留存与工作区留存分别建 fact；不要通过任意成员活动让所有工作区都活跃。

## 二十、来源归因窗口

广告点击到注册的归因窗口从 7 天改成 30 天会改变来源权重。指标版本记录窗口和优先级。重复触点保留 raw，first/last 模型在查询层实现，便于重算。

自然增长与 campaign 可能同时作用。来源分群用于诊断，不把 last-touch 直接当增量贡献。增量需实验或适当因果设计。

## 二十一、分群测试清单

- 同一互斥维度每个单位恰好一群。
- 群分母之和与总体守恒，除非明确 suppressed。
- 事件时属性 join 不产生零行或多行。
- unknown 率低于阈值且变化告警。
- 新枚举不会落入默认旧群。
- 结果发生后的属性不用于定义处理前群。
- 数据删除后分群仍满足隐私规则。

## 二十二、固定群与动态群

固定群用于比较同一批单位随时间的变化，例如 7 月 acquisition cohort。动态 segment 每天按当前条件重算，例如“当前企业计划”。把动态 segment 的月趋势解释为同一批用户行为会混入成员进出。

账户规模从 small 升为 medium 时，事件时分群会在升级后进入新群；首次规模分群则永久留在 small。前者回答当前规模体验，后者回答初始客户类型的长期结果。图表标题明确。

## 二十三、分群查询的守恒验证

对互斥维度运行 `SUM(group_count)=overall_count`。差额拆成 unknown、suppressed、late-arriving 和重复。对多值标签不要求守恒，但说明一个单位可贡献多群，不能把群数相加当总体。

版本规范化表保留 raw、canonical 与 mapping_version。新增 `8.2-hotfix` 时先更新映射并回填质量报告，不能自动截断字符串误归 8.2。

## 二十四、决策输出

分群分析最终写行动对象、证据、未排除解释和验证方法。例如“暂停8.2 rollout并修编码”比“paid用户质量差”更可操作也更少污名。若差异无法区分，输出下一步测量而不是强行命名原因。

## 来源

- [Google Analytics：Cohort exploration](https://support.google.com/analytics/answer/9670133?hl=en)（访问日期：2026-07-18）
- [Google Analytics Data API：Dimensions and Metrics](https://developers.google.com/analytics/devguides/reporting/data/v1/exploration-api-schema)（访问日期：2026-07-18）
- [NIST：Confidence Limits for the Mean](https://www.itl.nist.gov/div898/handbook/eda/section3/eda352.htm)（访问日期：2026-07-18）
