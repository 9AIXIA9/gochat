package http

import (
	"errors"
	"gochat/config"
	"gochat/internal/authorization/application/usecase"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const RefreshTokenCookieKey = "refresh_token"

type LoginRequest struct {
	Number   domain.UserNumber `json:"number" validate:"required,numeric"`
	Password domain.Password   `json:"password" validate:"required"`
}

type LoginResponseData struct {
	AccessToken domain.AccessToken `json:"access_token"`
}

func (r *LoginRequest) Bind(ginContext *gin.Context) error {
	return ginContext.ShouldBind(r)
}

func NewLoginHandler(useCase usecase.LoginUseCase, validator *ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[LoginRequest, *LoginRequest, *usecase.LoginInput, *usecase.LoginOutput](
		useCase,
		validator,
		func(request *LoginRequest) *usecase.LoginInput {
			return &usecase.LoginInput{
				Number:   request.Number,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *usecase.LoginOutput) {
			if time.Now().After(output.RefreshToken.ExpiredAt()) {
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
			case errors.Is(err, myErrors.ErrInvalidLength):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "password length is invalid"))
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "user does not exist"))
			case errors.Is(err, myErrors.ErrInvalidCredential):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "wrong account or password"))
			default:
				zap.L().Error("login handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		DefaultRequestTimeout,
	)
}
