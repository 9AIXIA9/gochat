package domain_test

import (
	"gochat/internal/profile/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	user := domain.LoadUser(fixedUserID)
	require.NotNil(t, user)
	assert.Equal(t, fixedUserID, user.ID())
}
