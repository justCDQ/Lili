---
stage: intermediate
direction: backend-data
topic: multipart-version-lifecycle
---

# Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag

大对象上传需要独立管理 upload ID、part number、校验和和完成清单。S3 ETag 不保证普遍等于对象MD5；版本和生命周期控制恢复与成本，CDN缓存又引入独立失效边界。

## 1. 数据流与责任边界

```mermaid
flowchart LR
 A["客户端/生产者"] --> B["授权与受控元数据"]
 B --> C["对象存储"]
 C --> D["异步处理/验证"]
 D --> E["业务发布与审计"]
```

对象存在不代表业务有效；数据库中的所有权、状态、object version 和校验结果决定它能否被引用。

## 2. CreateMultipartUpload

机制：初始化metadata/加密并返回uploadId。

使用：大文件和可恢复上传。

失败：创建后不完成导致parts收费。

验证：记录uploadId与业务upload状态。

取舍：增加请求与状态管理。

围绕 CreateMultipartUpload 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 3. UploadPart

机制：partNumber标识1..10000，重传同编号覆盖该part。

使用：并发/断点续传。

失败：客户端随意声明已完成part。

验证：ListParts和checksum核对。

取舍：并发高吞吐但内存/网络峰值。

围绕 UploadPart 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 4. CompleteMultipartUpload

机制：提交有序partNumber/ETag清单组合对象。

使用：全部part验证后原子完成。

失败：清单错序/缺part或多次竞争完成。

验证：记录最终version/checksum/size。

取舍：完成后才形成对象。

围绕 CompleteMultipartUpload 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 5. AbortMultipartUpload

机制：停止并清理未完成parts。

使用：取消/超时/失败。

失败：只删业务记录不abort。

验证：生命周期AbortIncomplete+主动清理。

取舍：清理最终可能需重试。

围绕 AbortMultipartUpload 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 6. 断点续传状态

机制：服务端保存upload owner、key、part size、expiry、completed parts。

使用：跨设备/进程恢复。

失败：信任客户端上传列表。

验证：从对象存储ListParts对账。

取舍：可恢复但状态与存储需一致。

围绕 断点续传状态 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 7. Versioning

机制：同key写入分配version ID，delete可产生delete marker。

使用：误覆盖恢复、并发引用固定版本。

失败：认为普通DELETE永久删除所有版本。

验证：列版本/恢复演练/权限分离。

取舍：提高恢复能力与存储成本。

围绕 Versioning 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 8. Lifecycle

机制：按prefix/tag/size/age transition/expire/abort。

使用：成本和保留自动化。

失败：规则重叠误删或把恢复数据转冷过早。

验证：policy模拟、inventory、restore演练。

取舍：自动省钱但执行异步且不可作精确定时器。

围绕 Lifecycle 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 9. CDN

机制：边缘缓存对象响应，origin与cache key独立。

使用：公开静态内容和大下载。

失败：私有对象使用公共cache key。

验证：签名cookie/URL与cache policy测试。

取舍：降低源流量但失效/权限复杂。

围绕 CDN 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 10. ETag

机制：表示对象特定版本的实体标签，格式依上传和加密。

使用：条件请求、缓存验证，在限定场景可校验MD5。

失败：把multipart/SSE-KMS ETag当内容MD5。

验证：使用明确checksum字段。

取舍：广泛可用但语义不能泛化。

围绕 ETag 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 11. 条件请求

机制：If-Match/If-None-Match等避免盲覆盖。

使用：并发写和只创建。

失败：检查后再PUT存在TOCTOU。

验证：直接发送条件header并处理412。

取舍：提升安全但S3兼容实现需核验。

围绕 条件请求 的接口必须保存受控 object ID、tenant、版本和状态；客户端提供的 bucket、key、类型或结果只能作为待验证输入。

## 12. 方案比较

|方案|收益|边界|
|---|---|---|
|单次PUT|简单|失败重传整个对象|
|Multipart|并发和续传|状态、parts成本|
|Versioning|误删恢复|所有旧版本计费|
|内容hash key|天然不可变CDN|需要引用切换|
|固定key+purge|URL稳定|传播延迟与并发覆盖|

## 13. 完整案例：50GB视频上传

### 输入与约束

浏览器弱网、可续传24小时、每part校验、取消需清理。

### 处理步骤

1. API创建业务upload与multipart，选择>=5MiB且part数<=10000。
2. 按part生成短期预签名，限制owner/key/partNumber/uploadId。
3. 客户端并发上传并上报checksum，服务端ListParts对账。
4. 完成时提交有序parts并HEAD最终size/checksum/version。
5. 失败/取消abort，生命周期7天兜底清理。

### 输出

断线只重传缺失parts，完成对象有固定version与校验记录。

### 验证

随机中断/重复part/错序清单/取消；未完成parts最终为零。

### 失败分支

只信客户端parts列表会完成混合/缺失对象；完成前必须与服务端存储状态对账。

### 恢复要求

失败后以数据库 upload/job 状态和对象 version 为准，重复任务使用 generation/条件更新；孤儿对象由清单与生命周期受控回收，不能根据用户文件名删除。

## 14. 完整案例：静态资源发布与回滚

### 输入与约束

Web资源全球CDN，发布后不可出现HTML引用旧/新混合，可快速回滚。

### 处理步骤

1. 构建产物用内容hash key并校验checksum。
2. 先上传全部不可变资源并验证可读。
3. 最后原子发布manifest/HTML版本指针。
4. CDN对hash资源长immutable缓存，HTML短缓存/验证。
5. 回滚切manifest到旧版本，生命周期延迟清理孤儿。

### 输出

新旧资源URL不覆盖，CDN无需大范围purge。

### 验证

发布中断时旧manifest仍完整；回滚只切指针；旧资源未过早生命周期删除。

### 失败分支

固定app.js覆盖会让边缘节点混合版本；使用内容寻址key。

### 恢复要求

失败后以数据库 upload/job 状态和对象 version 为准，重复任务使用 generation/条件更新；孤儿对象由清单与生命周期受控回收，不能根据用户文件名删除。

## 15. 故障注入矩阵

|注入|预期|禁止|
|---|---|---|
|签名后撤权|`multipart_initiated` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|上传中断|`incomplete_bytes` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|同key并发写|`part_retries` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|checksum错误|`complete_failures` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|扫描器崩溃|`abort_failures` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|生命周期延迟|`noncurrent_storage` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|对象存储429/503|`lifecycle_actions` 可观察且状态可恢复|越权发布、静默损坏或无限重试|
|数据库提交后断线|`cdn_hit_rate` 可观察且状态可恢复|越权发布、静默损坏或无限重试|

## 16. 调试与观测

1. `multipart_initiated`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
2. `incomplete_bytes`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
3. `part_retries`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
4. `complete_failures`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
5. `abort_failures`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
6. `noncurrent_storage`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
7. `lifecycle_actions`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
8. `cdn_hit_rate`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
9. `origin_bytes`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。
10. `checksum_mismatch`：按环境、操作和低基数结果分类，定义单位、采样点、SLO与告警窗口。

按 upload ID/object ID 从业务记录、签名、对象metadata/version、扫描任务到下载审计逐跳核对；签名query和敏感正文不进入普通日志。

## 17. S3兼容实现边界

1. Multipart 完成与版本发布依赖的一致性模型必须以实际对象存储产品和部署版本的官方文档/集成测试确认；S3 API兼容不自动表示语义完全一致。
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

1. bucket/key由服务端控制；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
2. 所有业务读取重新授权；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
3. 上传有大小/速率/checksum；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
4. 未验证对象隔离；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
5. 版本和删除策略明确；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
6. multipart孤儿可清理；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
7. 凭据最小权限且可轮换；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
8. 敏感URL不进日志；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
9. 故障有重试预算与幂等；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。
10. 备份/恢复和隐私删除已演练；Multipart Upload、断点续传、版本、生命周期、CDN 与 ETag 的负责人提供可重复验证证据。

## 19. 综合练习与验收

实现“50GB视频上传”并用“静态资源发布与回滚”验证不同约束。提交状态机、策略、对象清单、失败注入和观测面板。

- [ ] CreateMultipartUpload 的正常、边界、权限和失败路径均通过。
- [ ] UploadPart 的正常、边界、权限和失败路径均通过。
- [ ] CompleteMultipartUpload 的正常、边界、权限和失败路径均通过。
- [ ] AbortMultipartUpload 的正常、边界、权限和失败路径均通过。
- [ ] 断点续传状态 的正常、边界、权限和失败路径均通过。
- [ ] Versioning 的正常、边界、权限和失败路径均通过。
- [ ] Lifecycle 的正常、边界、权限和失败路径均通过。
- [ ] CDN 的正常、边界、权限和失败路径均通过。
- [ ] 两个完整案例都可在隔离环境重复运行。
- [ ] 对象存储故障不改变数据库业务不变量。
- [ ] 所有孤儿、版本、parts和审计记录有保留/清理策略。

## 来源

- [AWS S3 multipart upload](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html)（访问日期：2026-07-17）
- [AWS S3 Versioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Versioning.html)（访问日期：2026-07-17）
- [AWS S3 Lifecycle examples](https://docs.aws.amazon.com/AmazonS3/latest/userguide/lifecycle-configuration-examples.html)（访问日期：2026-07-17）
- [AWS S3 object metadata and ETag](https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingMetadata.html)（访问日期：2026-07-17）
- [RFC 9110 entity tags](https://www.rfc-editor.org/rfc/rfc9110.html#name-etag)（访问日期：2026-07-17）
