package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
	"gochat/internal/utils"
)

type JoinRoom struct {
	roomRepo     domain.RoomRepository
	userRoomRepo domain.UserRoomRepository
}

func NewJoinRoom(roomRepo domain.RoomRepository, userRoomRepo domain.UserRoomRepository) domain.JoinRoomUsecase {
	return &JoinRoom{roomRepo: roomRepo, userRoomRepo: userRoomRepo}
}

func (uc *JoinRoom) Execute(ctx context.Context, req *domain.JoinRoomRequest) (*domain.Response, error) {
	//查询房间信息
	room, err := uc.FindRoom(ctx, req.Number)
	if err != nil {
		if utils.IsNotFound(err) {
			return domain.RoomNotExistResponse, nil
		}
		return nil, err
	}

	if room == nil {
		return domain.RoomNotExistResponse, nil
	}

	//判断人数
	if room.CurrentUsers() >= room.MaxUsers() {
		return domain.RoomIsFullResponse, nil
	}

	// 判断密钥
	if err := uc.CheckSecret(ctx, req.Secret, room.SecretHash()); err != nil {
		return domain.WrongSecretResponse, nil
	}

	//加入房间(仓库)
	if err := uc.JoinRoom(ctx, req.UserNumber, req.Number); err != nil {
		if utils.IsDuplicate(err) {
			return domain.HasJoinedResponse, nil
		}
		return nil, err
	}

	return nil, nil
}

func (uc *JoinRoom) CheckSecret(ctx context.Context, origin, hash string) error {
	return encrypt.Compare(ctx, origin, hash)
}

func (uc *JoinRoom) JoinRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	return uc.userRoomRepo.Save(ctx, userNumber, roomNumber)
}

func (uc *JoinRoom) FindRoom(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	return uc.roomRepo.FindOneByNumber(ctx, number)
}
