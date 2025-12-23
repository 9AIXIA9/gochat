package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type GetMyProfileRequest struct {
	UserID kernel.UserID `json:"-" validate:"required"`
}

func (r *GetMyProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ginutils.GetUserID(ginContext)
	return nil
}

type GetMyProfileResponseData struct {
	Profile *dto.UserProfile `json:"profile"`
}

// NewGetMyProfileHandler 获取本用户资料
// @Summary      获取本用户资料
// @Description  根据token获取本用户资料
// @Tags         Profile
// @Produce      json
// @Success      200      {object}  api.Response{data=GetMyProfileResponseData}  "成功返回用户资料"
// @Failure      400      {object}  api.Response "请求参数错误"
// @Failure      404      {object}  api.Response "用户不存在"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /profile/user/{user_id} [get]
func NewGetMyProfileHandler(useCase application.GetUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetMyProfileRequest) *application.GetUserProfileInput {
			return &application.GetUserProfileInput{
				UserID: request.UserID,
			}
		},
		func(ginContext *gin.Context, output *application.GetUserProfileOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &GetMyProfileResponseData{
				Profile: dto.ToUserProfileDTO(output.Profile),
			})
		},
	)
}
