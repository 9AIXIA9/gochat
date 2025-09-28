package domain

import (
	"context"
)

type Usecase[request any] interface {
	Execute(ctx context.Context, request *request) (*Response, error)
}

type SignupUsecase interface {
	Execute(ctx context.Context, req *SignupRequest) (*Response, error)
	EncryptPwd(ctx context.Context, pwd string) (string, error)
	GenerateNumber(ctx context.Context) (UserNumber, error)
	CreateUser(ctx context.Context, user *User) error
}

type LoginUsecase interface {
	Execute(ctx context.Context, req *LoginRequest) (*Response, error)
	FindUser(ctx context.Context, number UserNumber) (*User, error)
	CheckPwd(ctx context.Context, origin, hash string) error
	GenerateRefreshToken() (RefreshToken, error)
	GenerateAuthToken(ctx context.Context, authInfo *AuthInfo) (AuthToken, error)
	SaveRefreshToken(ctx context.Context, token RefreshToken, info *RefreshInfo) error
}
type CreateRoomUsecase interface {
	Execute(ctx context.Context, req *CreateRoomRequest) (*Response, error)
	EncryptSecret(ctx context.Context, secret string) (string, error)
	GenerateNumber(ctx context.Context) (RoomNumber, error)
	CreateRoom(ctx context.Context, room *Room) error
}

type JoinRoomUsecase interface {
	Execute(ctx context.Context, req *JoinRoomRequest) (*Response, error)
	FindRoom(ctx context.Context, number RoomNumber) (*Room, error)
	CheckSecret(ctx context.Context, origin, hash string) error
	JoinRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
}

type LeaveRoomUsecase interface {
	Execute(ctx context.Context, req *LeaveRoomRequest) (*Response, error)
	LeaveRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
}

type UserConnectedUsecase interface {
	Execute(ctx context.Context, number UserNumber) error
	QueryUnsentMessages(ctx context.Context, number UserNumber) ([]*Message, error)
	SendUserManyMsgs(number UserNumber, msgs []*Message) []MessageID
	UpdateMessagesSentToOneUser(ctx context.Context, number UserNumber, msgIDs []MessageID) error
}

type SendMessageUsecase interface {
	Execute(ctx context.Context, req *SendMessageRequest) (*Response, error)
	SaveAndQueryUserNumberShouldSent(ctx context.Context, msg *Message) ([]UserNumber, error)
	SendMsgToManyUsers(msg *Message, numbers []UserNumber) []UserNumber
	UpdateMessageSentToManyUsers(ctx context.Context, msgID MessageID, userNumbers []UserNumber) error
}
