package id

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

type GocoreIDProvider struct {

}

func (provider *GocoreIDProvider) Register(container framework.Container) framework.NewInstance {
	return NewGocoreIDService
}

func (provider *GocoreIDProvider) Boot(container framework.Container) error {
	return nil
}

func (provider *GocoreIDProvider) IsDefer() bool {
	return false
}

func (provider *GocoreIDProvider) Params(container framework.Container) []interface{}  {
	return []interface{}{}
}

func (p *GocoreIDProvider) Name() string {
	return contract.IDKey
}
