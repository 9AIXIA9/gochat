package usecase

import "gochat/internal/domain"

type ExitRoom struct{}

func NewExitRoom() domain.ExitRoomUsecase {
	return &ExitRoom{}
}

func (uc *ExitRoom) ExitRoom(username, roomName string) error {
	//TODO implement me
	panic("implement me")
}
