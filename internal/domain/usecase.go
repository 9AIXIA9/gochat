package domain

import (
	"context"
)

//TODO: 逻辑部分很多还是不太精准 比如自己发送房间消息时不需要通知自己 房间不能多次加入 不能加入不存在房间 退出不在的房间 给不存在的人发消息 给不是房间号码的发送房间信息 给房间发送私人信息等等

type AuthUsecase interface {
	AuthTokenParser
}

type UserConnectedUsecase interface {
	Execute(ctx context.Context, number UserNumber) error
	MessageSender
	SendUnsentMessageAggregate
}

type Usecase[request any] interface {
	Execute(ctx context.Context, request *request) (*Response, error)
}

type RefreshTokenUsecase interface {
	Execute(ctx context.Context, req *RefreshTokenRequest) (*Response, error)
	RefreshTokenGenerator
	AuthTokenGenerator
	RefreshTokenAggregate
}

type SignupUsecase interface {
	Execute(ctx context.Context, req *SignupRequest) (*Response, error)
	Encryptor
	NumberGenerator
	UserSaver
}

type LoginUsecase interface {
	Execute(ctx context.Context, req *LoginRequest) (*Response, error)
	Comparator
	AuthTokenGenerator
	RefreshTokenGenerator
	UserFinder
	RefreshTokenSaver
}

type CreateRoomUsecase interface {
	Execute(ctx context.Context, req *CreateRoomRequest) (*Response, error)
	Encryptor
	NumberGenerator
	RoomSaver
}

type JoinRoomUsecase interface {
	Execute(ctx context.Context, req *JoinRoomRequest) (*Response, error)
	Comparator
	JoinRoomAggregate
}

type LeaveRoomUsecase interface {
	Execute(ctx context.Context, req *LeaveRoomRequest) (*Response, error)
	RoomLeaver
}

type SendPrivateMessageUsecase interface {
	Execute(ctx context.Context, req *SendMessageRequest) (*Response, error)
	MessageSender
	SendMessageAggregate
}

type SendRoomMessageUsecase interface {
	Execute(ctx context.Context, req *SendMessageRequest) (*Response, error)
	MessageSender
	RoomMemberFinder
	SendMessageAggregate
}
