package domain_test

import (
	"gochat/internal/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFriendship(t *testing.T) {
	friendship := domain.CreateFriendship(
		fixedFriendshipID,
		fixedUserID,
		fixedFriendID,
	)

	require.NotNil(t, friendship)
	assert.Equal(t, fixedFriendshipID, friendship.ID())
	assert.Equal(t, fixedUserID, friendship.UserID1())
	assert.Equal(t, fixedFriendID, friendship.UserID2())
}
