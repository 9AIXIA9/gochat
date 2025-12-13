package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.FullPath(), "/health_check") ||
			strings.Contains(c.FullPath(), "/swagger") ||
			strings.Contains(c.FullPath(), "/metrics") {
			c.Next()
			return
		}

		start := time.Now().UTC()
		c.Next()
		dur := time.Since(start)
		span := trace.SpanFromContext(c.Request.Context())
		fields := []zap.Field{zap.Duration("latency", dur)}
		if span != nil && span.SpanContext().IsValid() {
			fields = append(fields,
				zap.String("trace_id", span.SpanContext().TraceID().String()),
				zap.String("span_id", span.SpanContext().SpanID().String()),
			)
		}
		zap.L().Info("request completed", append(fields,
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
		)...)
	}
}
