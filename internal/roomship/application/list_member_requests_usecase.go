package application

import (
	"context"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
)

type ListMemberRequestsUseCase kernel.UseCase[*ListMemberRequestsInput, *ListMemberRequestsOutput]

const (
	defaultMemberRequestsLimit = 10
	maxMemberRequestsLimit     = 100
)

type ListMemberRequestsInput struct {
	UserID kernel.UserID
	BaseID kernel.OperationID
	Limit  int
}

func (r *ListMemberRequestsInput) Validate() error {
	if len(r.UserID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "user id is empty")
	}

	return nil
}

type ListMemberRequestsOutput struct {
	Requests []*domain.MemberRequest
}

type listMemberRequestsUseCase struct {
	finder domain.MemberRequestsFinderByUserID
}

func NewListMemberRequestsUseCase(
	finder domain.MemberRequestsFinderByUserID,
) (ListMemberRequestsUseCase, error) {
	if err := validate.NotNil(finder); err != nil {
		return nil, err
	}
	return &listMemberRequestsUseCase{
		finder: finder,
	}, nil
}

func (uc *listMemberRequestsUseCase) Execute(ctx context.Context, input *ListMemberRequestsInput) (*ListMemberRequestsOutput, error) {
	limit := input.Limit
	switch {
	case limit <= 0:
		limit = defaultMemberRequestsLimit
	case limit > maxMemberRequestsLimit:
		limit = maxMemberRequestsLimit
	}

	requests, err := uc.finder.FindsByUserID(ctx, input.UserID, limit, input.BaseID)
	if err != nil {
		return nil, err
	}

	return &ListMemberRequestsOutput{
		Requests: requests,
	}, nil
}
