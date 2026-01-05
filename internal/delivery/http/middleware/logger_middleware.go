package middleware

import (
	"context"
	"gochat/pkg/ctxutil"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
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

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", dur),
			zap.Error(c.Errors.Last()),
		}

		if requestID := ctxutil.RequestIDFrom(c.Request.Context()).String(); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		if traceID, spanID := getSpanIDAndTraceID(c.Request.Context()); traceID != "" && spanID != "" {
			fields = append(fields,
				zap.String("trace_id", traceID),
				zap.String("span_id", spanID),
			)
		}

		if userID := ctxutil.UserIDFrom(c.Request.Context()); userID != "" {
			fields = append(fields, zap.String("user_id", userID.String()))
		}

		zap.L().Info(
			"request completed",
			fields...,
		)
	}
}

func getSpanIDAndTraceID(ctx context.Context) (string, string) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return "", ""
	}

	spanCtx := span.SpanContext()
	if !spanCtx.IsValid() {
		return "", ""
	}

	return spanCtx.TraceID().String(), spanCtx.SpanID().String()
}
