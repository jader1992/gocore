package config

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type GocoreConfig struct {
	c framework.Container // 容器

	folder   string       // 文件夹
	keyDelim string       // 默认的分隔符
	lock     sync.RWMutex // 配置文件读写锁

	envMaps  map[string]string      // 所有的环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息
}

// find: 通过path获取某个元素
func (conf *GocoreConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

// IsExist check setting is exist
func (conf *GocoreConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

// Get 获取某个配置项
func (conf *GocoreConfig) Get(key string) interface{} {
	return conf.find(key)
}

// GetBool 获取bool类型配置
func (conf *GocoreConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

// GetInt 获取int类型配置
func (conf *GocoreConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *GocoreConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *GocoreConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string type
func (conf *GocoreConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *GocoreConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *GocoreConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *GocoreConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *GocoreConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *GocoreConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *GocoreConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml", // 默认的是：mapstructure
		Result:  val,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(conf.find(key)) // 此方法是 func Decode(input interface{}, output interface{}) error 具体实现
}

func (conf *GocoreConfig) loadConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	// 判断文件是否以yaml或者yml作为后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yal") {
		name := s[0]

		// 读取文件内容
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}

		// 直接对文本做文件替换
		bf = replace(bf, conf.envMaps)

		// 解析对应文件 【使用mapstructure.Decode的前提需要将内容转成map[string]interface{}】
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}

		conf.confMaps[name] = c
		conf.confRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && conf.c.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := conf.c.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

// 删除文件的操作
func (conf *GocoreConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()

	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

func NewGocoreConfig(params ...interface{}) (interface{}, error) {
	// 获取变量
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// 过程如下： 先检查配置文件夹是否存在，然后读取文件夹中的每个以 yaml 或者 yml 后缀的文件；读取之后，先用 replace 对环境变量进行一次
	// 替换；替换之后使用 go-yaml，对文件进行解析。

	// 示例化
	gocoreConf := &GocoreConfig{
		c:        container,
		folder:   envFolder,
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
		keyDelim: ".",
		lock:     sync.RWMutex{},
	}

	// 检查文件夹是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return gocoreConf, nil
	}

	//  读取每个文件
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		fileName := file.Name()
		err := gocoreConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// 监控文件夹文件

	watch, err := fsnotify.NewWatcher() // 初始化监察者
	if err != nil {
		return nil, err
	}

	err = watch.Add(envFolder) // 添加监察的内容
	if err != nil {
		return nil, err
	}

	go func() {
		// 防止进程失败
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					filename := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件： ", ev.Name)
						gocoreConf.loadConfigFile(folder, filename)
					}

					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件", ev.Name)
						gocoreConf.loadConfigFile(folder, filename)
					}

					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件", ev.Name)
						gocoreConf.removeConfigFile(folder, filename)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error: ", err)
					return
				}
			}
		}
	}()

	return gocoreConf, nil
}

// replace 用来替换配置文件中的环境变量占位符号，其格式为: env(xxx)
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	// 直接使用ReplaceAll替换。这个性能可能不是最优，但是配置文件加载，频率是比较低的，可以接受
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

// 还是用刚才的 database.mysql.password 举例，可以拆分为 3 个结构。database 去根 map 中寻找；如果有这个 key，就拿着
// mysql.password 的 path，去 database 这个 key 对应的 value 中进行寻找；而递归寻找到了最后一级 path 为 password，
// 发现这个 path 没有下一级了，就停止递归。
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// 判断是否有下个路径
	next, ok := source[path[0]]

	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next.(type) {
		case map[interface{}]interface{}: // 如果是interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}: // 如果是map[string]，直接循环调用
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil // 否则的话，返回nil
		}
	}
	return nil
}
