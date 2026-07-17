# Go 测试与基准：Test、表驱动测试、Fuzz 和 Benchmark

Go 的 `testing` 包与 `go test` 工具共同提供单元测试、子测试、示例测试、模糊测试和微基准。测试验证可观察契约；基准量化特定环境和输入下的性能，不代替正确性测试或生产观测。

## 测试文件与发现规则

测试文件以 `_test.go` 结尾，不参与普通包构建。测试函数签名必须是：

```go
func TestName(t *testing.T)
func BenchmarkName(b *testing.B)
func FuzzName(f *testing.F)
func ExampleName()
```

`go test` 为包构建测试二进制并运行。默认成功结果可能被缓存；需要强制重跑可用 `-count=1`。`go test ./...` 测试当前 module 下所有包。

黑盒测试可声明 `package examplesgo_test`，只能使用导出 API；白盒测试使用与被测代码相同的 package，可访问未导出标识符。优先测试公共契约，仅在复杂内部算法需要精确隔离时使用白盒测试。

## `testing.T` 的关键能力

| 方法 | 行为 | 使用边界 |
| --- | --- | --- |
| `Error/Errorf` | 标记失败并继续当前测试 | 希望一次收集多个独立断言 |
| `Fatal/Fatalf` | 标记失败并对当前 goroutine 调用 `runtime.Goexit` | 后续断言依赖当前结果 |
| `FailNow` | 立即停止当前测试 goroutine | 不能从辅助 goroutine 调用 |
| `Run` | 创建命名子测试并返回是否通过 | 表驱动、分组和选择执行 |
| `Helper` | 把调用函数标记为辅助函数 | 失败行定位到调用方 |
| `Cleanup` | 测试和所有子测试完成后 LIFO 清理 | 临时资源、恢复全局状态 |
| `TempDir` | 创建测试专属临时目录 | 测试结束自动删除 |
| `Setenv` | 设置环境变量并自动恢复 | 不能用于并行测试及其祖先 |
| `Context` | 返回测试结束前会取消的 context | 启动受测试生命周期管理的任务 |
| `Parallel` | 暂停并标记与其他并行测试并发运行 | 夹具必须隔离 |

`Fatalf` 只终止调用它的 goroutine；辅助 goroutine 应通过 channel 返回错误，由测试 goroutine调用 Fatal。测试函数返回或调用 FailNow 时，已注册的 Cleanup 仍会执行。

## 表驱动测试

表驱动不是语言特性，而是把输入、期望和名称组织为数据：

```go
func TestParsePositiveInt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "integer", input: "42", want: 42},
		{name: "trim spaces", input: " 7 ", want: 7},
		{name: "zero", input: "0", wantErr: true},
		{name: "syntax", input: "seven", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParsePositiveInt(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("value = %d, want %d", got, test.want)
			}
		})
	}
}
```

用稳定名称可执行单个 case：

```sh
go test -run 'TestParsePositiveInt/trim_spaces' -v
```

失败信息包含输入或 case 名、实际值与期望值。不要只写 `unexpected result`。错误测试优先用 `errors.Is/AsType` 验证类别和结构，不锁定可能调整的完整文案。

### 并行子测试

Go 1.22 起 `range` 每次迭代有独立迭代变量，但 case 中引用的 slice、map、server 或全局状态仍可能共享。调用 `t.Parallel()` 前应确保每个 case 有独立夹具。

并行测试在顺序测试及其父测试返回后运行。不要用 `Setenv`、`Chdir` 或其他进程级状态变化配合 Parallel。

## 测试替身与边界

Go 常用小接口替换外部依赖：

```go
type OrderStore interface {
	Load(ctx context.Context, id string) (Order, error)
}
```

fake 在内存中实现行为；stub 返回预设结果；spy 记录调用；mock 断言交互顺序。只替换真正的进程边界，不要为每个 struct 建同名接口。数据库约束、序列化、HTTP 头和事务行为还需要集成测试，因为内存 fake 无法证明这些性质。

时间与随机性应注入函数或接口：

```go
type Clock func() time.Time

func NewToken(now Clock) Token { /* ... */ }
```

固定等待 `time.Sleep` 会让测试慢且不稳定。并发测试用 channel、WaitGroup、context 和最终状态条件同步。

## HTTP 与文件测试

`httptest.NewServer` 启动真实本地 HTTP server，适合验证客户端超时、状态和 header；`httptest.NewRecorder` 在内存记录 handler 响应。测试结束调用 `Close`，或注册 `t.Cleanup(server.Close)`。

文件测试用 `t.TempDir()`，不依赖当前目录和开发机已有文件。测试夹具放 `testdata/`；Go 工具会忽略该目录作为 package，但测试可读取其中内容。

## 覆盖率的准确解释

```sh
go test -coverprofile=cover.out ./...
go tool cover -func=cover.out
go tool cover -html=cover.out
```

覆盖率表示被插桩语句块是否执行，不表示断言正确、所有输入组合或并发交错已验证。100% statement coverage 仍可能漏掉错误类别、顺序、授权和不变量。覆盖报告用于发现未触达路径，不能作为唯一质量目标。

`go test -cover` 默认分析被测包；`-coverpkg` 可扩大插桩包集合。跨包覆盖会增加运行开销，需明确报告口径。

## Fuzz 测试

模糊测试从 seed corpus 出发生成输入，寻找 panic、挂起或违反不变量的 case：

```go
func FuzzParsePositiveInt(f *testing.F) {
	f.Add("42")
	f.Add("0")
	f.Add("seven")
	f.Fuzz(func(t *testing.T, input string) {
		value, err := ParsePositiveInt(input)
		if err == nil && value < 1 {
			t.Fatalf("successful value = %d", value)
		}
	})
}
```

运行：

```sh
go test -fuzz=FuzzParsePositiveInt -fuzztime=30s
```

发现的最小失败输入会写入 `testdata/fuzz/...`，应提交为回归语料。fuzz 函数必须确定、快速且不能依赖共享可变全局状态。安全模糊测试还应设置资源上限，避免输入触发无限分配。

## Benchmark 的执行模型

Go 1.24 起推荐 `b.Loop()`：框架动态选择迭代次数，循环内调用被测操作。

```go
func BenchmarkParsePositiveInt(b *testing.B) {
	for b.Loop() {
		_, _ = ParsePositiveInt("42")
	}
}
```

运行：

```sh
go test -run '^$' -bench '^BenchmarkParsePositiveInt$' -benchmem -count=10
```

- `-run '^$'` 不运行普通测试。
- `-bench` 用正则选择基准；默认不运行 benchmark。
- `-benchmem` 报告每操作分配字节和次数。
- `-benchtime=3s` 增加采样时间；也可设固定次数如 `1000x`。
- `-count=10` 产生多个样本供统计比较。
- `-cpu=1,4,8` 改变 `GOMAXPROCS`，适用于并发基准。

`b.Loop` 会把循环体内调用视为存活，减少编译器删除无用工作的风险。仍应消费结果或验证一次输出，避免基准测到错误路径。

### 初始化与计时

耗时初始化不应计入操作：

```go
func BenchmarkLookup(b *testing.B) {
	index := buildLargeIndex()
	b.ResetTimer()
	for b.Loop() {
		result = index.Lookup("order-42")
	}
}
```

需要在循环中准备输入时可 `b.StopTimer()` / `b.StartTimer()`，但频繁切换也有开销。`b.ReportAllocs()` 等价于命令行 `-benchmem` 对该 benchmark 开启分配统计。自定义工作单位用 `b.SetBytes` 或 `b.ReportMetric`。

### 并行 Benchmark

```go
func BenchmarkCacheParallel(b *testing.B) {
	cache := newCache()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = cache.Get("key")
		}
	})
}
```

`RunParallel` 在多个 goroutine 中分配迭代，默认并行度受 `GOMAXPROCS` 影响，可用 `SetParallelism` 调整。它适合多 goroutine 可并行调用的 API，不适合本质顺序的操作。

## 用 benchstat 做 A/B 比较

在同一台空闲机器、相同 Go 版本、CPU 策略和输入下分别采样：

```sh
go test -run '^$' -bench Parse -benchmem -count=10 > old.txt
go test -run '^$' -bench Parse -benchmem -count=10 > new.txt
benchstat old.txt new.txt
```

`benchstat` 报告样本分布、变化百分比和统计检验。单次 `ns/op` 差异可能来自调频、后台进程、温度或调度。统计显著也不等于业务显著：每请求快 5ns 对网络服务可能没有价值，而一次分配减少可能降低 GC 压力。

## 完整案例：解析正整数

输入分四组：`"42"`、`" 7 "`、`"0"`、`"seven"`。

1. 函数先 `TrimSpace`，再调用 `strconv.Atoi`。
2. 语法错误用 `%w` 包装 `strconv.NumError`。
3. 解析成功但小于 1 时返回领域错误。
4. 表驱动测试逐项验证错误存在性和值。
5. benchmark 固定输入 `"42"`，量化热路径耗时与分配。

实现位于 [`../../examples/go/parse.go`](../../examples/go/parse.go)，测试与基准位于 [`../../examples/go/examples_test.go`](../../examples/go/examples_test.go)。

```sh
cd 05-backend-data/examples/go
go test -run TestParsePositiveInt -v
go test -run '^$' -bench BenchmarkParsePositiveInt -benchmem -count=3
```

失败分支：若忘记 `TrimSpace`，`" 7 "` case 失败；若只判断 Atoi 错误，`"0"` 会错误成功；若 benchmark 输入变成错误路径，结果不再代表正常解析。测试名称和失败输出会直接指出输入与实际值。

## 常见错误与修正

- 只测试成功路径：至少加入边界、无效输入、依赖错误、取消和超时。
- 比较完整错误字符串：改用 `errors.Is/AsType`，仅在文案本身是契约时比较文本。
- 测试依赖执行顺序：每个测试创建独立状态，清理用 Cleanup。
- 用 sleep 等 goroutine：用明确同步信号或最终条件加截止时间。
- 并行测试修改 env/工作目录：移除 Parallel 或注入配置。
- 基准包含初始化：移到循环外并 ResetTimer。
- 只跑一次 benchmark：多次采样，用 benchstat 比较。
- 把覆盖率当正确率：补充针对不变量和失败语义的断言。

## 验证清单

1. `go test ./...` 和 `go test -race ./...` 都通过。
2. `go test -shuffle=on -count=20 ./...` 能发现顺序依赖和偶发失败。
3. 测试失败信息包含 case、got、want 和 error。
4. 外部资源都有 Cleanup，测试无残留 goroutine。
5. benchmark 记录 Go 版本、OS/架构、CPU、GOMAXPROCS、输入和参数。
6. 性能变更保存修改前后多样本，并验证功能测试未退化。

## 练习

为 `ValidateOrder` 添加表驱动、fuzz 和 benchmark。完成标准：覆盖合法、空 ID、零数量、两个字段同时错误；用 `errors.Is` 和 `errors.AsType` 验证错误树；fuzz 保证任意输入不 panic；benchmark 分开测成功与双错误路径；使用 `-count=10` 和 benchstat 比较加入 `errors.Join` 前后的分配，但不以微基准替代错误完整性。

## 来源

- [Go：testing package](https://pkg.go.dev/testing)（访问日期：2026-07-17）
- [Go：go test command](https://pkg.go.dev/cmd/go#hdr-Test_packages)（访问日期：2026-07-17）
- [Go：Fuzzing](https://go.dev/doc/security/fuzz/)（访问日期：2026-07-17）
- [Go：benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)（访问日期：2026-07-17）
