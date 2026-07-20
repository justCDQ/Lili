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

Multipart Upload 的恢复事实是服务端保存的 upload ID 与 part 清单，发布事实是最终对象版本和 manifest 指针。ETag 在分段上传、加密或不同实现下不一定是内容 MD5，完整性必须使用明确的 checksum；CDN 缓存也不能代替源站版本记录。

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

断线恢复先用 upload ID 分页执行 `ListParts`，将服务端 part number、ETag/checksum 与本地清单比较，只重传缺失或校验不符的 part。完成请求必须提交有序清单并以 generation 条件更新一次；过期 generation 不得发布，取消后主动 abort，生命周期只负责兜底回收未完成 parts。

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

发布中断时保持旧 manifest 指针不变，已经上传的内容 hash 资源可安全复用。确认新版本全部可从源站读取后再切换指针；回滚只恢复旧 manifest。清理任务从仍被 manifest 引用的版本集合计算保留集，不能按上传时间直接删除旧资源。

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

大文件故障按业务 generation 和 upload ID 检查 `ListParts` 分页、重复 part、完成清单、最终 object version 与 checksum；静态发布故障再核对 manifest 指针、源站响应、CDN cache key、`Age` 和命中层级。先区分源站对象错误与边缘陈旧缓存，再决定续传、回滚或失效缓存。

## 17. S3兼容实现边界

1. 固化最小/最大 part 大小、最大 part 数与最后一段例外，按这些限制计算分段策略。
2. `ListParts` 可能分页；续传测试必须覆盖超过单页的上传，不能把首屏结果当完整清单。
3. 验证重复上传同一 part number 的覆盖规则，以及完成清单对顺序、ETag 和 checksum 的要求。
4. 明确 `CompleteMultipartUpload` 超时后的判定方式：先查询最终对象，不能直接创建第二个上传。
5. ETag 是否能作为 MD5 取决于分段与加密方式；跨产品校验统一使用显式 checksum。
6. 测试 versioning 开启后的覆盖、删除、delete marker 和指定 version 读取，保留业务发布的 version ID。
7. 生命周期规则分别验证未完成 multipart、当前版本和非当前版本；执行是异步的，不能承担同步取消。
8. CDN cache key 要明确是否包含 query、Host、压缩协商与选定 headers，避免不同资源误共享缓存。
9. 对 HTML/manifest 使用短缓存或条件请求，对内容 hash 资源使用 immutable；分别验证 `ETag`/`Last-Modified` 行为。
10. 演练 manifest 切换、边缘缓存未刷新和旧版本回滚，确保生命周期不会删除仍可被回滚指针引用的对象。

## 18. 生产检查

1. part 大小、并发数和最大 part 数由文件大小与产品限制计算，并设客户端内存/带宽上限。
2. 续传以服务端 `ListParts` 完整分页结果为准，缺失与 checksum 不符的 part 才重传。
3. complete 请求使用有序清单，超时后先 HEAD 最终对象，避免盲目创建第二个 multipart。
4. 最终对象保存 version 与显式 checksum，不把 multipart ETag 解释为内容 MD5。
5. 取消路径主动 abort，未完成 multipart 的数量、字节和最老年龄都有告警。
6. 生命周期分别覆盖未完成上传、非当前版本和孤儿 hash 资源，并验证不会删除回滚依赖。
7. 静态资源以内容 hash key 不可变发布，HTML/manifest 是唯一可切换发布指针。
8. CDN cache key、TTL、条件请求与压缩变体有集成测试，能区分 edge miss 与 origin error。
9. 发布中断保留旧 manifest；回滚测试证明无需覆盖资源或大范围 purge。
10. 版本恢复、delete marker、生命周期延迟与 CDN 陈旧缓存纳入定期故障演练。

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
