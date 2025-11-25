package domain_test

import (
	"gochat/internal/authorization/domain"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedUserIDEv  = kernel.UserID("user-1234-event")
	fixedEventIDEv = event.ID("event-1234")
)

func TestNewUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test successful creation
	idGenerator := eventMock.NewMockIDGenerator(ctrl)
	idGenerator.EXPECT().Generate().Return(fixedEventIDEv)

	ev, err := domain.NewUserCreatedEvent(fixedUserIDEv, idGenerator)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedEventIDEv, ev.ID())
	require.Equal(t, domain.TopicUserCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedUserIDEv), ev.AggregateID())
}

func TestToUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGenerator := eventMock.NewMockIDGenerator(ctrl)
	idGenerator.EXPECT().Generate().Return(fixedEventIDEv).AnyTimes()

	// Prepare a standard event with correct topic
	standardEvent := event.NewStandardEvent(
		kernel.ID(fixedUserIDEv),
		domain.TopicUserCreated,
		[]byte{},
		idGenerator,
	)

	// Test successful conversion
	ev, err := domain.ToUserCreatedEvent(standardEvent)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, standardEvent.ID(), ev.ID())
	require.Equal(t, standardEvent.Topic(), ev.Topic())
	require.Equal(t, standardEvent.AggregateID(), ev.AggregateID())

	// Test wrong topic error
	wrongEvent := event.NewStandardEvent(
		kernel.ID(fixedUserIDEv),
		"wrong.topic",
		[]byte{},
		idGenerator,
	)

	ev2, err := domain.ToUserCreatedEvent(wrongEvent)
	require.Error(t, err)
	require.Nil(t, ev2)
}

func TestAuthorizationUserCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGenerator := eventMock.NewMockIDGenerator(ctrl)
	idGenerator.EXPECT().Generate().Return(fixedEventIDEv).Times(2)
	ev, err := domain.NewUserCreatedEvent(fixedUserIDEv, idGenerator)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.Len(t, payload, 0)
	std := event.NewStandardEvent(kernel.ID(fixedUserIDEv), domain.TopicUserCreated, payload, idGenerator)
	converted, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
}
