package domain

type RoomNumber int64

type Room struct {
	Name         string
	Number       RoomNumber
	SecretHash   string
	Description  string
	CurrentUsers int
	MaxUsers     int
	Owner        UserNumber
}

type RoomRepository interface {
	Create(room *Room) (bool, error)
	JoinOne(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
	QueryByRoomNumber(roomNumber RoomNumber) (*Room, error)
	Delete(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
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

type ExitRoomUsecase interface {
	Logic(req *ExitRoomRequest) (*Response, error)
	ExitRoom(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}
