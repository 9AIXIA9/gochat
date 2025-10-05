package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/handler"
	"gochat/internal/infra/websocket/client"
	"gochat/internal/infra/websocket/upgrader"
	"strings"
)

const authQueryKey = "token"

// WebsocketJWTAuth HTTPS JWT认证中间件
func HTTPSJWTAuth(uc domain.AuthUsecase) gin.HandlerFunc {
	return JWTAuth(uc, func(c *gin.Context) {
		handler.ResponseSuccess(c, domain.InvalidTokenResponse)
	})
}

// WebsocketJWTAuth Websocket JWT认证中间件
func WebsocketJWTAuth(uc domain.AuthUsecase) gin.HandlerFunc {
	return JWTAuth(uc, func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.L().Error("upgrade connection failed", zap.Error(err))
			handler.ResponseError(c)
			return
		}

		conn.WriteJSON(&client.Message{
			Type: client.AuthType,
			Data: nil,
		})

		conn.Close()
	})
}

// JWTAuth JWT认证中间件
func JWTAuth(uc domain.AuthUsecase, response func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		var tokenStr string
		if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
			// Bearer token格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response(c)
				c.Abort()
				return
			}
			tokenStr = parts[1]
		} else if authQuery := c.Query(authQueryKey); authQuery != "" {
			tokenStr = authQuery
		} else {
			response(c)
			c.Abort()
			return
		}

		// 解析token
		if authInfo, err := uc.ParseAuthToken(c.Request.Context(), domain.AuthToken(tokenStr)); err != nil {
			response(c)
			c.Abort()
			return
		} else {
			c.Set(domain.AuthInfoKey, authInfo)
			c.Next()
		}
	}
}
