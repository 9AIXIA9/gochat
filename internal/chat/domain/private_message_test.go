package domain_test

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
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

	mockExister := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)

	mockExister.EXPECT().ExistByUserID(gomock.Any(), fixedUserID, fixedFriendID).Return(true, nil).Times(1)
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	start := time.Now().UTC()
	message, err := domain.CreatePrivateMessage(
		context.Background(),
		fixedFriendID,
		fixedUserID,
		"Hello, Friend!",
		mockEventIDGenerator,
		mockMessageIDGenerator,
		mockExister,
	)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedFriendID, message.RecipientID())
	assert.Equal(t, "Hello, Friend!", message.Content())
	evs := message.GetEvents()
	require.Len(t, evs, 1)
	assert.Equal(t, fixedEventID, evs[0].ID())
	assert.Equal(t, domain.TopicPrivateMessageCreated, evs[0].Topic())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)

	// content 为空
	message, err = domain.CreatePrivateMessage(
		context.Background(),
		fixedFriendID,
		fixedUserID,
		"",
		mockEventIDGenerator,
		mockMessageIDGenerator,
		mockExister,
	)
	require.ErrorIs(t, err, domain.ErrEmptyMessageContent)
	require.Nil(t, message)

	// 发送给自己
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1)

	start = time.Now().UTC()
	messageToSelf, err := domain.CreatePrivateMessage(
		context.Background(),
		fixedUserID,
		fixedUserID,
		"Hello, Self!",
		mockEventIDGenerator,
		mockMessageIDGenerator,
		mockExister,
	)
	require.NoError(t, err)
	require.NotNil(t, messageToSelf)
	assert.Equal(t, fixedMessageID, messageToSelf.ID())
	assert.Equal(t, fixedUserID, messageToSelf.SenderID())
	assert.Equal(t, fixedUserID, messageToSelf.RecipientID())
	assert.Equal(t, "Hello, Self!", messageToSelf.Content())
	evs = messageToSelf.GetEvents()
	require.Len(t, evs, 1)
	assert.Equal(t, fixedEventID, evs[0].ID())
	assert.Equal(t, domain.TopicPrivateMessageCreated, evs[0].Topic())
	assert.WithinDuration(t, start, messageToSelf.SentAt(), timeTolerance)
}
