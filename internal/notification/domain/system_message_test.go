package domain_test

import (
	"errors"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadSystemMessage(t *testing.T) {
	start := time.Now()
	message := domain.LoadSystemMessage(
		fixedMessageID,
		fixedRecipientID,
		domain.MessageStateUndelivered,
		fixedContent,
		time.Now().UTC(),
	)
	require.NotNil(t, message)

	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedRecipientID, message.RecipientID())
	assert.Equal(t, domain.MessageStateUndelivered, message.State())
	assert.Equal(t, fixedContent, message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)
}

func TestCreateSystemMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	start := time.Now()
	mockIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	//正常情况
	mockIDGenerator.EXPECT().
		Generate().
		Return(fixedMessageID).
		Times(1)

	mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil)

	message, err := domain.CreateSystemMessage(
		fixedRecipientID,
		fixedContent,
		mockIDGenerator,
		mockNotifier,
	)
	require.NoError(t, err)
	require.NotNil(t, message)

	assert.Equal(t, fixedMessageID, message.ID())
	assert.Equal(t, fixedRecipientID, message.RecipientID())
	assert.Equal(t, domain.MessageStateDelivered, message.State())
	assert.Equal(t, fixedContent, message.Content())
	assert.WithinDuration(t, start, message.SentAt(), timeTolerance)

	//通知失败但不影响消息创建
	mockIDGenerator.EXPECT().
		Generate().
		Return(fixedMessageID).
		Times(1)

	mockNotifier.EXPECT().
		Notify(gomock.Any()).
		Return(errors.New("test-notify-failed")).
		Times(1)

	messageWithFailedNotification, err := domain.CreateSystemMessage(
		fixedRecipientID,
		fixedContent,
		mockIDGenerator,
		mockNotifier,
	)
	require.NoError(t, err)
	require.NotNil(t, messageWithFailedNotification)
	assert.Equal(t, domain.MessageStateUndelivered, messageWithFailedNotification.State())

	//内容为空
	messageWithEmptyContent, err := domain.CreateSystemMessage(
		fixedRecipientID,
		"",
		mockIDGenerator,
		mockNotifier,
	)
	require.ErrorIs(t, err, domain.ErrEmptyContent)
	require.Nil(t, messageWithEmptyContent)
}

func TestSystemMessage_Deliver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifier := mocks.NewMockSystemMessageNotifier(ctrl)

	//正常情况
	message := domain.LoadSystemMessage(
		fixedMessageID,
		fixedRecipientID,
		domain.MessageStateUndelivered,
		fixedContent,
		time.Now().UTC(),
	)

	mockNotifier.EXPECT().
		Notify(message).
		Return(nil).
		Times(1)

	err := message.Deliver(mockNotifier)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStateDelivered, message.State())

	//消息状态不是未发送
	err = message.Deliver(mockNotifier)
	require.NoError(t, err)
}
