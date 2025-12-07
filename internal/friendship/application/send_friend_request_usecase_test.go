package application_test

import (
	"gochat/internal/friendship/application"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMocks "gochat/internal/shared/event/mocks"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendFriendRequestInput_Validate(t *testing.T) {
	//正常情况
	input := &application.SendFriendRequestInput{
		FromID:  fixedUserID,
		ToID:    fixedToID,
		Content: fixedContent,
	}

	err := input.Validate()
	require.NoError(t, err)

	//无内容情况
	inputWithNoContent := &application.SendFriendRequestInput{
		FromID: fixedUserID,
		ToID:   fixedToID,
	}

	err = inputWithNoContent.Validate()
	require.NoError(t, err)

	//空userID
	inputWithEmptyFromID := &application.SendFriendRequestInput{
		FromID:  "",
		ToID:    fixedToID,
		Content: fixedContent,
	}

	err = inputWithEmptyFromID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	//空ToID
	inputWithEmptyToID := &application.SendFriendRequestInput{
		FromID:  fixedUserID,
		ToID:    "",
		Content: fixedContent,
	}

	err = inputWithEmptyToID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewSendFriendRequestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//正常情况
	mockFriendshipExisterByUserID := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockFriendRequestExisterByUserIDAndState := mocks.NewMockFriendRequestExisterByUserIDAndState(ctrl)
	mockFriendRequestCreator := mocks.NewMockFriendRequestCreator(ctrl)
	mockIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockOperationIDGenerator := kernelmocks.NewMockOperationIDGenerator(ctrl)

	useCase, err := application.NewSendFriendRequestUseCase(
		mockFriendshipExisterByUserID,
		mockFriendRequestExisterByUserIDAndState,
		mockFriendRequestCreator,
		mockIDGenerator,
		mockOperationIDGenerator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	//空指针情况
	useCaseWithNil, err := application.NewSendFriendRequestUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSendFriendRequestUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFriendshipExisterByUserID := mocks.NewMockFriendshipExisterByUserID(ctrl)
	mockFriendRequestExisterByUserIDAndState := mocks.NewMockFriendRequestExisterByUserIDAndState(ctrl)
	mockFriendRequestCreator := mocks.NewMockFriendRequestCreator(ctrl)
	mockIDGenerator := eventMocks.NewMockIDGenerator(ctrl)
	mockOperationIDGenerator := kernelmocks.NewMockOperationIDGenerator(ctrl)

	//正常情况
	useCase, err := application.NewSendFriendRequestUseCase(
		mockFriendshipExisterByUserID,
		mockFriendRequestExisterByUserIDAndState,
		mockFriendRequestCreator,
		mockIDGenerator,
		mockOperationIDGenerator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockFriendshipExisterByUserID.EXPECT().ExistByUserID(nil, fixedUserID, fixedToID).Return(false, nil),
		mockFriendRequestExisterByUserIDAndState.EXPECT().ExistByUserIDAndState(nil, fixedUserID, fixedToID, domain.StatePending).Return(false, nil),
		mockOperationIDGenerator.EXPECT().Generate().Return(fixedOperationID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockFriendRequestCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.SendFriendRequestInput{
		FromID:  fixedUserID,
		ToID:    fixedToID,
		Content: fixedContent,
	})
	require.NoError(t, err)

	//已经是朋友情况
	mockFriendshipExisterByUserID.EXPECT().ExistByUserID(nil, fixedUserID, fixedToID).Return(true, nil)
	_, err = useCase.Execute(nil, &application.SendFriendRequestInput{
		FromID:  fixedUserID,
		ToID:    fixedToID,
		Content: fixedContent,
	})
	require.ErrorIs(t, err, domain.ErrAlreadyBeenFriends)

	//已有待处理的好友请求情况
	mockFriendshipExisterByUserID.EXPECT().ExistByUserID(nil, fixedUserID, fixedToID).Return(false, nil)
	mockFriendRequestExisterByUserIDAndState.EXPECT().ExistByUserIDAndState(nil, fixedUserID, fixedToID, domain.StatePending).Return(true, nil)
	_, err = useCase.Execute(nil, &application.SendFriendRequestInput{
		FromID:  fixedUserID,
		ToID:    fixedToID,
		Content: fixedContent,
	})
	require.ErrorIs(t, err, domain.ErrFriendRequestExists)
}
