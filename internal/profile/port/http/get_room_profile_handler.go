package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type GetRoomProfileRequest struct {
	RoomID kernel.RoomID `uri:"room_id" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
}

func (r *GetRoomProfileRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type GetRoomProfileResponseData struct {
	Profile *dto.RoomProfile `json:"profile"`
}

// NewGetRoomProfileHandler 获取房间资料
// @Summary      获取房间资料
// @Description  根据房间ID获取房间公开资料
// @Tags         Profile
// @Param        room_id  path      string     true  "房间ID"   example("019b593b-462e-74d6-bfda-0e103a172190")
// @Success      200      {object}  api.Response{data=GetRoomProfileResponseData} "成功返回房间资料"
// @Router       /profiles/rooms/{room_id} [get]
func NewGetRoomProfileHandler(useCase application.GetRoomProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetRoomProfileRequest) *application.GetRoomProfileInput {
			return &application.GetRoomProfileInput{
				RoomID: request.RoomID,
			}
		},
		func(ginContext *gin.Context, output *application.GetRoomProfileOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &GetRoomProfileResponseData{
				Profile: dto.ToRoomProfileDTO(output.Profile),
			})
		},
	)
}
