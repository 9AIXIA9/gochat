package presentation

import (
	"github.com/gin-gonic/gin"
	"gochat/internal/domain"
)

// SignupHandlerFunc 注册处理函数
func SignupHandlerFunc(signupUsecase domain.SignupUsecase) gin.HandlerFunc {
	return HandlerAdapter(func(req *domain.SignupRequest) (*domain.Response, error) {
		return signupUsecase.Logic(req)
	})
}

// LoginHandlerFunc 登录处理函数
func LoginHandlerFunc(loginUsecase domain.LoginUsecase) gin.HandlerFunc {
	return HandlerAdapter(func(req *domain.LoginRequest) (*domain.Response, error) {
		return loginUsecase.Logic(req)
	})
}

// CreateRoomHandlerFunc 创建房间处理函数
func CreateRoomHandlerFunc(CreateRoomUsecase domain.CreateRoomUsecase) gin.HandlerFunc {
	return HandlerAdapter(func(req *domain.CreateRoomRequest) (*domain.Response, error) {
		return CreateRoomUsecase.Logic(req)
	})
}

// ExitRoomHandlerFunc 退出房间处理函数
func ExitRoomHandlerFunc(exitRoomUsecase domain.ExitRoomUsecase) gin.HandlerFunc {
	return HandlerAdapter(func(req *domain.ExitRoomRequest) (*domain.Response, error) {
		return exitRoomUsecase.Logic(req)
	})
}

// JoinRoomHandlerFunc 加入房间处理函数
func JoinRoomHandlerFunc(joinRoomUsecase domain.JoinRoomUsecase) gin.HandlerFunc {
	return HandlerAdapter(func(req *domain.JoinRoomRequest) (*domain.Response, error) {
		return joinRoomUsecase.Logic(req)
	})
}
