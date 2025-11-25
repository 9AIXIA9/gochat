package domain_test

import (
	"testing"

	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomIDEv  kernel.RoomID = "social-room-123"
	fixedEventIDRC event.ID      = "social-event-room-created"
)

func TestNewRoomCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDRC)

	ev, err := domain.NewRoomCreatedEvent(fixedRoomIDEv, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedEventIDRC, ev.ID())
	require.Equal(t, domain.TopicRoomCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedRoomIDEv), ev.AggregateID())
}

func TestToRoomCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDRC).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedRoomIDEv), domain.TopicRoomCreated, []byte(""), idGen)
	converted, err := domain.ToRoomCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, converted)
	require.Equal(t, std.ID(), converted.ID())

	bad := event.NewStandardEvent(kernel.ID(fixedRoomIDEv), "wrong.topic", []byte(""), idGen)
	converted2, err := domain.ToRoomCreatedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestRoomCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDRC).Times(2)
	ev, err := domain.NewRoomCreatedEvent(fixedRoomIDEv, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.Equal(t, 0, len(payload))
	std := event.NewStandardEvent(kernel.ID(fixedRoomIDEv), domain.TopicRoomCreated, payload, idGen)
	converted, err := domain.ToRoomCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
}
