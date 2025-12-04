package application

import (
	"context"
	notificationDomain "gochat/internal/notification/domain"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type MemberRequestCreatedUseCase kernel.UseCase[*MemberRequestCreatedInput, *kernel.NoOutput]

type MemberRequestCreatedInput struct {
	RequestID kernel.OperationID
}

func (r *MemberRequestCreatedInput) Validate() error {
	if len(r.RequestID) == 0 {
		return myErrors.ErrEmptyInput
	}

	return nil
}

type memberRequestCreatedUseCase struct {
	requestFinder  domain.MemberRequestFinderByID
	roomshipFinder domain.RoomshipsFinderByRoomIDAndRole
	idGenerator    event.IDGenerator
	creator        event.UnpublishedEventsCreator
}

func NewMemberRequestCreatedUseCase(
	requestFinder domain.MemberRequestFinderByID,
	roomshipFinder domain.RoomshipsFinderByRoomIDAndRole,
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventsCreator,
) (MemberRequestCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		requestFinder,
		roomshipFinder,
		idGenerator,
		creator,
	); err != nil {
		return nil, err
	}

	return &memberRequestCreatedUseCase{
		requestFinder:  requestFinder,
		roomshipFinder: roomshipFinder,
		idGenerator:    idGenerator,
		creator:        creator,
	}, nil
}

func (uc *memberRequestCreatedUseCase) Execute(ctx context.Context, input *MemberRequestCreatedInput) (*kernel.NoOutput, error) {
	req, err := uc.requestFinder.FindByID(ctx, input.RequestID)
	if err != nil {
		return nil, err
	}

	ownerRoomships, err := uc.roomshipFinder.FindsByRoomIDAndRole(ctx, req.RoomID(), domain.OwnerRole)
	if err != nil {
		return nil, err
	}

	evs := make([]event.Event, 0, len(ownerRoomships))
	for _, roomship := range ownerRoomships {
		ev, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
			roomship.UserID(),
			uc.buildContent(req.ApplicantID(), req.Content(), req.RoomID()),
			uc.idGenerator,
		)
		if err != nil {
			return nil, err
		}
		evs = append(evs, ev)
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, evs); err != nil {
		return nil, err
	}

	return nil, nil
}

func (uc *memberRequestCreatedUseCase) buildContent(
	applicantID kernel.UserID,
	content string,
	roomID kernel.RoomID,
) string {
	return "User " + string(applicantID) + " has requested to join room " + string(roomID) + ". Message: " + content
}
