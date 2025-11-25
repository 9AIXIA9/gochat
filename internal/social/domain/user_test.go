package domain_test

import (
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	id := kernel.UserID("user-xyz")
	u := domain.CreateUser(id)
	require.NotNil(t, u)
	require.Equal(t, id, u.ID())
}
