package handler

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

func NewHealthCheckHandler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginutils.ResponseWithMessage(ginContext, api.CodeSuccess, "server is healthy")
	}
}
