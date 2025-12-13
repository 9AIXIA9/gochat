package middleware

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRecoverMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				zap.L().Error("panic recovered",
					zap.Any("panic", r),
					zap.String("stack", stack),
					zap.String("method", ginContext.Request.Method),
					zap.String("path", ginContext.Request.URL.Path),
				)
				ginutils.Response(ginContext, http.CodeServerError)
			}
		}()

		ginContext.Next()
	}
}
