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

type LoginByEmailRequest struct {
	Email    kernel.Email    `json:"email" validate:"required,email" example:"user12345@app.com"`
	Password domain.Password `json:"password" validate:"required" example:"your-password"`
}

type LoginByEmailResponseData struct {
	AccessToken domain.AccessToken `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIwMTliNTkyOS02NWJjLTc1NDktODljNi0zZjdkYzY4NzI1NzciLCJleHAiOjE3NjY4MjY5MTMsImlhdCI6MTc2NjgyMzMxM30.ogmbXIK85Eqnh3EP_8Ttj0kkuZxsxP5wfERPSc0vNgw"`
}

func (r *LoginByEmailRequest) Bind(ginContext *gin.Context) error {
	return ginContext.BindJSON(r)
}

// NewLoginByEmailHandler 通过Email进行用户登录
// @Summary      用户登录
// @Description  使用邮件和密码登录，成功后下发访问令牌与刷新令牌（刷新令牌存于 Cookie）
// @Tags         Authorization
// @Param        request  body      LoginByEmailRequest        true  "登录请求体"
// @Success      200      {object}  api.Response{data=LoginByEmailResponseData}   "登录成功，返回访问令牌"
// @Router       /auth/login/email [post]
func NewLoginByEmailHandler(useCase application.LoginByEmailUseCase, validator ginutils.Validator, cookieConfig *config.Cookie) gin.HandlerFunc {
	return ginutils.AdaptUseCaseToHandler(
		useCase,
		validator,
		func(request *LoginByEmailRequest) *application.LoginByEmailInput {
			return &application.LoginByEmailInput{
				Email:    request.Email,
				Password: request.Password,
			}
		},
		func(ginContext *gin.Context, output *application.LoginByEmailOutput) {
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
			ginutils.ResponseSuccessWithData(ginContext, &LoginByEmailResponseData{AccessToken: output.AccessToken})
		},
	)
}
