package application_test

import (
	"fmt"
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

const limit = 100

func TestUndeliveredMessagesPushRequestedInput_Validate(t *testing.T) {
	input := &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UndeliveredMessagesPushRequestedInput{
		UserID: "",
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewUndeliveredMessagesPushRequestedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPrivateMessagesFinder := mocks.NewMockPrivateMessagesFinderByRecipientIDAndState(ctrl)
	mockPrivateMessageNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	mockRoomMessagesFinder := mocks.NewMockRoomMessagesFinderByRecipientIDAndState(ctrl)
	mockRoomMessageNotifier := mocks.NewMockRoomMessageNotifier(ctrl)

	useCase, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		mockPrivateMessagesFinder,
		mockPrivateMessageNotifier,
		mockRoomMessagesFinder,
		mockRoomMessageNotifier,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUndeliveredMessagesPushRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPrivateMessagesFinder := mocks.NewMockPrivateMessagesFinderByRecipientIDAndState(ctrl)
	mockPrivateMessageNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	mockRoomMessagesFinder := mocks.NewMockRoomMessagesFinderByRecipientIDAndState(ctrl)
	mockRoomMessageNotifier := mocks.NewMockRoomMessageNotifier(ctrl)

	useCase, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		mockPrivateMessagesFinder,
		mockPrivateMessageNotifier,
		mockRoomMessagesFinder,
		mockRoomMessageNotifier,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	//正常情况
	mockPrivateMessages := make([]*domain.PrivateMessage, 0, limit-1)
	mockRoomMessages := make([]*domain.RoomMessage, 0, limit-1)
	for i := 0; i < limit-1; i++ {
		mockPrivateMessages = append(mockPrivateMessages, domain.LoadPrivateMessage(
			kernel.MessageID(fmt.Sprintf("message-%d", i)),
			kernel.UserID(fmt.Sprintf("sender-%d", i+1)),
			fixedUserID,
			fixedContent,
			domain.MessageStateUndelivered,
			time.Now().UTC(),
		))

		mockStates := make(map[kernel.UserID]domain.MessageState)
		mockStates[fixedUserID] = domain.MessageStateUndelivered
		for i := 0; i < fixedRoomMembersCount; i++ {
			mockStates[kernel.UserID(fmt.Sprintf("member-%d", i+1))] = domain.MessageStateDelivered
		}

		mockRoomMessages = append(mockRoomMessages, domain.LoadRoomMessage(
			kernel.MessageID(fmt.Sprintf("message-%d", i)),
			kernel.UserID(fmt.Sprintf("sender-%d", i+1)),
			mockStates,
			fixedRoomID,
			fixedContent,
			time.Now().UTC(),
		))
	}

	gomock.InOrder(
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockPrivateMessages, nil).Times(1),
		mockPrivateMessageNotifier.EXPECT().Notify(nil, gomock.Any()).Return(nil).Times(limit-1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockRoomMessages, nil).Times(1),
		mockRoomMessageNotifier.EXPECT().Notify(nil, gomock.Any(), []kernel.UserID{fixedUserID}).Return(nil).Times(limit-1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//private messages 为空
	mockRoomMessages = make([]*domain.RoomMessage, 0, limit-1)
	for i := 0; i < limit-1; i++ {
		mockStates := make(map[kernel.UserID]domain.MessageState)
		mockStates[fixedUserID] = domain.MessageStateUndelivered
		for i := 0; i < fixedRoomMembersCount; i++ {
			mockStates[kernel.UserID(fmt.Sprintf("member-%d", i+1))] = domain.MessageStateDelivered
		}

		mockRoomMessages = append(mockRoomMessages, domain.LoadRoomMessage(
			kernel.MessageID(fmt.Sprintf("message-%d", i)),
			kernel.UserID(fmt.Sprintf("sender-%d", i+1)),
			mockStates,
			fixedRoomID,
			fixedContent,
			time.Now().UTC(),
		))
	}

	gomock.InOrder(
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(nil, nil).Times(1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockRoomMessages, nil).Times(1),
		mockRoomMessageNotifier.EXPECT().Notify(nil, gomock.Any(), []kernel.UserID{fixedUserID}).Return(nil).Times(limit-1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//room messages 为空
	mockPrivateMessages = make([]*domain.PrivateMessage, 0, limit-1)
	for i := 0; i < limit-1; i++ {
		mockPrivateMessages = append(mockPrivateMessages, domain.LoadPrivateMessage(
			kernel.MessageID(fmt.Sprintf("message-%d", i)),
			kernel.UserID(fmt.Sprintf("sender-%d", i+1)),
			fixedUserID,
			fixedContent,
			domain.MessageStateUndelivered,
			time.Now().UTC(),
		))
	}

	gomock.InOrder(
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockPrivateMessages, nil).Times(1),
		mockPrivateMessageNotifier.EXPECT().Notify(nil, gomock.Any()).Return(nil).Times(limit-1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(nil, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//测试多次推送的情况
	mockPrivateMessages = make([]*domain.PrivateMessage, 0, limit)
	for i := 0; i < limit; i++ {
		mockPrivateMessages = append(mockPrivateMessages, domain.LoadPrivateMessage(
			kernel.MessageID(fmt.Sprintf("message-%d", i)),
			kernel.UserID(fmt.Sprintf("sender-%d", i+1)),
			fixedUserID,
			fixedContent,
			domain.MessageStateUndelivered,
			time.Now().UTC(),
		))
	}

	gomock.InOrder(
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(mockPrivateMessages, nil).Times(1),
		mockPrivateMessageNotifier.EXPECT().Notify(nil, gomock.Any()).Return(nil).Times(limit),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered, limit).Return(nil, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
