package app

import (
	"errors"
	"flag"
	"github.com/google/uuid"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/util"
	"path/filepath"
)

// GocoreApp 代表gocore框架的app实现
type GocoreApp struct {
	container  framework.Container // 服务容器
	baseFolder string              // 基础路径
	appId      string              // 表示当前这个app的唯一id, 可以用于分布式锁等
	configMap  map[string]string   // 配置加载
}

// Version 实现版本
func (app GocoreApp) Version() string {
	return "0.0.3"
}

func (app GocoreApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}

	// 如果参数也没有。表示默认的当前路径
	return util.GetExecDirectory()
}

// ConfigFolder 表示配置文件地址
func (app GocoreApp) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

func (app GocoreApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

func (app GocoreApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

func (app GocoreApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

func (app GocoreApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}

func (app GocoreApp) ConsoleFolder() string {
	if val, ok := app.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "console")
}

func (app GocoreApp) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

func (app GocoreApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

func (app GocoreApp) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (app GocoreApp) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "provider")
}

// AppID 表示这个App的唯一ID
func (app GocoreApp) AppId() string {
	return app.appId
}

// LoadAppConfig 加载配置文件map
func (app GocoreApp) LoadAppConfig(kv map[string]string)  {
	for key, val := range kv {
		app.configMap[key] = val
	}
}

func (app *GocoreApp) AppFolder() string {
  if val, ok := app.configMap["app_folder"]; ok {
    return val
  }
  return filepath.Join(app.BaseFolder(), "app") // 默认监控./app
}

func NewGocoreApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)

	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数，默认为当前路径")
		flag.Parse()
	}

	appId := uuid.New().String()
	configMap := map[string]string{}
	return &GocoreApp{baseFolder: baseFolder, container: container, appId: appId, configMap: configMap}, nil
}
