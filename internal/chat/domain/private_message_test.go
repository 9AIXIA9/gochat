package domain_test

import (
	"gochat/internal/chat/domain"
	eventMock "gochat/internal/shared/event/mocks"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadPrivateMessage(t *testing.T) {
	message := domain.LoadPrivateMessage(
		fixedMessageID,
		fixedUserID,
		fixedFriendID,
		"Hello, Friend!",
		time.Now().UTC(),
	)

	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedFriendID, message.RecipientID())
	assert.Equal(t, "Hello, Friend!", message.Content())
}

func TestCreatePrivateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	start := time.Now()

	message, err := domain.CreatePrivateMessage(
		fixedFriendID,
		fixedUserID,
		"Hello, Friend!",
		mockMessageIDGenerator,
		mockEventIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedFriendID, message.RecipientID())
	assert.Equal(t, "Hello, Friend!", message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)

	events := message.GetEvents()
	require.Len(t, events, 1)

	evCreated := events[0]
	assert.Equal(t, domain.TopicPrivateMessageCreated, evCreated.Topic())
	assert.Equal(t, message.ID().String(), evCreated.AggregateID().String())
}
