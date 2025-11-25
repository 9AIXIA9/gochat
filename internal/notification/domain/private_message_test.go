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
	pmID       kernel.MessageID = "pm-1"
	pmSenderID kernel.UserID    = "user-sender"
	pmRecID    kernel.UserID    = "user-recipient"
	pmContent  string           = "hello there"
)

func TestPrivateMessage_ReceiveAndDeliverAndRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sentAt := time.Now().UTC()
	m := domain.ReceivePrivateMessage(pmID, pmSenderID, pmRecID, pmContent, sentAt)
	require.NotNil(t, m)
	assert.Equal(t, pmID, m.ID())
	assert.Equal(t, domain.MessageStateUndelivered, m.State())
	assert.Equal(t, pmContent, m.Content())
	assert.Equal(t, pmSenderID, m.SenderID())
	assert.Equal(t, pmRecID, m.RecipientID())
	assert.WithinDuration(t, sentAt, m.SentAt(), 200*time.Millisecond)

	notifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	notifier.EXPECT().NotifyPrivateMessage(m).Return(nil)
	err := m.Deliver(notifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, m.State())

	// read changes to read state
	m.Read()
	assert.Equal(t, domain.MessageStateRead, m.State())
}

func TestPrivateMessage_DeliverError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := domain.ReceivePrivateMessage(pmID, pmSenderID, pmRecID, pmContent, time.Now().UTC())
	notifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	notifier.EXPECT().NotifyPrivateMessage(m).Return(assert.AnError)
	err := m.Deliver(notifier)
	require.Error(t, err)
	assert.Equal(t, domain.MessageStateUndelivered, m.State())
}

func TestPrivateMessage_LoadRetainsState(t *testing.T) {
	sentAt := time.Now().UTC().Add(-time.Minute)
	loaded := domain.LoadPrivateMessage(pmID, pmSenderID, pmRecID, domain.MessageStateRead, pmContent, sentAt)
	require.NotNil(t, loaded)
	assert.Equal(t, domain.MessageStateRead, loaded.State())
	assert.Equal(t, pmContent, loaded.Content())
	assert.Equal(t, pmSenderID, loaded.SenderID())
	assert.Equal(t, pmRecID, loaded.RecipientID())
	assert.Equal(t, sentAt, loaded.SentAt())
}
