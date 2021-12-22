package config

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"path/filepath"
)

type GocoreConfigProvider struct{}

func (provider *GocoreConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewGocoreConfig
}

func (provider *GocoreConfigProvider) Boot(c framework.Container) error {
	return nil
}

func (provider *GocoreConfigProvider) IsDefer() bool {
	return false
}

func (provider *GocoreConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv() // 获取env中的APP_ENV
	// 配置文件夹地址
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

func (provider *GocoreConfigProvider) Name() string {
	return contract.ConfigKey
}
