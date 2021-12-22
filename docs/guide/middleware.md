# 中间件

## 指南
gocore 的 HTTP 路由服务并没有自己开发，而是使用 gin。gin 生态已经有非常完善的(中间件体系)[https://github.com/gin-contrib]。

我们没有必要重新开发这些中间件。所以 gocore 框架提供了 middleware 命令来获取这些中间件。（不能直接使用go get 的方式来获取，因为 gocore 在 gin 基础上做了一些微调）

```
[~/Documents/workspace/gocore_workspace]$ ./gocore middleware
gocore middleware

Usage:
  gocore middleware [flags]
  gocore middleware [command]

Available Commands:
  add         add middleware to app, https://github.com/gin-contrib/[middleware].git
  list        list all installed middleware
  remove      remove middleware from app

Flags:
  -h, --help   help for middleware

Use "gocore middleware [command] --help" for more information about a command.
```

## 安装

可以安装 https://github.com/gin-contrib/ 项目中的任何中间件，使用命令 `./gocore middleware add gzip`

命令会从 https://github.com/gin-contrib/gzip.git 项目中下载中间件，并且安装到 `app/http/middleware` 中。

## 查询

检查目前已经安装了哪些中间件，可以使用命令 `./gocore middleware list`

## 删除

删除某个中间件，可以使用命令 `./gocore middleware remove`
