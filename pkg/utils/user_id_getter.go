package utils

import (
	"context"
	"gochat/internal/shared/kernel"
)

func GetUserIDFromCtx(ctx context.Context) kernel.UserID {
	userID, ok := ctx.Value("user_id").(kernel.UserID)
	if !ok {
		return ""
	}
	return userID
}
