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

// NewLoginHandler 用户登录
// @Summary      用户登录
// @Description  使用学号/账号和密码登录，成功后下发访问令牌与刷新令牌（刷新令牌存于 Cookie）
// @Tags         Authorization
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest        true  "登录请求体"
// @Success      200      {object}  api.Response{data=LoginResponseData}   "登录成功，返回访问令牌"
// @Failure      400      {object}  api.Response "请求参数错误或密码错误"
// @Failure      500      {object}  api.Response "服务器内部错误"
// @Router       /authorization/login [post]
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
			ginutils.ResponseSuccessWithData(ginContext, &LoginResponseData{AccessToken: output.AccessToken})
		},
	)
}
