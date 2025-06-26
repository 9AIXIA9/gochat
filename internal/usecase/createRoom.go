package usecase

import (
	"gochat/internal/domain"
)

type CreateRoom struct{}

func NewCreateRoom() domain.CreateRoomUsecase {
	return &CreateRoom{}
}

func (uc *CreateRoom) EncryptSecret(secret string) string {
	//TODO implement me
	panic("implement me")
}

func (uc *CreateRoom) GenerateNumber() domain.RoomNumber {
	//TODO implement me
	panic("implement me")
}

func (uc *CreateRoom) CreateRoom(owner domain.UserNumber, name, description string, maxUsers int) (*domain.Room, error) {
	//TODO implement me
	panic("implement me")
}
