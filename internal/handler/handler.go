package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
	"gochat/internal/domain"
	"time"
)

func Signup(usecase domain.SignupUsecase, timeout time.Duration) gin.HandlerFunc {
	return Adapter[SignupRequest, domain.SignupRequest](usecase, nil, timeout)
}
func Login(usecase domain.LoginUsecase, cookieConf *config.Cookie, tokenConf *config.RefreshToken, timeout time.Duration) gin.HandlerFunc {
	return Adapter[LoginRequest, domain.LoginRequest](
		usecase,
		func(c *gin.Context, resp *domain.Response) {
			// 从响应数据中提取刷新令牌
			if loginResp, ok := resp.Data.(domain.LoginResponse); ok {
				c.SetCookie("refresh_token", string(loginResp.RefreshToken), int(tokenConf.ExpireDuration.Seconds()), cookieConf.Path, cookieConf.Domain, cookieConf.Secure, cookieConf.HttpOnly)
			}
		},
		timeout,
	)
}
func RefreshToken(usecase domain.RefreshTokenUsecase, cookieConf *config.Cookie, tokenConf *config.RefreshToken, timeout time.Duration) gin.HandlerFunc {
	return Adapter[RefreshRequest, domain.RefreshTokenRequest](
		usecase,
		func(c *gin.Context, resp *domain.Response) {
			// 从响应数据中提取刷新令牌
			if refreshTokenResp, ok := resp.Data.(domain.RefreshTokenResponse); ok {
				c.SetCookie("refresh_token", string(refreshTokenResp.RefreshToken), int(tokenConf.ExpireDuration.Seconds()), cookieConf.Path, cookieConf.Domain, cookieConf.Secure, cookieConf.HttpOnly)
			}
		},
		timeout,
	)
}
func CreateRoom(usecase domain.CreateRoomUsecase, timeout time.Duration) gin.HandlerFunc {
	return Adapter[CreateRoomRequest, domain.CreateRoomRequest](usecase, nil, timeout)
}
func JoinRoom(usecase domain.JoinRoomUsecase, timeout time.Duration) gin.HandlerFunc {
	return Adapter[JoinRoomRequest, domain.JoinRoomRequest](usecase, nil, timeout)
}
func LeaveRoom(usecase domain.LeaveRoomUsecase, timeout time.Duration) gin.HandlerFunc {
	return Adapter[LeaveRoomRequest, domain.LeaveRoomRequest](usecase, nil, timeout)
}
func SendMessage(usecase domain.SendMessageUsecase, timeout time.Duration) gin.HandlerFunc {
	return Adapter[SendMessageRequest, domain.SendMessageRequest](usecase, nil, timeout)
}
