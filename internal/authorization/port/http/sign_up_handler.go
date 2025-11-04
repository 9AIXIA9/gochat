package http

import (
	"errors"
	"gochat/internal/authorization/application/usecase"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	http2 "gochat/internal/shared/http"
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
	UserNumber domain.UserNumber `json:"user_number"`
}

func NewSignUpHandler(useCase usecase.SignUpUseCase, validator *validator.Validator) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[SignUpRequest, *SignUpRequest, *usecase.SignUpInput, *usecase.SignUpOutput](
		useCase,
		validator,
		func(request *SignUpRequest) *usecase.SignUpInput {
			return &usecase.SignUpInput{
				Email:    request.Email,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *usecase.SignUpOutput) {
			ginutils.Response(ginContext, http2.NewApiResponseWithData(&SignUpResponseData{UserNumber: output.UserNumber}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.Response(ginContext, http2.NewApiResponseWithMessage(http2.CodeInvalidParam, "password length is invalid"))
			case errors.Is(err, myErrors.ErrDuplicatedKey):
				ginutils.Response(ginContext, http2.NewApiResponseWithMessage(http2.CodeInvalidParam, "the email address is used"))
			default:
				zap.L().Error("sign up handler failed", zap.Error(err))
				ginutils.Response(ginContext, http2.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
