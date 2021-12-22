package util

import (
	"io"
    "io/fs"
    "io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Exists 判断给定的文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) { // 检测已经报错的文件是否存在
			return true
		}
		return false
	}
	return true
}

// IsHiddenDIrectory 路径是否是隐藏路径
func IsHiddenDIrectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// Subdir 输出所有子目录，目录名
func Subdir(folder string) ([]string, error) {
	subs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, sub := range subs {
		if sub.IsDir() {
			ret = append(ret, sub.Name())
		}
	}
	return ret, nil
}

// DownloadFile 将文件下载到指定路径，并用该文件保存
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将内容写到你新文件中
	_, err = io.Copy(out, resp.Body)
	return err
}

// CopyFolder 将一个目录复制到另外一个目录中
func CopyFolder(source, destination string) error {
    var err error = filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
        var realPath string = strings.Replace(path, source, "", 1)
        if realPath == "" {
            return nil
        }
        if info.IsDir() {
            return os.Mkdir(filepath.Join(destination, realPath), 0755)
        } else {
            var data, err1 = ioutil.ReadFile(filepath.Join(source, realPath))
            if err1 != nil {
                return err1
            }
            return ioutil.WriteFile(filepath.Join(destination, realPath), data, 0777)
        }
    })
    return err
}

// CopyFile 将一个文件拷贝到另外一个目录中
func CopyFile(source, destination string) error  {
    var data, err1 = ioutil.ReadFile(source)
    if err1 != nil {
        return err1
    }
    return ioutil.WriteFile(destination, data, 0777)
}

