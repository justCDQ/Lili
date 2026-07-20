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

上传对象必须先处于存储层不可下载的 quarantine 状态。发布条件应绑定精确 object version、内容识别结果、扫描引擎与规则版本；扩展名、客户端 Content-Type 或一次“未发现恶意”都不能单独证明文件安全。

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

扫描任务以 object version 和扫描策略版本幂等；超时或引擎不可用时对象继续隔离，不降级为 clean。worker 只有在 generation 仍匹配时才能条件发布，旧结果不得覆盖重新上传的新版本；引擎或规则升级后可为仍保留的版本安排重扫。

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

导入任务绑定不可变压缩包 version，并持久化已验证 entry 清单与资源计数。失败后删除该 job ID 的临时解压空间并从原包重新验证；全部 entry 通过前不创建可见业务数据。失败清单保留受限审计信息，便于区分路径穿越、炸弹阈值、CRC 与恶意内容。

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

按 upload ID 依次核对 object version/checksum、magic/结构识别、quarantine policy、扫描 attempt、引擎/规则版本、verdict 与发布 CAS；压缩包再检查具体 entry 和触发的资源上限。审计记录主体、动作、对象 ID、策略版本和结果，不记录签名 query、正文或未经处理的敏感文件名。

## 17. S3兼容实现边界

1. Quarantine 必须由 bucket policy、独立 bucket 或拒绝读取的访问路径落实，不能只靠应用状态字段。
2. 类型识别需要读取足够的文件头和内部结构；验证对象存储的 range GET、最大扫描大小与流式读取行为。
3. 保存明确的 checksum 与 object version，避免扫描任务读取到同 key 后写入的新内容。
4. 对象事件可能重复、延迟或乱序；消费者以 version/generation 幂等，事件只触发工作而不作为发布事实。
5. 扫描服务账号只读 quarantine、只写受控结果；若对象加密，单独配置最小 KMS 解密权限。
6. 从 quarantine 到 clean 若通过 copy/tag 完成，要验证复制后的 checksum/version，并用数据库条件更新发布。
7. 生命周期不能在法务保留、人工复核或恶意样本调查结束前删除证据；不同 verdict 使用不同保留期。
8. 压缩包解析在隔离运行时执行，临时目录按 job ID 分配并在成功、失败和进程崩溃后都可回收。
9. 扫描和导入都要覆盖超时、资源耗尽、损坏文件与引擎崩溃，失败状态不能自动转为允许下载。
10. 审计导出与对象枚举必须正确处理分页，并定期对账 clean 对象、业务引用和扫描记录是否一一对应。

## 18. 生产检查

1. 存储策略直接拒绝读取 quarantine 对象，应用路由遗漏也不能绕过隔离。
2. 上传入口同时限制声明大小、实际读取字节、速率和租户配额，超限立即终止流。
3. 类型判定覆盖扩展名、Content-Type、magic 与内部结构，不一致结果按拒绝策略处理。
4. scanner 读取精确 object version，结果记录引擎、规则、签名库版本和完成时间。
5. timeout、崩溃或未知 verdict 保持 quarantine，只有匹配 generation 的 clean 结果可以发布。
6. EICAR、宏文档、伪扩展、损坏文件与 active content 测试集在扫描升级后回归。
7. 压缩包在隔离环境流式解析，文件数、层数、展开字节、压缩比与路径规则逐 entry 生效。
8. 下载每次重新授权；跨 tenant、旧 worker、旧 version 和被撤权用户的拒绝路径已验证。
9. 审计能关联上传、扫描、发布、下载与删除，同时避免记录正文、签名和敏感原始文件名。
10. quarantine 保留、人工复核、恶意证据和隐私删除的优先级及执行期限已经演练。

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
