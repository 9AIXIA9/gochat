package http

import (
	"errors"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SignUpRequest struct {
	Email    kernel.Email    `json:"email" validate:"required,email"`
	Password domain.Password `json:"password" validate:"required,min=6,max=50"`
}

func (r *SignUpRequest) Bind(ginContext *gin.Context) error {
	return ginContext.ShouldBind(r)
}

type SignUpResponseData struct {
	UserNumber kernel.UserNumber `json:"user_number"`
}

// NewSignUpHandler 用户注册
// @Summary      用户注册
// @Description  使用邮箱和密码注册账号，成功后返回系统分配的用户编号
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        request  body      SignUpRequest        true  "注册请求体"
// @Success      201      {object}  api.Response{data=SignUpResponseData}   "注册成功，返回用户编号"
// @Failure      400      {object}  api.Response "请求参数错误或密码不合法"
// @Failure      409      {object}  api.Response "邮箱已被占用"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /authorization/sign_up [post]
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
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, domain.ErrEmptyPassword) ||
				errors.Is(err, myErrors.ErrInvalidLength) ||
				errors.Is(err, domain.ErrInvalidPassword) ||
				errors.Is(err, myErrors.ErrInvalidFormat):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "password is invalid")
			case errors.Is(err, myErrors.ErrDuplicatedKey):
				ginutils.ResponseWithMessage(ginContext, api.CodeInvalidParam, "the email address is used")
			default:
				zap.L().Error("sign up handler failed", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
