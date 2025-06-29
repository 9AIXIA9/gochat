package usecase

import (
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

func (uc *JoinRoom) Logic(req *domain.JoinRoomRequest) (*domain.Response, error) {
	//查询房间信息
	room, err := uc.QueryRoom(req.URI.Number)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return domain.NewResponseWithDefaultMsg(domain.CodeRoomNotExist), nil
	}

	//判断人数
	if room.CurrentUsers >= room.MaxUsers {
		return domain.NewResponseWithDefaultMsg(domain.CodeRoomIsFull), nil
	}

	// 判断密钥
	if err := uc.CheckSecret(req.Body.Secret, room.SecretHash); err != nil {
		return domain.NewResponseWithDefaultMsg(domain.CodeWrongSecret), nil
	}

	//加入房间(仓库)
	if exist, err := uc.JoinRoom(req.UserNumber, req.URI.Number); err != nil {
		return nil, err
	} else if exist {
		return domain.NewResponseWithDefaultMsg(domain.CodeHasJoined), nil
	}

	return nil, nil
}

func (uc *JoinRoom) CheckSecret(origin, hash string) error {
	return encrypt.Compare(origin, hash)
}

func (uc *JoinRoom) JoinRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.userRoomRepo.Join(userNumber, roomNumber)
}

func (uc *JoinRoom) QueryRoom(number domain.RoomNumber) (*domain.Room, error) {
	return uc.roomRepo.QueryByRoomNumber(number)
}
