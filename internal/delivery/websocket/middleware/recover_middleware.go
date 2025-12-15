package middleware

import (
	"context"
	"gochat/internal/infrastructure/websocket"
	sharedHttp "gochat/internal/shared/api"
	"gochat/pkg/utils"
	"runtime/debug"

	"go.uber.org/zap"
)

func NewRecoverMiddleware() websocket.Middleware {
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *sharedHttp.Response {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					zap.L().Error("panic recovered",
						zap.Any("panic", r),
						zap.String("stack", stack),
						zap.String("user_id", utils.GetUserID(ctx).String()),
						zap.String("topic", websocket.GetTopic(ctx).String()),
					)
				}
			}()

			return next.Handle(ctx, data)
		})
	}
}
