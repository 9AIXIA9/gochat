package application

import (
	"context"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"

	"github.com/go-faster/errors"
)

type SendFriendRequestUseCase kernel.UseCase[*SendFriendRequestInput, *kernel.NoOutput]

type SendFriendRequestInput struct {
	FromID  kernel.UserID
	ToID    kernel.UserID
	Content string
}

func (r *SendFriendRequestInput) Validate() error {
	if len(r.FromID) == 0 || len(r.ToID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user ids can't be empty")
	}

	return nil
}

type sendFriendRequestUseCase struct {
	friendshipExisterByUserID    domain.FriendshipExisterByUserID
	friendRequestExisterByUserID domain.FriendRequestExisterByUserIDAndState
	friendRequestCreator         domain.FriendRequestCreator
	eventIDGenerator             event.IDGenerator
	operationIDGenerator         kernel.OperationIDGenerator
}

func NewSendFriendRequestUseCase(
	friendshipExisterByUserID domain.FriendshipExisterByUserID,
	friendRequestExisterByUserID domain.FriendRequestExisterByUserIDAndState,
	friendRequestCreator domain.FriendRequestCreator,
	eventIDGenerator event.IDGenerator,
	operationIDGenerator kernel.OperationIDGenerator,
) (SendFriendRequestUseCase, error) {
	if err := utils.CheckInterfaces(
		friendshipExisterByUserID,
		friendRequestExisterByUserID,
		friendRequestCreator,
		eventIDGenerator,
		operationIDGenerator,
	); err != nil {
		return nil, err
	}

	return &sendFriendRequestUseCase{
		friendshipExisterByUserID:    friendshipExisterByUserID,
		friendRequestExisterByUserID: friendRequestExisterByUserID,
		friendRequestCreator:         friendRequestCreator,
		eventIDGenerator:             eventIDGenerator,
		operationIDGenerator:         operationIDGenerator,
	}, nil
}

func (uc *sendFriendRequestUseCase) Execute(ctx context.Context, input *SendFriendRequestInput) (*kernel.NoOutput, error) {
	exist, err := uc.friendshipExisterByUserID.ExistByUserID(ctx, input.FromID, input.ToID)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, domain.ErrAlreadyBeenFriends
	}

	exist, err = uc.friendRequestExisterByUserID.ExistByUserIDAndState(ctx, input.FromID, input.ToID, domain.StatePending)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, domain.ErrFriendRequestExists
	}

	req, err := domain.CreateFriendRequest(input.FromID, input.ToID, input.Content, uc.operationIDGenerator, uc.eventIDGenerator)
	if err != nil {
		return nil, err
	}

	if err := uc.friendRequestCreator.Create(ctx, req); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, domain.ErrFriendRequestExists
		}
		return nil, err
	}

	return nil, nil
}
