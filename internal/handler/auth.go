package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
)

// injectAuthInfo 注入认证信息
func injectAuthInfo(c *gin.Context, req interface{}) {
	if authReq, ok := req.(domain.AuthContext); ok {
		info := GetAuthInfo(c)
		if info != nil {
			authReq.SetInfo(info)
		}
	}
}

// GetAuthInfo 从上下文中获取当前用户信息
func GetAuthInfo(c *gin.Context) *domain.AuthInfo {
	authInfoInterface, exists := c.Get(domain.AuthInfoKey)

	if !exists {
		return nil
	}

	authInfo, ok := authInfoInterface.(*domain.AuthInfo)

	if !ok {
		return nil
	}

	return authInfo
}
