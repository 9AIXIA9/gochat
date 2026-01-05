package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/roomship/application"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"github.com/gin-gonic/gin"
)

type LeaveRoomRequest struct {
	UserID kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RoomID kernel.RoomID `uri:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172192"`
}

func (r *LeaveRoomRequest) Bind(ginContext *gin.Context) error {
	userID := ctxutil.UserIDFrom(ginContext.Request.Context())
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
// @Param        room_id  path     string   true  "房间ID"    example(019b593b-462e-74d6-bfda-0e103a172190)
// @Success      200      {object}  api.Response "主动退出指定房间，幂等设计：不在房间也返回成功"
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
