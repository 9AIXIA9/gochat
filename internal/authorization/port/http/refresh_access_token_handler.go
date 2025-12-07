package http

import (
	"errors"
	"gochat/config"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/validator"
	myErrors "gochat/internal/shared/errors"
	sharedHttp "gochat/internal/shared/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefreshAccessTokenRequest struct {
	RefreshToken domain.RefreshToken `json:"-" validate:"required"`
}

type RefreshAccessTokenResponseData struct {
	AccessToken domain.AccessToken `json:"access_token"`
}

func (r *RefreshAccessTokenRequest) Bind(ginContext *gin.Context) error {
	refreshToken, err := ginContext.Cookie(RefreshTokenCookieKey)
	if err != nil {
		return err
	}
	r.RefreshToken = domain.RefreshToken(refreshToken)
	return nil
}

func NewRefreshAccessTokenHandler(useCase application.RefreshAccessTokenUseCase, validator *validator.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *RefreshAccessTokenRequest) *application.RefreshAccessTokenInput {
			return &application.RefreshAccessTokenInput{
				RefreshToken: request.RefreshToken,
			}
		},
		func(ginContext *gin.Context, output *application.RefreshAccessTokenOutput) {
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
			ginutils.Response(ginContext, sharedHttp.NewApiResponseWithData(RefreshAccessTokenResponseData{AccessToken: output.AccessToken}))
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, domain.ErrRefreshTokenExpired) ||
				errors.Is(err, domain.ErrAccessTokenGenerated) ||
				errors.Is(err, domain.ErrRefreshLimitExceeded) ||
				errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, sharedHttp.NewApiResponseWithMessage(sharedHttp.CodeInvalidParam, "invalid refresh token"))
			default:
				zap.L().Error("refresh access token handler failed", zap.Error(err))
				ginutils.Response(ginContext, sharedHttp.ResponseServerError)
			}
		},
		5*time.Second,
	)
}
