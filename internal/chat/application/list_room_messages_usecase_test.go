package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	fixedRoomMessageBaseID   kernel.MessageID = "message-123"
	maxRoomMessagesLimit                      = 100
	defaultRoomMessagesLimit                  = 50
)

func TestListRoomMessagesInput_Validate(t *testing.T) {
	input := &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  maxRoomMessagesLimit,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithNoBaseID := &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  maxRoomMessagesLimit,
	}

	err = inputWithNoBaseID.Validate()
	require.NoError(t, err)

	inputWithNoLimit := &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  0,
	}

	err = inputWithNoLimit.Validate()
	require.NoError(t, err)

	inputWithNegativeLimit := &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  -10,
	}

	err = inputWithNegativeLimit.Validate()
	require.NoError(t, err)

	inputWithNoBaseIDAndLimit := &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: "",
		Limit:  0,
	}

	err = inputWithNoBaseIDAndLimit.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.ListRoomMessagesInput{
		UserID: "",
		BaseID: fixedRoomMessageBaseID,
		Limit:  maxRoomMessagesLimit,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewListRoomMessagesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomMessagesFinderByRecipientID(ctrl)

	useCase, err := application.NewListRoomMessagesUseCase(
		finder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewListRoomMessagesUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestListRoomMessagesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	finder := mocks.NewMockRoomMessagesFinderByRecipientID(ctrl)

	useCase, err := application.NewListRoomMessagesUseCase(
		finder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 模拟正常情况，limit在范围内
	fixedLimit := maxRoomMessagesLimit - 1
	messages := getRoomMessages(fixedLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, fixedLimit, fixedRoomMessageBaseID).Return(messages, nil)

	output, err := useCase.Execute(nil, &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  maxRoomMessagesLimit - 1,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.RoomMessages, fixedLimit)

	// 模拟limit为0
	messages = getRoomMessages(defaultRoomMessagesLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, defaultRoomMessagesLimit, fixedRoomMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  0,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.RoomMessages, defaultRoomMessagesLimit)

	//模拟负数limit
	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, defaultRoomMessagesLimit, fixedRoomMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  -100,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.RoomMessages, defaultRoomMessagesLimit)

	// 模拟太大的limit
	messages = getRoomMessages(maxRoomMessagesLimit)

	finder.EXPECT().FindsByRecipientID(gomock.Any(), fixedUserID, maxRoomMessagesLimit, fixedRoomMessageBaseID).Return(messages, nil)

	output, err = useCase.Execute(nil, &application.ListRoomMessagesInput{
		UserID: fixedUserID,
		BaseID: fixedRoomMessageBaseID,
		Limit:  1000,
	})
	require.NoError(t, err)
	require.NotNil(t, output)

	require.Len(t, output.RoomMessages, maxRoomMessagesLimit)
}

func getRoomMessages(count int) []*domain.RoomMessage {
	messages := make([]*domain.RoomMessage, count)
	states := make(map[kernel.UserID]domain.MessageState, count+1)

	for i := 0; i < count; i++ {
		states[kernel.UserID("user-"+kernel.UserID(strconv.Itoa(i)).String())] = domain.MessageStateRead
	}
	states[fixedUserID] = domain.MessageStateDelivered

	for i := 0; i < count; i++ {
		messages[i] = domain.LoadRoomMessage(
			fixedMessageID,
			fixedFriendID,
			states,
			fixedRoomID,
			fixedContent,
			time.Now().UTC(),
		)
	}
	return messages
}
