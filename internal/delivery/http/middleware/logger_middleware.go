package middleware

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewLoggerMiddleware() gin.HandlerFunc {
	return ginzap.Ginzap(zap.L(), time.RFC3339, true)
}
