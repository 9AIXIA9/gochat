package application

import (
	"context"
	"gochat/internal/friendship/domain"
	notificationDomain "gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type FriendRequestCreatedUseCase kernel.UseCase[*FriendRequestCreatedInput, *kernel.NoOutput]

type FriendRequestCreatedInput struct {
	RequestID kernel.OperationID
}

func (r *FriendRequestCreatedInput) Validate() error {
	if len(r.RequestID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type friendRequestCreatedUseCase struct {
	requestIDFinder domain.FriendRequestFinderByID
	creator         event.UnpublishedEventsCreator
	idGenerator     event.IDGenerator
}

func NewFriendRequestCreatedUseCase(
	requestIDFinder domain.FriendRequestFinderByID,
	creator event.UnpublishedEventsCreator,
	idGenerator event.IDGenerator,
) (FriendRequestCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		idGenerator, requestIDFinder, creator,
	); err != nil {
		return nil, err
	}
	return &friendRequestCreatedUseCase{
		requestIDFinder: requestIDFinder,
		creator:         creator,
		idGenerator:     idGenerator,
	}, nil
}

func (uc *friendRequestCreatedUseCase) Execute(ctx context.Context, input *FriendRequestCreatedInput) (*kernel.NoOutput, error) {
	req, err := uc.requestIDFinder.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	if req.State() != domain.StatePending {
		//已处理则不发送通知
		return nil, nil
	}

	if req.From() == req.To() {
		return nil, domain.ErrAddYourselfAsFriend
	}

	ev, err := notificationDomain.NewFriendRequestCreatedNotificationRequestedEvent(
		input.RequestID,
		req.From(),
		req.To(),
		req.SentAt(),
		req.Content(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, []event.Event{ev}); err != nil {
		return nil, err
	}

	return nil, nil
}
