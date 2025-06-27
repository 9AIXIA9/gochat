package presentation

import "gochat/internal/domain"

// SignupRequest 注册请求
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Number   domain.UserNumber `json:"number" binding:"required"`
	Password string            `json:"password" binding:"required"`
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=200"`
	MaxUsers    int    `json:"max_users" binding:"required,min=2,max=100"`
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	RoomNumber domain.RoomNumber `uri:"number" binding:"required"`
}

// ExitRoomRequest 退出房间请求
type ExitRoomRequest struct {
	RoomNumber domain.RoomNumber `uri:"number" binding:"required"`
}

// ChatMessageRequest 聊天消息请求
type ChatMessageRequest struct {
	RoomNumber domain.RoomNumber `json:"room_number" binding:"required"`
	Content    string            `json:"content" binding:"required,min=1,max=1000"`
	Type       string            `json:"type" binding:"required,oneof=text image file"`
}
