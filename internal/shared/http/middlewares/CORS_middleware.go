package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域处理中间件
func CORS(config *CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.Origins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	})
}
