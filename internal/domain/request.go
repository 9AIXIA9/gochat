package domain

// SignupRequest 注册请求
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Number   UserNumber `json:"number" binding:"required"`
	Password string     `json:"password" binding:"required"`
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	*AuthInfo   `json:"-"`
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Secret      string `json:"secret" binding:"max=20"`
	Description string `json:"description" binding:"max=200"`
	MaxUsers    int    `json:"max_users" binding:"required,min=2,max=100"`
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	*AuthInfo  `json:"-"`
	RoomNumber RoomNumber `uri:"number" binding:"required"`
	Secret     string     `json:"secret" binding:"max=20"`
}

// ExitRoomRequest 退出房间请求
type ExitRoomRequest struct {
	*AuthInfo  `json:"-"`
	RoomNumber RoomNumber `uri:"number" binding:"required"`
}
