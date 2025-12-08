package application

import (
	"context"
	"gochat/internal/chat/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

const (
	defaultPrivateMessagesListLimit = 50
	maxPrivateMessagesListLimit     = 100
)

type ListPrivateMessagesUseCase kernel.UseCase[*ListPrivateMessagesInput, *ListPrivateMessagesOutput]

type ListPrivateMessagesInput struct {
	UserID kernel.UserID
	BaseID kernel.MessageID
	Limit  int
}

func (r *ListPrivateMessagesInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type ListPrivateMessagesOutput struct {
	PrivateMessages []*domain.PrivateMessage
}

type listPrivateMessagesUseCase struct {
	finder domain.PrivateMessagesFinderByRecipientID
}

func NewListPrivateMessagesUseCase(
	finder domain.PrivateMessagesFinderByRecipientID,
) (ListPrivateMessagesUseCase, error) {
	if err := utils.CheckInterfaces(finder); err != nil {
		return nil, err
	}
	return &listPrivateMessagesUseCase{
		finder: finder,
	}, nil
}

func (uc *listPrivateMessagesUseCase) Execute(ctx context.Context, input *ListPrivateMessagesInput) (*ListPrivateMessagesOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultPrivateMessagesListLimit
	case limit > maxPrivateMessagesListLimit:
		limit = maxPrivateMessagesListLimit
	}

	messages, err := uc.finder.FindsByRecipientID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListPrivateMessagesOutput{
		PrivateMessages: messages,
	}, nil
}
