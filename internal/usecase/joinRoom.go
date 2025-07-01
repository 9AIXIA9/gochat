package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
)

type JoinRoom struct {
	roomRepo     domain.RoomRepository
	userRoomRepo domain.UserRoomRepository
}

func NewJoinRoom(roomRepo domain.RoomRepository, userRoomRepo domain.UserRoomRepository) domain.JoinRoomUsecase {
	return &JoinRoom{roomRepo: roomRepo, userRoomRepo: userRoomRepo}
}

func (uc *JoinRoom) Logic(ctx context.Context, req *domain.JoinRoomRequest) (*domain.Message, error) {
	//查询房间信息
	room, err := uc.QueryRoom(ctx, req.Number)
	if err != nil {
		return nil, err
	}

	if room == nil {
		return domain.RoomNotExistResponse, nil
	}

	//判断人数
	if room.CurrentUsers >= room.MaxUsers {
		return domain.RoomIsFullResponse, nil
	}

	// 判断密钥
	if err := uc.CheckSecret(ctx, req.Secret, room.SecretHash); err != nil {
		return domain.WrongSecretResponse, nil
	}

	//加入房间(仓库)
	if exist, err := uc.JoinRoom(ctx, req.UserNumber, req.Number); err != nil {
		return nil, err
	} else if exist {
		return domain.HasJoinedResponse, nil
	}

	return domain.DefaultResponse, nil
}

func (uc *JoinRoom) CheckSecret(ctx context.Context, origin, hash string) error {
	return encrypt.Compare(ctx, origin, hash)
}

func (uc *JoinRoom) JoinRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.userRoomRepo.Join(ctx, userNumber, roomNumber)
}

func (uc *JoinRoom) QueryRoom(ctx context.Context, number domain.RoomNumber) (*domain.Room, error) {
	return uc.roomRepo.QueryByRoomNumber(ctx, number)
}
