package domain_test

import (
	"testing"
	"time"

	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	"gochat/internal/shared/kernel"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	rmID       kernel.MessageID = "rm-1"
	rmSenderID kernel.UserID    = "user-sender"
	rmRoomID   kernel.RoomID    = "room-777"
	rmContent  string           = "room msg"
)

func TestReceiveRoomMessageAndDeliver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	recipients := []kernel.UserID{"user-a", "user-b", "user-c"}
	m := domain.ReceiveRoomMessage(rmID, rmSenderID, rmRoomID, recipients, rmContent, time.Now().UTC())
	require.NotNil(t, m)
	for _, r := range recipients {
		assert.Equal(t, domain.MessageStateUndelivered, m.States()[r])
	}

	notifier := mocks.NewMockRoomMessageNotifier(ctrl)
	// Only deliver to undelivered recipients; we simulate delivering two of them
	undelivered := append([]kernel.UserID{}, recipients...)
	notifier.EXPECT().NotifyRoomMessage(m, undelivered).Return([]kernel.UserID{"user-a", "user-b"}, nil)
	err := m.Deliver(notifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, m.States()["user-a"])
	assert.Equal(t, domain.MessageStateDelivered, m.States()["user-b"])
	assert.Equal(t, domain.MessageStateUndelivered, m.States()["user-c"])

	// Mark last as read
	m.Read("user-c")
	assert.Equal(t, domain.MessageStateDelivered, m.States()["user-c"])
}

func TestRoomMessage_DeliverError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	recipients := []kernel.UserID{"user-a", "user-b"}
	m := domain.ReceiveRoomMessage(rmID, rmSenderID, rmRoomID, recipients, rmContent, time.Now().UTC())
	notifier := mocks.NewMockRoomMessageNotifier(ctrl)
	undelivered := append([]kernel.UserID{}, recipients...)
	notifier.EXPECT().NotifyRoomMessage(m, undelivered).Return(nil, assert.AnError)
	err := m.Deliver(notifier)
	require.Error(t, err)
	for _, r := range recipients {
		assert.Equal(t, domain.MessageStateUndelivered, m.States()[r])
	}
}

func TestRoomMessage_LoadRoomMessage(t *testing.T) {
	recipients := []kernel.UserID{"user-x", "user-y"}
	states := map[kernel.UserID]domain.MessageState{
		"user-x": domain.MessageStateDelivered,
		"user-y": domain.MessageStateUndelivered,
	}
	sentAt := time.Now().UTC()
	m := domain.LoadRoomMessage(rmID, rmSenderID, rmRoomID, recipients, states, rmContent, sentAt)
	require.NotNil(t, m)
	assert.Equal(t, rmID, m.ID())
	assert.Equal(t, rmSenderID, m.SenderID())
	assert.Equal(t, rmRoomID, m.RoomID())
	assert.Equal(t, recipients, m.RecipientIDs())
	assert.Equal(t, states, m.States())
	assert.Equal(t, rmContent, m.Content())
	assert.Equal(t, sentAt, m.SentAt())
}
