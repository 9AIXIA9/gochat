package utils

import (
	"math/rand"
	"time"
)

func BackoffWait(maxBackoffDuration time.Duration, times int) {
	if times <= 0 {
		return
	}
	base := 100 * time.Millisecond
	d := base << (times - 1)
	if d > maxBackoffDuration {
		d = maxBackoffDuration
	}
	// 加入 0~50% 抖动
	jitter := time.Duration(rand.Int63n(int64(d) / 2))
	time.Sleep(d/2 + jitter)
}
