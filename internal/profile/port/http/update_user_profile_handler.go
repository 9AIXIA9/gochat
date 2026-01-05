package http

import (
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/profile/application"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"gochat/pkg/ctxutil"

	"github.com/gin-gonic/gin"
)

type UpdateUserProfileRequest struct {
	UserID      kernel.UserID      `json:"-" validate:"required" example:"019b593b-462e-74d6-bfda-0e103a172191"`
	Name        string             `json:"name" validate:"omitempty,max=32" example:"Jack"`
	Gender      kernel.Gender      `json:"gender" validate:"omitempty,oneof=0 1 2" example:"1"`
	Email       kernel.Email       `json:"email" validate:"omitempty,email" example:"user@demo.com"`
	PhoneNumber kernel.PhoneNumber `json:"phone_number" validate:"omitempty,numeric,len=11" example:"13310001000"`
	Address     kernel.Address     `json:"address" validate:"omitempty,max=128" example:"China"`
	Sign        string             `json:"sign" validate:"omitempty,max=140" example:"I'm Jack!"`
}

func (r *UpdateUserProfileRequest) Bind(ginContext *gin.Context) error {
	r.UserID = ctxutil.UserIDFrom(ginContext.Request.Context())
	// 绑定查询参数
	if err := ginContext.BindJSON(r); err != nil {
		return err
	}
	return nil
}

// NewUpdateUserProfileHandler 更新当前用户资料
// @Summary      更新当前用户资料
// @Description  更新当前登录用户的基础资料（昵称、邮箱、电话等）
// @Tags         Profile
// @Security     BearerAuth
// @Param        request  body      UpdateUserProfileRequest  true  "更新用户资料请求体"
// @Success      200      {object}  api.Response    "更新成功"
// @Router       /profiles/me [put]
func NewUpdateUserProfileHandler(useCase application.UpdateUserProfileUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *UpdateUserProfileRequest) *application.UpdateUserProfileInput {
			return &application.UpdateUserProfileInput{
				UserID:      request.UserID,
				Name:        request.Name,
				Gender:      request.Gender,
				Email:       request.Email,
				PhoneNumber: request.PhoneNumber,
				Address:     request.Address,
				Sign:        request.Sign,
			}
		},
		func(ginContext *gin.Context, _ *kernel.NoOutput) {
			ginutils.ResponseSuccess(ginContext)
		},
	)
}
