package kernel

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/gin"
)

type GocoreKernelProvider struct {
	HttpEngine *gin.Engine
}

// Register 注册服务提供者
func(sp *GocoreKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewGocoreKernelService
}

// Boot 启动的时候判断是否由外界注入了Engine，如果注入的化，用注入的，如果没有，重新实例化
func (sp *GocoreKernelProvider) Boot(c framework.Container) error  {
	if sp.HttpEngine == nil {
		sp.HttpEngine = gin.Default()
	}
	sp.HttpEngine.SetContainer(c)
	return nil
}

// IsDefer 引擎的初始化我们希望开始就进行初始化
func (sp *GocoreKernelProvider) IsDefer() bool {
	return false
}

// Params 参数就是一个HttpEngine
func (sp *GocoreKernelProvider) Params(c framework.Container) []interface{}  {
	return []interface{}{sp.HttpEngine}
}

// Name 提供凭证
func (sp *GocoreKernelProvider) Name() string {
	return contract.KernelKey
}
