package usecase

import (
	"go.uber.org/zap"
	"gochat/internal/domain"
)

type ExitRoom struct {
	repo domain.RoomRepository
}

func NewExitRoom(repo domain.RoomRepository) domain.ExitRoomUsecase {
	return &ExitRoom{repo: repo}
}

func (uc *ExitRoom) Logic(req *domain.ExitRoomRequest) (*domain.Response, error) {
	// 退出房间
	if exist, err := uc.ExitRoom(req.UserNumber, req.RoomNumber); err != nil {
		return nil, err
	} else if !exist {
		return domain.NewResponseWithoutMsg(domain.CodeNotJoined), nil
	}

	zap.L().Info("exit room successfully",
		zap.Int64("userNumber", int64(req.UserNumber)),
		zap.Int64("roomNumber", int64(req.RoomNumber)))
	return nil, nil

}

func (uc *ExitRoom) ExitRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.repo.Delete(userNumber, roomNumber)
}
