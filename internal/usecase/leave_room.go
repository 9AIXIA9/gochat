package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

type LeaveRoom struct {
	domain.RoomLeaver
}

func NewLeaveRoom(leaver domain.RoomLeaver) domain.LeaveRoomUsecase {
	return &LeaveRoom{RoomLeaver: leaver}
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
