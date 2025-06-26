package presentation

import (
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type ExitRoomHandler struct {
	exitRoomUsecase domain.ExitRoomUsecase
	logger          domain.Logger
}

func ExitRoomHandlerFunc(exitRoomUsecase domain.ExitRoomUsecase, logger domain.Logger) gin.HandlerFunc {
	return (&ExitRoomHandler{
		exitRoomUsecase: exitRoomUsecase,
		logger:          logger,
	}).ExitRoom
}

func (h *ExitRoomHandler) ExitRoom(c *gin.Context) {
	var req ExitRoomRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("退出房间参数绑定失败: %v", err)
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

	// 退出房间
	if err := h.exitRoomUsecase.ExitRoom(userID, req.RoomNumber); err != nil {
		h.logger.Error("退出房间失败: %v", err)
		ResponseError(c, domain.CodeServerBusy, domain.CodeServerBusy.Msg())
		return
	}

	h.logger.Info("用户 %s(%d) 退出房间成功: %s", userID, userNumber, req.RoomNumber)
	ResponseSuccess(c, gin.H{
		"message": "退出房间成功",
	})
}
