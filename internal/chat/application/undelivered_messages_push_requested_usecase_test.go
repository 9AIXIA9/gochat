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

const fixedMessageLen = 10

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
	mockPrivateMessagesUpdater := mocks.NewMockPrivateMessagesUpdater(ctrl)
	mockRoomMessagesFinder := mocks.NewMockRoomMessagesFinderByRecipientIDAndState(ctrl)
	mockRoomMessageNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	mockRoomMessagesUpdater := mocks.NewMockRoomMessagesUpdater(ctrl)

	useCase, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		mockPrivateMessagesFinder,
		mockPrivateMessageNotifier,
		mockPrivateMessagesUpdater,
		mockRoomMessagesFinder,
		mockRoomMessageNotifier,
		mockRoomMessagesUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		nil, nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUndeliveredMessagesPushRequestedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPrivateMessagesFinder := mocks.NewMockPrivateMessagesFinderByRecipientIDAndState(ctrl)
	mockPrivateMessageNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	mockPrivateMessagesUpdater := mocks.NewMockPrivateMessagesUpdater(ctrl)
	mockRoomMessagesFinder := mocks.NewMockRoomMessagesFinderByRecipientIDAndState(ctrl)
	mockRoomMessageNotifier := mocks.NewMockRoomMessageNotifier(ctrl)
	mockRoomMessagesUpdater := mocks.NewMockRoomMessagesUpdater(ctrl)

	useCase, err := application.NewUndeliveredMessagesPushRequestedUseCase(
		mockPrivateMessagesFinder,
		mockPrivateMessageNotifier,
		mockPrivateMessagesUpdater,
		mockRoomMessagesFinder,
		mockRoomMessageNotifier,
		mockRoomMessagesUpdater,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	//正常情况
	mockPrivateMessages := make([]*domain.PrivateMessage, 0, fixedMessageLen)
	mockRoomMessages := make([]*domain.RoomMessage, 0, fixedMessageLen)
	for i := 0; i < fixedMessageLen; i++ {
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
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(mockPrivateMessages, nil).Times(1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(mockRoomMessages, nil).Times(1),
		mockPrivateMessageNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(fixedMessageLen),
		mockPrivateMessagesUpdater.EXPECT().Updates(nil, gomock.Any()).Return(nil).Times(1),
		mockRoomMessageNotifier.EXPECT().Notify(gomock.Any(), []kernel.UserID{fixedUserID}).Return([]kernel.UserID{fixedUserID}, nil).Times(fixedMessageLen),
		mockRoomMessagesUpdater.EXPECT().Updates(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//private messages 为空
	mockRoomMessages = make([]*domain.RoomMessage, 0, fixedMessageLen)
	for i := 0; i < fixedMessageLen; i++ {
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
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(nil, nil).Times(1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(mockRoomMessages, nil).Times(1),
		mockRoomMessageNotifier.EXPECT().Notify(gomock.Any(), []kernel.UserID{fixedUserID}).Return([]kernel.UserID{fixedUserID}, nil).Times(fixedMessageLen),
		mockRoomMessagesUpdater.EXPECT().Updates(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)

	//room messages 为空
	mockPrivateMessages = make([]*domain.PrivateMessage, 0, fixedMessageLen)
	for i := 0; i < fixedMessageLen; i++ {
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
		mockPrivateMessagesFinder.EXPECT().FindPrivateMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(mockPrivateMessages, nil).Times(1),
		mockRoomMessagesFinder.EXPECT().FindRoomMessagesByRecipientIDAndState(nil, fixedUserID, domain.MessageStateUndelivered).Return(nil, nil).Times(1),
		mockPrivateMessageNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(fixedMessageLen),
		mockPrivateMessagesUpdater.EXPECT().Updates(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UndeliveredMessagesPushRequestedInput{
		UserID: fixedUserID,
	})
	require.NoError(t, err)
}
