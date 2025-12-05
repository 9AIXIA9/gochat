package domain_test

import (
	"gochat/internal/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRoom(t *testing.T) {
	room := domain.CreateRoom(fixedRoomID)
	require.NotNil(t, room)
	assert.Equal(t, fixedRoomID, room.ID())
}
