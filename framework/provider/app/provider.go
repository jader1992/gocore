package app

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

// GocoreAppProvider 提供了App的具体实现方法
type GocoreAppProvider struct {
	BaseFolder string
}

func (p *GocoreAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewGocoreApp
}

func (p *GocoreAppProvider) Boot(container framework.Container) error {
	return nil
}

func (p *GocoreAppProvider) IsDefer() bool {
	return false
}

func (p *GocoreAppProvider) Params(container framework.Container) []interface{}  {
	return []interface{}{container, p.BaseFolder}
}

func (p *GocoreAppProvider) Name() string {
	return contract.AppKey
}


