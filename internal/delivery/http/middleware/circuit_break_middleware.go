package middleware

import (
	"fmt"
	"gochat/internal/infrastructure/breaker"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

func NewCircuitBreakMiddleware(conf *breaker.Config) gin.HandlerFunc {
	b := breaker.NewBreaker[any](conf)
	return func(ginContext *gin.Context) {
		_, err := b.Execute(func() (any, error) {
			ginContext.Next()
			// 将处理阶段的异常/5xx 映射为失败，参与熔断统计
			if len(ginContext.Errors) > 0 {
				return nil, ginContext.Errors.Last()
			}
			status := ginContext.Writer.Status()
			if status >= 500 {
				return nil, fmt.Errorf("upstream failed with status %d", status)
			}
			return nil, nil
		})

		if err != nil {
			ginutils.ResponseWithMessage(ginContext, api.CodeServiceUnavailable, "service unavailable")
			ginContext.Abort()
			return
		}
	}
}
