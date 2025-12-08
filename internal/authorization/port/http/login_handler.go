package http

import (
	"errors"
	"gochat/config"
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

const RefreshTokenCookieKey = "refresh_token"

type LoginRequest struct {
	Number   kernel.UserNumber `json:"number" validate:"required,numeric"`
	Password domain.Password   `json:"password" validate:"required"`
}

type LoginResponseData struct {
	AccessToken domain.AccessToken `json:"access_token"`
}

func (r *LoginRequest) Bind(ginContext *gin.Context) error {
	return ginContext.ShouldBind(r)
}

func NewLoginHandler(useCase application.LoginUseCase, validator ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LoginRequest) *application.LoginInput {
			return &application.LoginInput{
				Number:   request.Number,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *application.LoginOutput) {
			if time.Now().UTC().After(output.RefreshToken.ExpiredAt()) {
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
				return
			}
			ginContext.SetCookie(
				RefreshTokenCookieKey,
				output.RefreshToken.Token().String(),
				int(output.RefreshToken.ExpiredAt().Sub(time.Now()).Seconds()),
				cookieConfig.Path,
				cookieConfig.Domain,
				cookieConfig.Secure,
				cookieConfig.HttpOnly,
			)
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(&LoginResponseData{AccessToken: output.AccessToken}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, domain.ErrEmptyPassword) ||
				errors.Is(err, myErrors.ErrInvalidLength) ||
				errors.Is(err, domain.ErrInvalidPassword) ||
				errors.Is(err, myErrors.ErrInvalidNumber) ||
				errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "password is invalid"))
			default:
				zap.L().Error("login handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
