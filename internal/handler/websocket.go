package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/websocket/manager"
	"gochat/internal/infra/websocket/upgrader"
)

func Websocket(usecase domain.UserConnectedUsecase, m *manager.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			zap.L().Error("upgrade connection failed", zap.Error(err))
			ResponseError(c)
			return
		}

		info := GetAuthInfo(c)
		wait := m.AddClient(conn, info.UserNumber)

		// 即使出错也不要提前返回，保持连接并确保清理
		if err := usecase.Execute(c.Request.Context(), info.UserNumber); err != nil {
			zap.L().Error("user connected usecase execute failed", zap.Error(err))
		}

		// 阻塞直到连接关闭并完成清理
		wait()
	}
}
