package http

import (
	"gochat/config"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshAccessTokenRequest struct {
	RefreshToken domain.RefreshToken `json:"-" validate:"required" example:"omkGPw-xSCtdVWIxSRP852I9dL0jPzyicgUdRArmbGI"`
}

type RefreshAccessTokenResponseData struct {
	AccessToken domain.AccessToken `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMTliNTkyOS02NWJjLTc1NDktODljNi0zZjdkYzY4NzI1NzciLCJleHAiOjE3NjY4MjY5MTMsImlhdCI6MTc2NjgyMzMxM30.ogmbXIK85Eqnh3EP_8Ttj0kkuZxsxP5wfERPSc0vNgw"`
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
// @Success      200  {object}  api.Response{data=RefreshAccessTokenResponseData}  "刷新成功，返回新的访问令牌"
// @Router       /auth/tokens/refresh [put]
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
	)
}
