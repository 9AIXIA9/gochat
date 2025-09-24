package domain

import "time"

type ExternalRequest[request any] interface {
	ToDomain() *request
}

// SignupRequest 注册请求
type SignupRequest struct {
	Name     string
	Password string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Number   UserNumber
	Password string
}

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken RefreshToken
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	AuthInfo
	Name        string
	Secret      string
	Description string
	MaxUsers    int
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	AuthInfo
	Number RoomNumber
	Secret string
}

// LeaveRoomRequest 退出房间请求
type LeaveRoomRequest struct {
	AuthInfo
	Number RoomNumber
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	AuthInfo
	To      BaseNumber
	Content string
	SentAt  time.Time
}
