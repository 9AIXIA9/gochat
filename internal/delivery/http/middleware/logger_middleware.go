package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.FullPath(), "/health_check") ||
			strings.Contains(c.FullPath(), "/swagger") {
			c.Next()
			return
		}

		start := time.Now().UTC()
		c.Next()
		dur := time.Since(start)
		zap.L().Info(
			"request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", dur),
			zap.Error(c.Errors.Last()),
		)
	}
}
