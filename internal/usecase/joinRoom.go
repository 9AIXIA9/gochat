package usecase

import (
	"gochat/internal/domain"
	"gochat/internal/infra/ctrl"

	"net/http"
)

type JoinRoom struct {
}

func NewJoinRoom() domain.JoinRoomUsecase {
	return &JoinRoom{}
}

func (uc *JoinRoom) CheckSecret(number domain.RoomNumber, secret string) error {
	return nil
}

func (uc *JoinRoom) JoinRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber, w http.ResponseWriter, r *http.Request) error {
	return ctrl.JoinRoom(userNumber, roomNumber, w, r)
}
