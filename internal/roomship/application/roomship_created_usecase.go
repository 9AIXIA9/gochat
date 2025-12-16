package application

import (
	"context"
	chatDomain "gochat/internal/chat/domain"
	notificationDomain "gochat/internal/notification/domain"
	profileDomain "gochat/internal/profile/domain"
	"gochat/internal/roomship/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type RoomshipCreatedUseCase kernel.UseCase[*RoomshipCreatedInput, *kernel.NoOutput]

type RoomshipCreatedInput struct {
	RoomshipID domain.RoomshipID
}

func (r *RoomshipCreatedInput) Validate() error {
	if len(r.RoomshipID) == 0 {
		return myErrors.WrapBusiness(myErrors.ErrEmptyInput, "roomship id is empty")
	}
	return nil
}

type roomshipCreatedUseCase struct {
	roomshipFinderByID     domain.RoomshipFinderByID
	roomshipFinderByRoomID domain.RoomshipsFinderByRoomID
	idGenerator            event.IDGenerator
	creator                event.UnpublishedEventsCreator
}

func NewRoomshipCreatedUseCase(
	roomshipFinderByID domain.RoomshipFinderByID,
	roomshipFinderByRoomID domain.RoomshipsFinderByRoomID,
	idGenerator event.IDGenerator,
	creator event.UnpublishedEventsCreator,
) (RoomshipCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		roomshipFinderByID,
		roomshipFinderByRoomID,
		idGenerator,
		creator,
	); err != nil {
		return nil, err
	}

	return &roomshipCreatedUseCase{
		roomshipFinderByID:     roomshipFinderByID,
		roomshipFinderByRoomID: roomshipFinderByRoomID,
		idGenerator:            idGenerator,
		creator:                creator,
	}, nil
}

func (uc *roomshipCreatedUseCase) Execute(ctx context.Context, input *RoomshipCreatedInput) (*kernel.NoOutput, error) {
	newRoomship, err := uc.roomshipFinderByID.FindByID(ctx, input.RoomshipID)
	if err != nil {
		return nil, err
	}

	allRoomships, err := uc.roomshipFinderByRoomID.FindsByRoomID(ctx, newRoomship.RoomID())
	if err != nil {
		return nil, err
	}

	evs := make([]event.Event, 0, len(allRoomships)+1)
	for _, oldRoomship := range allRoomships {
		if oldRoomship.ID() == input.RoomshipID {
			continue
		}

		notificationEvRequested, err := notificationDomain.NewSystemMessageNotificationRequestedEvent(
			oldRoomship.UserID(),
			uc.buildContent(
				newRoomship.UserID(),
				newRoomship.RoomID(),
			),
			uc.idGenerator,
		)
		if err != nil {
			return nil, err
		}
		evs = append(evs, notificationEvRequested)
	}

	chatEvCreated, err := chatDomain.NewRoomshipCreatedEvent(
		chatDomain.RoomshipID(newRoomship.ID()),
		newRoomship.UserID(),
		newRoomship.RoomID(),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	profileEvCreated, err := profileDomain.NewRoomshipCreatedEvent(
		profileDomain.RoomshipID(newRoomship.ID()),
		newRoomship.UserID(),
		newRoomship.RoomID(),
		profileDomain.RoomshipRole(newRoomship.Role()),
		uc.idGenerator,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.creator.CreateUnpublishedEvents(ctx, append(evs, chatEvCreated, profileEvCreated)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (uc *roomshipCreatedUseCase) buildContent(
	userID kernel.UserID,
	roomID kernel.RoomID,
) string {
	return "User " + string(userID) + " has joined the room " + string(roomID) + "."
}
