package domain_test

import (
	"gochat/internal/chat/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatUser_CreateUser(t *testing.T) {
	user := domain.CreateUser(fixedUserID)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
}

func TestChatUser_LoadUser(t *testing.T) {
	user := domain.LoadUser(fixedUserID)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
}
