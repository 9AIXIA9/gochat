package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

// NewRateLimitMiddleware 限流中间件
func NewRateLimitMiddleware(client redis.Client, config *RateLimitConfig) gin.HandlerFunc {
	// 创建速率配置
	rate := limiter.Rate{
		Period: config.Period,
		Limit:  config.Limit,
	}

	var store limiter.Store
	if client != nil {
		s, err := redis.NewStore(client)
		if err != nil {
			log.Printf("ratelimit middleware redis store initialize failed, fallback to memory store, err: %v", err)
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
	return ginmiddleware.NewMiddleware(instance)
}
