package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"sync"
)

// GocoreRedis 代表gocore框架的redis实现
type GocoreRedis struct {
	container framework.Container
	clients   map[string]*redis.Client

	lock *sync.RWMutex
}

func (app *GocoreRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(app.container)

	// option对opt进行修改
	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有dsn,就生成dsn
	key := config.Uniqkey()

	// 判断是否已经实例化了redis.client
	app.lock.RLock()
	if db, ok := app.clients[key]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	// 没有实例化redis,就要进行实例化操作
	app.lock.Lock()
	defer app.lock.Unlock()

	client := redis.NewClient(config.Options)

	app.clients[key] = client

	return client, nil
}

// NewGocoreRedis 代表实例化GocoreRedis
func NewGocoreRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &GocoreRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}
