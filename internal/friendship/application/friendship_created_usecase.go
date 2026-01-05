package application

import (
	"context"
	"fmt"
	chatDomain "gochat/internal/chat/domain"
	"gochat/internal/friendship/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type FriendshipCreatedUseCase kernel.UseCase[*FriendshipCreatedInput, *kernel.NoOutput]

type FriendshipCreatedInput struct {
	FriendshipID domain.FriendshipID
}

func (r *FriendshipCreatedInput) Validate() error {
	if len(r.FriendshipID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "friendship id can't be empty")
	}

	return nil
}

type friendshipCreatedUseCase struct {
	idGenerator event.IDGenerator
	finderByID  domain.FriendshipFinderByID
	creator     event.UnpublishedEventsCreator
}

func NewFriendshipCreatedUseCase(
	finderByID domain.FriendshipFinderByID,
	creator event.UnpublishedEventsCreator,
	idGenerator event.IDGenerator,
) (FriendshipCreatedUseCase, error) {
	if err := validate.NotNil(
		idGenerator, finderByID, creator,
	); err != nil {
		return nil, err
	}
	return &friendshipCreatedUseCase{
		idGenerator: idGenerator,
		finderByID:  finderByID,
		creator:     creator,
	}, nil
}

func (uc *friendshipCreatedUseCase) Execute(ctx context.Context, input *FriendshipCreatedInput) (*kernel.NoOutput, error) {
	friendship, err := uc.finderByID.FindByID(ctx, input.FriendshipID)
	if err != nil {
		return nil, err
	}

	notificationEvRequested1, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
		friendship.UserID1(),
		uc.buildContent(friendship.UserID2()),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	notificationEvRequested2, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
		friendship.UserID2(),
		uc.buildContent(friendship.UserID1()),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	chatEvCreated, err := chatDomain.NewFriendshipCreatedEvent(
		chatDomain.FriendshipID(friendship.ID()),
		friendship.UserID1(),
		friendship.UserID2(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, []event.Event{
		notificationEvRequested1,
		notificationEvRequested2,
		chatEvCreated,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *friendshipCreatedUseCase) buildContent(userID kernel.UserID) string {
	return fmt.Sprintf("You are now friends with user %s.", userID)
}
