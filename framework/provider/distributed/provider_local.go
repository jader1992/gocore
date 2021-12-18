package distributed

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

// LocalDistributedProvider 提供App的具体实现方法
type LocalDistributedProvider struct {
}

// Register 注册HadeApp方法
func (sp *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistrubutedService
}

// Boot 启动调用
func (sp *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (sp *LocalDistributedProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (sp *LocalDistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (sp *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
