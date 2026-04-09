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
	if instance == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return ginmiddleware.NewMiddleware(instance, ginmiddleware.WithLimitReachedHandler(func(c *gin.Context) {
		ginutils.Response(c, api.CodeLimitExceeded)
	}))
}
