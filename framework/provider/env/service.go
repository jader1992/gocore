package env

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/jader1992/gocore/framework/contract"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GocoreEnv struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 保存所在的环境变量
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *GocoreEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

// Get 获取某个环境变量，如果没有设置，返回""
func (en *GocoreEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// IsExist 判断一个环境变量是否有被设置
func (en *GocoreEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// All 获取所有的环境变量，.env和运行环境变量融合后
func (en *GocoreEnv) All() map[string]string {
	return en.maps
}

// NewCoreEnv 有一个参数，.env文件所在的目录
// example: NewCoreEnv("/envfolder/") 会读取文件: /envfolder/.env
// .env的文件格式 FOO_ENV=BAR
func NewCoreEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewEnv param error")
	}

	// 读取folder文件
	folder := params[0].(string)

	// 实例化
	gocoreEnv := &GocoreEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// 解析folder/.env文件
	file := filepath.Join(folder, ".env")

	// 打开文件 .env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		br := bufio.NewReader(fi)
		for {
			// 按行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// 按照等号解析：形如： key=val
			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}
			// 保存map
			key := string(s[0])
			val := string(s[1])
			gocoreEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		gocoreEnv.maps[pair[0]] = pair[1]
	}

	// 返回实例
	return gocoreEnv, nil
}
