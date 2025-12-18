package middleware

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"gochat/pkg/utils"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// NewLoggerMiddleware returns a middleware that logs request latency and tracing info.
func NewLoggerMiddleware() websocket.Middleware {
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			start := time.Now().UTC()

			// Call the next handler in the chain
			resp := next.Handle(ctx, data)

			dur := time.Since(start)
			span := trace.SpanFromContext(ctx)
			fields := []zap.Field{zap.Duration("latency", dur)}
			if span != nil && span.SpanContext().IsValid() {
				fields = append(fields,
					zap.String("trace_id", span.SpanContext().TraceID().String()),
					zap.String("span_id", span.SpanContext().SpanID().String()),
				)
			}

			zap.L().Info("websocket request completed", append(fields,
				zap.String("user_id", utils.GetUserID(ctx).String()),
				zap.String("topic", websocket.GetTopic(ctx).String()),
				zap.Error(utils.GetError(ctx)),
			)...)

			return resp
		})
	}
}
