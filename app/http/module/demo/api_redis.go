package demo

import (
    "github.com/jader1992/gocore/framework/contract"
    "github.com/jader1992/gocore/framework/gin"
    "github.com/jader1992/gocore/framework/provider/redis"
    "time"
)

func (api *Api) DemoRedis(c *gin.Context)  {
    logger := c.MustMakeLog()
    logger.Info(c, "request start", nil)

    // 初始化一个redis
    redisService := c.MustMake(contract.RedisKey).(contract.RedisService)
    client, err := redisService.GetClient(redis.WithConfigPath("cache.default"), redis.WithRedisConfig(func(options *contract.RedisConfig) {
        options.MaxRetries = 3
    }))

    if err != nil {
        logger.Error(c, err.Error(), nil)
        c.AbortWithError(50001, err)
        return
    }

    if err := client.Set(c, "foo", "bar", 1 * time.Hour).Err(); err != nil {
        c.AbortWithError(500, err)
        return
    }

    val := client.Get(c, "foo").String()
    logger.Info(c, "redis get", map[string]interface{}{
        "val": val,
    })

    if err := client.Del(c, "foo").Err(); err != nil {
        c.AbortWithError(500, err)
        return
    }

    // 获取链路追踪信息
    traceService := c.MustMake(contract.TraceKey).(contract.Trace)
    traceContext := traceService.GetTrace(c)
    traceContextMap := traceService.ToMap(traceContext)

    obj := make(map[string]interface{})
    obj["result"] = "ok"
    obj["traceContextMap"] = traceContextMap

    c.JSON(200, obj)
}

func (api *Api) DemoCache(c *gin.Context)  {
    logger := c.MustMakeLog()
    logger.Info(c, "request start", nil)

    // 初始化cache服务
    cacheService :=  c.MustMake(contract.CacheKey).(contract.CacheService)

    // 设置key为foo
    err := cacheService.Set(c, "foo", "bar", 1 * time.Hour)
    if err != nil {
        c.AbortWithError(500, err)
        return
    }

    val, err := cacheService.Get(c, "foo")
    if err != nil {
        c.AbortWithError(500, err)
        return
    }

    logger.Info(c, "cache get", map[string]interface{}{
        val: val,
    })

    if err := cacheService.Del(c, "foo"); err != nil {
        c.AbortWithError(500, err)
        return
    }
    c.JSON(200, "ok")
}
