package middleware

import (
	"context"
	"time"

	"gochat/internal/infrastructure/websocket"
	"gochat/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTelemetryMiddleware(serviceName string) websocket.Middleware {
	tracer := otel.Tracer(serviceName + "/websocket")
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) ([]byte, error) {
			start := time.Now()
			ctx, span := tracer.Start(ctx, utils.GetWebsocketTopic(ctx), trace.WithSpanKind(trace.SpanKindServer))
			resp, err := next.Handle(ctx, data)

			dur := time.Since(start).Seconds()
			span.SetAttributes(
				attribute.String("websocket.topic", utils.GetWebsocketTopic(ctx)),
				attribute.String("websocket.user_id", utils.GetUserID(ctx).String()),
				attribute.Float64("websocket.duration_seconds", dur),
			)
			if err != nil {
				span.RecordError(err)
			}
			span.End()
			return resp, err
		})
	}
}
