package utils

import (
	"runtime/debug"

	"go.uber.org/zap"
)

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
