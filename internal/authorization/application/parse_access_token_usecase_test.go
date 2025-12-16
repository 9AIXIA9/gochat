package application_test

import (
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/domain/mocks"
	myErrors "gochat/internal/shared/errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestParseAccessTokenInput_Validate(t *testing.T) {
	input := &application.ParseAccessTokenInput{
		AccessToken: fixedAccessToken,
	}

	err := input.Validate()
	require.NoError(t, err)

	inputWithEmptyAccessToken := &application.ParseAccessTokenInput{
		AccessToken: "",
	}

	err = inputWithEmptyAccessToken.Validate()
	require.Error(t, err)
}

func TestNewParseAccessTokenUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccessTokenParser := mocks.NewMockAccessTokenParser(ctrl)

	useCase, err := application.NewParseAccessTokenUseCase(
		mockAccessTokenParser,
	)

	require.NoError(t, err)
	require.NotNil(t, useCase)

	useCaseWithNil, err := application.NewParseAccessTokenUseCase(
		nil,
	)

	require.ErrorIs(t, err, myErrors.ErrEmptyPointer)
	require.Nil(t, useCaseWithNil)
}

func TestParseAccessTokenUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccessTokenParser := mocks.NewMockAccessTokenParser(ctrl)

	useCase, err := application.NewParseAccessTokenUseCase(
		mockAccessTokenParser,
	)
	require.NoError(t, err)
	require.NotNil(t, useCase)

	gomock.InOrder(
		mockAccessTokenParser.EXPECT().Parse(fixedAccessToken).Return(fixedUserID, nil).Times(1),
	)

	_, err = useCase.Execute(nil, &application.ParseAccessTokenInput{
		AccessToken: fixedAccessToken,
	})
	require.NoError(t, err)
}
