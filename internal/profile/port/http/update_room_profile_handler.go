package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type UpdateRoomProfileRequest struct {
	UserID       kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	RoomID       kernel.RoomID `uri:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172192"`
	Name         string        `json:"name" example:"family"`
	Introduction string        `json:"introduction" example:"This is our family chat room."`
}

func (r *UpdateRoomProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	// 绑定查询参数
	if err := ginContext.BindJSON(r); err != nil {
		return err
	}
	if err := ginContext.BindUri(r); err != nil {
		return err
	}
	return nil
}

// NewUpdateRoomProfileHandler 更新房间资料
// @Summary      更新房间资料
// @Description  更新指定房间的资料（名称、简介等），需要房主身份
// @Tags         Profile
// @Security     BearerAuth
// @Param        request  body      UpdateRoomProfileRequest  true  "更新房间资料请求体"
// @Param   	 room_id  path      string                    true  "房间ID"  example("019b593b-462e-74d6-bfda-0e103a172192")
// @Success      200      {object}  api.Response    "更新成功"
// @Router       /profiles/rooms/{room_id} [put]
func NewUpdateRoomProfileHandler(useCase application.UpdateRoomProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *UpdateRoomProfileRequest) *application.UpdateRoomProfileInput {
			return &application.UpdateRoomProfileInput{
				UserID:       request.UserID,
				RoomID:       request.RoomID,
				Name:         request.Name,
				Introduction: request.Introduction,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
