package usecase

import (
	"context"
	"gochat/internal/domain"
)

type CreateRoom struct {
	domain.NumberGenerator
	domain.Encryptor
	domain.RoomSaver
}

func NewCreateRoom(
	generator domain.NumberGenerator,
	encryptor domain.Encryptor,
	saver domain.RoomSaver,
) domain.CreateRoomUsecase {
	return &CreateRoom{
		NumberGenerator: generator,
		Encryptor:       encryptor,
		RoomSaver:       saver,
	}
}

func (uc *CreateRoom) Execute(ctx context.Context, req *domain.CreateRoomRequest) (*domain.Response, error) {
	secretHash, err := uc.Encrypt(req.Secret)
	if err != nil {
		return nil, err
	}
	roomNumber := uc.GenerateNumber()

	room := domain.CreateRoom(roomNumber, secretHash, req.MaxUsers, req.UserNumber)

	if err := uc.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return domain.NewSuccessResponse(domain.CreateRoomResponse{
		RoomNumber: room.Number(),
		MaxUsers:   req.MaxUsers,
		Owner:      req.UserNumber,
	}), nil
}
