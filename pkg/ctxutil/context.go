package ctxutil

import (
	"context"
	"gochat/internal/shared/kernel"

	"go.opentelemetry.io/otel/trace"
)

const (
	userIDKey    = "user_id"
	requestIDKey = "request_id"
	errorKey     = "error"
)

// WithError returns a new context carrying an error value.
func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

// ErrorFrom extracts an error value from context, if present.
func ErrorFrom(ctx context.Context) error {
	if err, ok := ctx.Value(errorKey).(error); ok {
		return err
	}
	return nil
}

// WithUserID returns a new context carrying a user id.
func WithUserID(ctx context.Context, userID kernel.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFrom extracts user id from context if present.
func UserIDFrom(ctx context.Context) kernel.UserID {
	if userID, ok := ctx.Value(userIDKey).(kernel.UserID); ok {
		return userID
	}
	return ""
}

func WithRequestID(ctx context.Context, requestID kernel.OperationID) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFrom(ctx context.Context) kernel.OperationID {
	if requestID, ok := ctx.Value(requestIDKey).(kernel.OperationID); ok {
		return requestID
	}
	return ""
}

func SpanIDAndTraceIDFrom(ctx context.Context) (string, string) {
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
