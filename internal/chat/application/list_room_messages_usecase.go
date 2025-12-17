package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const (
	defaultRoomMessagesListLimit = 50
	maxRoomMessagesListLimit     = 100
)

type ListRoomMessagesUseCase kernel.UseCase[*ListRoomMessagesInput, *ListRoomMessagesOutput]

type ListRoomMessagesInput struct {
	UserID kernel.UserID
	BaseID kernel.MessageID
	Limit  int
}

func (r *ListRoomMessagesInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id can't be empty")
	}

	return nil
}

type ListRoomMessagesOutput struct {
	RoomMessages []*domain.RoomMessage
}

type listRoomMessagesUseCase struct {
	finder domain.RoomMessagesFinderByRecipientID
}

func NewListRoomMessagesUseCase(
	finder domain.RoomMessagesFinderByRecipientID,
) (ListRoomMessagesUseCase, error) {
	if err := utils.CheckInterfaces(finder); err != nil {
		return nil, err
	}
	return &listRoomMessagesUseCase{
		finder: finder,
	}, nil
}

func (uc *listRoomMessagesUseCase) Execute(ctx context.Context, input *ListRoomMessagesInput) (*ListRoomMessagesOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultRoomMessagesListLimit
	case limit > maxRoomMessagesListLimit:
		limit = maxRoomMessagesListLimit
	}

	messages, err := uc.finder.FindsByRecipientID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListRoomMessagesOutput{
		RoomMessages: messages,
	}, nil
}
