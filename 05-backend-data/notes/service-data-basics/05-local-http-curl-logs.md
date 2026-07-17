# 本地运行 HTTP 服务、使用 curl 请求与查看日志

## 学习目标

本文把前四篇概念组成一个可操作闭环：启动有限暴露面的本地服务，用 curl 精确构造和检查请求，以结构化日志关联状态与耗时，并按 DNS、连接、HTTP、应用顺序诊断失败。

## 1. 本地服务的最小运行条件

启动前要知道可执行程序、工作目录、配置、监听地址、端口和关闭方式。开发服务绑定 `127.0.0.1` 可降低意外暴露；绑定 `0.0.0.0` 会监听所有本地 IPv4 接口，不应在不可信网络默认使用。

```sh
go run ./cmd/service
```

`go run` 会编译并启动临时二进制，适合开发。调试退出码和信号时更可靠的是先 `go build` 再运行，因为 go 命令本身可能包装子进程状态。

服务“启动日志已打印”不等于监听已成功，日志应在 `net.Listen` 成功后记录实际地址。固定端口冲突要明确失败，测试用端口 0 由系统分配。

## 2. curl 的请求模型

curl 是 URL 数据传输工具。它的默认输出通常是响应正文，进度与诊断写 stderr。自动化要分别捕获正文、Header、状态、时间和进程退出码。

```sh
curl --silent --show-error \
  --output response.json \
  --dump-header response.headers \
  --write-out '%{http_code}\n' \
  http://127.0.0.1:8080/health
```

`--include` 把响应 Header 放到 stdout 与正文一起，适合人工观察但不利于把正文直接交 JSON 解析。`--dump-header` 分离保存。`--verbose` 把连接与协议诊断写 stderr，可能包含敏感 Header，分享前脱敏。

`--fail-with-body` 在 HTTP 4xx/5xx 让 curl 返回失败，同时保留正文；`--fail` 可能不显示正文。curl 成功退出通常只证明传输按其规则完成，HTTP 500 默认不一定是进程失败，脚本必须选择相应选项并检查状态。

## 3. 构造方法、Header 与正文

```sh
curl --silent --show-error --fail-with-body \
  --request POST \
  --header 'Content-Type: application/json' \
  --header 'Accept: application/json' \
  --data-binary @task.json \
  http://127.0.0.1:8080/tasks
```

`--data` 会隐式选择 POST 并有数据处理规则，`--data-binary @file` 更适合按文件字节发送。`--json @file` 是 curl 提供的便利选项，会设置相关 JSON Header，最低支持版本要在团队环境确认。

不要把令牌直接写命令行或分享完整 shell history。可从受限文件读取 Header，或使用环境/秘密注入后确保 `set -x` 关闭；进程参数仍可能被本机其他工具看到。

URL 查询参数用 `--get --data-urlencode` 等安全编码方式，避免手拼空格、`&` 和 Unicode。Shell 中 URL 含 `&` 时必须引用，否则会被解释为后台控制符。

## 4. 超时与重试

curl `--connect-timeout` 限制连接阶段，`--max-time` 限制整体操作。两者处理不同故障。

```sh
curl --connect-timeout 2 --max-time 5 http://127.0.0.1:8080/health
```

重试选项要谨慎：只有方法与业务操作可安全重试时启用，设置最大次数和总时间。POST 创建在响应丢失时可能已执行，客户端重试会重复创建；需服务端幂等键。

本地测试也应设置超时，避免 CI 永久挂住。超时后服务端可能继续处理，handler 必须观察请求 context，并让下游调用可取消。

## 5. 日志的最小结构

请求日志应能关联一次请求，至少包含时间、severity、事件名、request_id、方法、路由模板、状态、响应字节和持续时间。路径中可能含用户数据；应记录低基数路由模板用于聚合，并按隐私策略处理实际目标。

```json
{
  "time":"2026-07-17T02:30:00Z",
  "severity":"INFO",
  "event":"http_request_completed",
  "request_id":"req-7f3a",
  "method":"GET",
  "route":"/health",
  "status":200,
  "duration_ms":1.7,
  "response_bytes":16
}
```

不记录 Authorization、Cookie、Set-Cookie、密码、令牌和完整请求正文。必要业务字段使用白名单、脱敏和访问控制。请求 ID 用于关联，不是认证凭据；若接受上游 ID，只信任受控代理或校验格式并另存内部 trace ID。

日志严重性要稳定。正常 404 是否 WARN 取决于服务语义；健康探针 404 可能配置错误，公网扫描 404 可能常态。用指标统计总体率，用日志保留具体上下文。

## 6. 状态和耗时捕获

Go 中包装 `ResponseWriter` 可记录第一个状态和写入字节。接口可能包含 Flusher、Hijacker、Pusher 等扩展；简单 wrapper 若不转发会改变流式、WebSocket 或 HTTP/2 行为。入门服务只支持普通响应时应明确限制。

```go
type captureWriter struct {
    http.ResponseWriter
    status int
    bytes  int
}

func (w *captureWriter) WriteHeader(status int) {
    if w.status != 0 { return }
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

func (w *captureWriter) Write(data []byte) (int, error) {
    if w.status == 0 { w.WriteHeader(http.StatusOK) }
    count, err := w.ResponseWriter.Write(data)
    w.bytes += count
    return count, err
}
```

重复 WriteHeader 只以第一次为准，wrapper 也保持该语义。写入错误可能在客户端断开时出现，handler 应处理能处理的错误；请求完成日志记录实际已写字节。

## 7. 完整服务

### 7.1 Handler 与日志中间件

```go
package localhttp

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "time"
)

type captureWriter struct {
    http.ResponseWriter
    status int
    bytes  int
}

func (w *captureWriter) WriteHeader(status int) {
    if w.status != 0 { return }
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

func (w *captureWriter) Write(data []byte) (int, error) {
    if w.status == 0 { w.WriteHeader(http.StatusOK) }
    n, err := w.ResponseWriter.Write(data)
    w.bytes += n
    return n, err
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        captured := &captureWriter{ResponseWriter: w}
        next.ServeHTTP(captured, r)
        status := captured.status
        if status == 0 { status = http.StatusOK }
        logger.Info("http request completed",
            "event", "http_request_completed",
            "method", r.Method,
            "path", r.URL.Path,
            "status", status,
            "response_bytes", captured.bytes,
            "duration", time.Since(start),
        )
    })
}

func Handler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })
    mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
        select {
        case <-time.After(200 * time.Millisecond):
            w.Header().Set("Content-Type", "application/json")
            _ = json.NewEncoder(w).Encode(map[string]string{"status": "finished"})
        case <-r.Context().Done():
            return
        }
    })
    return mux
}
```

生产服务不应为每次 slow 调用创建不受控计时器负载；案例用于取消演示。日志 middleware 记录实际路径，真实高流量服务应使用路由模板和请求 ID，并处理 panic 恢复的日志策略。

### 7.2 启动程序

```go
func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    server := &http.Server{
        Addr:              "127.0.0.1:8080",
        Handler:           Logging(logger, Handler()),
        ReadHeaderTimeout: 5 * time.Second,
        IdleTimeout:       30 * time.Second,
    }
    logger.Info("server listening", "address", server.Addr)
    if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
        logger.Error("server failed", "error", err)
        os.Exit(1)
    }
}
```

严格实现应先 net.Listen 成功，再记录 listener.Addr；上例固定地址日志在端口冲突时会先打印“listening”再失败，是需要在练习修复的可观察问题。

## 8. curl 验证闭环

正常请求：

```sh
curl --silent --show-error --fail-with-body \
  --dump-header /tmp/health.headers \
  --output /tmp/health.json \
  --write-out 'status=%{http_code} total=%{time_total}\n' \
  http://127.0.0.1:8080/health
```

预期终端输出类似 `status=200 total=0.001234`，Header 含 `Content-Type: application/json`，正文可解析且 status 字段为 ok。日志出现同一路径、200、响应字节与持续时间。

方法失败：

```sh
curl --silent --show-error --include --request POST \
  http://127.0.0.1:8080/health
```

Go ServeMux 方法模式对不匹配方法返回 405，并可能设置 Allow。curl 默认仍可能退出 0；加 `--fail-with-body` 才让 405 成为 curl 失败。

超时失败：

```sh
curl --silent --show-error --max-time 0.05 \
  http://127.0.0.1:8080/slow
echo "$?"
```

客户端约 50ms 超时并非 HTTP 状态响应；curl 返回超时相关非零码。请求 context 被取消，handler 返回不写正文。middleware 可能记录默认 200/0 bytes，因为 Go handler 没有感知并写状态；生产观测需另设 canceled outcome，不能把它算正常 200。

## 9. httptest 自动验证

```go
func TestHealth(t *testing.T) {
    var logs bytes.Buffer
    logger := slog.New(slog.NewJSONHandler(&logs, nil))
    server := httptest.NewServer(Logging(logger, Handler()))
    defer server.Close()

    client := &http.Client{Timeout: time.Second}
    response, err := client.Get(server.URL + "/health")
    if err != nil { t.Fatal(err) }
    defer response.Body.Close()
    if response.StatusCode != 200 { t.Fatalf("status=%d", response.StatusCode) }
    var body map[string]string
    if err := json.NewDecoder(response.Body).Decode(&body); err != nil { t.Fatal(err) }
    if body["status"] != "ok" { t.Fatalf("body=%v", body) }
    if !strings.Contains(logs.String(), `"event":"http_request_completed"`) {
        t.Fatalf("logs=%s", logs.String())
    }
}
```

测试使用系统临时端口，不与开发 8080 冲突。它验证真实 HTTP 客户端/服务器连接、状态、JSON 和日志事件。更精确测试应解析每行 JSON 并断言字段类型，不只 Contains。

## 10. 分层诊断流程

1. 进程：是否仍运行，启动错误与退出码是什么。
2. 监听：实际地址族、IP、端口和网络命名空间。
3. 名称：应用实际 resolver 对目标名称的 A/AAAA 结果。
4. 连接：curl verbose 的目标地址、连接耗时和拒绝/超时。
5. TLS：证书名称、信任链、版本和握手错误。
6. HTTP：方法、目标、Header、状态与媒体类型。
7. 应用：请求 ID 对应日志、业务错误码、下游耗时。
8. 数据：事务提交、约束、最终行与幂等结果。

不要一开始删除缓存、重启全部组件或关闭 TLS 验证，这会销毁证据并扩大风险。每步记录命令、时间和环境。

## 11. 日志与指标的关系

日志记录单个事件上下文；指标聚合请求率、错误率和延迟分布；trace 连接跨服务操作。高流量下不应靠逐行日志计算所有 SLO，日志采样又不能丢失必要审计事件。

OpenTelemetry Logs Data Model 定义时间、observed timestamp、severity、body、attributes、trace/span ID 等通用字段。采用时映射已有字段，避免同一含义多种名字。

## 12. 调试清单

- curl 无输出：同时看 stderr 和退出码，检查是否 `--silent` 隐藏进度但保留错误。
- 返回 HTML 而预期 JSON：可能命中代理/默认 404，检查状态、Content-Type 和 Server。
- 日志显示 200 但客户端超时：捕获取消/写错误，不把未写响应默认为成功。
- 8080 被占用：找监听进程，不随意终止不属于任务的服务。
- localhost 一台成功一台失败：比较解析地址族与服务绑定。
- curl 命令泄露令牌：清理日志并轮换凭据，不能只删除终端历史。
- 日志量激增：检查重试、健康探针、错误循环和高基数字段。

## 13. 练习

1. 改用 `net.Listen("tcp","127.0.0.1:0")` 后再打印实际地址，消除误导启动日志。
2. 为 middleware 区分 canceled、panic、写失败与正常完成 outcome。
3. 用 curl 分离保存 Header/body/status，写脚本断言 JSON 字段与退出码。
4. 给 `/slow` 添加服务端 100ms deadline，比较客户端 50ms 与 500ms 的结果。
5. 为日志添加受控 request ID，测试不记录 Authorization 与 Cookie。

## 来源

- [curl 官方手册](https://curl.se/docs/manpage.html)（访问日期：2026-07-17）
- [Go 标准库：net/http](https://pkg.go.dev/net/http)（访问日期：2026-07-17）
- [Go 标准库：net/http/httptest](https://pkg.go.dev/net/http/httptest)（访问日期：2026-07-17）
- [Go 标准库：log/slog](https://pkg.go.dev/log/slog)（访问日期：2026-07-17）
- [OpenTelemetry：Logs Data Model](https://opentelemetry.io/docs/specs/otel/logs/data-model/)（访问日期：2026-07-17）
