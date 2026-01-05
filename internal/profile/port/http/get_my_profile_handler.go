package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	"gochat/internal/profile/dto"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	_ "gochat/internal/shared/api"

	"github.com/gin-gonic/gin"
)

type GetMyProfileRequest struct {
	UserID kernel.UserID `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
}

func (r *GetMyProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ctxutil.UserIDFrom(ginContext.Request.Context())
	return nil
}

type GetMyProfileResponseData struct {
	Profile *dto.UserProfile `json:"profile"`
}

// NewGetMyProfileHandler 获取本用户资料
// @Summary      获取本用户资料
// @Description  根据token获取本用户资料
// @Tags         Profile
// @Security     BearerAuth
// @Success      200      {object}  api.Response{data=GetMyProfileResponseData}  "成功返回用户资料"
// @Router       /profiles/me [get]
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
