package handler

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"gochat/pkg/ctxutil"

	"go.uber.org/zap"
)

func NewNotFoundHandler() websocket.HandlerFunc {
	return func(ctx context.Context, data []byte) *api.Response {
		zap.L().Warn("websocket topic not found",
			zap.String("user_id", ctxutil.UserIDFrom(ctx).String()),
			zap.String("topic", websocket.GetTopic(ctx).String()),
		)
		return api.NewResponse(api.CodeNotFound)
	}
}
