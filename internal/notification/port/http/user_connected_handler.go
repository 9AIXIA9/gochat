package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/notification/application/usecase"
	"gochat/internal/notification/infrastructure/websocket"
	sharedHttp "gochat/internal/shared/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewUserConnectedHandler(
	uc usecase.UserConnectedUseCase,
	manager *websocket.Manager,
) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		userID := ginutils.GetUserID(ginContext)

		client, err := manager.Connect(userID, ginContext.Writer, ginContext.Request, nil)
		if err != nil {
			zap.L().Error("websocket connect failed", zap.Error(err))
			ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			return
		}

		client.Start()
		defer client.Close()

		if _, err := uc.Execute(ginContext.Request.Context(), &usecase.UserConnectedInput{
			UserID: userID,
		}); err != nil {
			zap.L().Error("user connected usecase execute failed", zap.Error(err))
			ginutils.Response(ginContext, sharedHttp.ResponseServerError)
		}

		client.Wait()
	}
}
