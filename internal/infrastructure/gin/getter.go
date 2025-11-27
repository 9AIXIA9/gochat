package gin

import (
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

const (
	userIDKey = "user_id"
)

func SetUserID(ginContext *gin.Context, userID kernel.UserID) {
	ginContext.Set(userIDKey, userID)
}

func GetUserID(ginContext *gin.Context) kernel.UserID {
	idAny, ok := ginContext.Get(userIDKey)
	if !ok {
		return kernel.EmptyUserID
	}

	id, ok := idAny.(kernel.UserID)
	if !ok {
		return kernel.EmptyUserID
	}
	return id
}
