package http

import (
	"gochat/config"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/shared/api"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/gin-gonic/gin"
)

const RefreshTokenCookieKey = "refresh_token"

type LoginByNumberRequest struct {
	Number   kernel.UserNumber `json:"number" validate:"required,numeric" example:"2004426295315795968"`
	Password domain.Password   `json:"password" validate:"required" example:"your-password"`
}

type LoginByNumberResponseData struct {
	AccessToken domain.AccessToken `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMTliNTkyOS02NWJjLTc1NDktODljNi0zZjdkYzY4NzI1NzciLCJleHAiOjE3NjY4MjY5MTMsImlhdCI6MTc2NjgyMzMxM30.ogmbXIK85Eqnh3EP_8Ttj0kkuZxsxP5wfERPSc0vNgw"`
}

func (r *LoginByNumberRequest) Bind(ginContext *gin.Context) error {
	return ginContext.BindJSON(r)
}

// NewLoginByNumberHandler 通过Number进行用户登录
// @Summary      用户登录
// @Description  使用账号和密码登录，成功后下发访问令牌与刷新令牌（刷新令牌存于 Cookie）
// @Tags         Authorization
// @Param        request  body      LoginByNumberRequest        true  "登录请求体"
// @Success      200      {object}  api.Response{data=LoginByNumberResponseData}   "登录成功，返回访问令牌"
// @Router       /auth/login/user_number [post]
func NewLoginByNumberHandler(useCase application.LoginByNumberUseCase, validator ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LoginByNumberRequest) *application.LoginByNumberInput {
			return &application.LoginByNumberInput{
				Number:   request.Number,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *application.LoginByNumberOutput) {
			if time.Now().UTC().After(output.RefreshToken.ExpiredAt()) {
				ginutils.Response(ginContext, api.CodeServerError)
				return
			}
			ginContext.SetCookie(
				RefreshTokenCookieKey,
				output.RefreshToken.Token().String(),
				int(output.RefreshToken.ExpiredAt().Sub(time.Now().UTC()).Seconds()),
				cookieConfig.Path,
				cookieConfig.Domain,
				cookieConfig.Secure,
				cookieConfig.HttpOnly,
			)
			ginutils.ResponseSuccessWithData(ginContext, &LoginByNumberResponseData{AccessToken: output.AccessToken})
		},
	)
}
