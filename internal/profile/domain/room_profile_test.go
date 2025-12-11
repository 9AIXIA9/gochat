package domain_test

import (
	"gochat/internal/profile/domain"
	"gochat/internal/profile/domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadRoomProfile(t *testing.T) {
	profile := domain.LoadRoomProfile(
		fixedProfileID,
		fixedName,
		fixedIntroduction,
		time.Now().UTC(),
	)

	require.NotNil(t, profile)
	require.Equal(t, fixedProfileID, profile.ID())
	require.Equal(t, fixedName, profile.Name())
	require.Equal(t, fixedIntroduction, profile.Introduction())
}

func TestCreateRoomProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDGenerator := mocks.NewMockProfileIDGenerator(ctrl)
	mockIDGenerator.EXPECT().
		Generate().
		Return(fixedProfileID).
		Times(1)

	profile := domain.CreateRoomProfile(
		fixedName,
		fixedIntroduction,
		mockIDGenerator,
	)

	require.NotNil(t, profile)
	require.Equal(t, fixedProfileID, profile.ID())
	require.Equal(t, fixedName, profile.Name())
	require.Equal(t, fixedIntroduction, profile.Introduction())
}
