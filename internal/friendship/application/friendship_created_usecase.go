package application

import (
	"context"
	"fmt"
	"gochat/internal/friendship/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type FriendshipCreatedUseCase kernel.UseCase[*FriendshipCreatedInput, *kernel.NoOutput]

type FriendshipCreatedInput struct {
	FriendshipID domain.FriendshipID
}

func (r *FriendshipCreatedInput) Validate() error {
	if len(r.FriendshipID) == 0 {
		return myErrors.ErrEmptyInput
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
	if err := utils.CheckInterfaces(
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
	//TODO 同步到 chat上下文

	friendship, err := uc.finderByID.FindByID(ctx, input.FriendshipID)
	if err != nil {
		return nil, err
	}

	ev1, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
		friendship.UserID1(),
		uc.buildContent(friendship.UserID2()),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	ev2, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
		friendship.UserID2(),
		uc.buildContent(friendship.UserID1()),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, []event.Event{ev1, ev2}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *friendshipCreatedUseCase) buildContent(userID kernel.UserID) string {
	return fmt.Sprintf("You are now friends with user %s.", userID)
}
