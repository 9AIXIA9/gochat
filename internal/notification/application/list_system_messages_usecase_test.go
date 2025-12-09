package application_test

import (
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedSystemMessageBaseID   kernel.MessageID = "message-123"
	maxSystemMessagesLimit                      = 100
	defaultSystemMessagesLimit                  = 50
)

func TestListSystemMessagesInput_Validate(t *testing.T) {
	input := &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  maxSystemMessagesLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxSystemMessagesLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListSystemMessagesInput{
		UserID: "",
		BaseID: fixedSystemMessageBaseID,
		Limit:  maxSystemMessagesLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListSystemMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockSystemMessageFinderByUserID(ctrl)

	useCase, err := application.NewListSystemMessagesUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListSystemMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListSystemMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockSystemMessageFinderByUserID(ctrl)

	useCase, err := application.NewListSystemMessagesUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxSystemMessagesLimit - 1
	messages := getSystemMessages(fixedLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, fixedLimit, fixedSystemMessageBaseID).Return(messages, nil)

	output, err := useCase.Execute(nil, &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  maxSystemMessagesLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.SystemMessages, fixedLimit)

	// 模拟limit为0
	messages = getSystemMessages(defaultSystemMessagesLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultSystemMessagesLimit, fixedSystemMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.SystemMessages, defaultSystemMessagesLimit)

	//模拟负数limit
	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, defaultSystemMessagesLimit, fixedSystemMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.SystemMessages, defaultSystemMessagesLimit)

	// 模拟太大的limit
	messages = getSystemMessages(maxSystemMessagesLimit)

	finder.EXPECT().FindsByUserID(gomock.Any(), fixedUserID, maxSystemMessagesLimit, fixedSystemMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListSystemMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedSystemMessageBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.SystemMessages, maxSystemMessagesLimit)
}

func getSystemMessages(count int) []*domain.SystemMessage {
	messages := make([]*domain.SystemMessage, count)
	for i := 0; i < count; i++ {
		messages[i] = domain.LoadSystemMessage(
			fixedMessageID,
			fixedUserID,
			domain.MessageStateDelivered,
			fixedContent,
			time.Now().UTC(),
		)
	}
	return messages
}
