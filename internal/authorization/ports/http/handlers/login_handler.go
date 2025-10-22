package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/config"
	myErrors "gochat/internal/shared/errors"
	httputils "gochat/internal/shared/http"
	ginutils "gochat/internal/shared/infrastructure/gin"
	"gochat/internal/shared/kernel"
	"time"
)

const RefreshTokenCookieKey = "refresh_token"

type LoginRequest struct {
	Number   domain.UserNumber `json:"number" validate:"required,numeric"`
	Password domain.Password   `json:"password" validate:"required"`
}

type LoginResponseData struct {
	AccessToken kernel.AccessToken `json:"access_token"`
}

func (r *LoginRequest) Bind(ginContext *gin.Context) error {
	return ginContext.ShouldBind(r)
}

func NewLogin(useCase domain.LoginUseCase, validator *ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[LoginRequest, *LoginRequest, *domain.LoginInput, *domain.LoginOutput](
		useCase,
		validator,
		func(request *LoginRequest) *domain.LoginInput {
			return &domain.LoginInput{
				Number:   request.Number,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *domain.LoginOutput) {
			if time.Now().After(output.RefreshToken.ExpiredAt()) {
				ginutils.Response(ginContext, httputils.ResponseServerError)
				return
			}
			ginContext.SetCookie(
				RefreshTokenCookieKey,
				output.RefreshToken.String(),
				int(output.RefreshToken.ExpiredAt().Sub(time.Now()).Seconds()),
				cookieConfig.Path,
				cookieConfig.Domain,
				cookieConfig.Secure,
				cookieConfig.HttpOnly,
			)
			ginutils.Response(ginContext, httputils.NewApiResponseWithData(&LoginResponseData{AccessToken: output.AccessToken}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "password length is invalid"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "user does not exist"))
			case errors.Is(err, myErrors.ErrInvalidCredential):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "wrong account or password"))
			default:
				zap.L().Error("login handler failed", zap.Error(err))
				ginutils.Response(ginContext, httputils.ResponseServerError)
			}
		},
		DefaultRequestTimeout,
	)
}
