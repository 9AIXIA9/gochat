package http

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	_ "gochat/internal/shared/api"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Email    kernel.Email    `json:"email" validate:"required,email" example:"youremail@demo.com"`
	Password domain.Password `json:"password" validate:"required,min=6,max=50" example:"your-password"`
}

func (r *SignUpRequest) Bind(ginContext *gin.Context) error {
	return ginContext.BindJSON(r)
}

type SignUpResponseData struct {
	UserNumber kernel.UserNumber `json:"user_number" example:"2004426295315795968"`
}

// NewSignUpHandler 用户注册
// @Summary      用户注册
// @Description  使用邮箱和密码注册账号，成功后返回系统分配的账号
// @Tags         Authorization
// @Param        request  body      SignUpRequest        true  "注册请求体"
// @Success      200      {object}  api.Response{data=SignUpResponseData}   "注册成功，返回账号"
// @Router       /auth/sign-up [post]
func NewSignUpHandler(useCase application.SignUpUseCase, validator ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *SignUpRequest) *application.SignUpInput {
			return &application.SignUpInput{
				Email:    request.Email,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *application.SignUpOutput) {
			ginutils.ResponseSuccessWithData(ginContext, &SignUpResponseData{UserNumber: output.UserNumber})
		},
	)
}
