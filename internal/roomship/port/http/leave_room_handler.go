package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type LeaveRoomRequest struct {
	UserID kernel.UserID `json:"-" validate:"required"`
	RoomID kernel.RoomID `uri:"room_id" validate:"required"`
}

func (r *LeaveRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ginutils.GetUserID(ginContext)
	r.UserID = userID
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

// NewLeaveRoomHandler 退出房间
// @Summary      退出房间
// @Description  当前登录用户退出指定房间
// @Tags         Roomship
// @Security     BearerAuth
// @Param        room_id  path      kernel.RoomID                   true  "房间ID"
// @Success      200      {object}  api.Response "退出成功，若本身不在房间视为成功"
// @Router       /rooms/{room_id}/members/me [delete]
func NewLeaveRoomHandler(useCase application.LeaveRoomUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LeaveRoomRequest) *application.LeaveRoomInput {
			return &application.LeaveRoomInput{
				UserID: request.UserID,
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
