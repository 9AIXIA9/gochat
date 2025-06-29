package usecase

import (
	"gochat/internal/domain"
)

type LeaveRoom struct {
	repo domain.UserRoomRepository
}

func NewLeaveRoom(repo domain.UserRoomRepository) domain.LeaveRoomUsecase {
	return &LeaveRoom{repo}
}

func (uc *LeaveRoom) Logic(req *domain.LeaveRoomRequest) (*domain.Response, error) {
	// 退出房间
	if exist, err := uc.LeaveRoom(req.UserNumber, req.URI.RoomNumber); err != nil {
		return nil, err
	} else if !exist {
		return domain.NewResponseWithDefaultMsg(domain.CodeNotJoined), nil
	}

	return nil, nil
}

func (uc *LeaveRoom) LeaveRoom(userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.repo.Leave(userNumber, roomNumber)
}
