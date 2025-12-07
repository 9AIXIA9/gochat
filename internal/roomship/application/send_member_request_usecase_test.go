package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	kernelmocks "gochat/internal/shared/kernel/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendMemberRequestInput_Validate(t *testing.T) {
	input := &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: fixedPassword,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.SendMemberRequestInput{
		UserID:   "",
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: fixedPassword,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRoomID := &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   "",
		Content:  fixedContent,
		Password: fixedPassword,
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyPassword := &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: "",
	}

	err = inputWithEmptyPassword.Validate()
	require.NoError(t, err)

	inputWithEmptyContent := &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  "",
		Password: fixedPassword,
	}

	err = inputWithEmptyContent.Validate()
	require.NoError(t, err)
}

func TestNewSendMemberRequestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRequestExister := mocks.NewMockMemberRequestExisterByUserIDAndRoomIDAndState(ctrl)
	mockRoomshipExister := mocks.NewMockRoomshipExisterByUserIDAndRoomID(ctrl)
	mockFinder := mocks.NewMockRoomFinderByID(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockOperationIDGenerator := kernelmocks.NewMockOperationIDGenerator(ctrl)
	mockComparator := mocks.NewMockComparator(ctrl)
	mockCreator := mocks.NewMockMemberRequestCreator(ctrl)

	useCase, err := application.NewSendMemberRequestUseCase(
		mockMemberRequestExister,
		mockRoomshipExister,
		mockFinder,
		mockEventIDGenerator,
		mockOperationIDGenerator,
		mockComparator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSendMemberRequestUseCase(
		nil, nil, nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSendMemberRequestUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemberRequestExister := mocks.NewMockMemberRequestExisterByUserIDAndRoomIDAndState(ctrl)
	mockRoomshipExister := mocks.NewMockRoomshipExisterByUserIDAndRoomID(ctrl)
	mockFinder := mocks.NewMockRoomFinderByID(ctrl)
	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockOperationIDGenerator := kernelmocks.NewMockOperationIDGenerator(ctrl)
	mockComparator := mocks.NewMockComparator(ctrl)
	mockCreator := mocks.NewMockMemberRequestCreator(ctrl)

	useCase, err := application.NewSendMemberRequestUseCase(
		mockMemberRequestExister,
		mockRoomshipExister,
		mockFinder,
		mockEventIDGenerator,
		mockOperationIDGenerator,
		mockComparator,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	//正常情况
	mockRoom := domain.LoadRoom(
		fixedRoomID,
		fixedOperatorID,
		fixedRoomNumber,
		fixedPasswordEncrypted,
		fixedMaxMemberCount,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockMemberRequestExister.EXPECT().ExistByUserIDAndRoomIDAndState(nil, fixedUserID, fixedRoomID, domain.StatePending).Return(false, nil),
		mockRoomshipExister.EXPECT().ExistByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(false, nil),
		mockFinder.EXPECT().FindByID(nil, fixedRoomID).Return(mockRoom, nil),
		mockComparator.EXPECT().Compare(fixedPasswordEncrypted.String(), fixedPassword.String()).Return(nil),
		mockOperationIDGenerator.EXPECT().Generate().Return(fixedOperationID),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: fixedPassword,
	})

	require.NoError(t, err)

	//已存在申请
	gomock.InOrder(
		mockMemberRequestExister.EXPECT().ExistByUserIDAndRoomIDAndState(nil, fixedUserID, fixedRoomID, domain.StatePending).Return(true, nil),
	)

	_, err = useCase.Execute(nil, &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: fixedPassword,
	})

	require.ErrorIs(t, err, domain.ErrMemberRequestAlreadyExists)

	//已是成员
	gomock.InOrder(
		mockMemberRequestExister.EXPECT().ExistByUserIDAndRoomIDAndState(nil, fixedUserID, fixedRoomID, domain.StatePending).Return(false, nil),
		mockRoomshipExister.EXPECT().ExistByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(true, nil),
	)

	_, err = useCase.Execute(nil, &application.SendMemberRequestInput{
		UserID:   fixedUserID,
		RoomID:   fixedRoomID,
		Content:  fixedContent,
		Password: fixedPassword,
	})

	require.ErrorIs(t, err, domain.ErrMemberRequestAlreadyExists)
}
