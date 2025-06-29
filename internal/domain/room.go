package domain

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

type LeaveRoomUsecase interface {
	Logic(req *LeaveRoomRequest) (*Response, error)
	LeaveRoom(userNumber UserNumber, roomNumber RoomNumber) (bool, error)
}
