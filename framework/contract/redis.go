package contract

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jader1992/gocore/framework"
)

const RedisKey = "gocore:redis"

// RedisOption 代表初始化时候的选项
type RedisOption func(container framework.Container, config *RedisConfig) error

// RedisService 标识一个redis服务
type RedisService interface {
	// GetClient 获取redis连接实例
	GetClient(option ...RedisOption) (*redis.Client, error)
}

// RedisConfig 为框架定义Redis配置接口
type RedisConfig struct {
	*redis.Options
}

// Uniqkey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) Uniqkey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
