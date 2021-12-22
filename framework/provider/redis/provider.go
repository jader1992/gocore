package redis

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
)

// GocoreRedisProvider 提供Redis的具体实现方法
type GocoreRedisProvider struct {
}

// Register 注册方法
func (h *GocoreRedisProvider) Register(container framework.Container) framework.NewInstance {
	return NewGocoreRedis
}

// Boot 启动调用
func (h *GocoreRedisProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *GocoreRedisProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (h *GocoreRedisProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (h *GocoreRedisProvider) Name() string {
	return contract.RedisKey
}
