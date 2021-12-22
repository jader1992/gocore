package env

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

type GocoreTestingEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *GocoreTestingEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewGocoreTestingEnv
}

// Boot will called when the service instantiate
func (provider *GocoreTestingEnvProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *GocoreTestingEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *GocoreTestingEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

// Name define the name for this service
func (provider *GocoreTestingEnvProvider) Name() string {
	return contract.EnvKey
}
