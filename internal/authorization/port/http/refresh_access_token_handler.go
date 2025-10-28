package http

import (
	"errors"
	"gochat/config"
	"gochat/internal/authorization/application/usecase"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	myErrors "gochat/internal/shared/errors"
	http2 "gochat/internal/shared/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RefreshAccessTokenRequest struct {
	RefreshToken domain.RefreshToken `json:"-" validate:"required"`
}

func (r *RefreshAccessTokenRequest) Bind(ginContext *gin.Context) error {
	refreshToken, err := ginContext.Cookie(RefreshTokenCookieKey)
	if err != nil {
		return err
	}
	r.RefreshToken = domain.RefreshToken(refreshToken)
	return nil
}

func NewRefreshAccessTokenHandler(useCase usecase.RefreshAccessTokenUseCase, validator *ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler[RefreshAccessTokenRequest, *RefreshAccessTokenRequest, *usecase.RefreshAccessTokenInput, *usecase.RefreshAccessTokenOutput](
		useCase,
		validator,
		func(request *RefreshAccessTokenRequest) *usecase.RefreshAccessTokenInput {
			return &usecase.RefreshAccessTokenInput{
				RefreshToken: "",
			}
		},
		func(ginContext *gin.Context, output *usecase.RefreshAccessTokenOutput) {
			if time.Now().After(output.RefreshToken.ExpiredAt()) {
				ginutils.Response(ginContext, http2.ResponseServerError)
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
			ginutils.Response(ginContext, http2.ResponseSuccess)
		},
		func(ginContext *gin.Context, err error) {
			switch {
			case errors.Is(err, myErrors.ErrNotFound):
				ginutils.Response(ginContext, http2.NewApiResponseWithMessage(http2.CodeInvalidParam, "refresh token does not exist"))
			case errors.Is(err, myErrors.ErrExpired):
				ginutils.Response(ginContext, http2.NewApiResponseWithMessage(http2.CodeInvalidParam, "refresh token expired"))
			case errors.Is(err, myErrors.ErrExceedMaxValue):
				ginutils.Response(ginContext, http2.NewApiResponseWithMessage(http2.CodeInvalidParam, "the number of token refreshes is exhausted"))
			default:
				zap.L().Error("refresh access token handler failed", zap.Error(err))
				ginutils.Response(ginContext, http2.ResponseServerError)
			}
		},
		DefaultRequestTimeout,
	)
}
