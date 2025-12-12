package domain_test

import (
	"gochat/internal/profile/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadRoomship(t *testing.T) {
	roomship := domain.LoadRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
		domain.OwnerRole,
	)
	require.NotNil(t, roomship)
	assert.Equal(t, fixedRoomshipID, roomship.ID())
	assert.Equal(t, fixedUserID, roomship.UserID())
	assert.Equal(t, fixedRoomID, roomship.RoomID())
	assert.Equal(t, domain.OwnerRole, roomship.Role())
}
