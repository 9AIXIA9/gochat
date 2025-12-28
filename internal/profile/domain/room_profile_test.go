package domain_test

import (
	"gochat/internal/profile/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	maxIntroductionLen = 200
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

func TestRoomProfile_Updates(t *testing.T) {
	profile := domain.LoadRoomProfile(
		fixedRoomID,
		"",
		"",
		time.Now().UTC(),
	)

	err := profile.UpdateName(fixedName)
	require.NoError(t, err)
	assert.Equal(t, fixedName, profile.Name())

	err = profile.UpdateIntroduction(fixedIntroduction)
	require.NoError(t, err)
	assert.Equal(t, fixedIntroduction, profile.Introduction())

	// Test exceeding max lengths
	longName := ""
	for i := 0; i < maxNameLen+1; i++ {
		longName += "a"
	}
	err = profile.UpdateName(longName)
	require.Error(t, err)

	longIntroduction := ""
	for i := 0; i < maxIntroductionLen+1; i++ {
		longIntroduction += "a"
	}
	err = profile.UpdateIntroduction(longIntroduction)
	require.Error(t, err)
}
