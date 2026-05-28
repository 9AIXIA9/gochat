package domain_test

import (
	"encoding/json"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	eventMock "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedEventID     event.ID         = "event-123"
	fixedMessageID   kernel.MessageID = "notification-123"
	fixedRecipientID kernel.UserID    = "user-456"
	fixedSenderID    kernel.UserID    = "user-789"
	fixedContent     string           = "Hello, this is a notification!"
)

type AliasPayload struct {
	Content  string
	SenderID kernel.UserID
	SendAt   time.Time
}

func TestLoadNotification(t *testing.T) {
	fixedPayload, err := json.Marshal(&AliasPayload{
		Content:  fixedContent,
		SenderID: fixedSenderID,
		SendAt:   time.Now().UTC(),
	})
	require.NoError(t, err)

	notification := domain.LoadNotification(
		fixedMessageID,
		fixedRecipientID,
		domain.StateDelivered,
		fixedPayload,
	)

	require.NotNil(t, notification)
	assert.Equal(t, fixedMessageID, notification.ID())
	assert.Equal(t, fixedRecipientID, notification.RecipientID())
	assert.Equal(t, domain.StateDelivered, notification.State())
	assert.Equal(t, fixedPayload, []byte(notification.RawPayload()))

	var tempPayload AliasPayload
	err = json.Unmarshal(notification.RawPayload(), &tempPayload)
	require.NoError(t, err)
	require.NoError(t, err)
	assert.Equal(t, fixedContent, tempPayload.Content)
	assert.Equal(t, fixedSenderID, tempPayload.SenderID)
	assert.WithinDuration(t, time.Now().UTC(), tempPayload.SendAt, time.Minute)
}

func TestCreateNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	// 正常情况
	mockIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)
	fixedPayload, err := json.Marshal(&AliasPayload{
		Content:  fixedContent,
		SenderID: fixedSenderID,
		SendAt:   time.Now().UTC(),
	})
	require.NoError(t, err)

	notification, err := domain.CreateNotification(
		fixedMessageID,
		fixedRecipientID,
		fixedPayload,
		mockIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, notification)
	assert.Equal(t, fixedMessageID, notification.ID())
	assert.Equal(t, fixedRecipientID, notification.RecipientID())
	assert.Equal(t, fixedPayload, []byte(notification.RawPayload()))
	evs := notification.GetEvents()
	require.Len(t, evs, 1)
	assert.Equal(t, fixedEventID, evs[0].ID())
	assert.Equal(t, contract.TopicPushRequested, evs[0].Topic())

	// 错误情况：没有 recipient ID
	_, err = domain.CreateNotification(
		fixedMessageID,
		"",
		fixedPayload,
		mockIDGenerator,
	)
	require.ErrorIs(t, err, domain.ErrNotificationNoRecipient)
}
