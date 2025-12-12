package domain_test

import (
	"gochat/internal/profile/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoadRoomProfile(t *testing.T) {
	profile := domain.LoadRoomProfile(
		fixedRoomID,
		fixedName,
		fixedIntroduction,
		time.Now().UTC(),
	)

	require.NotNil(t, profile)
	require.Equal(t, fixedRoomID, profile.ID())
	require.Equal(t, fixedName, profile.Name())
	require.Equal(t, fixedIntroduction, profile.Introduction())
}

func TestCreateRoomProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	signedAt := time.Now().UTC()
	profile := domain.CreateRoomProfile(
		fixedRoomID,
		signedAt,
	)

	require.NotNil(t, profile)
	require.Equal(t, fixedRoomID, profile.ID())
	require.Equal(t, signedAt, profile.CreatedAt())
}
