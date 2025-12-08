package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRefuseMemberRequestInput_Validate(t *testing.T) {
	input := &application.RefuseMemberRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.RefuseMemberRequestInput{
		UserID:    "",
		RequestID: fixedOperationID,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRequestID := &application.RefuseMemberRequestInput{
		UserID:    fixedUserID,
		RequestID: "",
	}

	err = inputWithEmptyRequestID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)
}

func TestNewRefuseMemberRequestUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)
	mockMemberRequestFinderByRequestID := mocks.NewMockMemberRequestFinderByID(ctrl)
	mockMemberRequestUpdater := mocks.NewMockMemberRequestUpdater(ctrl)

	useCase, err := application.NewRefuseMemberRequestUseCase(
		mockMemberRequestFinderByRequestID,
		mockRoomshipFinder,
		mockMemberRequestUpdater,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewRefuseMemberRequestUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestRefuseMemberRequestUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)
	mockMemberRequestFinderByRequestID := mocks.NewMockMemberRequestFinderByID(ctrl)
	mockMemberRequestUpdater := mocks.NewMockMemberRequestUpdater(ctrl)

	useCase, err := application.NewRefuseMemberRequestUseCase(
		mockMemberRequestFinderByRequestID,
		mockRoomshipFinder,
		mockMemberRequestUpdater,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	mockRequest := domain.LoadMemberRequest(
		fixedOperationID,
		domain.StatePending,
		fixedApplicantID,
		fixedRoomID,
		fixedContent,
		"",
		time.Time{},
		time.Now().UTC(),
	)

	mockRoomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.OwnerRole,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockMemberRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequest, nil),
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(mockRoomship, nil),
		mockMemberRequestUpdater.EXPECT().Update(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.RefuseMemberRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.StateRefused, mockRequest.State())
	assert.Equal(t, fixedUserID, mockRequest.OperatorID())

	// 权限不足
	mockRoomshipNotOwner := domain.LoadRoomship(
		fixedRoomshipID,
		fixedRoomID,
		fixedUserID,
		domain.MemberRole,
		time.Now().UTC(),
	)

	gomock.InOrder(
		mockMemberRequestFinderByRequestID.EXPECT().FindByID(nil, fixedOperationID).Return(mockRequest, nil),
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(mockRoomshipNotOwner, nil),
	)

	_, err = useCase.Execute(nil, &application.RefuseMemberRequestInput{
		UserID:    fixedUserID,
		RequestID: fixedOperationID,
	})
	require.ErrorIs(t, err, domain.ErrNotAdmin)
}
