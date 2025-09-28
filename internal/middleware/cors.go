package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
)

// CORS 跨域处理中间件
func CORS(conf *config.CORS) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     conf.Origins,
		AllowMethods:     conf.AllowMethods,
		AllowHeaders:     conf.AllowHeaders,
		ExposeHeaders:    conf.ExposeHeaders,
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           conf.MaxAge,
	})
}
