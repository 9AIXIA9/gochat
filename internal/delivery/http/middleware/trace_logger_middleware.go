package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TraceLoggerMiddleware enriches existing logs with trace/span IDs if present.
func TraceLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		traceID, _ := c.Get("trace_id")
		spanID, _ := c.Get("span_id")
		if traceID != nil && spanID != nil {
			zap.L().With(zap.String("trace_id", traceID.(string)), zap.String("span_id", spanID.(string))).Debug("request traced")
		}
	}
}
