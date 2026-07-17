# 统一 API 错误模型与可恢复契约

错误契约把失败从人类可读字符串转换为客户端可分类、可恢复、可追踪的稳定数据。它必须同时表达 HTTP 层结果、领域问题、具体发生实例和字段位置，但不能泄露服务端内部实现。

## 1. 错误的四层信息

```mermaid
flowchart LR
    A["HTTP 状态码：通用分类"] --> B["问题类型：稳定机器语义"]
    B --> C["发生实例：本次请求细节"]
    C --> D["追踪标识：服务端定位"]
```

- HTTP 状态码让浏览器、代理和通用客户端知道成功、认证、限流或服务故障类别。
- 问题类型是长期稳定的机器标识，例如 `https://api.example.com/problems/insufficient-stock`。
- 发生实例描述本次冲突涉及哪个资源或字段。
- `trace_id`/`request_id` 用于关联服务端观测数据，不等于向客户端暴露堆栈。

客户端业务分支应读取问题类型或明确扩展字段，不应解析 `detail` 文案。文案可能本地化或改写。

## 2. RFC 9457 Problem Details

JSON 使用媒体类型 `application/problem+json`。标准成员如下：

| 成员 | 类型与默认 | 作用 | 使用边界 |
|---|---|---|---|
| `type` | URI 引用；缺失等同 `about:blank` | 标识问题类型 | 自定义类型应稳定并最好可解析到文档 |
| `status` | JSON 数字；可选 | 记录本次 HTTP 状态 | 只是提示，必须与实际响应状态一致 |
| `title` | 字符串；可选 | 问题类型的人类摘要 | 除本地化外不随实例变化 |
| `detail` | 字符串；可选 | 本次问题的人类说明 | 帮助纠正请求，不放堆栈或 SQL |
| `instance` | URI 引用；可选 | 标识本次发生 | 对客户端应视为不透明标识 |

扩展成员可以承载机器需要的数据，如 `errors`、`current_version`、`retryable`。扩展名应从字母开始，最好至少三个字符，并只用字母、数字和下划线，以便跨格式处理。

```json
{
  "type": "https://api.example.com/problems/validation-error",
  "title": "Request validation failed",
  "status": 422,
  "detail": "Two fields are invalid.",
  "instance": "/problem-instances/01JZ8T4Q",
  "trace_id": "7e32f6a14c7340a7",
  "errors": [
    {"pointer": "/items/0/quantity", "code": "minimum", "minimum": 1},
    {"pointer": "/shipping/address_id", "code": "not_found"}
  ]
}
```

字段位置可采用 JSON Pointer。错误对象要给机器稳定的 `code` 和相关参数，显示文案由客户端或服务端的本地化层生成。

## 3. 分类与状态映射

同一个内部错误不应直接映射；先根据客户端可采取的动作分类：

| 情况 | 状态 | 稳定问题类型 | 客户端动作 |
|---|---:|---|---|
| JSON 无法解析 | 400 | malformed-request | 修正文档语法 |
| 字段违反业务约束 | 422 | validation-error | 标注字段并修改 |
| 未提供或 token 无效 | 401 | invalid-credentials | 登录或刷新凭据 |
| 已认证但策略拒绝 | 403 | forbidden | 停止或申请权限 |
| 资源不存在/不可发现 | 404 | resource-not-found | 刷新列表 |
| 唯一键或状态冲突 | 409 | conflict | 修改输入或刷新状态 |
| 条件更新版本过期 | 412 | stale-version | 重新读取后合并 |
| 请求过多 | 429 | rate-limit-exceeded | 按 `Retry-After` 等待 |
| 依赖暂时失败 | 503 | dependency-unavailable | 有界退避重试 |
| 未知服务端故障 | 500 | internal-error | 报告 instance，不盲目重试写请求 |

不要用 409 表示任何验证错误，也不要把数据库错误码原样暴露。唯一约束冲突可以映射 409，但应转成领域语义，例如“邮箱已注册”，并避免泄露一个敏感账号是否存在。

## 4. 可重试性不是状态码的简单函数

重试决策至少依赖操作是否安全/幂等、服务端是否可能已执行、问题类型、`Retry-After` 和重试预算。

- GET 遇到暂时 503 通常可退避重试。
- 带幂等键的创建请求超时后，可用同一键重试并查询已有结果。
- 未做幂等保护的支付 POST 超时，执行结果未知，直接重试可能重复扣款。
- 422 是输入确定性错误，原样重试无意义。
- 429/503 可给 `Retry-After`，其值可以是秒数或 HTTP 日期。

服务端可返回 `retryable` 作为领域提示，但客户端仍必须受总 deadline、最大次数、指数退避和随机抖动约束。所有客户端同时按固定间隔重试会形成同步重试风暴。

## 5. Go 1.26 的领域错误映射

领域层返回结构化类别，HTTP 层统一编码，handler 不应在每个分支自行拼 JSON：

```go
type Problem struct {
    Type     string       `json:"type"`
    Title    string       `json:"title"`
    Status   int          `json:"status"`
    Detail   string       `json:"detail,omitempty"`
    Instance string       `json:"instance,omitempty"`
    TraceID  string       `json:"trace_id,omitempty"`
    Errors   []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
    Pointer string `json:"pointer"`
    Code    string `json:"code"`
}

func writeProblem(w http.ResponseWriter, p Problem) {
    w.Header().Set("Content-Type", "application/problem+json")
    w.Header().Set("Cache-Control", "no-store")
    w.WriteHeader(p.Status)
    _ = json.NewEncoder(w).Encode(p)
}
```

必须先设置响应头和状态再写正文；第一次 `Write` 会隐式发送 200。编码失败时响应可能已部分写出，不能再可靠更换状态，因此问题结构应保持可编码，并在测试中覆盖。

领域错误使用 `errors.Is/As` 识别，而不是比较字符串：

```go
func problemFor(err error, base, traceID string, logger *slog.Logger) Problem {
    switch {
    case errors.Is(err, ErrOrderNotFound):
        return Problem{Type: base + "/resource-not-found", Title: "Resource not found", Status: 404}
    case errors.Is(err, ErrVersionConflict):
        return Problem{Type: base + "/stale-version", Title: "Resource version is stale", Status: 412}
    default:
        logger.Error("request failed", "trace_id", traceID, "error", err)
        return Problem{Type: base + "/internal-error", Title: "Internal server error", Status: 500, TraceID: traceID}
    }
}
```

外部响应隐藏内部错误，完整 cause 只进入受权限保护的日志/trace。日志也要脱敏，不能记录密码、Cookie、Authorization、身份证号或完整支付数据。

## 6. 批量与部分成功

批量请求有三种不同契约，必须提前选择：

1. 原子批量：任一项失败则全部回滚，返回最相关的问题。
2. 独立批量：每项都有成功/失败结果，整体 200 只表示批处理协议执行完成。
3. 异步批量：返回 202 和任务资源，任务内记录各项结果。

不要在普通单资源接口返回多个无关问题类型。字段验证可以在同一 `validation-error` 类型的 `errors` 扩展中列出多个位置，因为它们共享同一问题语义。

## 7. 安全边界

错误差异可能形成侧信道。例如登录时“账号不存在”和“密码错误”会帮助枚举账号；重置密码接口宜对存在与不存在的邮箱返回相同外部结果。对资源权限可在需要隐藏存在性时统一返回 404。

不得返回：

- SQL 文本、表名、约束内部名和连接字符串；
- 文件系统路径、源代码行、堆栈和依赖版本；
- token、Cookie、密钥、内部服务地址；
- 可推断其他租户资源存在性的 ID 或计数；
- 客户端无法安全展示的未转义用户输入。

`detail` 是数据，不是可信 HTML。前端用文本节点显示，不能直接插入 `innerHTML`。

## 8. 完整案例：创建订单

### 输入

请求包含两个商品，其中一项数量为 0；随后客户端修正数量，但库存已不足；第三次请求提交成功后响应在网络中丢失。每次请求使用同一业务意图对应的幂等键。

### 步骤

1. 网关生成 trace ID，验证正文大小和 JSON 语法。
2. handler 校验字段，第一次返回 422，`errors[0].pointer=/items/1/quantity`。
3. 客户端修正后重试；事务锁定库存，发现不足，返回 409 `insufficient-stock` 和可用数量。
4. 客户端减少数量并以新业务意图的新幂等键提交；服务创建订单并缓存状态/响应摘要。
5. 响应丢失后，客户端用同一幂等键重试，服务返回原订单而不是再创建。

### 输出

```http
HTTP/1.1 409 Conflict
Content-Type: application/problem+json

{"type":"https://api.example.com/problems/insufficient-stock","title":"Insufficient stock","status":409,"instance":"/problem-instances/p_91","item_id":"sku_8","available":2,"trace_id":"01JZ8T4Q"}
```

### 验证

- JSON Schema 验证每种问题正文。
- 同一 `type` 的 `title` 除语言外保持不变。
- 实际 HTTP 状态与 `status` 一致。
- 日志可按 trace ID 找到内部 cause，外部正文不存在内部表名。
- 同一幂等键重放返回同一个订单 ID。

### 失败分支

若数据库唯一约束错误直接返回 `duplicate key value violates unique constraint orders_idempotency_key_key`，既泄露内部结构，也没有告诉客户端可恢复动作。修正为识别约束对应的领域冲突，查询已保存结果并返回原响应；其他未知数据库错误统一转为 500，内部保留 cause。

## 9. 契约测试与观测

对每个公开问题类型维护：状态码、必需扩展、可重试性、是否会暴露资源存在性、客户端行为。测试至少包括：

```text
given: 无效 quantity 与固定 trace ID
when:  POST /orders
then:  HTTP 422
and:   Content-Type = application/problem+json
and:   type 以稳定 URI 表示 validation-error
and:   errors 含 /items/0/quantity
and:   响应不含 SQL、stack、Authorization
```

指标按低基数问题类型和规范化路由聚合。不要把完整 instance、用户 ID 或自由文本 detail 作为指标 label，否则时序库基数失控。高基数 trace ID 进入日志和 trace。

## 10. 常见错误

- 所有错误都用 400：客户端无法区分认证、冲突、限流和服务故障。
- 返回 200 加 `success:false`：代理和通用监控把失败当成功。
- 让客户端解析中文 `message`：改文案即破坏程序。
- `status` 字段与实际状态不一致：中间件和客户端产生不同判断。
- 对 500 自动无限重试：放大故障，并可能重复执行非幂等写入。
- 将请求 ID 当认证凭据：追踪标识可被看到，不具备授权能力。
- 为了统一格式吞掉 `WWW-Authenticate`、`Allow`、`Retry-After`：丢失标准 HTTP 恢复信息。

## 11. 练习

为注册、登录、更新资料、提交订单定义错误目录。每个问题类型写出状态码、扩展成员、客户端动作、日志内容和安全风险，并实现统一 Go encoder 的表驱动测试。

完成标准：至少覆盖 400/401/403/404/409/412/422/429/500/503；客户端不解析 detail；敏感账号不可枚举；可重试写操作有幂等保护；每个响应能用 trace ID 在服务端定位。

## 来源

- [RFC 9457: Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457.html)（访问日期：2026-07-17）
- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110.html)（访问日期：2026-07-17）
- [RFC 6901: JSON Pointer](https://www.rfc-editor.org/rfc/rfc6901.html)（访问日期：2026-07-17）
- [Go 1.26 Package errors](https://pkg.go.dev/errors)（访问日期：2026-07-17）
