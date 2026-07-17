# HTTP 请求响应、方法、状态码、Header、JSON 与 Cookie

## 学习目标

本文说明 HTTP 消息的语义结构、常用方法与状态码、表示元数据、JSON 数值边界和 Cookie 属性，并实现一个严格验证输入的 Go HTTP 创建接口。

## 1. HTTP 消息与资源

HTTP 是无状态的请求/响应协议。请求目标标识资源，方法表达对资源的语义；响应状态码表达本次请求结果。连接可复用，但服务器不能仅因两个请求使用同一连接就假设来自同一用户。

概念请求包含方法、目标 URI、协议字段和可选内容；响应包含状态码、字段和可选内容。HTTP/1.1、HTTP/2、HTTP/3 在线路 framing 上不同，但共享 RFC 9110 定义的核心语义。

```text
POST /v1/tasks HTTP/1.1
Host: api.example.test
Content-Type: application/json
Accept: application/json
Content-Length: 24

{"title":"read RFC 9110"}
```

应用应使用协议库，不手工按换行解析 HTTP。消息 framing、重复字段、连接管理和版本差异包含安全边界。

## 2. 方法语义

| 方法 | 典型语义 | 安全 | 幂等 |
| --- | --- | --- | --- |
| GET | 获取目标资源当前表示 | 是 | 是 |
| HEAD | 与 GET 相同但响应无内容 | 是 | 是 |
| POST | 按资源语义处理所附表示 | 否 | 不保证 |
| PUT | 用请求表示创建或替换目标状态 | 否 | 是 |
| PATCH | 按补丁语义部分修改 | 否 | 取决于补丁格式/操作 |
| DELETE | 请求删除目标资源关联 | 否 | 是 |
| OPTIONS | 获取通信选项 | 是 | 是 |

安全表示客户端不请求改变服务器状态；日志和计费等附带影响不改变该定义。幂等表示重复相同请求的预期服务器状态效果与一次相同，不代表响应码、时间或日志完全相同。

GET 不应用于执行转账或删除。爬虫、预取和缓存可按安全语义自动发起 GET。POST 默认不幂等，网络失败后盲目重试可能重复创建；可通过业务幂等键建立安全重试契约。

## 3. 常用状态码

状态码第一位划分类别：1xx 信息，2xx 成功，3xx 重定向，4xx 客户端可归因问题，5xx 服务器未履约。具体状态必须按定义选择。

| 状态 | 用途 | 常见边界 |
| --- | --- | --- |
| 200 OK | 请求成功并返回表示 | 不用于所有创建场景 |
| 201 Created | 已创建资源 | 通常用 Location 指向新资源 |
| 204 No Content | 成功且无响应内容 | 不发送消息正文 |
| 400 Bad Request | 请求语法或通用格式非法 | 与字段语义错误需稳定约定 |
| 401 Unauthorized | 缺少/无效认证凭据 | 语义实际是未认证 |
| 403 Forbidden | 已理解但拒绝授权 | 不必泄露资源是否存在 |
| 404 Not Found | 未找到或不愿披露 | 可用于权限隐藏策略 |
| 405 Method Not Allowed | 资源不支持该方法 | 响应应含 Allow |
| 409 Conflict | 与目标当前状态冲突 | 如唯一键或版本冲突 |
| 413 Content Too Large | 请求内容超限 | 限制应尽早执行 |
| 415 Unsupported Media Type | 内容媒体类型不支持 | 不等于 JSON 语法错误 |
| 422 Unprocessable Content | 语法可处理但指令语义非法 | RFC 9110 使用此名称 |
| 500 Internal Server Error | 未分类服务端失败 | 对外不泄露内部栈 |
| 503 Service Unavailable | 临时不可服务 | 可结合 Retry-After |

返回 200 但错误写在 JSON `success:false` 会破坏代理、监控和通用客户端语义。反过来，4xx/5xx 仍应提供稳定的机器可读错误结构。

## 4. Header 字段

字段名大小写不敏感；字段值语法由字段定义决定。不是所有重复字段都能用逗号合并，`Set-Cookie` 是重要特例。应用不要自行把所有同名字段拼接。

常用字段：

- `Content-Type` 描述当前消息内容的媒体类型与参数。
- `Accept` 表示客户端可接受响应媒体类型。
- `Content-Length` 是 framing 相关长度；库通常处理。
- `Authorization` 携带认证凭据，禁止完整记录。
- `Location` 在创建或重定向响应中指示 URI。
- `ETag` 与条件请求可用于缓存验证和并发控制。
- `Cache-Control` 指示缓存行为，不能只用自定义字段代替。

服务器必须设置字段和状态后再写正文。在 Go 中第一次 `Write` 若未调用 `WriteHeader` 会隐式发送 200；随后再设置 Header 或状态已太晚。

## 5. 媒体类型与内容协商

JSON 常用媒体类型 `application/json`。比较 Content-Type 时使用媒体类型解析器，因为参数、大小写与空白有语法规则；不要仅做完整字符串相等。

没有请求内容时 Content-Type 可缺失；有 JSON 正文时服务器可要求 `application/json`，不支持时返回 415。`Accept` 不接受可生成格式时可返回 406，但实际 API 要清晰规定默认值和兼容范围。

字符编码对 `application/json` 由 JSON 规范定义互操作要求，不应随意接受非标准 charset 转码后当同一签名内容。

## 6. JSON 数据模型与边界

JSON 值包括 object、array、number、string、true、false、null。对象成员名是字符串；规范指出对象成员名应唯一以获得可互操作行为，但不同实现遇到重复名可能取前、取后或报告错误。

JSON number 语法不规定 IEEE 754，也没有单独整数类型。跨 JavaScript、Go 和数据库传递大整数时，必须约定范围或用十进制字符串。Go `encoding/json` 解码到 `any` 默认用 float64；解到有类型 struct 或使用 `Decoder.UseNumber` 可避免部分意外，但仍需范围验证。

字符串可以转义 Unicode，但不代表业务文本已规范化。解析成功后仍要验证长度、允许字符和规范化策略。

严格边界通常执行：限制总字节；检查媒体类型；解码一个顶层值；拒绝未知字段；拒绝尾随第二个值；验证必填、范围和跨字段不变量。

## 7. Cookie

服务器通过 `Set-Cookie` 响应字段让用户代理保存 cookie；后续匹配请求通过 `Cookie` 字段回送。Cookie 是客户端提供的状态，必须验证、签名或仅存不透明会话 ID，不能因“服务器曾设置”就信任内容。

关键属性：

- `Secure`：只在安全传输上发送；本地开发策略要单独处理。
- `HttpOnly`：禁止通过浏览器非 HTTP API 访问，降低脚本窃取风险，但不阻止请求携带。
- `SameSite`：限制跨站请求携带，值与默认行为要按当前用户代理规范验证。
- `Path`/`Domain`：决定发送范围，不是访问控制边界。
- `Max-Age`/`Expires`：持久期；服务端会话仍需独立失效机制。

删除 cookie 通常用相同名称、Path/Domain 并设置过去 Expires 或非正 Max-Age。属性不匹配可能留下另一个同名 cookie。

Cookie 容量和数量受用户代理限制，不适合保存大型状态。认证 cookie 还需要 CSRF 策略、会话轮换、撤销和 TLS。

## 8. 完整案例：创建任务接口

### 8.1 契约

`POST /tasks` 接受最多 64 KiB 的 `application/json`：

```json
{"title":"read RFC 9110","priority":2}
```

`title` 去除首尾空白后 1–100 个 Unicode 码点；priority 为 1–3。成功返回 201、Location 与任务 JSON。失败统一返回 `{code,message}` JSON。

### 8.2 Handler

```go
package taskapi

import (
    "encoding/json"
    "errors"
    "io"
    "mime"
    "net/http"
    "strings"
    "unicode/utf8"
)

type createInput struct {
    Title    string `json:"title"`
    Priority int    `json:"priority"`
}

type problem struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(value)
}

func Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        writeJSON(w, http.StatusMethodNotAllowed, problem{"METHOD_NOT_ALLOWED", "use POST"})
        return
    }
    mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
    if err != nil || mediaType != "application/json" {
        writeJSON(w, http.StatusUnsupportedMediaType, problem{"UNSUPPORTED_MEDIA_TYPE", "use application/json"})
        return
    }

    r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields()
    var input createInput
    if err := decoder.Decode(&input); err != nil {
        writeJSON(w, http.StatusBadRequest, problem{"INVALID_JSON", "request JSON is invalid"})
        return
    }
    var extra any
    if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
        writeJSON(w, http.StatusBadRequest, problem{"INVALID_JSON", "one JSON value is required"})
        return
    }

    input.Title = strings.TrimSpace(input.Title)
    length := utf8.RuneCountInString(input.Title)
    if length < 1 || length > 100 || input.Priority < 1 || input.Priority > 3 {
        writeJSON(w, http.StatusUnprocessableEntity, problem{"INVALID_TASK", "title or priority is invalid"})
        return
    }
    output := struct {
        ID       string `json:"id"`
        Title    string `json:"title"`
        Priority int    `json:"priority"`
    }{"task-1", input.Title, input.Priority}
    w.Header().Set("Location", "/tasks/task-1")
    writeJSON(w, http.StatusCreated, output)
}
```

真实存储生成 ID 时要在事务/存储层完成，并处理唯一约束。示例固定 ID 仅用于协议验证。

### 8.3 正常测试

```go
func TestCreate(t *testing.T) {
    request := httptest.NewRequest(http.MethodPost, "/tasks",
        strings.NewReader(`{"title":"read RFC 9110","priority":2}`))
    request.Header.Set("Content-Type", "application/json; charset=utf-8")
    response := httptest.NewRecorder()
    Create(response, request)

    if response.Code != http.StatusCreated { t.Fatalf("status=%d", response.Code) }
    if got := response.Header().Get("Location"); got != "/tasks/task-1" {
        t.Fatalf("Location=%q", got)
    }
    var body map[string]any
    if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil { t.Fatal(err) }
    if body["id"] != "task-1" { t.Fatalf("body=%v", body) }
}
```

输入经过媒体类型解析、限流读取、严格解码、尾随检查和语义验证，输出 201；测试验证状态、Location 与可解析正文，不只匹配字符串空白。

### 8.4 失败分支

| 输入 | 预期 |
| --- | --- |
| GET `/tasks` | 405，Allow: POST |
| Content-Type `text/plain` | 415 |
| `{"title":` | 400 INVALID_JSON |
| 含未知字段 `prio` | 400 INVALID_JSON |
| 两个连续 JSON 对象 | 400 INVALID_JSON |
| title 空白或 priority 4 | 422 INVALID_TASK |
| 正文超过 64 KiB | 400 当前示例；可细分 413 |

`MaxBytesReader` 的超限错误可用类型判断映射为 413；生产代码应实现该稳定分支。任何内部编码/存储错误要返回 500 并记录内部错误，不能把数据库文本回传客户端。

仓库中的[可运行 Task API 示例](../../examples/service-data-basics/taskapi/)保存了正常响应与方法、媒体类型、语法、未知字段、语义失败测试。

## 9. Cookie 示例

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    opaqueSessionID,
    Path:     "/",
    Secure:   true,
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   3600,
})
```

`opaqueSessionID` 应不可预测，服务端保存会话状态并验证过期/撤销。日志只记录经过设计的会话关联 ID，不能记录 cookie Value。

## 10. 调试清单

- Handler 返回 200 而预期 4xx：检查是否在 `WriteHeader` 前已写正文。
- JSON 多余字段未报错：是否使用 typed struct 与 `DisallowUnknownFields`。
- 大整数改变：是否先解码到 float64，改为明确字段类型或字符串协议。
- Cookie 未发送：检查 Secure、Domain、Path、SameSite、过期和请求 scheme。
- 重复 Set-Cookie 丢失：不要按普通逗号列表合并该字段。
- POST 重试产生重复：引入幂等键和存储唯一约束，不只客户端去重。
- 代理后状态异常：保留端到端状态/字段证据，检查代理改写和大小限制。

## 11. 练习

1. 为 Create 写表驱动测试覆盖表中所有失败分支，并断言 JSON 错误码。
2. 识别 `MaxBytesError` 并返回 413，保证超限测试稳定。
3. 添加 GET `/tasks/{id}` 和 ETag 条件请求，明确 304 无正文语义。
4. 设计 POST 幂等键表和冲突响应，模拟客户端超时后重试。
5. 用 httptest CookieJar 验证 Secure、Path 和删除属性行为。

## 来源

- [RFC 9110：HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)（访问日期：2026-07-17）
- [RFC 8259：The JavaScript Object Notation Data Interchange Format](https://www.rfc-editor.org/rfc/rfc8259)（访问日期：2026-07-17）
- [RFC 6265：HTTP State Management Mechanism](https://www.rfc-editor.org/rfc/rfc6265)（访问日期：2026-07-17）
- [Go 标准库：encoding/json](https://pkg.go.dev/encoding/json)（访问日期：2026-07-17）
- [Go 标准库：net/http](https://pkg.go.dev/net/http)（访问日期：2026-07-17）
