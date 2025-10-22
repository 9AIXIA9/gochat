package gin

import (
	"github.com/gin-gonic/gin"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

const UserIDKey = "user_id"

func GetUserID(ginContext *gin.Context) (kernel.UserID, error) {
	unknownUserID, ok := ginContext.Get(UserIDKey)
	if !ok {
		return "", myErrors.ErrNotFound
	}

	userID, ok := unknownUserID.(string)
	if !ok {
		return "", myErrors.ErrInvalidCredential
	}
	return kernel.UserID(userID), nil
}
