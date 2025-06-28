package usecase

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
	"gochat/internal/infra/encrypt"
)

type JoinRoom struct {
	repo domain.RoomRepository
}

func NewJoinRoom(repo domain.RoomRepository) domain.JoinRoomUsecase {
	return &JoinRoom{repo: repo}
}

func (uc *JoinRoom) Logic(req *domain.JoinRoomRequest) (*domain.Response, error) {
	//查询房间信息
	room, err := uc.QueryRoom(req.RoomNumber)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return domain.NewResponseWithoutMsg(domain.CodeRoomNotExist), nil
	}

	//判断人数
	if room.CurrentUsers >= room.MaxUsers {
		return domain.NewResponseWithoutMsg(domain.CodeRoomIsFull), nil
	}

	// 判断密钥
	if err := uc.CheckSecret(req.Secret, room.SecretHash); err != nil {
		return domain.NewResponseWithoutMsg(domain.CodeWrongSecret), nil
	}

	//加入房间(仓库)
	if exist, err := uc.JoinRoom(req.UserNumber, req.RoomNumber); err != nil {
		return nil, err
	} else if exist {
		return domain.NewResponseWithoutMsg(domain.CodeHasJoined), nil
	}

	zap.L().Info("join room successfully",
		zap.Int64("userNumber", int64(req.UserNumber)),
		zap.Int64("roomNumber", int64(req.RoomNumber)))

	return nil, nil
}

func (uc *JoinRoom) CheckSecret(origin, hash string) error {
	return encrypt.Compare(origin, hash)
}

func (uc *JoinRoom) JoinRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.repo.JoinOne(userNumber, roomNumber)
}

func (uc *JoinRoom) QueryRoom(number domain.RoomNumber) (*domain.Room, error) {
	return uc.repo.QueryByRoomNumber(number)
}
