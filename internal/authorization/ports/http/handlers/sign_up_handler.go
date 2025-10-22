package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/authorization/domain"
	myErrors "gochat/internal/shared/errors"
	httputils "gochat/internal/shared/http"
	ginutils "gochat/internal/shared/infrastructure/gin"
	"gochat/internal/shared/kernel"
	"time"
)

const DefaultRequestTimeout = 5 * time.Second

type SignUpRequest struct {
	Email    kernel.Email    `json:"email" validate:"required,email"`
	Password domain.Password `json:"password" validate:"required,min=6,max=50"`
}

func (r *SignUpRequest) Bind(ginContext *gin.Context) error {
	return ginContext.ShouldBind(r)
}

type SignUpResponseData struct {
	UserNumber domain.UserNumber `json:"user_number"`
}

func NewSignUp(useCase domain.SignUpUseCase, validator *ginutils.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[SignUpRequest, *SignUpRequest, *domain.SignUpInput, *domain.SignUpOutput](
		useCase,
		validator,
		func(request *SignUpRequest) *domain.SignUpInput {
			return &domain.SignUpInput{
				Email:    request.Email,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *domain.SignUpOutput) {
			ginutils.Response(ginContext, httputils.NewApiResponseWithData(&SignUpResponseData{UserNumber: output.UserNumber}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "password length is invalid"))
			case errors.Is(err, myErrors.ErrDuplicatedKey):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "the email address is used"))
			default:
				zap.L().Error("sign up handler failed", zap.Error(err))
				ginutils.Response(ginContext, httputils.ResponseServerError)
			}
		},
		DefaultRequestTimeout,
	)
}
