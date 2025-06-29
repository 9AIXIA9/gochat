package handler

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
)

func Signup(usecase domain.SignupUsecase) gin.HandlerFunc {
	return Adapter[SignupRequest, domain.SignupRequest](usecase)
}
func Login(usecase domain.LoginUsecase) gin.HandlerFunc {
	return Adapter[LoginRequest, domain.LoginRequest](usecase)
}
func CreateRoom(usecase domain.CreateRoomUsecase) gin.HandlerFunc {
	return Adapter[CreateRoomRequest, domain.CreateRoomRequest](usecase)
}
func JoinRoom(usecase domain.JoinRoomUsecase) gin.HandlerFunc {
	return Adapter[JoinRoomRequest, domain.JoinRoomRequest](usecase)
}
func LeaveRoom(usecase domain.LeaveRoomUsecase) gin.HandlerFunc {
	return Adapter[LeaveRoomRequest, domain.LeaveRoomRequest](usecase)
}
