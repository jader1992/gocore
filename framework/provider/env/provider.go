package env

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

type GocoreEnvProvider struct {
	Folder string
}

func (provider *GocoreEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewCoreEnv
}

func (provider *GocoreEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

func (provider *GocoreEnvProvider) IsDefer() bool {
	return false
}

func (provider *GocoreEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

func (provider *GocoreEnvProvider) Name() string {
	return contract.EnvKey
}
