package presentation

import (
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type JoinRoomHandler struct {
	joinRoomUsecase domain.JoinRoomUsecase
	logger          domain.Logger
}

func JoinRoomHandlerFunc(joinRoomUsecase domain.JoinRoomUsecase, logger domain.Logger) gin.HandlerFunc {
	return (&JoinRoomHandler{
		joinRoomUsecase: joinRoomUsecase,
		logger:          logger,
	}).JoinRoom
}

func (h *JoinRoomHandler) JoinRoom(c *gin.Context) {
	var req JoinRoomRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("加入房间参数绑定失败: %v", err)
		ResponseError(c, domain.CodeInvalidParam, err.Error())
		return
	}

	// 获取当前用户信息
	userID, userNumber, exists := utils.GetCurrentUser(c)
	if !exists {
		h.logger.Error("获取用户信息失败")
		ResponseError(c, domain.CodeUnauthorized, domain.CodeUnauthorized.Msg())
		return
	}

	// 解析房间号码
	roomNumber := domain.RoomNumber(0) // 这里需要根据实际情况转换字符串到RoomNumber

	// 加入房间（WebSocket连接）
	if err := h.joinRoomUsecase.JoinRoom(userNumber, roomNumber, c.Writer, c.Request); err != nil {
		h.logger.Error("加入房间失败: %v", err)
		ResponseError(c, domain.CodeRoomNotExist, domain.CodeRoomNotExist.Msg())
		return
	}

	h.logger.Info("用户 %s(%d) 加入房间成功: %s", userID, userNumber, req.RoomNumber)
}
