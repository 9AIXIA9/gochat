package application

import (
	"context"
	"errors"
	"gochat/internal/friendship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type FriendRequestAgreedUseCase kernel.UseCase[*FriendRequestAgreedInput, *kernel.NoOutput]

type FriendRequestAgreedInput struct {
	RequestID kernel.OperationID
}

func (r *FriendRequestAgreedInput) Validate() error {
	if len(r.RequestID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "request id can't be empty")
	}

	return nil
}

type friendRequestAgreedUseCase struct {
	friendRequestFinderByID domain.FriendRequestFinderByID
	creator                 domain.FriendshipCreator
	friendshipIDGenerator   domain.FriendshipIDGenerator
	idGenerator             event.IDGenerator
}

func NewFriendRequestAgreedUseCase(
	friendRequestFinderByID domain.FriendRequestFinderByID,
	creator domain.FriendshipCreator,
	friendshipIDGenerator domain.FriendshipIDGenerator,
	idGenerator event.IDGenerator,
) (FriendRequestAgreedUseCase, error) {
	if err := validate.NotNil(
		friendRequestFinderByID, creator, idGenerator, friendshipIDGenerator,
	); err != nil {
		return nil, err
	}
	return &friendRequestAgreedUseCase{
		friendRequestFinderByID: friendRequestFinderByID,
		creator:                 creator,
		friendshipIDGenerator:   friendshipIDGenerator,
		idGenerator:             idGenerator,
	}, nil
}

func (uc *friendRequestAgreedUseCase) Execute(ctx context.Context, input *FriendRequestAgreedInput) (*kernel.NoOutput, error) {
	req, err := uc.friendRequestFinderByID.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	if req.State() != domain.StateAgreed {
		return nil, domain.ErrFriendRequestNotAgreed
	}

	friendship, err := domain.CreateFriendship(
		req.From(),
		req.To(),
		uc.friendshipIDGenerator,
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.Create(ctx, friendship); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, nil
		}
		return nil, err
	}

	return nil, nil
}
