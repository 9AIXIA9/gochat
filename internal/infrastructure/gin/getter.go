package gin

import (
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

func GetUserID(ginContext *gin.Context) kernel.UserID {
	idAny, ok := ginContext.Get("user_id")
	if !ok {
		return kernel.EmptyUserID
	}

	id, ok := idAny.(kernel.UserID)
	if !ok {
		return kernel.EmptyUserID
	}
	return id
}
