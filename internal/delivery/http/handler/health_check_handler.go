package handler

import (
	ginutils "gochat/internal/infrastructure/gin"

	"github.com/gin-gonic/gin"
)

func NewHealthCheckHandler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginutils.ResponseSuccess(ginContext)
	}
}
