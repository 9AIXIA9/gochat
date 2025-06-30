package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/types"
)

type LeaveRoom struct {
	repo domain.UserRoomRepository
}

func NewLeaveRoom(repo domain.UserRoomRepository) domain.LeaveRoomUsecase {
	return &LeaveRoom{repo}
}

func (uc *LeaveRoom) Logic(ctx context.Context, req *domain.LeaveRoomRequest) (*domain.Message, error) {
	// 退出房间
	if exist, err := uc.LeaveRoom(ctx, req.UserNumber, req.Number); err != nil {
		return nil, err
	} else if !exist {
		return types.NotJoinedResponse, nil
	}

	return types.DefaultResponse, nil
}

func (uc *LeaveRoom) LeaveRoom(ctx context.Context, userNumber domain.UserNumber, roomNumber domain.RoomNumber) (bool, error) {
	return uc.repo.Leave(ctx, userNumber, roomNumber)
}
