package utils

import (
	"gochat/internal/domain"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser 从上下文中获取当前用户信息
func GetCurrentUser(c *gin.Context) (userID uint64, userNumber domain.UserNumber, exists bool) {
	userIDInterface, exists1 := c.Get(domain.UserIDKey)
	userNumberInterface, exists2 := c.Get(domain.UserNumberKey)

	if !exists1 || !exists2 {
		return 0, 0, false
	}

	userID, ok1 := userIDInterface.(uint64)
	userNumber, ok2 := userNumberInterface.(domain.UserNumber)

	if !ok1 || !ok2 {
		return 0, 0, false
	}

	return userID, userNumber, true
}
