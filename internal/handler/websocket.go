package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/client"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/infra/websocket/upgrader"
)

func Websocket(usecase domain.UserConnectedUsecase, m *manager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := GetUserInfo(c)

		//升级连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.L().Error("upgrade connection failed", zap.Error(err))
			ResponseError(c)
			return
		}

		userClient := client.NewClient(conn, info.UserNumber, m.DropClient)

		m.AddClient(userClient)

		//执行逻辑
		if err := usecase.Execute(c.Request.Context(), info.UserNumber); err != nil {
			zap.L().Error("user connected usecase execute failed", zap.Error(err))
			return
		}

		//保持连接
		userClient.Start()
	}
}
