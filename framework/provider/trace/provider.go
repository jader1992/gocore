package trace

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

type GocoreTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *GocoreTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewGocoreTraceService
}

// Boot will called when the service instantiate
func (provider *GocoreTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *GocoreTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *GocoreTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this service
func (provider *GocoreTraceProvider) Name() string {
	return contract.TraceKey
}
