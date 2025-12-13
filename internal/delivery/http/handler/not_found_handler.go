package handler

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
)

func NewNotFoundHandler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginutils.Response(ginContext, http.CodeNotFound)
	}
}
