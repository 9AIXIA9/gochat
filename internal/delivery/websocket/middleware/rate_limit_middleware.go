package middleware

import (
	"context"
	"fmt"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"gochat/pkg/ctxutil"

	"github.com/ulule/limiter/v3"
	"go.uber.org/zap"
)

// NewRateLimitMiddleware applies rate limiting to websocket routes with the shared ulule limiter.
func NewRateLimitMiddleware(limiter *limiter.Limiter) websocket.Middleware {
	if limiter == nil {
		return func(next websocket.Handler) websocket.Handler { return next }
	}

	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			key := fmt.Sprintf("%s:%s", ctxutil.UserIDFrom(ctx).String(), websocket.GetTopic(ctx).String())
			result, err := limiter.Get(ctx, key)
			if err != nil {
				zap.L().Warn("websocket ratelimit check failed, fail-open", zap.Error(err), zap.String("key", key))
				return next.Handle(ctx, data)
			}
			if result.Reached {
				return api.ResponseLimitExceeded
			}

			return next.Handle(ctx, data)
		})
	}
}
