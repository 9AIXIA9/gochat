package http

import (
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	sharedHttp "gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GetRoomProfileRequest struct {
	RoomID kernel.RoomID `uri:"room_id" validate:"required"`
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
// @Produce      json
// @Param        room_id  path      int                        true  "房间ID"
// @Success      200      {object}  GetRoomProfileResponseData "成功返回房间资料"
// @Failure      400      {object}  sharedHttp.Response "请求参数错误"
// @Failure      404      {object}  sharedHttp.Response "房间不存在"
// @Failure      500      {object}  sharedHttp.Response "服务器内部错误"
// @Router       /profile/room/{room_id} [get]
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
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrEmptyInput):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "input is empty")
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeNotFound, "room is not found")
			default:
				zap.L().Error("GetRoomProfileHandler error", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
	)
}
