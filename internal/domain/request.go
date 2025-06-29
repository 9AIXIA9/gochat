package domain

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
		Number   UserNumber `json:"number,string" binding:"required"`
		Password string     `json:"password" binding:"required"`
	}
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	AuthInfo `json:"-"`
	Body     struct {
		Name        string `json:"name" binding:"required,min=1,max=50"`
		Secret      string `json:"secret" binding:"required,max=20"`
		Description string `json:"description" binding:"max=200"`
		MaxUsers    int    `json:"max_users" binding:"required,min=2,max=100"`
	}
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	AuthInfo `json:"-"`
	URI      struct {
		Number RoomNumber `uri:"number" binding:"required"`
	}
	Body struct {
		Secret string `json:"secret" binding:"required,max=20"`
	}
}

// LeaveRoomRequest 退出房间请求
type LeaveRoomRequest struct {
	AuthInfo `json:"-"`
	URI      struct {
		RoomNumber RoomNumber `uri:"number" binding:"required"`
	}
}
