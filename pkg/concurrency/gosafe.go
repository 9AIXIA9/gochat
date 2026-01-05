package concurrency

import (
	"runtime/debug"

	"go.uber.org/zap"
)

// GoSafe runs fn in a goroutine and logs any panic.
func GoSafe(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				zap.L().Error("panic recovered",
					zap.Any("panic", err),
					zap.String("stack", stack),
				)
			}
		}()
		fn()
	}()
}
