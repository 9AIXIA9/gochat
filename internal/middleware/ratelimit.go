package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	rRedis "github.com/ulule/limiter/v3/drivers/store/redis"
	"gochat/internal/config"
	"log"
)

// RateLimit 限流中间件
func RateLimit(rdb *redis.Client, conf *config.RateLimit) gin.HandlerFunc {
	// 创建速率配置
	rate := limiter.Rate{
		Period: conf.Period,
		Limit:  conf.Limit,
	}

	var store limiter.Store
	if rdb != nil {
		s, err := rRedis.NewStore(rdb)
		if err != nil {
			log.Printf("redis store init failed, fallback to memory store, err: %v", err)
			store = memory.NewStore()
		} else {
			store = s
		}
	} else {
		log.Println("redis client is nil, using memory store for rate limiting")
		store = memory.NewStore()
	}

	// 创建限流实例并返回 gin 中间件
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}
