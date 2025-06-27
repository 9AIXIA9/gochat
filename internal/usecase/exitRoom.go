package usecase

import "gochat/internal/domain"

type ExitRoom struct{}

func NewExitRoom() domain.ExitRoomUsecase {
	return &ExitRoom{}
}

func (uc *ExitRoom) ExitRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	//TODO implement me
	panic("implement me")
}
