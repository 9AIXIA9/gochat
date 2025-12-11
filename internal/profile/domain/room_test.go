package domain_test

import (
	"gochat/internal/profile/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadRoom(t *testing.T) {
	room := domain.LoadRoom(fixedRoomID)
	require.NotNil(t, room)
	assert.Equal(t, fixedRoomID, room.ID())
}
