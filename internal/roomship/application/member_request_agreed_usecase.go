package application

import (
	"context"
	"errors"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type MemberRequestAgreedUseCase kernel.UseCase[*MemberRequestAgreedInput, *kernel.NoOutput]

type MemberRequestAgreedInput struct {
	RequestID kernel.OperationID
}

func (r *MemberRequestAgreedInput) Validate() error {
	if len(r.RequestID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "request id is empty")
	}

	return nil
}

type memberRequestAgreedUseCase struct {
	roomshipIDGenerator domain.RoomshipIDGenerator
	idGenerator         event.IDGenerator
	creator             domain.RoomshipCreator
	finder              domain.MemberRequestFinderByID
}

func NewMemberRequestAgreedUseCase(
	roomshipIDGenerator domain.RoomshipIDGenerator,
	idGenerator event.IDGenerator,
	finder domain.MemberRequestFinderByID,
	creator domain.RoomshipCreator,
) (MemberRequestAgreedUseCase, error) {
	if err := utils.CheckInterfaces(
		roomshipIDGenerator,
		idGenerator,
		creator,
		finder,
	); err != nil {
		return nil, err
	}
	return &memberRequestAgreedUseCase{
		roomshipIDGenerator: roomshipIDGenerator,
		idGenerator:         idGenerator,
		creator:             creator,
		finder:              finder,
	}, nil
}

func (uc *memberRequestAgreedUseCase) Execute(ctx context.Context, input *MemberRequestAgreedInput) (*kernel.NoOutput, error) {
	req, err := uc.finder.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	roomship, err := domain.CreateRoomship(
		req.ApplicantID(),
		req.RoomID(),
		domain.MemberRole,
		uc.roomshipIDGenerator,
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.Create(ctx, roomship); err != nil {
		if errors.Is(err, myErrors.ErrDuplicatedKey) {
			return nil, nil
		}
		return nil, err
	}

	return nil, nil
}
