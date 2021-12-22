package services

import (
	"context"
	"errors"
	redisV8 "github.com/go-redis/redis/v8"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/redis"
	"sync"
	"time"
)

type GocoreRedisCache struct {
	container framework.Container
	client    *redisV8.Client
	lock      sync.RWMutex
}

func (r *GocoreRedisCache) Get(ctx context.Context, key string) (string, error) {
    var val string
    if err := r.GetObj(ctx, key, &val); err != nil {
        return "", err
    }
    return val, nil
}

// GetObj 获取某个key对应的对象, 对象必须实现 https://pkg.go.dev/encoding#BinaryUnMarshaler
func (r *GocoreRedisCache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := r.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisV8.Nil) {
		return ErrKeyNotFound
	}

	err := cmd.Scan(model)
	if err != nil {
		return err
	}
	return nil
}

// GetMany 获取某些key对应的值
func (r *GocoreRedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline() // 开启管道获取数据

	vals := make(map[string]string)
	cmds := make([]*redisV8.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}

	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	return vals, nil
}

// Set 设置某个key和值到缓存，带超时时间
func (r *GocoreRedisCache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

// SetObj 设置某个key和对象到缓存, 对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (r *GocoreRedisCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return r.client.Set(ctx, key, val, timeout).Err()
}

// SetMany 设置多个key和值到缓存
func (r *GocoreRedisCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipeline := r.client.Pipeline()
	cmds := make([]*redisV8.StatusCmd, 0, len(data))

	for k, v := range data {
		cmds = append(cmds, pipeline.Set(ctx, k, v, timeout))
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (r *GocoreRedisCache) SetForever(ctx context.Context, key string, val string) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *GocoreRedisCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return r.client.Set(ctx, key, val, NoneDuration).Err()
}

func (r *GocoreRedisCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

func (r *GocoreRedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *GocoreRedisCache) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc contract.RememberFunc, obj interface{}) error {
	err := r.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	objNew, err := rememberFunc(ctx, r.container)
	if err != nil {
		return err
	}

	if err := r.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := r.GetObj(ctx, key, obj); err != nil {
		return err
	}

	return nil
}

func (r *GocoreRedisCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return r.client.IncrBy(ctx, key, step).Result()
}

func (r *GocoreRedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, 1).Result()
}

func (r *GocoreRedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, -1).Result()
}

func (r *GocoreRedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *GocoreRedisCache) DelMany(ctx context.Context, keys []string) error {
	pipeline := r.client.Pipeline()
	cmds := make([]*redisV8.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipeline.Del(ctx, key))
	}
	_, err := pipeline.Exec(ctx)
	return err
}

func NewRedisCache(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	if !container.IsBind(contract.RedisKey) {
		err := container.Bind(&redis.GocoreRedisProvider{})
		if err != nil {
			return nil, err
		}
	}

	// 获取redis服务配置，并且实例化redis.Client
	redisService := container.MustMake(contract.RedisKey).(contract.RedisService)
	client, err := redisService.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil, err
	}

	// 返回RedisCache实例
	obj := &GocoreRedisCache{
		container: container,
		client:    client,
		lock:      sync.RWMutex{},
	}

	return obj, nil
}
