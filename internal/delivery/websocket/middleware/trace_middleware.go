package middleware

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"gochat/pkg/ctxutil"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// NewTraceMiddleware starts an OTel span for each websocket message routed to a handler.
// The span name is the topic, and attributes include the topic and user_id.
func NewTraceMiddleware(serviceName string) websocket.Middleware {
	tracer := otel.Tracer(serviceName)
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			topic := websocket.GetTopic(ctx).String()
			userID := ctxutil.UserIDFrom(ctx).String()

			ctxWithSpan, span := tracer.Start(ctx, topic, trace.WithSpanKind(trace.SpanKindServer))
			defer span.End()

			span.SetAttributes(
				attribute.String("rpc.system", "websocket"),
				attribute.String("rpc.method", topic),
				attribute.String("enduser.id", userID),
				attribute.Int("message.payload_size", len(data)),
			)

			resp := next.Handle(ctxWithSpan, data)
			if err := ctxutil.ErrorFrom(ctxWithSpan); err != nil {
				span.SetStatus(codes.Error, err.Error())
			}
			return resp
		})
	}
}
