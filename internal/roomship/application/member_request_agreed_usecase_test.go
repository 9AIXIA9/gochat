package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMemberRequestAgreedInput_Validate(t *testing.T) {
	input := &application.MemberRequestAgreedInput{
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyRequestID := &application.MemberRequestAgreedInput{
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewMemberRequestAgreedUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockRoomshipCreator(ctrl)
	mockFinder := mocks.NewMockMemberRequestFinderByID(ctrl)

	useCase, err := application.NewMemberRequestAgreedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockFinder,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewMemberRequestAgreedUseCase(
		nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestMemberRequestAgreedUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipIDGenerator := mocks.NewMockRoomshipIDGenerator(ctrl)
	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockCreator := mocks.NewMockRoomshipCreator(ctrl)
	mockFinder := mocks.NewMockMemberRequestFinderByID(ctrl)

	useCase, err := application.NewMemberRequestAgreedUseCase(
		mockRoomshipIDGenerator,
		mockIDGenerator,
		mockFinder,
		mockCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	mockMemberRequest := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StatePending,
		fixedApplicantID,
		fixedRoomID,
		fixedContent,
		fixedOperatorID,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockMemberRequest, nil),
		mockRoomshipIDGenerator.EXPECT().Generate().Return(fixedRoomshipID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.MemberRequestAgreedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)

	// 已经创建
	gomock.InOrder(
		mockFinder.EXPECT().FindByID(nil, fixedOperationID).Return(mockMemberRequest, nil),
		mockRoomshipIDGenerator.EXPECT().Generate().Return(fixedRoomshipID),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockCreator.EXPECT().Create(nil, gomock.Any()).Return(myErrors.ErrDuplicatedKey),
	)

	_, err = useCase.Execute(nil, &application.MemberRequestAgreedInput{
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
}
