package domain

type UserRepository interface {
	Create(user *User) (bool, error)
	QueryByNumber(number UserNumber) (*User, error)
}

type RoomRepository interface {
	Create(room *Room) (bool, error)
	QueryByRoomNumber(roomNumber RoomNumber) (*Room, error)
}

type UserRoomRepository interface {
	Join(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
	Leave(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}

type MessageRepository interface {
	Create(message Message) (bool, error)
	QueryAllByRoomNumber(number RoomNumber) ([]Message, error)
}
