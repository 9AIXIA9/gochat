package presentation

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type ExitRoomHandler struct {
	exitRoomUsecase domain.ExitRoomUsecase
}

func ExitRoomHandlerFunc(exitRoomUsecase domain.ExitRoomUsecase) gin.HandlerFunc {
	return (&ExitRoomHandler{
		exitRoomUsecase: exitRoomUsecase,
	}).ExitRoom
}

func (h *ExitRoomHandler) ExitRoom(c *gin.Context) {
	var req ExitRoomRequest
	if err := c.ShouldBindUri(&req); err != nil {
		ResponseError(c, domain.CodeInvalidParam, "")
		return
	}

	// 获取当前用户信息
	_, userNumber, exists := utils.GetCurrentUser(c)
	if !exists {
		zap.L().Error("获取用户信息失败")
		ResponseError(c, domain.CodeUnauthorized, domain.CodeUnauthorized.Msg())
		return
	}

	// 退出房间
	if err := h.exitRoomUsecase.ExitRoom(userNumber, req.RoomNumber); err != nil {
		ResponseError(c, domain.CodeServerBusy, "exit room failed",
			zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(req.RoomNumber)), zap.Error(err))
		return
	}

	zap.L().Info("exit room successfully", zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(req.RoomNumber)))
	ResponseSuccess(c, gin.H{
		"message": "退出房间成功",
	})
}
