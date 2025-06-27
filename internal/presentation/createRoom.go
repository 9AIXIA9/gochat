package presentation

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateRoomHandler struct {
	createRoomUsecase domain.CreateRoomUsecase
}

func CreateRoomHandlerFunc(createRoomUsecase domain.CreateRoomUsecase) gin.HandlerFunc {
	return (&CreateRoomHandler{
		createRoomUsecase: createRoomUsecase,
	}).CreateRoom
}

func (h *CreateRoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ResponseError(c, domain.CodeInvalidParam, "")
		return
	}

	// 获取当前用户信息
	_, userNumber, exists := utils.GetCurrentUser(c)
	if !exists {
		ResponseError(c, domain.CodeUnauthorized, "get current user info failed")
		return
	}

	// 创建房间
	room, err := h.createRoomUsecase.CreateRoom(userNumber, req.Name, req.Description, req.MaxUsers)
	if err != nil {
		ResponseError(c, domain.CodeServerBusy, "create room failed", zap.Error(err),
			zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(room.Number)))
		return
	}

	zap.L().Info("create room successfully",
		zap.Int64("userNumber", int64(userNumber)), zap.Int64("roomNumber", int64(room.Number)))

	// 返回房间信息
	ResponseSuccess(c, gin.H{
		"room_number": room.Number,
		"room_name":   room.Name,
		"description": room.Description,
		"max_users":   room.MaxUsers,
		"owner":       userNumber,
	})
}
