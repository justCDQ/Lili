---
stage: intermediate
direction: backend-data
topic: upload-security
---

# 上传类型与内容校验、权限、病毒扫描与审计

文件上传是不可信字节进入存储和解析器的入口。安全链路要限制大小与速率、验证magic和结构、隔离未扫描对象、使用最小权限扫描器、阻止主动内容执行，并保留可追踪审计。

## 1. 数据流与责任边界

```mermaid
flowchart LR
 A["客户端/生产者"] --> B["授权与受控元数据"]
 B --> C["对象存储"]
 C --> D["异步处理/验证"]
 D --> E["业务发布与审计"]
```

对象存在不代表业务有效；数据库中的所有权、状态、object version 和校验结果决定它能否被引用。

## 2. 大小限制

机制：网关、预签policy、对象HEAD和解压后分别限制。

使用：防内存/存储/压缩炸弹。

失败：只信Content-Length。

验证：超限/分块/缺header测试。

取舍：多层校验增加实现但封闭绕过。

围绕 大小限制 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 3. 文件名

机制：只作显示metadata，服务端生成随机key。

使用：避免路径、控制字符和覆盖。

失败：把../../name作为object key/Content-Disposition。

验证：规范化和安全编码测试。

取舍：保留原名利于用户但需最小化。

围绕 文件名 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 4. Content-Type

机制：客户端声明仅作提示。

使用：下载表现和初步路由。

失败：把image/png当真实性证明。

验证：magic bytes+结构解析。

取舍：严格类型可能拒绝边界格式。

围绕 Content-Type 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 5. Magic/结构

机制：读取受限头和安全解析器验证格式。

使用：识别伪装文件与损坏。

失败：仅扩展名白名单。

验证：构造polyglot/truncated样本。

取舍：解析器本身有漏洞面。

围绕 Magic/结构 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 6. Quarantine

机制：未扫描对象与公开/业务可见区隔离。

使用：阻止异步窗口暴露。

失败：上传后立即生成公共URL。

验证：策略证明app不能读quarantine。

取舍：多一次copy/状态管理。

围绕 Quarantine 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 7. Malware scan

机制：隔离worker读取并扫描，结果带engine/signature version。

使用：发现已知恶意内容。

失败：扫描成功等于永久安全。

验证：EICAR/超时/引擎更新后重扫。

取舍：不能发现全部未知威胁。

围绕 Malware scan 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 8. Archive safety

机制：限制层数、文件数、总解压大小、压缩比和路径。

使用：防zip bomb/zip slip。

失败：直接解压到共享目录。

验证：恶意archive corpus和资源限制。

取舍：限制可能拒绝合法大包。

围绕 Archive safety 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 9. Active content

机制：SVG/HTML/PDF/Office可包含脚本、链接或宏。

使用：隔离域下载、转码/消毒。

失败：同源inline展示用户HTML/SVG。

验证：CSP、Content-Disposition、独立origin测试。

取舍：转码损失特性但降低风险。

围绕 Active content 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 10. 权限

机制：上传、扫描、发布、下载使用不同身份和prefix。

使用：最小权限与职责分离。

失败：扫描器拥有全桶删除。

验证：策略测试和CloudTrail/审计。

取舍：角色更多但缩小爆炸半径。

围绕 权限 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 11. 审计

机制：记录actor、upload、object version、checksum、scan和发布决策。

使用：追责、重新扫描和删除。

失败：记录签名URL/token/完整PII。

验证：按upload ID关联并脱敏。

取舍：存储与隐私保留需平衡。

围绕 审计 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 12. 方案比较

|方案|收益|边界|
|---|---|---|
|同步扫描|用户立即得结果|请求超时/大文件|
|异步quarantine|可扩展隔离|状态机复杂|
|原文件下载|保真|主动内容风险|
|转码后发布|统一安全格式|质量/功能损失|
|杀毒签名|已知威胁|不能证明绝对安全|

## 13. 完整案例：企业附件上传

### 输入与约束

PDF/Office最大100MiB，扫描期间不可下载，tenant严格隔离。

### 处理步骤

1. API创建upload记录与quarantine随机key。
2. 直传后服务端HEAD校验size/checksum并识别真实类型。
3. scanner用只读quarantine+写clean权限扫描，宏文件按策略隔离。
4. 通过后数据库条件迁移state并发布固定object version。
5. 下载API重新授权并签clean version，审计全链路。

### 输出

只有clean且授权对象可下载，扫描结果可按引擎版本重查。

### 验证

EICAR、伪扩展、超大、扫描timeout、跨tenant、旧worker发布均被阻止。

### 失败分支

扫描前业务表直接引用对象会被猜测下载；quarantine policy必须从存储层拒绝。

### 恢复要求

失败后以数据库 upload/job 状态和对象 version 为准，重复任务使用 generation/条件更新；孤儿对象由清单与生命周期受控回收，不能根据用户文件名删除。

## 14. 完整案例：压缩包导入

### 输入与约束

ZIP包含CSV/图片，最多1000文件、解压2GiB、层数2，不能路径穿越。

### 处理步骤

1. 隔离容器无网络/低权限解析。
2. 逐entry规范化路径，拒绝绝对路径、..、symlink和设备文件。
3. 流式累计文件数、未压缩字节、压缩比和深度。
4. 每个文件单独类型/病毒校验，不执行宏/脚本。
5. 验证全部后事务创建导入任务，临时文件生命周期清理。

### 输出

恶意包在资源上限内失败，合法包产生明确清单。

### 验证

zip slip、zip bomb、嵌套、重名、Unicode路径和损坏CRC样本。

### 失败分支

先完整解压再检查上限会先耗尽磁盘；限制必须在流式写入前/过程中。

### 恢复要求

失败后以数据库 upload/job 状态和对象 version 为准，重复任务使用 generation/条件更新；孤儿对象由清单与生命周期受控回收，不能根据用户文件名删除。

## 15. 故障注入矩阵

|注入|预期|禁止|
|---|---|---|
|签名后撤权|`upload_rejected_size` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|上传中断|`type_mismatch` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|同key并发写|`scan_queue_age` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|checksum错误|`scan_timeout` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|扫描器崩溃|`malware_detected` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|生命周期延迟|`quarantine_bytes` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|对象存储429/503|`publish_conflict` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|数据库提交后断线|`download_denied` 可观察且状态可恢复|越权发布、静默损坏或无限重试|

## 16. 调试与观测

1. `upload_rejected_size`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
2. `type_mismatch`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
3. `scan_queue_age`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
4. `scan_timeout`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
5. `malware_detected`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
6. `quarantine_bytes`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
7. `publish_conflict`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
8. `download_denied`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
9. `archive_limit`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
10. `audit_gap`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。

按 upload ID/object ID 从业务记录、签名、对象metadata/version、扫描任务到下载审计逐跳核对；签名query和敏感正文不进入普通日志。

## 17. S3兼容实现边界

1. Quarantine 与 clean 发布依赖的一致性模型必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
2. 条件请求 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
3. ETag/Checksum 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
4. 版本与delete marker 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
5. multipart最小part/part数 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
6. 生命周期执行时间 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
7. IAM/Policy语法 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
8. 加密/KMS 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
9. 事件通知 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
10. List分页 必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。

## 18. 生产检查

1. bucket/key由服务端控制；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
2. 所有业务读取重新授权；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
3. 上传有大小/速率/checksum；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
4. 未验证对象隔离；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
5. 版本和删除策略明确；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
6. multipart孤儿可清理；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
7. 凭据最小权限且可轮换；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
8. 敏感URL不进日志；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
9. 故障有重试预算与幂等；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。
10. 备份/恢复和隐私删除已演练；上传类型与内容校验、权限、病毒扫描与审计 的负责人提供可重复验证证据。

## 19. 综合练习与验收

实现“企业附件上传”并用“压缩包导入”验证不同约束。提交状态机、策略、对象清单、失败注入和观测面板。

- [ ] 大小限制 的正常、边界、权限和失败路径均通过。
- [ ] 文件名 的正常、边界、权限和失败路径均通过。
- [ ] Content-Type 的正常、边界、权限和失败路径均通过。
- [ ] Magic/结构 的正常、边界、权限和失败路径均通过。
- [ ] Quarantine 的正常、边界、权限和失败路径均通过。
- [ ] Malware scan 的正常、边界、权限和失败路径均通过。
- [ ] Archive safety 的正常、边界、权限和失败路径均通过。
- [ ] Active content 的正常、边界、权限和失败路径均通过。
- [ ] 两个完整案例都可在隔离环境重复运行。
- [ ] 对象存储故障不改变数据库业务不变量。
- [ ] 所有孤儿、版本、parts和审计记录有保留/清理策略。

## 来源

- [OWASP File Upload Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html)（访问日期：2026-07-17）
- [AWS S3 security best practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/security-best-practices.html)（访问日期：2026-07-17）
- [AWS S3 presigned URL best practices](https://docs.aws.amazon.com/prescriptive-guidance/latest/presigned-url-best-practices/introduction.html)（访问日期：2026-07-17）
- [CISA Malware Analysis](https://www.cisa.gov/resources-tools/services/malware-analysis)（访问日期：2026-07-17）
- [RFC 9110 Content-Disposition references](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-17）
