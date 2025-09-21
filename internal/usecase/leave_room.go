package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

type LeaveRoom struct {
	repo domain.UserRoomRepository
}

func NewLeaveRoom(repo domain.UserRoomRepository) domain.LeaveRoomUsecase {
	return &LeaveRoom{repo}
}

func (uc *LeaveRoom) Execute(ctx context.Context, req *domain.LeaveRoomRequest) (*domain.Response, error) {
	// 退出房间
	if err := uc.LeaveRoom(ctx, req.UserNumber, req.Number); err != nil {
		if utils.IsNotFound(err) {
			return domain.NotJoinedResponse, nil
		}
		return nil, err
	}

	return nil, nil
}

func (uc *LeaveRoom) LeaveRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) error {
	return uc.repo.Delete(ctx, userNumber, roomNumber)
}
