package cache

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/cache/services"
	"strings"
)

// GocoreCacheProvider 缓存服务提供者
type GocoreCacheProvider struct {
	framework.Container

	Driver string // Driver
}

func (p *GocoreCacheProvider) Register(c framework.Container) framework.NewInstance {
	if p.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewMemoryCache
		}

		cs := tcs.(contract.Config)
		p.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	switch p.Driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return services.NewMemoryCache
	default:
		return services.NewMemoryCache
	}
}

func (p *GocoreCacheProvider) Boot(c framework.Container) error {
	return nil
}

func (p *GocoreCacheProvider) IsDefer() bool {
	return false
}

func (p *GocoreCacheProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (p *GocoreCacheProvider) Name() string {
	return contract.CacheKey
}
