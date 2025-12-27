package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	"gochat/internal/shared/kernel"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type GetUserProfileRequest struct {
	UserID kernel.UserID `uri:"user_id" validate:"required"`
}

func (r *GetUserProfileRequest) Bind(ginContext *gin.Context) error {
	if err := ginContext.ShouldBindUri(r); err != nil {
		return err
	}
	return nil
}

type GetUserProfileResponseData struct {
	Profile *dto.UserProfile `json:"profile"`
}

// NewGetUserProfileHandler 获取用户资料
// @Summary      获取用户资料
// @Description  根据用户ID获取用户公开资料
// @Tags         Profile
// @Param        user_id  path      kernel.UserID                         true  "用户ID"
// @Success      200      {object}  api.Response{data=GetUserProfileResponseData}  "成功返回用户资料"
// @Router       /profiles/users/{user_id} [get]
func NewGetUserProfileHandler(useCase application.GetUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *GetUserProfileRequest) *application.GetUserProfileInput {
			return &application.GetUserProfileInput{
				UserID: request.UserID,
			}
		},
		func(ginContext *gin.Context, output *application.GetUserProfileOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &GetUserProfileResponseData{
				Profile: dto.ToUserProfileDTO(output.Profile),
			})
		},
	)
}
