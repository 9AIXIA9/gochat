package application_test

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	eventMock "gochat/internal/shared/event/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSignUpInput_Validate(t *testing.T) {
	input := &application.SignUpInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyPassword := &application.SignUpInput{
		Email:    fixedEmail,
		Password: "",
	}

	err = inputWithEmptyPassword.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidLength)

	inputWithEmptyEmail := &application.SignUpInput{
		Email:    "",
		Password: fixedPassword,
	}

	err = inputWithEmptyEmail.Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidFormat)
}

func TestNewSignUpUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockUserIDGenerator := mocks.NewMockUserIDGenerator(ctrl)
	mockNumberGenerator := mocks.NewMockUserNumberGenerator(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)
	mockUserCreator := mocks.NewMockUserCreator(ctrl)

	useCase, err := application.NewSignUpUseCase(
		mockEventIDGenerator,
		mockUserIDGenerator,
		mockNumberGenerator,
		mockEncryptor,
		mockUserCreator,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewSignUpUseCase(
		nil, nil, nil, nil, nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestSignUpUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventIDGenerator := eventMock.NewMockIDGenerator(ctrl)
	mockUserIDGenerator := mocks.NewMockUserIDGenerator(ctrl)
	mockNumberGenerator := mocks.NewMockUserNumberGenerator(ctrl)
	mockEncryptor := mocks.NewMockEncryptor(ctrl)
	mockUserCreator := mocks.NewMockUserCreator(ctrl)

	useCase, err := application.NewSignUpUseCase(
		mockEventIDGenerator,
		mockUserIDGenerator,
		mockNumberGenerator,
		mockEncryptor,
		mockUserCreator,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockEncryptor.EXPECT().Encrypt(fixedPassword.String()).Return(fixedEncryptedPassword.String(), nil).Times(1),
		mockUserIDGenerator.EXPECT().Generate().Return(fixedUserID).Times(1),
		mockNumberGenerator.EXPECT().Generate().Return(fixedUserNumber).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockUserCreator.EXPECT().Create(nil, gomock.Any()).Return(nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.SignUpInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	})
	require.NoError(t, err)

	// 邮件已注册
	gomock.InOrder(
		mockEncryptor.EXPECT().Encrypt(fixedPassword.String()).Return(fixedEncryptedPassword.String(), nil).Times(1),
		mockUserIDGenerator.EXPECT().Generate().Return(fixedUserID).Times(1),
		mockNumberGenerator.EXPECT().Generate().Return(fixedUserNumber).Times(1),
		mockEventIDGenerator.EXPECT().Generate().Return(fixedEventID).Times(1),
		mockUserCreator.EXPECT().Create(nil, gomock.Any()).Return(myErrors.ErrDuplicatedKey).Times(1),
	)
	_, err = useCase.Execute(nil, &application.SignUpInput{
		Email:    fixedEmail,
		Password: fixedPassword,
	})
	require.ErrorContains(t, err, "email is already used")
}
