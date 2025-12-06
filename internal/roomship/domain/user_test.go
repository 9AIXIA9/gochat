package domain_test

import (
	"gochat/internal/roomship/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	u := domain.LoadUser(fixedUserID)
	require.NotNil(t, u)
	require.Equal(t, fixedUserID, u.ID())
}
