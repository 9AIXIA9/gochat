package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedPrivateMessageBaseID   kernel.MessageID = "message-123"
	maxPrivateMessagesLimit                      = 100
	defaultPrivateMessagesLimit                  = 50
)

func TestListPrivateMessagesInput_Validate(t *testing.T) {
	input := &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  maxPrivateMessagesLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxPrivateMessagesLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListPrivateMessagesInput{
		UserID: "",
		BaseID: fixedPrivateMessageBaseID,
		Limit:  maxPrivateMessagesLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListPrivateMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockPrivateMessagesFinderByRecipientID(ctrl)

	useCase, err := application.NewListPrivateMessagesUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListPrivateMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListPrivateMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockPrivateMessagesFinderByRecipientID(ctrl)

	useCase, err := application.NewListPrivateMessagesUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxPrivateMessagesLimit - 1
	messages := getPrivateMessages(fixedLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, fixedLimit, fixedPrivateMessageBaseID).Return(messages, nil)

	output, err := useCase.Execute(nil, &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  maxPrivateMessagesLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.PrivateMessages, fixedLimit)

	// 模拟limit为0
	messages = getPrivateMessages(defaultPrivateMessagesLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, defaultPrivateMessagesLimit, fixedPrivateMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.PrivateMessages, defaultPrivateMessagesLimit)

	//模拟负数limit
	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, defaultPrivateMessagesLimit, fixedPrivateMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.PrivateMessages, defaultPrivateMessagesLimit)

	// 模拟太大的limit
	messages = getPrivateMessages(maxPrivateMessagesLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, maxPrivateMessagesLimit, fixedPrivateMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListPrivateMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedPrivateMessageBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.PrivateMessages, maxPrivateMessagesLimit)
}

func getPrivateMessages(count int) []*domain.PrivateMessage {
	messages := make([]*domain.PrivateMessage, count)
	for i := 0; i < count; i++ {
		messages[i] = domain.LoadPrivateMessage(
			fixedMessageID,
			fixedFriendID,
			fixedUserID,
			fixedContent,
			domain.MessageStateDelivered,
			time.Now().UTC(),
		)
	}
	return messages
}
