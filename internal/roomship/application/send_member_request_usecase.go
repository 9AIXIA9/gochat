package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type SendMemberRequestUseCase kernel.UseCase[*SendMemberRequestInput, *kernel.NoOutput]

type SendMemberRequestInput struct {
	UserID   kernel.UserID
	RoomID   kernel.RoomID
	Content  string
	Password domain.Password
}

func (r *SendMemberRequestInput) Validate() error {
	if len(r.UserID) == 0 || len(r.RoomID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id or room id is empty")
	}
	return nil
}

type sendMemberRequestUseCase struct {
	memberRequestExister domain.MemberRequestExisterByUserIDAndRoomIDAndState
	roomshipExister      domain.RoomshipExisterByUserIDAndRoomID
	finder               domain.RoomFinderByID
	eventIDGenerator     event.IDGenerator
	operationIDGenerator kernel.OperationIDGenerator
	comparator           domain.Comparator
	creator              domain.MemberRequestCreator
}

func NewSendMemberRequestUseCase(
	memberRequestExister domain.MemberRequestExisterByUserIDAndRoomIDAndState,
	roomshipExister domain.RoomshipExisterByUserIDAndRoomID,
	finder domain.RoomFinderByID,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator kernel.OperationIDGenerator,
	comparator domain.Comparator,
	creator domain.MemberRequestCreator,
) (SendMemberRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		memberRequestExister,
		roomshipExister,
		finder,
		eventIDGenerator,
		operationIDGenerator,
		comparator,
		creator,
	); err != nil {
		return nil, err
	}
	return &sendMemberRequestUseCase{
		memberRequestExister: memberRequestExister,
		roomshipExister:      roomshipExister,
		finder:               finder,
		eventIDGenerator:     eventIDGenerator,
		operationIDGenerator: operationIDGenerator,
		comparator:           comparator,
		creator:              creator,
	}, nil
}

func (uc *sendMemberRequestUseCase) Execute(ctx context.Context, input *SendMemberRequestInput) (*kernel.NoOutput, error) {
	exist, err := uc.memberRequestExister.ExistByUserIDAndRoomIDAndState(ctx, input.UserID, input.RoomID, domain.StatePending)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, domain.ErrMemberRequestAlreadyExists
	}

	exist, err = uc.roomshipExister.ExistByUserIDAndRoomID(ctx, input.UserID, input.RoomID)
	if err != nil {
		return nil, domain.ErrIsAlreadyMember
	}

	if exist {
		return nil, domain.ErrMemberRequestAlreadyExists
	}

	room, err := uc.finder.FindByID(ctx, input.RoomID)
	if err != nil {
		return nil, err
	}

	if err := room.PasswordEncrypted().Compare(input.Password, uc.comparator); err != nil {
		return nil, err
	}

	req, err := domain.CreateMemberRequest(
		input.UserID,
		input.RoomID,
		input.Content,
		uc.operationIDGenerator,
		uc.eventIDGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.Create(ctx, req); err != nil {
		return nil, err
	}

	return nil, nil
}
