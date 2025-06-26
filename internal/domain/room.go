package domain

import (
	"net/http"
)

type RoomNumber int64

type Room struct {
	Name        string
	Number      RoomNumber
	SecretHash  string
	Description string
	MaxUsers    int
	Owner       UserNumber
	CreatedAt   int64
}

type CreateRoomUsecase interface {
	EncryptSecret(secret string) string
	GenerateNumber() RoomNumber
	CreateRoom(owner UserNumber, name, description string, maxUsers int) (*Room, error)
}

type JoinRoomUsecase interface {
	CheckSecret(number RoomNumber, secret string) error
	JoinRoom(userNumber UserNumber, roomNumber RoomNumber, w http.ResponseWriter, r *http.Request) error
}

type ExitRoomUsecase interface {
	ExitRoom(username, roomName string) error
}
