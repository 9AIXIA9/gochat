package domain

import (
	"context"
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
	Save(ctx context.Context, message *Message) error
	QueryMessages(ctx context.Context, number BaseNumber, count int) ([]*Message, error)
	QueryAllMessages(ctx context.Context, number BaseNumber) ([]*Message, error)
}
