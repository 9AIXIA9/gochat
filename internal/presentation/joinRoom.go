package presentation

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type JoinRoomHandler struct {
	joinRoomUsecase domain.JoinRoomUsecase
}

func JoinRoomHandlerFunc(joinRoomUsecase domain.JoinRoomUsecase) gin.HandlerFunc {
	return (&JoinRoomHandler{
		joinRoomUsecase: joinRoomUsecase,
	}).JoinRoom
}

func (h *JoinRoomHandler) JoinRoom(c *gin.Context) {
	var req JoinRoomRequest
	if err := c.ShouldBindUri(&req); err != nil {
		ResponseError(c, domain.CodeInvalidParam, "")
		return
	}

	// 获取当前用户信息
	_, userNumber, exists := utils.GetCurrentUser(c)
	if !exists {
		ResponseError(c, domain.CodeUnauthorized, "get current user info failed")
		return
	}

	// 加入房间（WebSocket连接）
	if err := h.joinRoomUsecase.JoinRoom(userNumber, req.RoomNumber, c.Writer, c.Request); err != nil {
		ResponseError(c, domain.CodeRoomNotExist, "join room failed",
			zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(req.RoomNumber)), zap.Error(err))
		return
	}

	zap.L().Info("join room successfully", zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(req.RoomNumber)))
}
