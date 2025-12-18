package utils

import (
	"context"
	"gochat/internal/shared/kernel"
)

const (
	userIDKey = "user_id"
	errorKey  = "error"
)

func SetError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

func GetError(ctx context.Context) error {
	err, ok := ctx.Value(errorKey).(error)
	if !ok {
		return nil
	}
	return err
}

func SetUserID(ctx context.Context, userID kernel.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) kernel.UserID {
	userID, ok := ctx.Value(userIDKey).(kernel.UserID)
	if !ok {
		return ""
	}
	return userID
}
