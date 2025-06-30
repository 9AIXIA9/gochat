package domain

import "context"

//todo 加上retry

type UserRepository interface {
	Create(ctx context.Context, user *User) (bool, error)
	QueryByNumber(ctx context.Context, number UserNumber) (*User, error)
}

type RoomRepository interface {
	Create(ctx context.Context, room *Room) (bool, error)
	QueryByRoomNumber(ctx context.Context, roomNumber RoomNumber) (*Room, error)
}

type UserRoomRepository interface {
	Join(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) (bool, error)
	Leave(ctx context.Context, userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}

type ChatRepository interface {
	Create(ctx context.Context, chat Chat) (bool, error)
	QueryAllByRoomNumber(ctx context.Context, number RoomNumber) ([]Chat, error)
}
