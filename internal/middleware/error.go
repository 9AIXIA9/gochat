package middleware

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
)

// Error 错误处理中间件
func Error(logger domain.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo
	}
}
