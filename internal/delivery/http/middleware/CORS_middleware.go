package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORSMiddleware 跨域处理中间件
func NewCORSMiddleware(config *CORSConfig) gin.HandlerFunc {
	allowOrigins := config.AllowOrigins
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	})
}
