package ctxutil

import (
	"context"
	"gochat/internal/shared/kernel"
)

const (
	userIDKey = "user_id"
	errorKey  = "error"
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
