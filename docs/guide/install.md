# 安装

---
## 可执行文件

我们有两种方式来获取可执行的gocore文件，第一种是直接下载对应操作系统的gocore文件，另外一种是下载源码自己编译

### 直接下载

下载地址：
xxx

将生成的可执行文件 gocore 放到 $PATH 目录中：
`cp gocore /usr/local/bin/`

### 源码编译

下载 git 地址：`git@github.com/jianfengye/gocore:cloud/gocore.git` 到目录 gocore

在 gocore 目录中运行命令 `go run main.go build self`

将生成的可执行文件 gocore 放到 $PATH 目录中：
`cp gocore /usr/local/bin/`


## 初始化项目

使用命令 `gocore new [app]` 在当前目录创建子项目

```
[~/Documents/workspace/gocore_workspace]$ gocore new --help
create a new app

Usage:
  gocore new [app] [flags]

Aliases:
  new, create, init

Flags:
  -f, --force        if app exist, overwrite app, default: false
  -h, --help         help for new
  -m, --mod string   go mod name, default: folder name
```

你可以通过 --mod 来定义项目名字的模块名字，默认项目名，目录名，模块名是相同的

接下来，可以通过命令 `go run main.go` 看到如下信息：

```
[~/Documents/workspace/gocore_workspace]$ go run main.go
gocore commands

Usage:
  gocore [command]

Available Commands:
  app         start app serve
  build       build gocore command
  command     all about commond
  cron        about cron command
  deploy      deploy app by ssh
  dev         dev mode
  env         get current environment
  help        get help info
  middleware  gocore middleware
  new         create a new app
  provider    about gocore service provider
  swagger     swagger operator

Flags:
  -h, --help   help for gocore

Use "gocore [command] --help" for more information about a command.
```

至此，项目安装成功。
