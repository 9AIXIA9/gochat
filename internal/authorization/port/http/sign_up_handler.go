package http

import (
	"errors"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"gochat/internal/shared/kernel"
	"time"

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
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "password is invalid")
			case errors.Is(err, myErrors.ErrDuplicatedKey):
				ginutils.ResponseWithMessage(ginContext, sharedHttp.CodeInvalidParam, "the email address is used")
			default:
				zap.L().Error("sign up handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.CodeServerError)
			}
		},
		5*time.Second,
	)
}
