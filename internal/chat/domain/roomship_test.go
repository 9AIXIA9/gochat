package domain_test

import (
	"gochat/internal/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRoomship(t *testing.T) {
	friendship := domain.CreateRoomship(
		fixedRoomshipID,
		fixedUserID,
		fixedRoomID,
	)

	require.NotNil(t, friendship)
	assert.Equal(t, fixedRoomshipID, friendship.ID())
	assert.Equal(t, fixedUserID, friendship.UserID())
	assert.Equal(t, fixedRoomID, friendship.RoomID())
}
