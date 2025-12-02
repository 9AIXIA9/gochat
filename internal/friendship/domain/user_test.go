package domain_test

import (
	"gochat/internal/friendship/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user := domain.CreateUser(fixedUserID)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
}
