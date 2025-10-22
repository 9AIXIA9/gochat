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
	"time"
)

type RefreshAccessTokenRequest struct {
	RefreshTokenString string `json:"-" validate:"required"`
}

func (r *RefreshAccessTokenRequest) Bind(ginContext *gin.Context) error {
	refreshTokenString, err := ginContext.Cookie(RefreshTokenCookieKey)
	if err != nil {
		return err
	}
	r.RefreshTokenString = refreshTokenString
	return nil
}

func NewRefreshAccessToken(useCase domain.RefreshAccessTokenUseCase, validator *ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[RefreshAccessTokenRequest, *RefreshAccessTokenRequest, *domain.RefreshAccessTokenInput, *domain.RefreshAccessTokenOutput](
		useCase,
		validator,
		func(request *RefreshAccessTokenRequest) *domain.RefreshAccessTokenInput {
			return &domain.RefreshAccessTokenInput{
				RefreshTokenString: request.RefreshTokenString,
			}
		},
		func(ginContext *gin.Context, output *domain.RefreshAccessTokenOutput) {
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
			ginutils.Response(ginContext, httputils.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "refresh token does not exist"))
			case errors.Is(err, myErrors.ErrExpired):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "refresh token expired"))
			case errors.Is(err, myErrors.ErrExceedMaxValue):
				ginutils.Response(ginContext, httputils.NewApiResponseWithMessage(httputils.CodeInvalidParam, "the number of token refreshes is exhausted"))
			default:
				zap.L().Error("refresh access token handler failed", zap.Error(err))
				ginutils.Response(ginContext, httputils.ResponseServerError)
			}
		},
		DefaultRequestTimeout,
	)
}
