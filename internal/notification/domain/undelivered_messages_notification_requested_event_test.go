package domain_test

import (
	"testing"

	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	undeliveredEventID event.ID      = "event-undelivered-request"
	undeliveredUserID  kernel.UserID = "user-undelivered"
)

func TestNewUndeliveredMessagesNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(undeliveredEventID)
	ev, err := domain.NewUndeliveredMessagesNotificationRequestedEvent(undeliveredUserID, idGen)
	require.NoError(t, err)
	require.Equal(t, undeliveredEventID, ev.ID())
	require.Equal(t, domain.TopicUndeliveredMessagesNotificationRequested, ev.Topic())
	require.Equal(t, kernel.ID(undeliveredUserID), ev.AggregateID())
}

func TestToUndeliveredMessagesNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(undeliveredEventID).AnyTimes()
	std := event.NewStandardEvent(kernel.ID(undeliveredUserID), domain.TopicUndeliveredMessagesNotificationRequested, []byte(""), idGen)
	converted, err := domain.ToUndeliveredMessagesNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, std.ID(), converted.ID())
	bad := event.NewStandardEvent(kernel.ID(undeliveredUserID), "wrong.topic", []byte(""), idGen)
	converted2, err := domain.ToUndeliveredMessagesNotificationRequestedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestUndeliveredMessagesNotificationRequestedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(undeliveredEventID).Times(2)
	ev, err := domain.NewUndeliveredMessagesNotificationRequestedEvent(undeliveredUserID, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	require.Equal(t, 0, len(payload))
	std := event.NewStandardEvent(kernel.ID(undeliveredUserID), domain.TopicUndeliveredMessagesNotificationRequested, payload, idGen)
	converted, err := domain.ToUndeliveredMessagesNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.ID(), converted.ID())
}
