package domain_test

import (
	"gochat/internal/chat/domain"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fixedChatUserID  kernel.UserID     = "chat-user-123"
	fixedChatUserNum kernel.UserNumber = "20001"
)

func TestChatUser_CreateUser(t *testing.T) {
	user := domain.CreateUser(fixedChatUserID, fixedChatUserNum)
	require.NotNil(t, user)
	assert.Equal(t, fixedChatUserID, user.ID())
	assert.Equal(t, fixedChatUserNum, user.Number())
}

func TestChatUser_LoadUser(t *testing.T) {
	user := domain.LoadUser(fixedChatUserID, fixedChatUserNum)
	require.NotNil(t, user)
	assert.Equal(t, fixedChatUserID, user.ID())
	assert.Equal(t, fixedChatUserNum, user.Number())
}
