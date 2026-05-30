package timeout

import (
	"context"
)

func CheckCtxTimeout(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
