package utils

import (
	"context"
	"math/rand"
	"time"
)

func BackoffWait(ctx context.Context, retry int, baseDelay time.Duration) {
	maxDelay := baseDelay * (1 << retry)
	jitter := time.Duration(rand.Int63n(int64(maxDelay) / 2))
	delay := maxDelay + jitter

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}
