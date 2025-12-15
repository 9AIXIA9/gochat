package middleware

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	sharedHttp "gochat/internal/shared/api"
	"gochat/pkg/utils"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTelemetryMiddleware(serviceName string) websocket.Middleware {
	tracer := otel.Tracer(serviceName + "/websocket")
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *sharedHttp.Response {
			start := time.Now().UTC()
			ctx, span := tracer.Start(ctx, websocket.GetTopic(ctx).String(), trace.WithSpanKind(trace.SpanKindServer))
			resp := next.Handle(ctx, data)

			dur := time.Since(start).Seconds()
			span.SetAttributes(
				attribute.String("websocket.topic", websocket.GetTopic(ctx).String()),
				attribute.String("websocket.user_id", utils.GetUserID(ctx).String()),
				attribute.Float64("websocket.duration_seconds", dur),
			)
			span.End()
			return resp
		})
	}
}
