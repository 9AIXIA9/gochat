package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/infra/websocket/client"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/infra/websocket/upgrader"
)

func Websocket(m *manager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := GetUserInfo(c)

		//升级连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.L().Error("upgrade connection failed", zap.Error(err))
			ResponseError(c)
		}

		userClient := client.NewClient(conn, info.UserNumber)

		m.AddClient(userClient)

		//启动服务 保持连接
		userClient.Start()
		defer func() {
			if err := m.DropClient(userClient.Number()); err != nil {
				zap.L().Error("close client failed", zap.Int64("user_number", int64(info.UserNumber)),
					zap.Error(err))
			}
		}()
	}
}
