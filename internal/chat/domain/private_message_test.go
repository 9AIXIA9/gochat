package domain_test

import (
	"context"
	"errors"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
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
		domain.MessageStateRead,
		time.Now().UTC(),
	)

	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedFriendID, message.RecipientID())
	assert.Equal(t, "Hello, Friend!", message.Content())
	assert.Equal(t, domain.MessageStateRead, message.State())
}

func TestCreatePrivateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExister := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)

	mockExister.EXPECT().ExistByUserID(gomock.Any(), fixedUserID, fixedFriendID).Return(true, nil).Times(1)
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(1)

	start := time.Now().UTC()
	message, err := domain.CreatePrivateMessage(
		context.Background(),
		fixedFriendID,
		fixedUserID,
		"Hello, Friend!",
		mockMessageIDGenerator,
		mockNotifier,
		mockExister,
	)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedUserID, message.SenderID())
	assert.Equal(t, fixedFriendID, message.RecipientID())
	assert.Equal(t, "Hello, Friend!", message.Content())
	assert.Equal(t, domain.MessageStateUndelivered, message.State())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)

	// content 为空
	message, err = domain.CreatePrivateMessage(
		context.Background(),
		fixedFriendID,
		fixedUserID,
		"",
		mockMessageIDGenerator,
		mockNotifier,
		mockExister,
	)
	require.ErrorIs(t, err, domain.ErrEmptyMessageContent)
	require.Nil(t, message)

	// 发送失败
	mockExister.EXPECT().ExistByUserID(gomock.Any(), fixedUserID, fixedFriendID).Return(true, nil).Times(1)
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockNotifier.EXPECT().Notify(gomock.Any()).Return(errors.New("test")).Times(1)

	start = time.Now().UTC()
	messageWithFailedDeliver, err := domain.CreatePrivateMessage(
		context.Background(),
		fixedFriendID,
		fixedUserID,
		"Hello, Friend!",
		mockMessageIDGenerator,
		mockNotifier,
		mockExister,
	)
	require.NoError(t, err)
	require.NotNil(t, messageWithFailedDeliver)
	assert.Equal(t, fixedMessageID, messageWithFailedDeliver.ID())
	assert.Equal(t, fixedUserID, messageWithFailedDeliver.SenderID())
	assert.Equal(t, fixedFriendID, messageWithFailedDeliver.RecipientID())
	assert.Equal(t, "Hello, Friend!", messageWithFailedDeliver.Content())
	assert.Equal(t, domain.MessageStateUndelivered, messageWithFailedDeliver.State())
	assert.WithinDuration(t, start, messageWithFailedDeliver.SentAt(), timeTolerance)

	// 发送给自己，也尝试投递，但创建态仍然保持未投递
	mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1)
	mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(1)

	start = time.Now().UTC()
	messageToSelf, err := domain.CreatePrivateMessage(
		context.Background(),
		fixedUserID,
		fixedUserID,
		"Hello, Self!",
		mockMessageIDGenerator,
		mockNotifier,
		mockExister,
	)
	require.NoError(t, err)
	require.NotNil(t, messageToSelf)
	assert.Equal(t, fixedMessageID, messageToSelf.ID())
	assert.Equal(t, fixedUserID, messageToSelf.SenderID())
	assert.Equal(t, fixedUserID, messageToSelf.RecipientID())
	assert.Equal(t, "Hello, Self!", messageToSelf.Content())
	assert.Equal(t, domain.MessageStateUndelivered, messageToSelf.State())
	assert.WithinDuration(t, start, messageToSelf.SentAt(), timeTolerance)
}

func TestPrivateMessage_Deliver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)

	message := domain.LoadPrivateMessage(
		fixedMessageID,
		fixedUserID,
		fixedFriendID,
		"Hello, Friend!",
		domain.MessageStateUndelivered,
		time.Now().UTC(),
	)

	// 成功投递
	mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(1)

	err := message.Deliver(mockNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, message.State())

	// 已经是已投递状态，跳过投递
	err = message.Deliver(mockNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, message.State())

	// 投递失败
	messageUndelivered := domain.LoadPrivateMessage(
		fixedMessageID,
		fixedUserID,
		fixedFriendID,
		"Hello, Friend!",
		domain.MessageStateUndelivered,
		time.Now().UTC(),
	)

	mockNotifier.EXPECT().Notify(gomock.Any()).Return(errors.New("test")).Times(1)

	err = messageUndelivered.Deliver(mockNotifier)
	require.Error(t, err)
	assert.Equal(t, domain.MessageStateUndelivered, messageUndelivered.State())
}
