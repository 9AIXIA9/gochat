package domain

type Usecase[request any] interface {
	Logic(request *request) (*Response, error)
}

type SignupUsecase interface {
	Logic(req *SignupRequest) (*Response, error)
	EncryptPwd(pwd string) (string, error)
	GenerateNumber() UserNumber
	CreateUser(user *User) (bool, error)
}

type LoginUsecase interface {
	Logic(req *LoginRequest) (*Response, error)
	QueryUser(number UserNumber) (*User, error)
	CheckPwd(origin, hash string) error
	GenerateToken(authInfo *AuthInfo) (string, error)
}
type CreateRoomUsecase interface {
	Logic(req *CreateRoomRequest) (*Response, error)
	EncryptSecret(secret string) (string, error)
	GenerateNumber() RoomNumber
	CreateRoom(room *Room) (bool, error)
}

type JoinRoomUsecase interface {
	Logic(req *JoinRoomRequest) (*Response, error)
	QueryRoom(number RoomNumber) (*Room, error)
	CheckSecret(origin, hash string) error
	JoinRoom(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}

type LeaveRoomUsecase interface {
	Logic(req *LeaveRoomRequest) (*Response, error)
	LeaveRoom(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}
