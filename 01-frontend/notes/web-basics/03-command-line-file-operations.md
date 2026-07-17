# 命令行文件与目录操作

## 是什么与为什么需要

shell 读取命令并调用操作系统程序。当前工作目录决定相对路径的起点。前端开发的安装、构建、测试、Git 和服务器通常从命令行执行，因此必须能安全地定位和操作文件。

## 关键命令

```sh
pwd                    # 显示当前位置
ls -la                 # 包含隐藏项的详细列表
cd path/to/project     # 切换目录
mkdir -p src/components
cp source.txt copy.txt
cp -R assets assets-backup
mv old.txt new.txt     # 移动或重命名
rm file.txt
rm -R directory        # 递归删除，通常不可撤销
```

路径含空格时使用引号：`cd "My Project"`。以 `-` 开头的文件名可用 `rm -- -name`。执行修改前用 `pwd`、`ls` 确认目标。

## shell 怎样执行一条命令

```mermaid
flowchart LR
    A["读取命令行"] --> B["展开变量、通配符和引用"]
    B --> C["处理重定向与管道"]
    C --> D["查找命令或内建命令"]
    D --> E["创建进程并等待"]
    E --> F["得到标准输出、标准错误和退出状态"]
```

空格用于分隔参数，因此引用会直接影响程序收到的参数。单引号通常保留字面字符；双引号仍允许 `$变量` 等展开。不要把未知文本拼接成 shell 命令；在脚本或程序中应使用参数数组调用子进程。

| 数据通道 | 编号 | 默认目标 | 常见操作 |
| --- | ---: | --- | --- |
| 标准输入 | 0 | 终端键盘 | `< input.txt` |
| 标准输出 | 1 | 终端 | `> out.txt` 覆盖，`>> out.txt` 追加 |
| 标准错误 | 2 | 终端 | `2> error.log` |

管道 `|` 只把前一命令的标准输出连接到后一命令的标准输入，标准错误默认不会进入管道。

## 最小文件操作序列

```sh
mkdir -p hello-web/assets
cd hello-web
pwd
touch index.html
cp index.html index.backup.html
mv index.backup.html backup.html
ls -la
```

验证每步后再继续：

```sh
test -f index.html && printf '%s\n' 'index exists'
find . -maxdepth 2 -type f -print
printf '%s\n' 'alpha' 'beta' | sort
```

在脚本中，退出状态决定后续控制流。`&&` 只在左侧成功时运行右侧，`||` 只在左侧失败时运行右侧。交互终端中的结果不能直接证明脚本可移植；macOS 的 BSD 工具与 GNU 工具在部分选项上不同。

## 删除、覆盖与平台差异

不要在不确定目录执行 `rm -R`，不要照抄带通配符的删除命令。`mv` 到已存在目标可能覆盖；不同系统选项和默认行为有差异。`cp` 不等于版本控制，无法表达历史和合并。

## 退出状态、重定向与管道

退出码 `0` 表示成功，非零表示失败；`command --help` 和 `man command` 查看当前系统文档。重定向 `>` 会覆盖文件，`>>` 追加；管道 `|` 把前一程序输出交给后一程序。

## 批量操作前的安全检查

删除前先用只读命令列出精确目标，不以通配符猜测。复制目录时确认目标是否已存在，因为不同命令可能创建 `target/source` 或合并进 `target`。练习：建立包含空格和以连字符开头文件名的目录，完成复制、重命名、查找和清理。完成标准：每条命令都能说明当前目录、参数列表、副作用和退出状态；不使用 `sudo`；删除前有可恢复副本。

## 完整案例：整理一批静态页面而不误删文件

输入目录如下，目标是把 Markdown 源文件复制到备份目录，把 HTML 产物移动到 `dist`，同时保留日志证明实际处理了哪些文件：

```text
release candidate/
├── about.html
├── about.md
├── index.html
├── index.md
├── notes old.md
└── -draft.md
```

### 1. 建立安全前提

先进入父目录，再用引号处理空格：

```sh
cd "release candidate"
pwd
find . -maxdepth 1 -type f -print
```

`pwd` 必须显示预期目录。`find` 输出作为输入证据；如果出现子目录或符号链接，应先决定是否纳入，不能直接扩大删除范围。

### 2. 创建目标并逐项复制

```sh
mkdir -p backup dist
cp -- about.md index.md "notes old.md" ./backup/
cp -- -draft.md ./backup/
```

`--` 结束选项解析，使 `-draft.md` 被当作文件名。并非所有非 POSIX 工具都支持相同选项，执行前查本机手册。复制后检查：

```sh
find backup -maxdepth 1 -type f -print
```

预期输出有四个 Markdown 文件。若目标目录已经含同名文件，`cp` 的覆盖行为取决于实现和选项；重要数据应先比较或使用版本控制，不依赖交互提示。

### 3. 移动构建产物

```sh
mv -- about.html index.html ./dist/
find dist -maxdepth 1 -type f -print
```

同一文件系统内移动通常是重命名目录项，跨文件系统可能表现为复制后删除；中途失败时结果可能不同。移动后应验证源目录和目标目录，而不是只看命令没有报错。

### 4. 保存标准输出与标准错误

```sh
find . -maxdepth 2 -type f -print > manifest.txt 2> manifest-error.log
wc -l manifest.txt manifest-error.log
```

`>` 在命令启动前打开并截断目标文件。不要把输出重定向到同时作为输入读取的同一文件，否则可能在程序读取前清空。错误日志为空只表示该命令没有写标准错误，不证明业务结果正确。

### 5. 使用退出状态控制后续步骤

```sh
test -f dist/index.html && test -f backup/index.md && printf '%s\n' 'layout verified'
```

只有两个 `test` 都返回 0 才打印成功信息。要查看上一条命令状态，可在交互 shell 立即运行 `printf '%s\n' "$?"`；其他命令会覆盖该状态。

失败分支：`cd` 失败后继续运行会在错误目录操作，因此脚本必须在切换目录失败时停止；未引用 `notes old.md` 会变成两个参数；`*.md` 在无匹配时的行为因 shell 设置不同；`rm -R dist` 无回收站保障。修复方式是确认目录、引用变量、先列出匹配，并尽量用版本控制或临时目录提供恢复路径。

### 6. 案例输出与验收

最终结构应包含 `backup` 的四个 Markdown、`dist` 的两个 HTML 和清单日志。验收命令：

```sh
find . -maxdepth 2 -type f -print | sort
```

逐行核对输出，不把 `sort` 成功当作文件正确。完成后才考虑清理原目录；清理动作应单独执行和审查，不与复制、验证写成难以中断的一行。

## 管道与引用的进一步边界

双引号把变量展开结果保持为一个参数，例如 `cp -- "$source" "$target"`；不加引号可能发生字段分割和路径名展开。单引号完全保留 `$` 等普通字符，适合固定文本。管道中各命令可能在子进程环境执行，变量副作用不应依赖具体 shell 实现。

命令替换会去掉末尾换行，不适合安全保存任意二进制或包含换行的文件名。批量处理复杂文件名时，优先使用能以 NUL 分隔的工具组合或使用 Python 等语言的路径 API。

## 可移植脚本的最小纪律

交互终端允许人工观察并中止，脚本则会重复执行，所以要把前置条件写成检查。脚本开头应明确解释器，使用只依赖目标环境的命令，并让失败产生非零退出状态。不要假设 GNU 与 BSD 工具拥有相同长选项。

```sh
#!/bin/sh

source_dir=${1-}
target_dir=${2-}

if [ -z "$source_dir" ] || [ -z "$target_dir" ]; then
  printf '%s\n' 'usage: copy-html SOURCE_DIR TARGET_DIR' >&2
  exit 64
fi

if [ ! -d "$source_dir" ]; then
  printf 'source directory not found: %s\n' "$source_dir" >&2
  exit 66
fi

mkdir -p "$target_dir" || exit 1
find "$source_dir" -type f -name '*.html' -print
```

这个示例只列出文件，没有执行复制，因此可先核对集合。`$1` 在未提供参数时可能触发脚本策略差异，`${1-}` 为缺失参数提供空值；所有路径变量均加双引号。真正增加复制步骤前，还要决定是否递归保留目录、如何处理同名文件和符号链接。

验证脚本时至少覆盖：正常目录、带空格目录、不存在目录、空目录、不可写目标和中断恢复。只测试“理想路径”无法证明批处理安全。

## 来源

- [POSIX：Shell Command Language](https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html) — 访问日期：2026-07-17
- [GNU Coreutils：Directory listing](https://www.gnu.org/software/coreutils/manual/html_node/Directory-listing.html) — 访问日期：2026-07-17
- [GNU Coreutils：Basic operations](https://www.gnu.org/software/coreutils/manual/html_node/Basic-operations.html) — 访问日期：2026-07-17
