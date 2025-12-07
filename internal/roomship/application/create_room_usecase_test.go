package application_test

import (
	"gochat/internal/roomship/application"
	"gochat/internal/roomship/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateRoomInput_Validate(t *testing.T) {
	input := &application.CreateRoomInput{
		UserID:         fixedUserID,
		MaxMemberCount: fixedMaxMemberCount,
		Password:       fixedPassword,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyUserID := &application.CreateRoomInput{
		UserID:         "",
		MaxMemberCount: fixedMaxMemberCount,
		Password:       fixedPassword,
	}

	err = inputWithEmptyUserID.Validate()
	require.ErrorIs(t, err, myErrors.ErrEmptyInput)

	inputWithEmptyMaxMemberCount := &application.CreateRoomInput{
		UserID:         fixedUserID,
		MaxMemberCount: 0,
		Password:       fixedPassword,
	}

	err = inputWithEmptyMaxMemberCount.Validate()
	require.NoError(t, err)
	assert.Equal(t, defaultMemberCount, inputWithEmptyMaxMemberCount.MaxMemberCount)

	inputWithEmptyPassword := &application.CreateRoomInput{
		UserID:         fixedUserID,
		MaxMemberCount: fixedMaxMemberCount,
		Password:       "",
	}

	err = inputWithEmptyPassword.Validate()
	require.NoError(t, err)
}

func TestNewCreateRoomUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockRoomIDGenerator := mocks.NewMockRoomIDGenerator(ctrl)
	mockRoomNumberGenerator := mocks.NewMockRoomNumberGenerator(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)
	mockRoomCreator := mocks.NewMockRoomCreator(ctrl)

	useCase, err := application.NewCreateRoomUseCase(
		mockIDGenerator,
		mockRoomIDGenerator,
		mockRoomNumberGenerator,
		mockEncryptor,
		mockRoomCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewCreateRoomUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestCreateRoomUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockRoomIDGenerator := mocks.NewMockRoomIDGenerator(ctrl)
	mockRoomNumberGenerator := mocks.NewMockRoomNumberGenerator(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)
	mockRoomCreator := mocks.NewMockRoomCreator(ctrl)

	useCase, err := application.NewCreateRoomUseCase(
		mockIDGenerator,
		mockRoomIDGenerator,
		mockRoomNumberGenerator,
		mockEncryptor,
		mockRoomCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockEncryptor.EXPECT().Encrypt(fixedPassword.String()).Return(fixedPasswordEncrypted.String(), nil),
		mockRoomIDGenerator.EXPECT().Generate().Return(fixedRoomID),
		mockRoomNumberGenerator.EXPECT().Generate().Return(fixedRoomNumber),
		mockIDGenerator.EXPECT().Generate().Return(fixedEventID),
		mockRoomCreator.EXPECT().Create(nil, gomock.Any()).Return(nil),
	)

	_, err = useCase.Execute(nil, &application.CreateRoomInput{
		UserID:         fixedUserID,
		MaxMemberCount: fixedMaxMemberCount,
		Password:       fixedPassword,
	})

	require.NoError(t, err)
}
