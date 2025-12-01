package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type SendFriendRequestUseCase kernel.UseCase[*SendFriendRequestInput, *kernel.NoOutput]

type SendFriendRequestInput struct {
	FromID   kernel.UserID
	ToNumber kernel.UserNumber
	Content  string
}

func (r *SendFriendRequestInput) Validate() error {
	if len(r.FromID) == 0 {
		return myErrors.ErrEmptyInput
	}

	if err := r.ToNumber.Validate(); err != nil {
		return err
	}

	return nil
}

type sendFriendRequestUseCase struct {
	finderByID           domain.UserFinderByID
	finderByNumber       domain.UserFinderByNumber
	saver                domain.UserSaver
	eventIDGenerator     event.IDGenerator
	operationIDGenerator domain.OperationIDGenerator
}

func NewSendFriendRequestUseCase(
	finderByID domain.UserFinderByID,
	finderByNumber domain.UserFinderByNumber,
	saver domain.UserSaver,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator domain.OperationIDGenerator,
) (SendFriendRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		finderByID,
		finderByNumber,
		saver,
		eventIDGenerator,
		operationIDGenerator,
	); err != nil {
		return nil, err
	}

	return &sendFriendRequestUseCase{
		finderByID:           finderByID,
		finderByNumber:       finderByNumber,
		saver:                saver,
		eventIDGenerator:     eventIDGenerator,
		operationIDGenerator: operationIDGenerator,
	}, nil
}

func (uc *sendFriendRequestUseCase) Execute(ctx context.Context, input *SendFriendRequestInput) (*kernel.NoOutput, error) {
	//TODO 事务处理
	fromUser, err := uc.finderByID.FindByID(ctx, input.FromID)
	if err != nil {
		return nil, err
	}

	toUser, err := uc.finderByNumber.FindByNumber(ctx, input.ToNumber)
	if err != nil {
		return nil, err
	}

	req, err := fromUser.SendFriendRequest(toUser.ID(), input.Content, uc.operationIDGenerator)
	if err != nil {
		return nil, err
	}

	if err := toUser.ReceiveFriendRequest(req, uc.eventIDGenerator); err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, fromUser); err != nil {
		return nil, err
	}

	if err := uc.saver.Save(ctx, toUser); err != nil {
		return nil, err
	}

	return nil, nil
}
