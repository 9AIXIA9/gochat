package websocket

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewTelemetryMiddleware(serviceName string) websocket.Middleware {
	tracer := otel.Tracer(serviceName + "/websocket")
	return websocket.MiddlewareFunc(func(reqCtx context.Context, req *websocket.Request, next websocket.Handler) (*websocket.Response, error) {
		start := time.Now()
		ctx, span := tracer.Start(reqCtx, req.RequestTopic.String(), trace.WithSpanKind(trace.SpanKindServer))

		if _, err := next.Handle(ctx, req); err != nil {
			span.RecordError(err)
		}

		dur := time.Since(start).Seconds()
		span.SetAttributes(
			attribute.String("websocket.topic", req.RequestTopic.String()),
			attribute.Float64("websocket.duration_seconds", dur),
		)
		span.End()
		return nil, nil
	})
}
