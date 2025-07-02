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

type ChatRepository interface {
	Save(ctx context.Context, chat *Chat) error
	FindAllByRoomNumber(ctx context.Context, number RoomNumber) ([]*Chat, error)
}
