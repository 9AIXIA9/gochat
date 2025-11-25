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
	fixedUserIDEv  kernel.UserID = "social-user-123"
	fixedEventIDEv event.ID      = "social-event-123"
)

func TestNewUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDEv)

	ev, err := domain.NewUserCreatedEvent(fixedUserIDEv, idGen)
	require.NoError(t, err)
	require.NotNil(t, ev)
	require.Equal(t, fixedEventIDEv, ev.ID())
	require.Equal(t, domain.TopicUserCreated, ev.Topic())
	require.Equal(t, kernel.ID(fixedUserIDEv), ev.AggregateID())
}

func TestToUserCreatedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDEv).AnyTimes()

	std := event.NewStandardEvent(kernel.ID(fixedUserIDEv), domain.TopicUserCreated, []byte(""), idGen)
	converted, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.NotNil(t, converted)
	require.Equal(t, std.ID(), converted.ID())

	bad := event.NewStandardEvent(kernel.ID(fixedUserIDEv), "wrong.topic", []byte(""), idGen)
	converted2, err := domain.ToUserCreatedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestUserCreatedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(fixedEventIDEv).Times(2)
	// create event
	ev, err := domain.NewUserCreatedEvent(fixedUserIDEv, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	// empty payload expected
	require.Equal(t, 0, len(payload))
	// build standard event with same payload and convert
	std := event.NewStandardEvent(kernel.ID(fixedUserIDEv), domain.TopicUserCreated, payload, idGen)
	converted, err := domain.ToUserCreatedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
	// unmarshal on empty should be no-op
	require.Equal(t, kernel.ID(fixedUserIDEv), converted.AggregateID())
}
