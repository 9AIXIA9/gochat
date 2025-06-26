package middleware

import (
	"gochat/internal/domain"
	"gochat/internal/presentation"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(uc domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			presentation.ResponseError(c, domain.CodeUnauthorized, domain.CodeUnauthorized.Msg())
			c.Abort()
			return
		}

		// Bearer token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			presentation.ResponseError(c, domain.CodeInvalidToken, domain.CodeInvalidToken.Msg())
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		if uc.ParseToken(c, tokenString) {
			c.Next()
		} else {
			presentation.ResponseError(c, domain.CodeInvalidToken, domain.CodeInvalidToken.Msg())
			c.Abort()
			return
		}
	}
}
