package handler

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/websocket"
	sharedHttp "gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewWebsocketHandler(server *websocket.Server) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		userID := ginutils.GetUserID(ginContext)

		if err := server.ServeWS(ginContext.Writer, ginContext.Request, userID); err != nil {
			zap.L().Error("serve websocket failed", zap.Error(err))
			ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			return
		}
	}
}
