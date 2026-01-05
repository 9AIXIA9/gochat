package middleware

import (
	"gochat/pkg/ctxutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		c.Next()
		dur := time.Since(start)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		zap.L().Info(
			"request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", dur),
			zap.String("request_id", ctxutil.RequestIDFrom(c.Request.Context()).String()),
			zap.Error(c.Errors.Last()),
		)
	}
}
