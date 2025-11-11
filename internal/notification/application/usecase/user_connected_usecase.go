package usecase

import (
	"context"
	"errors"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type UserConnectedUseCase kernel.UseCase[*UserConnectedInput, *kernel.NoOutput]

type UserConnectedInput struct {
	UserID kernel.UserID
}

func (r *UserConnectedInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type userConnectedUseCase struct {
	messageFinder        application.MessageFinder
	messageNotifier      application.MessageNotifier
	messageStatesUpdater application.MessageStatesUpdater
}

func NewUserConnectedUseCase(
	messageFinder application.MessageFinder,
	messageNotifier application.MessageNotifier,
	messageStatesUpdater application.MessageStatesUpdater,
) UserConnectedUseCase {
	return &userConnectedUseCase{
		messageFinder:        messageFinder,
		messageNotifier:      messageNotifier,
		messageStatesUpdater: messageStatesUpdater,
	}
}

func (uc *userConnectedUseCase) Execute(ctx context.Context, input *UserConnectedInput) (*kernel.NoOutput, error) {
	messages, err := uc.messageFinder.FindMessagesByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	ids := make([]domain.MessageID, 0, len(messages))
	for _, message := range messages {
		if err := uc.messageNotifier.Notify(input.UserID, message); err != nil {
			if errors.Is(err, myErrors.ErrNotFound) {
				return nil, nil
			}
			return nil, err
		}
		ids = append(ids, message.ID())
	}
	if err := uc.messageStatesUpdater.UpdateMessageStates(ctx, input.UserID, ids, domain.MessageStateDelivered); err != nil {
		return nil, err
	}

	return nil, nil
}
