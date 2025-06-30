package middleware

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"gochat/internal/types"
	"strings"
)

// JWTAuth JWT认证中间件
func JWTAuth(uc domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			handler.ResponseSuccess(c, types.UnauthorizedResponse)
			c.Abort()
			return
		}

		// Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.ResponseSuccess(c, types.InvalidTokenResponse)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		if authInfo, err := uc.ParseToken(c.Request.Context(), tokenString); err != nil {
			handler.ResponseSuccess(c, types.InvalidTokenResponse)
			c.Abort()
			return
		} else {
			c.Set(domain.AuthInfoKey, authInfo)
			c.Next()
		}
	}
}
