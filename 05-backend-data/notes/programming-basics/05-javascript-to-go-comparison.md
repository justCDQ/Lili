# 从 JavaScript 到 Go：类型、数据结构与错误模型对照

## 学习目标

本文不是语法替换表，而是解释同一需求在 JavaScript 与 Go 中如何建立等价契约。完成后应能识别类型、数值、集合、Unicode、错误和模块边界的关键差异，并在迁移时保持可观察行为一致。

## 1. 运行与类型模型

JavaScript 的类型属于运行时值；同一变量可先后引用不同类型。Go 是静态类型语言，变量类型在编译期确定，赋值和调用必须满足类型规则。Go 仍可通过接口表达多种具体类型，但这不是取消类型检查。

```js
let value = 1;
value = "one"; // 合法
```

```go
var value int = 1
// value = "one" // 编译错误
```

静态类型会在构建阶段拒绝一部分错误，但不能自动保证数值范围、业务状态、文件内容或网络输入正确。两种语言都需要运行时边界验证。

### 1.1 零值与缺失

Go 变量未显式初始化时获得其类型的零值：布尔为 `false`，数值为 0，字符串为 `""`，指针、函数、接口、切片、map 和 channel 为 `nil`。JavaScript 的未初始化 `let` 绑定在声明执行后值为 `undefined`；对象缺失属性读取也得到 `undefined`。

Go 的普通 `int` 字段无法区分“未提供”和“明确提供 0”。需要该区别时可用指针、`value + bool`、专用可选类型或解码器提供的字段存在信息。不要在所有字段上机械使用指针；应由接口语义决定。

| 语义 | JavaScript 常见表示 | Go 常见表示 |
| --- | --- | --- |
| 必填整数 | Number + 运行时验证 | `int64` + 解码后范围验证 |
| 可选整数 | `undefined` / 缺失属性 | `*int64` 或专用 option |
| 明确空值 | `null` | 指针为 nil，需结合协议约定 |
| 固定状态集合 | String + 检查 | 自定义 string 类型 + 检查 |

## 2. 数值与转换

JavaScript Number 是 IEEE 754 双精度浮点数，安全整数上限为 `2^53-1`；BigInt 表示任意精度整数但不能与 Number 直接混算。Go 提供多种有符号、无符号整数及浮点类型，不同已定义数值类型通常需要显式转换。

Go 的转换语法不等于安全校验。`int32(3_000_000_000)` 会按目标宽度产生截断结果，而不是自动返回错误。外部输入先验证范围，再转换。

```go
func toInt32(n int64) (int32, error) {
    if n < math.MinInt32 || n > math.MaxInt32 {
        return 0, fmt.Errorf("out of int32 range: %d", n)
    }
    return int32(n), nil
}
```

金额跨语言传输时，最小货币单位的十进制整数仍受 JavaScript 安全整数范围约束。可能超过范围的 ID、计数或金额应在 JSON 中约定为十进制字符串，并在两端严格解析；JSON 数字文本本身不规定接收语言必须保持任意精度。

## 3. 字符串与 Unicode

JavaScript String 是 UTF-16 码元序列，`length` 返回码元数。Go string 是只读字节序列，通常保存 UTF-8，`len` 返回字节数。两者直接按索引都不保证得到完整 Unicode 码点。

```js
const text = "A😀";
console.log(text.length);       // 3 个 UTF-16 码元
console.log([...text].length);  // 2 个码点
```

```go
text := "A😀"
fmt.Println(len(text))                    // 5 个 UTF-8 字节
fmt.Println(utf8.RuneCountInString(text)) // 2 个码点
```

跨语言契约不能写“最多 20 字符”后各自使用 `length`。要明确按 UTF-8 字节、Unicode 码点还是扩展字素簇，并使用相同 Unicode 版本和规范化策略。

## 4. Object、struct、Array、slice 与 map

JavaScript Object 是动态属性集合，运行时可增删字段；Go struct 的字段集合与类型在编译期确定。JSON 解码到 Go struct 时，未知字段默认可能被忽略，严格边界可使用 `Decoder.DisallowUnknownFields`。

JavaScript Array 是对象，可动态增长并允许混合值。Go 数组长度属于类型，切片提供动态长度但元素类型固定。Go 切片描述符可能共享底层数组；JavaScript Array 变量也共享对象引用，但具体复制和增长规则不同，不能把一种语言的直觉直接套到另一种。

JavaScript Map 的键可为对象，按引用身份区分；Go map 键必须是可比较类型，slice、map 和 function 不能直接作为键。两者都不应在无同步情况下进行实现不支持的并发写入；JavaScript 单一事件循环也不能保证跨 worker 或共享内存不存在竞态。

Go map 遍历顺序未指定。JavaScript Map 明确按插入顺序迭代，普通 Object 的属性枚举也有规范化顺序规则。若输出需要跨语言稳定顺序，应显式排序，不依赖容器默认。

## 5. 函数、方法与接口

JavaScript 函数是一等对象，可捕获词法环境，参数数量和类型在运行时不自动强制。Go 函数也可作为值并形成闭包，但函数签名规定参数和返回类型，支持多个返回值。

JavaScript 方法调用中的 `this` 取决于调用形式；箭头函数捕获外层 `this`。Go 方法有显式接收者，调用不会依据动态调用语法重新绑定接收者。

Go 接口按方法集隐式实现，不需要 `implements` 声明。接口适合在消费者侧描述所需最小行为，例如函数只需要 `io.Reader`，无需知道输入来自文件、网络还是内存。

```go
func DecodeOrders(r io.Reader) ([]Order, error) {
    var orders []Order
    err := json.NewDecoder(r).Decode(&orders)
    return orders, err
}
```

空接口值与“包含 nil 指针的接口值”不同：接口包含动态类型和动态值，只要动态类型存在，接口就不等于 nil。这是迁移错误处理中常见边界。

## 6. 错误与异常

JavaScript 同步代码用 `throw` 中断当前控制流，Promise 用拒绝表示异步失败，`await` 会把拒绝转为抛出。应抛出 Error 对象并保留 `cause`。

Go 通常把 `error` 作为最后一个返回值显式传播。调用者用 `if err != nil` 处理。`fmt.Errorf("context: %w", err)` 包装并保留错误链，`errors.Is` 按目标错误匹配，`errors.As` 提取类型。

panic 表示普通返回机制难以继续的异常状态，不是日常输入验证的替代。`recover` 只在延迟函数中生效，通常用于进程边界隔离并把 panic 转成受控失败；它不能恢复已经破坏的业务不变量。

| 场景 | JavaScript | Go |
| --- | --- | --- |
| 参数业务非法 | 抛自定义 Error 或返回显式结果 | 返回非 nil error |
| 异步 I/O 失败 | Promise rejection / `await` 抛出 | 返回 error；异步由 goroutine/channel 协调 |
| 保留原因 | `new Error(msg, {cause})` | `%w` + `errors.Is/As` |
| 程序不变量破坏 | 可能 throw | 可能 panic，但应谨慎 |

错误迁移不能只把 `throw` 改成第二返回值。必须检查每个调用点：是否添加上下文、是否能处理、失败时是否仍产生部分输出、资源是否释放。

## 7. 异步模型概览

JavaScript Promise 表示未来完成值，事件循环安排任务和微任务；异步函数返回 Promise。Go 用 goroutine 运行并发函数，用 channel、锁、context 等协调。两者都不保证业务操作自动有超时、取消、顺序或幂等性。

把 `Promise.all` 逐字翻译成每项启动 goroutine 可能导致无界并发。迁移时先记录原实现的最大并发、失败策略、结果顺序、取消方式和资源上限，再选择 worker pool、信号量或顺序执行。

## 8. 模块与依赖边界

ESM 使用 import/export；Go 源文件声明 package 并通过 import path 引用包。Node 包与 Go 模块的版本语义、解析算法和发布约束不同，不应假设 `package.json` 与 `go.mod` 一一对应。

迁移时先画出外部可观察边界：HTTP 请求响应、JSON、数据库 schema、文件格式、命令行 stdout/stderr/exit code。内部目录可按 Go 惯例重组，只要契约有兼容测试保护。

## 9. 完整案例：词频统计的等价实现

### 9.1 契约

输入是字符串数组。每项先执行 Unicode NFC，再转小写；空字符串忽略。输出按出现次数降序、次数相同按标准化后的字符串升序。结果字段为 `word` 和 `count`。

输入：

```json
["Go", "go", "API", "api", "api", "é", "é", ""]
```

预期步骤：`Go/go` 合并为 `go` 两次；`API/api/api` 合并为 `api` 三次；两种 `é` 经 NFC 合并为两次；空串忽略。预期输出：

```json
[
  {"word":"api","count":3},
  {"word":"go","count":2},
  {"word":"é","count":2}
]
```

这里的小写仅用于教学契约，并不构成适用于所有语言的身份规范。生产系统需明确 locale、Unicode 版本与安全需求。

### 9.2 JavaScript 实现

```js
export function frequency(words) {
  if (!Array.isArray(words)) throw new TypeError("words must be an array");
  const counts = new Map();
  for (const [index, raw] of words.entries()) {
    if (typeof raw !== "string") {
      throw new TypeError(`words[${index}] must be a string`);
    }
    const word = raw.normalize("NFC").toLowerCase();
    if (word === "") continue;
    counts.set(word, (counts.get(word) ?? 0) + 1);
  }
  return [...counts.entries()]
    .map(([word, count]) => ({ word, count }))
    .sort((a, b) => b.count - a.count || a.word.localeCompare(b.word, "und"));
}
```

`localeCompare` 的具体排序取决于实现提供的国际化数据。若协议要求跨运行时逐字节一致，应定义码点或 UTF-8 字节排序，而不是把默认 locale 排序当成固定协议。

可用以下断言验证核心结果：

```js
const input = ["Go", "go", "API", "api", "api", "é", "e\u0301", ""];
const actual = frequency(input);
console.assert(JSON.stringify(actual) === JSON.stringify([
  { word: "api", count: 3 },
  { word: "go", count: 2 },
  { word: "é", count: 2 },
]));
```

### 9.3 Go 实现

Go 标准库没有 Unicode 规范化包；官方扩展模块 `golang.org/x/text/unicode/norm` 提供 NFC。引入前应在 `go.mod` 固定并审查版本。

```go
package frequency

import (
    "fmt"
    "sort"
    "strings"

    "golang.org/x/text/unicode/norm"
)

type Entry struct {
    Word  string `json:"word"`
    Count int    `json:"count"`
}

func Count(words []string) ([]Entry, error) {
    counts := make(map[string]int, len(words))
    for _, raw := range words {
        word := strings.ToLower(norm.NFC.String(raw))
        if word == "" {
            continue
        }
        counts[word]++
    }

    out := make([]Entry, 0, len(counts))
    for word, count := range counts {
        out = append(out, Entry{Word: word, Count: count})
    }
    sort.Slice(out, func(i, j int) bool {
        if out[i].Count != out[j].Count {
            return out[i].Count > out[j].Count
        }
        return out[i].Word < out[j].Word
    })
    return out, nil
}

func ValidateInput(values []any) ([]string, error) {
    words := make([]string, len(values))
    for i, value := range values {
        word, ok := value.(string)
        if !ok {
            return nil, fmt.Errorf("words[%d] must be a string", i)
        }
        words[i] = word
    }
    return words, nil
}
```

Go `[]string` 已在调用边界提供元素静态类型，因此 `Count` 无需逐项动态检查。若输入来自解码到 `[]any` 的非类型化 JSON，则 `ValidateInput` 明确转换并返回带索引错误。更直接的做法是把 JSON 解码到 `[]string`，让解码器拒绝非字符串元素。

### 9.4 Go 测试

```go
func TestCount(t *testing.T) {
    input := []string{"Go", "go", "API", "api", "api", "é", "e\u0301", ""}
    got, err := Count(input)
    if err != nil {
        t.Fatal(err)
    }
    want := []Entry{{"api", 3}, {"go", 2}, {"é", 2}}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("Count() = %#v, want %#v", got, want)
    }
}
```

运行步骤：

```bash
go mod init example.com/frequency
go get golang.org/x/text/unicode/norm
go test ./...
```

`go get` 会修改 `go.mod` 和 `go.sum`，应审查并提交这两个文件。测试输出 `ok` 表示样本契约通过。

### 9.5 排序差异与调整

上述 JavaScript 使用 locale collation，Go 使用字符串字节词法比较，遇到非 ASCII 同频词时顺序可能不同。要得到跨语言完全一致结果，应把契约改为 UTF-8 字节升序，并在 JavaScript 实现相同比较器，或两端采用同版本 Unicode Collation Algorithm 与相同 locale/options。

这说明迁移的正确标准是可观察契约一致，不是代码形状相似。示例的预期三项恰好不暴露全部排序差异，必须加入 `ä`、`z`、中文和组合字符等同频样本。

### 9.6 失败分支

JavaScript 输入 `["go", 7]` 在索引 1 抛 TypeError。Go 若 JSON 直接解码到 `[]string`，解码阶段返回类型错误；若调用 `ValidateInput`，返回 `words[1] must be a string`。两端都不应静默把数字转成字符串。

超大数组会使 map 与输出占用 O(k) 空间；输入无界时改用流式解析、设最大单词数与不同词数量。计数可能溢出 Go `int`，长期流应选择明确宽度并检查上限；JavaScript 计数超过安全整数也必须失败或改用 BigInt。

## 10. 迁移检查表

- 列出外部接口、样本与错误，而不是从目录结构开始翻译。
- 标记每个 Number 是整数、浮点、金额、时间还是不透明 ID。
- 为 `null`、`undefined`、缺失、零值和空集合分别定义语义。
- 明确字符串长度、规范化、大小写与排序算法。
- 检查对象与切片/map 的共享修改以及复制深度。
- 不依赖 Go map 顺序，协议输出显式排序。
- 把 throw/rejection 的每条失败路径映射为 error，并保留原因。
- 为异步任务定义并发上限、超时、取消和部分成功策略。
- 先建立兼容测试，再重构为 Go 惯用包边界。
- 对性能结论使用基准与 profile，不凭语言标签判断。

## 11. 调试常见迁移问题

- JSON 中大整数改变：检查 JavaScript 解析前是否已失真，改用字符串协议。
- 中文截断乱码：检查 Go 是否按字节切片、JavaScript 是否切断代理对。
- 输出顺序不稳定：排序 map 键或输出条目，不比较未定义顺序。
- 字段缺失被当成零：使用可选表示并在解码后验证存在性。
- error 看似 nil 却进入失败：检查接口是否包含 nil 指针动态值。
- goroutine 数持续增加：检查从 Promise 迁移后是否缺少取消、结果接收或并发上限。
- 测试通过但协议不兼容：比较真实序列化字节、状态码和错误结构。

## 12. 练习

1. 为词频案例实现统一 UTF-8 字节比较器，并加入非 ASCII 同频测试。
2. 把可能超过安全整数的计数改为十进制字符串输出，分别实现 JS BigInt 与 Go uint64 边界。
3. 设计 JSON patch 字段，使“缺失”“null”“0”在 Go 中可以区分并测试。
4. 找一个使用 `Promise.all` 的 JS 函数，写出其失败和顺序契约，再设计有限并发 Go 版本。
5. 构造包含 nil 指针的 Go error 接口，解释为何 `err != nil`，并修复返回方式。

## 来源

- [ECMA-262：ECMAScript Language Types](https://tc39.es/ecma262/#sec-ecmascript-language-types)（访问日期：2026-07-17）
- [ECMA-262：Map Objects](https://tc39.es/ecma262/#sec-map-objects)（访问日期：2026-07-17）
- [Go 语言规范](https://go.dev/ref/spec)（访问日期：2026-07-17）
- [Go 官方博客：Error handling and Go](https://go.dev/blog/error-handling-and-go)（访问日期：2026-07-17）
- [Go 扩展模块：golang.org/x/text/unicode/norm](https://pkg.go.dev/golang.org/x/text/unicode/norm)（访问日期：2026-07-17）
