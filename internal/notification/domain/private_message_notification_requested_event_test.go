package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	"gochat/internal/notification/domain"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	pmEventID     event.ID         = "event-pm-notify"
	pmMessageID   kernel.MessageID = "pm-msg-1"
	pmSenderIDEv  kernel.UserID    = "pm-sender"
	pmRecipientID kernel.UserID    = "pm-recipient"
)

func TestNewPrivateMessageNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(pmEventID)
	sentAt := time.Now().UTC()
	ev, err := domain.NewPrivateMessageNotificationRequestedEvent(pmMessageID, pmRecipientID, pmSenderIDEv, "hello", sentAt, idGen)
	require.NoError(t, err)
	require.Equal(t, pmEventID, ev.ID())
	require.Equal(t, domain.TopicPrivateMessageNotificationRequested, ev.Topic())
	require.Equal(t, kernel.ID(pmMessageID), ev.AggregateID())
	require.Equal(t, pmSenderIDEv, ev.SenderID())
	require.Equal(t, pmRecipientID, ev.RecipientID())
	require.Equal(t, "hello", ev.Content())
	require.WithinDuration(t, sentAt, ev.SentAt(), 200*time.Millisecond)
}

func TestToPrivateMessageNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(pmEventID).AnyTimes()
	sentAt := time.Now().UTC().Add(-time.Minute)
	payload, _ := json.Marshal(struct {
		SenderID    kernel.UserID
		RecipientID kernel.UserID
		Content     string
		SentAt      time.Time
	}{pmSenderIDEv, pmRecipientID, "content", sentAt})
	std := event.NewStandardEvent(kernel.ID(pmMessageID), domain.TopicPrivateMessageNotificationRequested, payload, idGen)
	converted, err := domain.ToPrivateMessageNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, pmSenderIDEv, converted.SenderID())
	require.Equal(t, pmRecipientID, converted.RecipientID())
	require.Equal(t, "content", converted.Content())
	require.WithinDuration(t, sentAt, converted.SentAt(), 200*time.Millisecond)
	bad := event.NewStandardEvent(kernel.ID(pmMessageID), "wrong.topic", payload, idGen)
	converted2, err := domain.ToPrivateMessageNotificationRequestedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}
