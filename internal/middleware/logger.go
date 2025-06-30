package middleware

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logger() gin.HandlerFunc {
	return ginzap.Ginzap(zap.L(), time.RFC3339, true)
}

func Recover() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(zap.L(), true)
}
