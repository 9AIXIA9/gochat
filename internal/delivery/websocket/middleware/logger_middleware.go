package middleware

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"gochat/pkg/ctxutil"
	"time"

	"go.uber.org/zap"
)

// NewLoggerMiddleware returns a middleware that logs request latency
func NewLoggerMiddleware() websocket.Middleware {
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			start := time.Now().UTC()

			// Call the next handler in the chain
			resp := next.Handle(ctx, data)

			dur := time.Since(start)

			fields := []zap.Field{
				zap.String("user_id", ctxutil.UserIDFrom(ctx).String()),
				zap.String("topic", websocket.GetTopic(ctx).String()),
				zap.Duration("latency", dur),
				zap.Error(ctxutil.ErrorFrom(ctx)),
			}

			if requestID := ctxutil.RequestIDFrom(ctx).String(); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			if traceID, spanID := ctxutil.SpanIDAndTraceIDFrom(ctx); traceID != "" && spanID != "" {
				fields = append(fields,
					zap.String("trace_id", traceID),
					zap.String("span_id", spanID),
				)
			}

			zap.L().Info(
				"request completed",
				fields...,
			)
			return resp
		})
	}
}
