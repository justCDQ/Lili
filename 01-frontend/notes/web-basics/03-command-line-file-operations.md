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

## 实际使用

```sh
mkdir -p hello-web/assets
cd hello-web
pwd
touch index.html
cp index.html index.backup.html
mv index.backup.html backup.html
ls -la
```

## 常见错误与边界

不要在不确定目录执行 `rm -R`，不要照抄带通配符的删除命令。`mv` 到已存在目标可能覆盖；不同系统选项和默认行为有差异。`cp` 不等于版本控制，无法表达历史和合并。

## 补充知识

退出码 `0` 表示成功，非零表示失败；`command --help` 和 `man command` 查看当前系统文档。重定向 `>` 会覆盖文件，`>>` 追加；管道 `|` 把前一程序输出交给后一程序。

## 来源

- [GNU Coreutils：目录操作](https://www.gnu.org/software/coreutils/manual/html_node/Directory-listing.html)
- [GNU Coreutils：文件操作](https://www.gnu.org/software/coreutils/manual/html_node/Basic-operations.html)
- [VS Code：Getting started with the terminal](https://code.visualstudio.com/docs/terminal/getting-started)

访问日期：2026-07-16。
