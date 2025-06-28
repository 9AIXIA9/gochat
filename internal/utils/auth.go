package utils

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser 从上下文中获取当前用户信息
func GetCurrentUser(c *gin.Context) (userNumber domain.UserNumber) {
	userNumberInterface, exists := c.Get(domain.AuthKey)

	if !exists {
		return 0
	}

	userNumber, ok := userNumberInterface.(domain.UserNumber)

	if !ok {
		return 0
	}

	return userNumber
}
