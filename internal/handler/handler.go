package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/config"
	"gochat/internal/domain"
)

func Signup(usecase domain.SignupUsecase) gin.HandlerFunc {
	return Adapter[SignupRequest, domain.SignupRequest](usecase, nil)
}
func Login(conf *config.RefreshToken, usecase domain.LoginUsecase) gin.HandlerFunc {
	return Adapter[LoginRequest, domain.LoginRequest](
		usecase,
		func(c *gin.Context, resp *domain.Response) {
			// 从响应数据中提取刷新令牌
			if loginResp, ok := resp.Data.(domain.LoginResponse); ok {
				c.SetCookie("refresh_token", string(loginResp.RefreshToken), int(conf.ExpireDuration.Seconds()), "/", "", true, true)
			}
		},
	)
}
func RefreshToken(conf *config.RefreshToken, usecase domain.RefreshTokenUsecase) gin.HandlerFunc {
	return Adapter[RefreshRequest, domain.RefreshTokenRequest](
		usecase,
		func(c *gin.Context, resp *domain.Response) {
			// 从响应数据中提取刷新令牌
			if refreshTokenResp, ok := resp.Data.(domain.RefreshTokenResponse); ok {
				c.SetCookie("refresh_token", string(refreshTokenResp.RefreshToken), int(conf.ExpireDuration.Seconds()), "/", "", true, true)
			}
		},
	)
}
func CreateRoom(usecase domain.CreateRoomUsecase) gin.HandlerFunc {
	return Adapter[CreateRoomRequest, domain.CreateRoomRequest](usecase, nil)
}
func JoinRoom(usecase domain.JoinRoomUsecase) gin.HandlerFunc {
	return Adapter[JoinRoomRequest, domain.JoinRoomRequest](usecase, nil)
}
func LeaveRoom(usecase domain.LeaveRoomUsecase) gin.HandlerFunc {
	return Adapter[LeaveRoomRequest, domain.LeaveRoomRequest](usecase, nil)
}
func SendMessage(usecase domain.SendMessageUsecase) gin.HandlerFunc {
	return Adapter[SendMessageRequest, domain.SendMessageRequest](usecase, nil)
}
