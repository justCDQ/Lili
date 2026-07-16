# HTTP 请求响应、方法、状态码、Header、JSON 与 Cookie

## 是什么

HTTP 请求由方法、目标、Header 和可选内容组成；响应由状态码、Header 和可选内容组成。GET 获取表示，POST 提交处理，PUT 替换目标状态，PATCH 部分修改，DELETE 请求删除。Header 携带元数据，JSON 是常用表示格式，Cookie 让用户代理在后续请求中回送服务器设置的状态。

```go
func create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { http.Error(w, "method", 405); return }
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
    var in struct { Name string `json:"name"` }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil { http.Error(w, "bad json", 400); return }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"name": in.Name})
}
```

## 常见错误与边界

- 状态码表达请求结果：4xx 是客户端可归因问题，5xx 是服务器未履约；错误体仍需稳定结构。
- GET、HEAD、OPTIONS、TRACE 语义上安全；PUT、DELETE 和安全方法是幂等语义，但实现仍需正确。
- 校验 `Content-Type`、正文大小和未知字段；JSON Number 精度跨语言需约定。
- Cookie 设置 `Secure`、`HttpOnly`、`SameSite`、Path、期限；Cookie 不是可信输入。
- Header 名不区分大小写，不能将敏感值写日志。

## 为什么需要

这些概念组成一次服务请求从寻址、传输、处理到持久化的最小闭环。缺少其中任一层的明确契约，都会让连接失败、协议错误、并发问题或数据不一致难以定位。

## 关键特性或规则

本文已有的规则、选择条件与复杂度约束共同构成判断标准。使用前必须明确输入类型、规模、资源所有权、失败语义和可观察结果；任何依赖实现细节的结论都需要测试或 Profile 验证。

## 实际怎么使用

运行本文 Go 服务或数据库示例，使用 curl 发出正常、非法方法、错误 JSON、超大正文和并发请求。逐层记录 DNS/地址、端口、请求 Header、状态码、日志、数据变化和错误恢复，并为核心处理函数添加测试。

## 补充知识

本地成功只验证单进程与本机网络条件。进入容器、代理或远端数据库后，还要显式处理超时、连接池、取消、重试、幂等、事务边界和敏感日志。

## 来源

- [RFC 9110：HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)（访问日期：2026-07-16）
- [RFC 8259：JSON](https://www.rfc-editor.org/rfc/rfc8259)（访问日期：2026-07-16）
- [RFC 6265：HTTP State Management](https://www.rfc-editor.org/rfc/rfc6265)（访问日期：2026-07-16）
