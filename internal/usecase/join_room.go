package usecase

import (
	"context"
	"gochat/internal/domain"
	"gochat/internal/utils"
)

type JoinRoom struct {
	domain.Comparator
	domain.JoinRoomAggregateUOW
}

func NewJoinRoom(
	comparator domain.Comparator,
	joinRoomAggregateUOW domain.JoinRoomAggregateUOW,
) domain.JoinRoomUsecase {
	return &JoinRoom{
		Comparator:           comparator,
		JoinRoomAggregateUOW: joinRoomAggregateUOW,
	}
}

func (uc *JoinRoom) Execute(ctx context.Context, req *domain.JoinRoomRequest) (resp *domain.Response, err error) {
	//查询房间信息
	err = uc.DoJoinRoomUOW(ctx, func(aggregate domain.JoinRoomAggregate) error {
		room, err := aggregate.FindRoom(ctx, req.Number)
		if err != nil {
			if utils.IsNotFound(err) {
				resp = domain.RoomNotExistResponse
				return nil
			}
			return err
		}

		if room == nil {
			resp = domain.RoomNotExistResponse
			return nil
		}

		//判断人数
		if room.CurrentUsers() >= room.MaxUsers() {
			resp = domain.RoomIsFullResponse
			return nil
		}

		// 判断密钥
		if err := uc.Compare(req.Secret, room.SecretHash()); err != nil {
			resp = domain.WrongSecretResponse
			return nil
		}

		//加入房间(仓库)
		if err := aggregate.JoinRoom(ctx, req.UserNumber, req.Number); err != nil {
			if utils.IsDuplicate(err) {
				resp = domain.HasJoinedResponse
				return nil
			}
			return err
		}

		return nil
	})
	return resp, err
}
