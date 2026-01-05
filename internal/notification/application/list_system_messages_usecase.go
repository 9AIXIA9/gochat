package application

import (
	"context"
	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

const (
	defaultSystemMessagesListLimit = 50
	maxSystemMessagesListLimit     = 100
)

type ListSystemMessagesUseCase kernel.UseCase[*ListSystemMessagesInput, *ListSystemMessagesOutput]

type ListSystemMessagesInput struct {
	UserID kernel.UserID
	BaseID kernel.MessageID
	Limit  int
}

func (r *ListSystemMessagesInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}

	return nil
}

type ListSystemMessagesOutput struct {
	SystemMessages []*domain.SystemMessage
}

type listSystemMessagesUseCase struct {
	finder domain.SystemMessageFinderByUserID
}

func NewListSystemMessagesUseCase(
	finder domain.SystemMessageFinderByUserID,
) (ListSystemMessagesUseCase, error) {
	if err := validate.NotNil(finder); err != nil {
		return nil, err
	}
	return &listSystemMessagesUseCase{
		finder: finder,
	}, nil
}

func (uc *listSystemMessagesUseCase) Execute(ctx context.Context, input *ListSystemMessagesInput) (*ListSystemMessagesOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultSystemMessagesListLimit
	case limit > maxSystemMessagesListLimit:
		limit = maxSystemMessagesListLimit
	}

	messages, err := uc.finder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListSystemMessagesOutput{
		SystemMessages: messages,
	}, nil
}
