package handler

import (
	"gochat/internal/domain"
)

// SignupRequest 注册请求
type SignupRequest struct {
	Body struct {
		Name     string `json:"name" binding:"required,min=2,max=20"`
		Password string `json:"password" binding:"required,min=6,max=50"`
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Body struct {
		Number   domain.UserNumber `json:"number,string" binding:"required"`
		Password string            `json:"password" binding:"required"`
	}
}

// RefreshRequest 更新令牌请求
type RefreshRequest struct {
	Cookie struct {
		RefreshToken domain.RefreshToken `json:"refresh_token"`
	}
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	domain.AuthInfo `json:"-"`
	Body            struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Secret      string `json:"secret" binding:"required,max=20"`
		Description string `json:"description" binding:"max=200"`
		MaxUsers    int    `json:"max_users" binding:"required,min=2,max=100"`
	}
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	domain.AuthInfo `json:"-"`
	URI             struct {
		Number domain.RoomNumber `uri:"number,string" binding:"required"`
	}
	Body struct {
		Secret string `json:"secret" binding:"required,max=20"`
	}
}

// LeaveRoomRequest 退出房间请求
type LeaveRoomRequest struct {
	domain.AuthInfo `json:"-"`
	URI             struct {
		RoomNumber domain.RoomNumber `uri:"number,string" binding:"required"`
	}
}

// SendMessageRequest 发送信息请求
type SendMessageRequest struct {
	domain.AuthInfo `json:"-"`
	URI             struct {
		To domain.BaseNumber `uri:"to,string" binding:"required"`
	}
	Body struct {
		Content string `json:"content" binding:"required"`
		SentAt  int64  `json:"sent_at" binding:"required"`
	}
}
