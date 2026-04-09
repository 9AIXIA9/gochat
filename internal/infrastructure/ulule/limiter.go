package ulule

import (
	"log"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewLimiter(client redis.Client, config *Config) *limiter.Limiter {
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

	return limiter.New(store, rate)
}
