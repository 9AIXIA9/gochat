package domain_test

import (
	"encoding/json"
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
	welcomeEventID event.ID          = "event-welcome-email"
	welcomeUserID  kernel.UserID     = "welcome-user"
	welcomeEmail   kernel.Email      = "welcome@example.com"
	welcomeNumber  kernel.UserNumber = "123456"
)

func TestNewWelcomeEmailNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(welcomeEventID)
	ev, err := domain.NewWelcomeEmailNotificationRequestedEvent(welcomeUserID, welcomeEmail, welcomeNumber, idGen)
	require.NoError(t, err)
	require.Equal(t, welcomeEventID, ev.ID())
	require.Equal(t, domain.TopicWelcomeEmailNotificationRequested, ev.Topic())
	require.Equal(t, kernel.ID(welcomeUserID), ev.AggregateID())
	require.Equal(t, welcomeEmail, ev.Email())
	require.Equal(t, welcomeNumber, ev.Number())
}

func TestToWelcomeEmailNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(welcomeEventID).AnyTimes()
	payload, _ := json.Marshal(struct {
		Email  kernel.Email
		Number kernel.UserNumber
	}{welcomeEmail, welcomeNumber})
	std := event.NewStandardEvent(kernel.ID(welcomeUserID), domain.TopicWelcomeEmailNotificationRequested, payload, idGen)
	converted, err := domain.ToWelcomeEmailNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, welcomeEmail, converted.Email())
	require.Equal(t, welcomeNumber, converted.Number())
	bad := event.NewStandardEvent(kernel.ID(welcomeUserID), "wrong.topic", payload, idGen)
	converted2, err := domain.ToWelcomeEmailNotificationRequestedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}

func TestWelcomeEmailNotificationRequestedEvent_MarshalUnmarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(welcomeEventID).Times(2)
	ev, err := domain.NewWelcomeEmailNotificationRequestedEvent(welcomeUserID, welcomeEmail, welcomeNumber, idGen)
	require.NoError(t, err)
	payload, err := ev.Marshal()
	require.NoError(t, err)
	var decoded struct {
		Email  kernel.Email
		Number kernel.UserNumber
	}
	require.NoError(t, json.Unmarshal(payload, &decoded))
	require.Equal(t, welcomeEmail, decoded.Email)
	require.Equal(t, welcomeNumber, decoded.Number)
	std := event.NewStandardEvent(kernel.ID(welcomeUserID), domain.TopicWelcomeEmailNotificationRequested, payload, idGen)
	converted, err := domain.ToWelcomeEmailNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, ev.Email(), converted.Email())
}
