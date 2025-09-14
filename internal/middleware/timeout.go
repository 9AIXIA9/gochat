package middleware

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"time"
)

//todo 看一下timeout的中间件

// Timeout 超时控制中间件
func Timeout(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(timeoutResponse),
	)
}

func timeoutResponse(c *gin.Context) {
	handler.ResponseSuccess(c, domain.TimeoutResponse)
}
