# 文件、目录、路径、权限、进程与环境变量

## 学习目标

本文解释操作系统如何命名和访问文件、Unix 权限位如何参与授权、程序如何成为进程，以及环境变量如何随进程创建而传递。案例将建立一个限制在指定目录内的文件读取命令。

## 1. 文件与目录

普通文件保存字节序列及元数据，例如类型、大小、所有者、权限和时间戳。文件名不是文件内容的一部分；目录保存“名称到文件系统对象”的映射，因此同一对象可能通过多个硬链接名称访问。

目录也是文件系统对象，但读、写、执行权限的含义与普通文件不同。对普通文件，读允许读取内容，写允许修改内容，执行允许把它作为程序执行。对目录：

- 读允许列出目录项名称。
- 写允许在目录中创建、删除或重命名目录项，通常还需要执行权限。
- 执行又称 search，允许查找并穿过目录中的路径组件。

能读取一个文件不等于能列出其父目录；已知精确名称且对目录有 search 权限时，可能在没有目录 read 权限的情况下打开文件。反过来，能列出名称但没有 search 权限，不能通过名称访问条目。

删除文件通常是修改父目录的目录项，不是写文件内容。因此只读文件仍可能被对父目录有足够权限的进程删除，具体还受 sticky bit、ACL、扩展安全机制和平台规则影响。

## 2. 路径与解析

路径由路径组件组成。POSIX 中绝对路径以 `/` 开头，从进程根目录开始；相对路径从进程当前工作目录，或特定 API 提供的目录文件描述符开始。

```text
/srv/app/data/orders.json   绝对路径
data/orders.json            相对当前工作目录
./data/../data/input.json   包含 . 和 ..
```

路径解析逐组件查找目录项。`.` 指当前目录，`..` 指父目录；符号链接会把解析切换到链接内容指示的路径。链接形成循环或跟随过多时会失败。路径变长、组件不存在、途中组件不是目录、缺少 search 权限也会失败。

相对路径基于运行时工作目录，不基于源码文件或可执行文件所在目录。IDE、测试、systemd、容器与 shell 可能设置不同工作目录，所以配置文件路径应显式传入，或使用文档明确的配置根目录。

Go 用 `path/filepath` 处理本地操作系统路径，`path` 用于 `/` 分隔的 URL 等路径。`filepath.Join` 负责分隔符并清理词法组件；`filepath.Clean` 只做词法清理，不访问文件系统，也不会解析符号链接。

```go
joined := filepath.Join("data", "daily", "..", "orders.json")
// Unix 上为 data/orders.json
```

## 3. 规范化不等于安全约束

防止路径穿越不能只检查输入是否包含 `..`。编码、绝对路径、符号链接、挂载点以及检查后被替换的目录项都可能绕过字符串规则。

词法层的最低检查流程是：拒绝绝对输入；清理路径；拼到固定根；计算相对路径；拒绝结果为 `..` 或以 `../` 开头。但如果攻击者能在根目录内创建或替换符号链接，这仍不能保证最终对象留在根中。

更强的 Unix 实现使用以可信目录文件描述符为起点的 `openat` 系列操作，并按平台能力限制符号链接与跨挂载解析。安全目标必须明确：只防 `../`，还是还要防不可信用户并发替换目录结构。

检查 `os.Stat` 后再 `os.Open` 存在 TOCTOU：两次系统调用之间，名称指向对象可能变化。应直接执行所需打开操作并处理结果；若打开后必须验证对象，使用打开的文件描述符读取元数据。

## 4. Unix 权限位

传统模式位分 owner、group、other 三类，每类包含 read(4)、write(2)、execute/search(1)。`0750` 表示 owner 为 7，即 rwx；group 为 5，即 r-x；other 为 0，即 ---。前导 0 表示八进制写法。

```text
-rw-r-----  0640  普通文件：owner 读写，group 只读
drwxr-x---  0750  目录：owner 全部，group 可列出和穿过
```

系统先根据进程有效用户 ID、有效组 ID 与补充组确定匹配的权限类别，再检查请求操作。它不是把 owner/group/other 三组权限全部相加。ACL、capability、MAC 等扩展机制可进一步影响结果。

`umask` 清除新建对象请求模式中的位。例如程序请求创建 `0666` 文件，umask `0027` 通常得到 `0640`。程序仍应为敏感文件请求最小模式，并验证部署环境，而不是假设某个全局 umask。

特殊位包括 set-user-ID、set-group-ID 和 sticky。它们有平台与对象类型相关语义，不应在不了解安全影响时启用。共享临时目录常用 sticky 限制用户删除他人条目。

## 5. 进程及其资源

程序是可执行代码和静态数据；进程是一次运行实例。进程通常有进程 ID、父进程 ID、地址空间、线程、当前工作目录、根目录、用户与组身份、环境、打开的文件描述符和信号处理状态。

创建子进程时，父进程安排参数、环境、工作目录和标准流。POSIX `exec` 家族用新程序映像替换当前进程映像；进程 ID 通常保持，但代码、数据和栈被替换。文件描述符是否跨 exec 保留取决于 close-on-exec 标志。

进程退出会释放其私有地址空间和内核资源引用，但写入的持久性不能只由“退出了”推断。应用必须检查写入、关闭和必要的同步错误。父进程还要回收子进程退出状态；未回收的已退出子进程在 Unix 中会保留最小内核记录。

信号是异步通知机制。程序收到终止信号时可开始优雅关闭，但处理必须有截止时间，因为强制终止无法执行清理。不能只依赖退出钩子保存关键数据。

## 6. 环境变量

环境是进程启动时可见的名称和值集合。POSIX 表示为 `name=value` 字符串数组；新程序通常继承父进程提供的环境。修改当前进程环境不会追溯修改父进程，也不会自动更新已启动子进程。

```go
value, found := os.LookupEnv("APP_ENV")
if !found {
    value = "development"
}
```

`LookupEnv` 可区分变量未设置和已设置为空字符串，`Getenv` 对两者都返回空字符串。每个变量都要定义：是否必填、空字符串是否合法、类型、范围、默认值和是否允许运行中变化。

`PATH` 是可执行文件搜索目录列表。运行没有斜杠的命令名时，相关 API 可能按 PATH 依次查找。高权限或安全敏感程序不应依赖不可信 PATH；应使用可信绝对路径、受控环境，并避免 shell 解释。

环境变量不是秘密保险箱。它可能被子进程继承、诊断工具读取、崩溃报告收集或错误日志输出。密钥应通过受控秘密管理注入，限制子进程环境，永不打印完整值，并支持轮换。

环境变量只保存字符串。解析布尔值时不要把任意非空值都当 true；应接受明确集合并拒绝未知拼写。数值要检查解析错误和范围，时间要明确单位。

## 7. Go 文件 API 的关键行为

`os.Open` 只读打开，返回 `*os.File` 与 error。`os.OpenFile` 用标志明确读写、创建、截断或追加。创建模式仍会受 umask 和平台权限模型影响。

```go
file, err := os.Open(path)
if err != nil {
    return fmt.Errorf("open %q: %w", path, err)
}
defer file.Close()
```

只写 `defer file.Close()` 会忽略关闭错误。读取文件时关闭错误通常不影响已经读取的数据；写文件时关闭可能暴露延迟写入错误，应显式检查。长循环中每轮 defer 会积累到函数返回，可提取单次处理函数。

判断错误用 `errors.Is(err, fs.ErrNotExist)`、`errors.Is(err, fs.ErrPermission)` 等，不要解析系统错误文本。错误消息要附加操作和安全路径上下文，同时用 `%w` 保留底层错误。

## 8. 完整案例：受限目录文件查看器

### 8.1 契约

命令接收环境变量 `VIEWER_ROOT` 和一个相对文件路径。只允许词法上位于 root 内的普通文件，最大 1 MiB；stdout 输出内容，stderr 输出诊断。缺配置或非法路径退出 2，打开/读取失败退出 1，成功退出 0。

这份入门实现假设攻击者不能在 `VIEWER_ROOT` 中创建或替换符号链接。若根目录对不可信用户可写，必须使用目标平台的目录文件描述符安全 API 加固。

### 8.2 可运行程序

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

func resolveUnder(root, userPath string) (string, error) {
    if root == "" {
        return "", errors.New("VIEWER_ROOT is required")
    }
    if userPath == "" || filepath.IsAbs(userPath) {
        return "", errors.New("path must be a non-empty relative path")
    }
    clean := filepath.Clean(userPath)
    if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
        return "", errors.New("path escapes root")
    }
    rootAbs, err := filepath.Abs(root)
    if err != nil {
        return "", fmt.Errorf("resolve root: %w", err)
    }
    candidate := filepath.Join(rootAbs, clean)
    relative, err := filepath.Rel(rootAbs, candidate)
    if err != nil {
        return "", fmt.Errorf("compare path: %w", err)
    }
    if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
        return "", errors.New("path escapes root")
    }
    return candidate, nil
}

func copyFile(root, userPath string, output io.Writer) error {
    path, err := resolveUnder(root, userPath)
    if err != nil {
        return err
    }
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open %q: %w", userPath, err)
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return fmt.Errorf("stat opened file: %w", err)
    }
    if !info.Mode().IsRegular() {
        return fmt.Errorf("%q is not a regular file", userPath)
    }
    const max = int64(1 << 20)
    if info.Size() > max {
        return fmt.Errorf("%q exceeds %d bytes", userPath, max)
    }
    written, err := io.CopyN(output, file, max+1)
    if err != nil && !errors.Is(err, io.EOF) {
        return fmt.Errorf("read %q: %w", userPath, err)
    }
    if written > max {
        return fmt.Errorf("%q grew beyond %d bytes", userPath, max)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: viewer RELATIVE_PATH")
        os.Exit(2)
    }
    root, found := os.LookupEnv("VIEWER_ROOT")
    if !found || root == "" {
        fmt.Fprintln(os.Stderr, "viewer: VIEWER_ROOT is required")
        os.Exit(2)
    }
    if err := copyFile(root, os.Args[1], os.Stdout); err != nil {
        fmt.Fprintln(os.Stderr, "viewer:", err)
        os.Exit(1)
    }
}
```

### 8.3 正常输入与验证

```bash
mkdir -p sandbox/docs
printf 'status=ready\n' > sandbox/docs/status.txt
chmod 0640 sandbox/docs/status.txt
VIEWER_ROOT="$PWD/sandbox" go run main.go docs/status.txt
```

解析步骤：确认相对路径；清理并拼接 root；再次计算相对关系；打开对象；从已打开描述符检查是普通文件且初始大小不超限；最多复制 1 MiB+1；输出 `status=ready`，退出 0。

用 `stat` 或 `ls -l` 验证模式，但不要假设不同平台输出格式相同。用 `wc -c` 验证 stdout 精确字节数为 13。

### 8.4 失败分支

执行 `VIEWER_ROOT="$PWD/sandbox" go run main.go ../secret.txt`，在打开前返回 `path escapes root`。传 `/etc/passwd` 会因绝对路径失败。把文件权限改为 `0000` 后，以无特权且非所有者用户运行应得到 permission denied；特权用户可能仍能读取，因此测试必须记录身份。

传入目录 `docs` 会在 `Mode().IsRegular()` 失败。文件初始小于上限、读取时增长超过上限时，`CopyN` 的 `written > max` 分支阻止超限内容被当成成功；stdout 可能已有部分数据，因此调用者必须以退出码为最终成功信号，严格系统可先缓冲到受限临时区域再一次输出。

### 8.5 符号链接边界

若 `sandbox/docs/link` 指向 root 外文件，当前词法检查仍可能读取目标。可以在打开前用 `EvalSymlinks` 检查最终路径，但检查与打开之间仍有竞态。根目录不可信时，应采用 Linux `openat2` 的解析限制或同等级平台机制；这是实现约束，不是所有 POSIX 系统的通用 API。

## 9. 诊断清单

- 路径不存在：记录当前工作目录、原始输入、清理后路径和操作名称。
- 权限不足：检查进程有效 UID/GID、补充组、每级目录 search 权限、文件位与 ACL。
- 命令找错版本：查看 PATH 顺序和 `command -v` 结果，避免只看文件名。
- 文件偶发变化：检查是否按名称多次 stat/open，改为围绕已打开描述符操作。
- 配置为空：用 `LookupEnv` 区分缺失与空值，并检查父进程实际传入环境。
- 子进程泄密：构造最小环境白名单，不把整个父环境无条件转交。
- 文件描述符耗尽：检查打开/关闭路径、异常分支和长循环中的 defer。

## 10. 练习

1. 为 `resolveUnder` 写表驱动测试，覆盖空值、绝对路径、`..`、正常嵌套路径和相似前缀目录。
2. 写程序打印自身 PID、父 PID、工作目录和指定环境变量，再由父程序以不同环境启动。
3. 创建 `0640` 文件和 `0750` 目录，分别改变 owner/group/other 位并记录实际访问结果。
4. 将 viewer 改成先缓冲后输出，保证失败时 stdout 为空，并限制内存占用。
5. 在 Linux 上研究 `openat2(2)` 的 `RESOLVE_BENEATH` 与 `RESOLVE_NO_SYMLINKS`，说明它们不是 POSIX 通则。

## 来源

- [POSIX.1-2024：General Concepts—File Access Permissions 与 Pathname Resolution](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html)（访问日期：2026-07-17）
- [POSIX.1-2024：Environment Variables](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap08.html)（访问日期：2026-07-17）
- [POSIX.1-2024：open 与 openat](https://pubs.opengroup.org/onlinepubs/9799919799/functions/open.html)（访问日期：2026-07-17）
- [Go 标准库：os](https://pkg.go.dev/os)（访问日期：2026-07-17）
- [Go 标准库：path/filepath](https://pkg.go.dev/path/filepath)（访问日期：2026-07-17）
