package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetGlobalEnv(env string) {
	// 设置Gin模式，兼容 dev/production 等别名
	switch strings.ToLower(env) {
	case "dev", "development", "debug":
		gin.SetMode(gin.DebugMode)
	case "prod", "production", "release":
		gin.SetMode(gin.ReleaseMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
		zap.L().Warn(
			"unknown gin mode, fallback to debug",
			zap.String("env", env),
		)
	}
}
