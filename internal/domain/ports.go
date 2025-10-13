package domain

import (
	"context"
	"time"
)

type SendUnsentMessageAggregate interface {
	UnsentMessageFinder
	SentMessageUpdater
}

type SendMessageAggregate interface {
	MessageSaver
	SentMessageUpdater
}

type RefreshTokenAggregate interface {
	RefreshTokenFinder
	RefreshTokenSaver
}

type JoinRoomAggregate interface {
	RoomFinder
	RoomJoiner
}

type UserSaver interface {
	SaveUser(ctx context.Context, user *User) error
}

type UserFinder interface {
	FindUser(ctx context.Context, number UserNumber) (*User, error)
}

type RoomSaver interface {
	SaveRoom(ctx context.Context, room *Room) error
}

type RoomJoiner interface {
	JoinRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
}

type RoomFinder interface {
	FindRoom(ctx context.Context, number RoomNumber) (*Room, error)
}

type RoomLeaver interface {
	LeaveRoom(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
}

type RoomMemberFinder interface {
	FindRoomMembers(ctx context.Context, number RoomNumber) ([]UserNumber, error)
}

type UnsentMessageFinder interface {
	FindUnsentMessages(ctx context.Context, number UserNumber) ([]*Message, error)
}

type SentMessageUpdater interface {
	UpdateMessageSent(ctx context.Context, userNumber UserNumber, msgID MessageID) error
}

type MessageSaver interface {
	SaveMessage(ctx context.Context, msg *Message) error
}

type RefreshTokenSaver interface {
	SaveRefreshToken(ctx context.Context, token RefreshToken, info *RefreshInfo, expireDuration time.Duration) error
}
type RefreshTokenFinder interface {
	FindRefreshToken(ctx context.Context, token RefreshToken) (*RefreshInfo, error)
}

type Encryptor interface {
	Encrypt(raw string) (string, error)
}

type Comparator interface {
	Compare(origin, hash string) error
}

type AuthTokenGenerator interface {
	GenerateAuthToken(authInfo *AuthInfo) (AuthToken, error)
}
type AuthTokenParser interface {
	ParseAuthToken(token AuthToken) (*AuthInfo, error)
}

type RefreshTokenGenerator interface {
	GenerateRefreshToken() (RefreshToken, error)
}

type NumberGenerator interface {
	GenerateNumber() BaseNumber
}

type MessageSender interface {
	SendMessage(number UserNumber, msg *Message) error
}
