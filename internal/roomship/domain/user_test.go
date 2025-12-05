package domain_test

import (
	"gochat/internal/roomship/domain"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	id := kernel.UserID("user-xyz")
	u := domain.LoadUser(id)
	require.NotNil(t, u)
	require.Equal(t, id, u.ID())
}
