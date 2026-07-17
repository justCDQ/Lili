# 二进制、十六进制、字符编码、时间与时区

## 学习目标

本文解释位、字节、整数编码、UTF-8、时间点、持续时间、UTC 偏移和 IANA 时区的区别，并实现一个带长度、时间戳和 UTF-8 载荷的可验证二进制记录格式。

## 1. 位、字节与进制

位只有 0 或 1。字节通常为 8 位，可表示 0 到 255。二进制直接展示位模式；十六进制每一位对应 4 bit，因此两个十六进制数字正好表示一个字节。

```text
十进制  172
二进制  10101100
十六进制 AC
```

进制只改变书写，不改变值。Go 整数字面量可写 `0b10101100`、`0xac` 和 `172`；格式化 `%08b` 输出 8 位二进制，`%02x` 输出两位小写十六进制。

位运算：`&` 按位与，用于清除或测试位；`|` 按位或，用于设置位；`^` 异或；`&^` 位清除；`<<`、`>>` 移位。使用位字段时为每一位定义名称，验证保留位为零，避免散落魔法数字。

```go
const (
    FlagCompressed byte = 1 << 0
    FlagEncrypted  byte = 1 << 1
)
flags := FlagCompressed | FlagEncrypted
compressed := flags&FlagCompressed != 0
```

## 2. 有符号整数、宽度与溢出

无符号 n 位整数范围是 0 到 `2^n-1`。Go 的有符号整数使用二进制补码，n 位范围为 `-2^(n-1)` 到 `2^(n-1)-1`。固定协议字段要使用 `uint16`、`int64` 等确定位宽类型，不用宽度依平台的 `int`。

Go 整数运行时运算按类型宽度截断，不自动返回溢出错误。处理长度、金额和时间戳时，在运算与转换前检查范围。

```go
var x uint8 = 255
x++
fmt.Println(x) // 0
```

文本中的 `"123"` 与二进制整数 123 不同：前者 UTF-8 字节为 `31 32 33` 十六进制，后者若用 32 位大端是 `00 00 00 7b`。协议必须规定表示、位宽、符号和字节序。

## 3. 字节序

多字节整数需要规定高低字节顺序。大端把最高有效字节放在最低地址或流的前面，小端相反。

数值 `0x12345678`：

```text
大端字节：12 34 56 78
小端字节：78 56 34 12
```

网络协议常使用大端，但这是具体协议约定，不是所有网络负载的自动规则。Go 的 `encoding/binary.BigEndian` 和 `LittleEndian` 显式读写整数。不要用 `unsafe` 把结构体内存直接当跨平台协议，因为填充、对齐、架构字节序和类型布局会变化。

## 4. 字符、码点与编码

Unicode 为抽象字符分配码点，例如 `狸` 为 U+72F8。字符编码把码点序列映射为字节。UTF-8 使用 1 到 4 字节编码 Unicode 标量值，ASCII U+0000–U+007F 保持单字节相同。

```go
text := "狸力"
fmt.Printf("% x\n", []byte(text)) // e7 8b b8 e5 8a 9b
```

Go string 是任意字节序列，可能包含无效 UTF-8。`utf8.ValidString` 验证；`range` 遇到无效编码会产生 U+FFFD，宽度为 1 字节。若协议要求有效 UTF-8，应在边界拒绝无效输入，不能解码后默默替换再签名或存储。

编码与序列化不同。UTF-8 只定义文本到字节；JSON 还定义结构、字符串转义和数字语法；Base64 把任意字节编码成有限 ASCII 字符集合，体积通常增加，且不提供加密或完整性保护。

十六进制也只是编码表示。摘要或密钥的十六进制字符串长度是原字节数两倍；比较前先明确大小写、前缀和允许长度。

## 5. 时间相关的四个概念

时间点是时间线上的瞬间；持续时间是两个事件间隔；UTC 偏移表示某本地时间与 UTC 的数值差；时区是包含历史与未来民用时间规则的地区标识，如 `Asia/Shanghai`。

RFC 3339 时间戳 `2026-07-17T10:30:00+08:00` 表示一个时间点，`+08:00` 是该表示使用的偏移，不自动携带地区时区 ID。相同时间点可写成 `2026-07-17T02:30:00Z`。

日历日期 `2026-07-17` 可能只表示生日、账期或当地营业日，不一定是 UTC 零点。把纯日期强制转时间点会引入无根据的时区选择。

## 6. 墙上时钟与单调时钟

墙上时钟用于回答当前日期时间，可能被 NTP 校正、管理员修改或虚拟化环境调整。单调时钟只保证向前测量经过时间，不对应日历日期。

Go `time.Now` 返回的 Time 可携带单调时钟读数；用 `time.Since(start)` 或 `end.Sub(start)` 测量同一进程内持续时间时可使用该读数，避免墙钟跳变。序列化 Time 会去除单调部分，所以跨进程持续时间不能依赖它。

超时通常基于持续时间与单调计时；审计事件记录基于 UTC 墙钟，并可额外记录请求顺序或单调耗时。单个墙钟时间戳不能严格排序分布式系统所有事件。

## 7. RFC 3339、Unix 时间戳与精度

RFC 3339 是互联网时间戳文本格式的规范子集，包含日期、时间、秒和 UTC 关系。Go `time.RFC3339` 是布局常量，布局通过参考时间 `Mon Jan 2 15:04:05 MST 2006` 描述格式。

```go
t, err := time.Parse(time.RFC3339, "2026-07-17T10:30:00+08:00")
fmt.Println(t.UTC().Format(time.RFC3339)) // 2026-07-17T02:30:00Z
```

Unix 时间通常指自 1970-01-01T00:00:00Z 起的秒数，但 API 可能使用秒、毫秒、微秒或纳秒。字段名应带单位，例如 `created_at_ms`，并检查范围。Unix 时间戳本身不保存原始时区。

浮点 Unix 秒可能丢失精度。协议可用整数秒加纳秒字段，或 RFC 3339 小数秒文本。数据库时间类型的精度和时区语义是实现相关的，必须看所选数据库文档。

## 8. IANA 时区与夏令时

IANA Time Zone Database 记录地区偏移规则和历史变更。地区时区不是固定偏移。某些本地时间在夏令时切换时不存在，另一些会重复两次；解析“当地 01:30”必须定义歧义策略。

Go 用 `time.LoadLocation("America/New_York")` 读取时区规则。运行环境需要提供时区数据库，可由系统、Go 安装或嵌入的 `time/tzdata` 提供。容器精简镜像缺少 tzdata 时加载可能失败，部署要测试。

未来民用规则可能改变，保存预约时通常需要保存用户选择的当地日期时间、时区 ID 与业务重算策略，而不只是当前换算出的 UTC 时间点。

## 9. 完整案例：二进制事件记录

### 9.1 格式

定义版本 1 记录：

| 偏移 | 长度 | 字段 | 编码 |
| --- | --- | --- | --- |
| 0 | 2 | magic | 固定字节 `4c 49`，即 LI |
| 2 | 1 | version | `01` |
| 3 | 1 | flags | bit0 表示文本已 NFC；其余必须为 0 |
| 4 | 8 | unix nanoseconds | 有符号 int64，大端 |
| 12 | 2 | payload length | uint16，大端 |
| 14 | N | payload | 有效 UTF-8 |

最大载荷 4096 字节，虽然字段能表示 65535，业务限制更小。时间范围必须能由 Go `time.Unix(0,nanos)` 表示并满足产品允许区间。

### 9.2 编码与解码

```go
package record

import (
    "encoding/binary"
    "errors"
    "fmt"
    "time"
    "unicode/utf8"
)

const headerSize = 14
const maxPayload = 4096

type Record struct {
    Time       time.Time
    Normalized bool
    Text       string
}

func Encode(record Record) ([]byte, error) {
    payload := []byte(record.Text)
    if !utf8.Valid(payload) {
        return nil, errors.New("payload is not valid UTF-8")
    }
    if len(payload) > maxPayload {
        return nil, fmt.Errorf("payload is %d bytes, maximum is %d", len(payload), maxPayload)
    }
    out := make([]byte, headerSize+len(payload))
    copy(out[0:2], []byte{'L', 'I'})
    out[2] = 1
    if record.Normalized {
        out[3] = 1
    }
    binary.BigEndian.PutUint64(out[4:12], uint64(record.Time.UnixNano()))
    binary.BigEndian.PutUint16(out[12:14], uint16(len(payload)))
    copy(out[14:], payload)
    return out, nil
}

func Decode(data []byte) (Record, error) {
    if len(data) < headerSize {
        return Record{}, errors.New("record is shorter than header")
    }
    if data[0] != 'L' || data[1] != 'I' {
        return Record{}, errors.New("invalid magic")
    }
    if data[2] != 1 {
        return Record{}, fmt.Errorf("unsupported version %d", data[2])
    }
    if data[3]&^byte(1) != 0 {
        return Record{}, fmt.Errorf("reserved flag bits set: 0x%02x", data[3])
    }
    size := int(binary.BigEndian.Uint16(data[12:14]))
    if size > maxPayload {
        return Record{}, fmt.Errorf("payload length %d exceeds maximum", size)
    }
    if len(data) != headerSize+size {
        return Record{}, fmt.Errorf("length mismatch: header=%d actual=%d", size, len(data)-headerSize)
    }
    payload := data[14:]
    if !utf8.Valid(payload) {
        return Record{}, errors.New("payload is not valid UTF-8")
    }
    nanos := int64(binary.BigEndian.Uint64(data[4:12]))
    return Record{
        Time:       time.Unix(0, nanos).UTC(),
        Normalized: data[3]&1 != 0,
        Text:       string(payload),
    }, nil
}
```

把 int64 转为 uint64 写入再读回 int64，可以保持相同 64 位二进制补码位模式。协议明确有符号语义，不能在其他语言中把字段当无符号时间。

### 9.3 输入、输出和验证

输入时间 `2026-07-17T10:30:00+08:00`、文本 `狸力`、Normalized=true。解析时间后 Encode，再以十六进制打印：

```go
parsed, _ := time.Parse(time.RFC3339, "2026-07-17T10:30:00+08:00")
encoded, err := Encode(Record{Time: parsed, Normalized: true, Text: "狸力"})
if err != nil { panic(err) }
fmt.Printf("%x\n", encoded)

decoded, err := Decode(encoded)
if err != nil { panic(err) }
fmt.Println(decoded.Time.Format(time.RFC3339Nano), decoded.Text)
```

步骤是写 magic/version/flags；把时间转 Unix 纳秒并按大端写 8 字节；写载荷字节长度 6，即 `0006`；追加 UTF-8 `e78bb8e58a9b`。输出时间为同一瞬间的 UTC 表示，文本仍为 `狸力`。

测试应断言 Decode(Encode(x)) 的时间点 `Equal`、文本和 flags 相同，并断言编码的最后六字节精确为预期 UTF-8。

仓库中的[可运行 Record 示例](../../examples/computer-systems/record/)包含 round-trip、长度、版本、保留位与非法 UTF-8 测试。

### 9.4 失败分支

- 把 magic 改成 `00 00`，返回 invalid magic。
- version 改为 2，返回 unsupported version，不能按版本 1 猜测。
- flags 设置 `0x80`，因保留位非零失败。
- 长度字段写 7 而实际 6，返回 length mismatch，不读取越界。
- 载荷写无效 UTF-8 `ff`，返回编码错误。
- 文本超过 4096 字节，Encode 在分配结果前拒绝。

格式未包含校验和或认证码，随机损坏可能恰好仍通过结构验证；它也不提供真实性和防篡改。需要这些属性时使用明确的 CRC 或密码学认证方案，并把覆盖范围写入协议。

## 10. 时间调试清单

- 相差 8 小时：确认值是时间点还是当地时间，是否重复应用偏移。
- 相差 1000 倍：检查秒、毫秒、微秒、纳秒单位字段。
- 夏令时附近重复/缺失：保留 IANA 时区 ID并定义歧义策略。
- 测试只在部分机器失败：固定时区、locale 和 tzdata 版本或注入 Location。
- 超时出现负数：使用单调持续时间，不以墙钟时间戳相减。
- JSON 时间解析不一致：只接受明确 RFC 3339 形式，拒绝模糊本地字符串。
- 二进制值异常：以十六进制逐字段标偏移，核对宽度、有符号和字节序。

## 11. 练习

1. 手工编码 `0x01020304` 的大小端字节，再用 `encoding/binary` 验证。
2. 为 Record 编写 round-trip 与六个失败分支测试，并运行 fuzz 测试保证 Decode 不 panic。
3. 比较 `len("狸力")`、rune 数与终端显示宽度，说明三者用途。
4. 解析同一 RFC 3339 时间点的 `Z` 和 `+08:00` 表示，用 `Time.Equal` 验证。
5. 构造一个夏令时重复当地时间，研究 Go `time.Date` 的选择并写业务显式策略。

## 来源

- [Unicode Standard 17.0：Chapter 3 Conformance](https://www.unicode.org/versions/Unicode17.0.0/core-spec/chapter-3/)（访问日期：2026-07-17）
- [RFC 3629：UTF-8](https://www.rfc-editor.org/rfc/rfc3629)（访问日期：2026-07-17）
- [RFC 3339：Date and Time on the Internet](https://www.rfc-editor.org/rfc/rfc3339)（访问日期：2026-07-17）
- [IANA Time Zone Database](https://www.iana.org/time-zones)（访问日期：2026-07-17）
- [Go 标准库：time](https://pkg.go.dev/time)（访问日期：2026-07-17）
