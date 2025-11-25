package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomCreatedEventID event.ID          = "room-created-ev-1"
	fixedRoomCreatedRoomID  kernel.RoomID     = "room-created-123"
	fixedRoomCreatedNumber  kernel.RoomNumber = "80001"
	fixedRoomCreatedOwnerID kernel.UserID     = "owner-abc"
)

func TestRoomCreatedEvent_NewRoomCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomCreatedEventID)

	ev, err := domain.NewRoomCreatedEvent(fixedRoomCreatedRoomID, fixedRoomCreatedNumber, fixedRoomCreatedOwnerID, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedRoomCreatedEventID, ev.ID())
	require.Equal(t, domain.TopicRoomCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomCreatedRoomID), ev.AggregateID())
	require.Equal(t, fixedRoomCreatedNumber, ev.Number())
	require.Equal(t, fixedRoomCreatedOwnerID, ev.OwnerID())
	// payload should be json with number & ownerID
	require.NotEmpty(t, ev.Payload())
}

func TestRoomCreatedEvent_ToRoomCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomCreatedEventID).AnyTimes()

	payload := []byte("{\"OwnerID\":\"" + fixedRoomCreatedOwnerID.String() + "\",\"Number\":\"" + fixedRoomCreatedNumber.String() + "\"}")
	std := event.NewStandardEvent(kernel.ID(fixedRoomCreatedRoomID), domain.TopicRoomCreated, payload, idGen)
	parsed, err := domain.ToRoomCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, parsed)
	require.Equal(t, fixedRoomCreatedOwnerID, parsed.OwnerID())
	require.Equal(t, fixedRoomCreatedNumber, parsed.Number())

	wrong := event.NewStandardEvent(kernel.ID(fixedRoomCreatedRoomID), "wrong.topic", payload, idGen)
	parsed2, err := domain.ToRoomCreatedEvent(wrong)
	require.ErrorIs(t, err, errors.ErrWrongEventType)
	require.Nil(t, parsed2)
}

func TestChatRoomCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedRoomCreatedEventID).Times(2)
	ev, err := domain.NewRoomCreatedEvent(fixedRoomCreatedRoomID, fixedRoomCreatedNumber, fixedRoomCreatedOwnerID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	std := event.NewStandardEvent(kernel.ID(fixedRoomCreatedRoomID), domain.TopicRoomCreated, payload, idGen)
	converted, err := domain.ToRoomCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.OwnerID(), converted.OwnerID())
	require.Equal(t, ev.Number(), converted.Number())
}
