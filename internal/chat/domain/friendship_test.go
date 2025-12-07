package domain_test

import (
	"gochat/internal/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFriendship(t *testing.T) {
	friendship := domain.LoadFriendship(
		fixedFriendshipID,
		fixedUserID,
		fixedFriendID,
	)

	require.NotNil(t, friendship)
	assert.Equal(t, fixedFriendshipID, friendship.ID())
	assert.Equal(t, fixedUserID, friendship.UserID1())
	assert.Equal(t, fixedFriendID, friendship.UserID2())
}
