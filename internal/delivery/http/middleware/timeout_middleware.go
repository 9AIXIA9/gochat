package middleware

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func NewTimeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithResponse(func(ginContext *gin.Context) {
			ginutils.Response(ginContext, http.CodeTimeout)
		}),
	)
}
