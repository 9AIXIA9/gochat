package middleware

import (
	"gochat/internal/handler"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				zap.L().Error("panic recovered",
					zap.Any("panic", r),
					zap.String("stack", stack),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)

				// 直接返回标准错误响应
				handler.ResponseError(c)
			}
		}()

		c.Next()
	}
}
