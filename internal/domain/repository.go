package domain

import (
	"context"
	"time"
)

//todo 加上retry

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindOneByNumber(ctx context.Context, number UserNumber) (*User, error)
}

type RoomRepository interface {
	Save(ctx context.Context, room *Room) error
	FindOneByNumber(ctx context.Context, number RoomNumber) (*Room, error)
}

type UserRoomRepository interface {
	Save(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
	Delete(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) error
}

type MessageRepository interface {
	SaveAndQueryUserNumberShouldSent(ctx context.Context, message *Message) ([]UserNumber, error)
	UpdateMessagesSentToOneUser(ctx context.Context, number UserNumber, msgIDs []MessageID) error
	UpdateMessageSentToManyUsers(ctx context.Context, msgID MessageID, userNumbers []UserNumber) error
	QueryUnsentMessages(ctx context.Context, number UserNumber) ([]*Message, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token RefreshToken, info *RefreshInfo, expireDuration time.Duration) error
	FindByToken(ctx context.Context, token RefreshToken) (*RefreshInfo, error)
	FindByUserNumber(ctx context.Context, userNumber UserNumber) (*RefreshInfo, error)
	DeleteByToken(ctx context.Context, token RefreshToken) error
	DeleteByUserNumber(ctx context.Context, userNumber UserNumber) error
}
