package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendRoomMessageInput_Validate(t *testing.T) {
	input := &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptySenderID := &application.SendRoomMessageInput{
		SenderID: "",
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	}

	err = inputWithEmptySenderID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRecipientID := &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   "",
		Content:  fixedContent,
	}

	err = inputWithEmptyRecipientID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewSendRoomMessageUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockExister := mocks.NewMockRoomshipExisterByUserIDAndRoomID(ctrl)
	mockMessageCreator := mocks.NewMockRoomMessageCreator(ctrl)

	useCase, err := application.NewSendRoomMessageUseCase(
		mockMessageIDGenerator,
		mockEventIDGenerator,
		mockExister,
		mockMessageCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSendRoomMessageUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSendRoomMessageUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockExister := mocks.NewMockRoomshipExisterByUserIDAndRoomID(ctrl)
	mockMessageCreator := mocks.NewMockRoomMessageCreator(ctrl)

	useCase, err := application.NewSendRoomMessageUseCase(
		mockMessageIDGenerator,
		mockEventIDGenerator,
		mockExister,
		mockMessageCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockExister.EXPECT().ExistByUserIDAndRoomID(nil, fixedRoomID, fixedUserID).Return(true, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.NoError(t, err)

	// 不是成员
	gomock.InOrder(
		mockExister.EXPECT().ExistByUserIDAndRoomID(nil, fixedRoomID, fixedUserID).Return(false, nil).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SendRoomMessageInput{
		SenderID: fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
	})
	require.ErrorIs(t, err, domain.ErrNotMember)
}
