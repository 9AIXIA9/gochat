package application_test

import (
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendPrivateMessageInput_Validate(t *testing.T) {
	input := &application.SendPrivateMessageInput{
		SenderID:    fixedUserID,
		RecipientID: fixedFriendID,
		Content:     fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptySenderID := &application.SendPrivateMessageInput{
		SenderID:    "",
		RecipientID: fixedFriendID,
		Content:     fixedContent,
	}

	err = inputWithEmptySenderID.Validate()
	require.Error(t, err)

	inputWithEmptyRecipientID := &application.SendPrivateMessageInput{
		SenderID:    fixedUserID,
		RecipientID: "",
		Content:     fixedContent,
	}

	err = inputWithEmptyRecipientID.Validate()
	require.Error(t, err)
}

func TestNewSendPrivateMessageUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExister := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	mockMessageCreator := mocks.NewMockPrivateMessageCreator(ctrl)

	useCase, err := application.NewSendPrivateMessageUseCase(
		mockExister,
		mockMessageIDGenerator,
		mockNotifier,
		mockMessageCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSendPrivateMessageUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSendPrivateMessageUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExister := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockMessageIDGenerator := kernelmocks.NewMockMessageIDGenerator(ctrl)
	mockNotifier := mocks.NewMockPrivateMessageNotifier(ctrl)
	mockMessageCreator := mocks.NewMockPrivateMessageCreator(ctrl)

	useCase, err := application.NewSendPrivateMessageUseCase(
		mockExister,
		mockMessageIDGenerator,
		mockNotifier,
		mockMessageCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockExister.EXPECT().ExistByUserID(nil, fixedUserID, fixedFriendID).Return(true, nil).Times(1),
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockNotifier.EXPECT().Notify(gomock.Any()).Return(nil).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SendPrivateMessageInput{
		SenderID:    fixedUserID,
		RecipientID: fixedFriendID,
		Content:     fixedContent,
	})
	require.NoError(t, err)

	// 不是好友
	gomock.InOrder(
		mockExister.EXPECT().ExistByUserID(nil, fixedUserID, fixedFriendID).Return(false, nil).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SendPrivateMessageInput{
		SenderID:    fixedUserID,
		RecipientID: fixedFriendID,
		Content:     fixedContent,
	})
	require.Error(t, err)

	// 发送给自己
	gomock.InOrder(
		mockMessageIDGenerator.EXPECT().Generate().Return(fixedMessageID).Times(1),
		mockMessageCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SendPrivateMessageInput{
		SenderID:    fixedUserID,
		RecipientID: fixedUserID,
		Content:     fixedContent,
	})
	require.NoError(t, err)
}
