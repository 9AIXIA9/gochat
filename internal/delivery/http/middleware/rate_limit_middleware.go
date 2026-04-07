package middleware

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

// NewRateLimitMiddleware 限流中间件
func NewRateLimitMiddleware(instance *limiter.Limiter) gin.HandlerFunc {
	// 创建限流实例并返回 gin 中间件
	return ginmiddleware.NewMiddleware(instance, ginmiddleware.WithLimitReachedHandler(func(c *gin.Context) {
		ginutils.Response(c, api.CodeLimitExceeded)
	}))
}
