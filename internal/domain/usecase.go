package domain

import "context"

type Usecase[request any] interface {
	Logic(ctx context.Context, request *request) (*Message, error)
}

type SignupUsecase interface {
	Logic(ctx context.Context, req *SignupRequest) (*Message, error)
	EncryptPwd(ctx context.Context, pwd string) (string, error)
	GenerateNumber(ctx context.Context) (UserNumber, error)
	CreateUser(ctx context.Context, user *User) (bool, error)
}

type LoginUsecase interface {
	Logic(ctx context.Context, req *LoginRequest) (*Message, error)
	QueryUser(ctx context.Context, number UserNumber) (*User, error)
	CheckPwd(ctx context.Context, origin, hash string) error
	GenerateToken(ctx context.Context, authInfo *AuthInfo) (string, error)
}
type CreateRoomUsecase interface {
	Logic(ctx context.Context, req *CreateRoomRequest) (*Message, error)
	EncryptSecret(ctx context.Context, secret string) (string, error)
	GenerateNumber(ctx context.Context) (RoomNumber, error)
	CreateRoom(ctx context.Context, room *Room) (bool, error)
}

type JoinRoomUsecase interface {
	Logic(ctx context.Context, req *JoinRoomRequest) (*Message, error)
	QueryRoom(ctx context.Context, number RoomNumber) (*Room, error)
	CheckSecret(ctx context.Context, origin, hash string) error
	JoinRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}

type LeaveRoomUsecase interface {
	Logic(ctx context.Context, req *LeaveRoomRequest) (*Message, error)
	LeaveRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}
