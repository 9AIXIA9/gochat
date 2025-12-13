package utils

import (
	"context"
	"gochat/internal/shared/kernel"
)

const (
	userIDKey = "user_id"
)

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
