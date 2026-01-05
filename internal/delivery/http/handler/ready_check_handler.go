package handler

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

func NewReadyCheckHandler(isReady func() bool) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		if !isReady() {
			ginutils.ResponseWithMessage(ginContext, api.CodeServiceUnavailable, "server is not ready")
			return
		}
		ginutils.ResponseWithMessage(ginContext, api.CodeSuccess, "server is ready")
	}
}
