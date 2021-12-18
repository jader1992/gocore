package static

import (
  "github.com/jader1992/gocore/framework/gin"
  "net/http"
  "os"
  "path"
  "strings"
)

// https://github.com/gin-contrib/static/blob/master/static.go

const INDEX = "index.html"

type ServeFileSystem interface {
   http.FileSystem
   Exists(prefix string, path string) bool
}

// localFileSystem 本地换件系统结构
type localFileSystem struct {
  http.FileSystem
  root string
  indexes bool
}

// LocalFile 初始化本地文件系统
// root 根路径
// indexes
func LocalFile(root string, indexes bool) *localFileSystem {
  return &localFileSystem{
    FileSystem: gin.Dir(root, indexes),
    root: root,
    indexes: indexes,
  }
}

// Exists 检查文件是否存在
func (l *localFileSystem) Exists(prefix string, filepath string) bool {
  // 去除前缀 && 长度变短
  if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
    // 获取文件路径
    name := path.Join(l.root, p)
    // 检查路径状态
    stats, err := os.Stat(name)
    if err != nil {
      return false
    }

    // 如果是一个目录
    if stats.IsDir() {
      // 寻找目录下的index.html
      if !l.indexes {
        index := path.Join(name, INDEX)
        _, err := os.Stat(index)
        if err != nil {
          return false
        }
      }
    }

    // 如果文件，则直接返回
    return true
  }
  return false
}

// 返回一个以root为根目录的本地服务
func ServerRoot(urlPrefix, root string) gin.HandlerFunc {
  return Serve(urlPrefix, LocalFile(root, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
  // 初始化文件服务
  fileserver := http.FileServer(fs)
  // 去调http.request.URL 前缀
  if urlPrefix != "" {
    fileserver = http.StripPrefix(urlPrefix, fileserver)
  }
  return func(c *gin.Context) {
    // 判断文件是否存在
    if fs.Exists(urlPrefix, c.Request.URL.Path) {
      // 如果找到，调用文件服务器 ServeHTTP 来处理这个请求
      fileserver.ServeHTTP(c.Writer, c.Request)
      // 终止后续的请求
      c.Abort()
    }
  }
}
