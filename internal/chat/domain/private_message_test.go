package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	"gochat/internal/shared/event"
	eventMocks "gochat/internal/shared/event/mocks"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedPrivateMessageID   kernel.MessageID = "priv-msg-123"
	fixedPrivateEventID     event.ID         = "event-priv-1"
	fixedPrivateSenderID    kernel.UserID    = "user-sender"
	fixedPrivateRecipient   kernel.UserID    = "user-recipient"
	privateMessageContent                    = "hello private"
	privateMsgTimeTolerance                  = 150 * time.Millisecond
)

func TestPrivateMessage_CreatePrivateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recipient := domain.CreateUser(fixedPrivateRecipient, "50001")

	msgIDGen := mocks.NewMockMessageIDGenerator(ctrl)
	msgIDGen.EXPECT().Generate().Return(fixedPrivateMessageID)

	evIDGen := eventMocks.NewMockIDGenerator(ctrl)
	evIDGen.EXPECT().Generate().Return(fixedPrivateEventID)

	start := time.Now().UTC()
	pm, err := domain.CreatePrivateMessage(recipient, fixedPrivateSenderID, privateMessageContent, msgIDGen, evIDGen)
	require.NoError(t, err)
	require.NotNil(t, pm)
	assert.Equal(t, fixedPrivateMessageID, pm.ID())
	assert.Equal(t, fixedPrivateSenderID, pm.SenderID())
	assert.Equal(t, fixedPrivateRecipient, pm.RecipientID())
	assert.Equal(t, privateMessageContent, pm.Content())
	assert.WithinDuration(t, start, pm.SentAt(), privateMsgTimeTolerance)

	events := pm.GetEvents()
	assert.Len(t, events, 1)
	assert.Empty(t, pm.GetEvents()) // drained

	createdEv := events[0]
	assert.Equal(t, fixedPrivateEventID, createdEv.ID())
	assert.Equal(t, domain.TopicPrivateMessageCreated, createdEv.Topic())
	assert.Equal(t, kernel.ID(fixedPrivateMessageID), createdEv.AggregateID())
	assert.Equal(t, []byte(""), createdEv.Payload())
}

func TestPrivateMessage_LoadPrivateMessage(t *testing.T) {
	pm := domain.LoadPrivateMessage(fixedPrivateMessageID, fixedPrivateSenderID, fixedPrivateRecipient, privateMessageContent, time.Now().UTC())
	require.NotNil(t, pm)
	assert.Equal(t, fixedPrivateMessageID, pm.ID())
	assert.Equal(t, fixedPrivateSenderID, pm.SenderID())
	assert.Equal(t, fixedPrivateRecipient, pm.RecipientID())
	assert.Equal(t, privateMessageContent, pm.Content())
}
