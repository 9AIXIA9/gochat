package presentation

import (
	"gochat/internal/domain"
	"gochat/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateRoomHandler struct {
	createRoomUsecase domain.CreateRoomUsecase
	logger            domain.Logger
}

func CreateRoomHandlerFunc(createRoomUsecase domain.CreateRoomUsecase, logger domain.Logger) gin.HandlerFunc {
	return (&CreateRoomHandler{
		createRoomUsecase: createRoomUsecase,
		logger:            logger,
	}).CreateRoom
}

func (h *CreateRoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("创建房间参数绑定失败: %v", err)
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

	// 创建房间
	room, err := h.createRoomUsecase.CreateRoom(userNumber, req.Name, req.Description, req.MaxUsers)
	if err != nil {
		h.logger.Error("创建房间失败: %v", err)
		ResponseError(c, domain.CodeServerBusy, domain.CodeServerBusy.Msg())
		return
	}

	h.logger.Info("用户 %s(%d) 创建房间成功: %s", userID, userNumber, room.Number)

	// 返回房间信息
	ResponseSuccess(c, gin.H{
		"room_number": room.Number,
		"room_name":   room.Name,
		"description": room.Description,
		"max_users":   room.MaxUsers,
		"owner":       userNumber,
	})
}
