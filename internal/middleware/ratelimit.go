package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gochat/internal/config"
)

// RateLimit 限流中间件
func RateLimit(conf *config.RateLimit) gin.HandlerFunc {
	//todo redis存储

	// 创建速率配置：每分钟20次请求
	rate := limiter.Rate{
		Period: conf.Period,
		Limit:  conf.Limit,
	}

	// 使用内存存储
	store := memory.NewStore()

	// 创建限流实例
	instance := limiter.New(store, rate)

	// 返回gin中间件
	return mgin.NewMiddleware(instance)
}
