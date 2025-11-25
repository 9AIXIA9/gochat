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
	roomMsgEventID event.ID         = "event-room-msg-notify"
	roomMsgID      kernel.MessageID = "room-msg-1"
	roomMsgSender  kernel.UserID    = "room-sender"
	roomIDEvent    kernel.RoomID    = "room-1"
)

func TestNewRoomMessageNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(roomMsgEventID)
	sentAt := time.Now().UTC()
	recipients := []kernel.UserID{"user-a", "user-b"}
	ev, err := domain.NewRoomMessageNotificationRequestedEvent(roomMsgID, roomMsgSender, roomIDEvent, recipients, "payload", sentAt, idGen)
	require.NoError(t, err)
	require.Equal(t, roomMsgEventID, ev.ID())
	require.Equal(t, domain.TopicRoomMessageNotificationRequested, ev.Topic())
	require.Equal(t, kernel.ID(roomMsgID), ev.AggregateID())
	require.Equal(t, roomMsgSender, ev.SenderID())
	require.Equal(t, roomIDEvent, ev.RoomID())
	require.ElementsMatch(t, recipients, ev.RecipientIDs())
}

func TestToRoomMessageNotificationRequestedEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	idGen := eventMocks.NewMockIDGenerator(ctrl)
	idGen.EXPECT().Generate().Return(roomMsgEventID).AnyTimes()
	recipients := []kernel.UserID{"user-a", "user-b"}
	sentAt := time.Now().UTC().Add(-time.Minute)
	payload, _ := json.Marshal(struct {
		MessageID    kernel.MessageID
		SenderID     kernel.UserID
		RoomID       kernel.RoomID
		RecipientIDs []kernel.UserID
		Content      string
		SentAt       time.Time
	}{roomMsgID, roomMsgSender, roomIDEvent, recipients, "content", sentAt})
	std := event.NewStandardEvent(kernel.ID(roomMsgID), domain.TopicRoomMessageNotificationRequested, payload, idGen)
	converted, err := domain.ToRoomMessageNotificationRequestedEvent(std)
	require.NoError(t, err)
	require.Equal(t, roomMsgSender, converted.SenderID())
	require.Equal(t, roomIDEvent, converted.RoomID())
	require.ElementsMatch(t, recipients, converted.RecipientIDs())
	require.Equal(t, "content", converted.Content())
	require.WithinDuration(t, sentAt, converted.SentAt(), 200*time.Millisecond)
	bad := event.NewStandardEvent(kernel.ID(roomMsgID), "wrong.topic", payload, idGen)
	converted2, err := domain.ToRoomMessageNotificationRequestedEvent(bad)
	require.ErrorIs(t, err, myErrors.ErrWrongEventType)
	require.Nil(t, converted2)
}
