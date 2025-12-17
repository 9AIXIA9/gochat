package application_test

import (
	"gochat/internal/profile/application"
	"gochat/internal/profile/domain"
	"gochat/internal/profile/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateRoomProfileInput_Validate(t *testing.T) {
	input := &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.UpdateRoomProfileInput{
		UserID:       "",
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyRoomID := &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       "",
		Name:         fixedName,
		Introduction: fixedIntroduction,
	}

	err = inputWithEmptyRoomID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithOnlyRoomIDAndUserID := &application.UpdateRoomProfileInput{
		UserID: fixedUserID,
		RoomID: fixedRoomID,
	}

	err = inputWithOnlyRoomIDAndUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

}

func TestNewUpdateRoomProfileUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockRoomProfileUpdater(ctrl)
	mockProfileFinder := mocks.NewMockRoomProfileFinder(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)

	useCase, err := application.NewUpdateRoomProfileUseCase(
		mockUpdater,
		mockProfileFinder,
		mockRoomshipFinder,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewUpdateRoomProfileUseCase(
		nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestUpdateRoomProfileUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUpdater := mocks.NewMockRoomProfileUpdater(ctrl)
	mockProfileFinder := mocks.NewMockRoomProfileFinder(ctrl)
	mockRoomshipFinder := mocks.NewMockRoomshipFinderByUserIDAndRoomID(ctrl)

	useCase, err := application.NewUpdateRoomProfileUseCase(
		mockUpdater,
		mockProfileFinder,
		mockRoomshipFinder,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	// 正常情况
	gomock.InOrder(
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(domain.LoadRoomship(
			fixedRoomshipID,
			fixedUserID,
			fixedRoomID,
			domain.OwnerRole,
		), nil).Times(1),
		mockProfileFinder.EXPECT().FindByID(nil, fixedRoomID).Return(domain.LoadRoomProfile(
			fixedRoomID,
			fixedName,
			fixedSign,
			time.Now().UTC(),
		), nil).Times(1),
		mockUpdater.EXPECT().Update(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	})
	require.NoError(t, err)

	//权限不足
	gomock.InOrder(
		mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(domain.LoadRoomship(
			fixedRoomshipID,
			fixedUserID,
			fixedRoomID,
			domain.MemberRole,
		), nil),
	)

	_, err = useCase.Execute(nil, &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	})
	require.ErrorIs(t, err, domain.ErrNoPermission)

	//不在该房间
	mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(nil, myErrors.ErrNotFound).Times(1)
	_, err = useCase.Execute(nil, &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	})
	require.ErrorIs(t, err, domain.ErrNoPermission)

	// 不存在该档案
	mockRoomshipFinder.EXPECT().FindByUserIDAndRoomID(nil, fixedUserID, fixedRoomID).Return(domain.LoadRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
		domain.OwnerRole,
	), nil).Times(1)
	mockProfileFinder.EXPECT().FindByID(nil, fixedRoomID).Return(nil, myErrors.ErrNotFound).Times(1)

	_, err = useCase.Execute(nil, &application.UpdateRoomProfileInput{
		UserID:       fixedUserID,
		RoomID:       fixedRoomID,
		Name:         fixedName,
		Introduction: fixedIntroduction,
	})
	require.ErrorContains(t, err, "room profile not found")
}
