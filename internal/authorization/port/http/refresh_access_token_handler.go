package http

import (
	"gochat/config"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
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

// NewRefreshAccessTokenHandler 刷新访问令牌
// @Summary      刷新访问令牌
// @Description  使用刷新令牌 Cookie 刷新访问令牌，并重新设置刷新令牌 Cookie
// @Tags         Authorization
// @Produce      json
// @Success      200  {object}  api.Response{data=RefreshAccessTokenResponseData}  "刷新成功，返回新的访问令牌"
// @Failure      400  {object}  api.Response         "刷新令牌无效或已过期"
// @Failure      500  {object}  api.Response         "服务器内部错误"
// @Router       /authorization/refresh_access_token [get]
func NewRefreshAccessTokenHandler(useCase application.RefreshAccessTokenUseCase, validator ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *RefreshAccessTokenRequest) *application.RefreshAccessTokenInput {
			return &application.RefreshAccessTokenInput{
				RefreshToken: request.RefreshToken,
			}
		},
		func(ginContext *gin.Context, output *application.RefreshAccessTokenOutput) {
			if time.Now().UTC().After(output.RefreshToken.ExpiredAt()) {
				ginutils.Response(ginContext, api.CodeServerError)
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
			ginutils.ResponseSuccessWithData(
				ginContext,
				RefreshAccessTokenResponseData{
					AccessToken: output.AccessToken,
				},
			)
		},
		func(ginContext *gin.Context, err error) {
			if myErrors.IsBusinessError(err) {
				ginutils.ResponseWithMessage(ginContext, api.CodeSuccess, err.Error())
			} else {
				zap.L().Error("RefreshAccessToken failed", zap.Error(err))
				ginutils.Response(ginContext, api.CodeServerError)
			}
		},
	)
}
