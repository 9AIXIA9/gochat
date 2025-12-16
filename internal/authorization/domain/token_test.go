package domain_test

import (
	"gochat/internal/authorization/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessToken_Validate(t *testing.T) {
	// Test valid token
	validToken := domain.AccessToken("valid-token")
	err := validToken.Validate()
	require.NoError(t, err)

	// Test invalid token(empty string)
	invalidToken := domain.AccessToken("")
	err = invalidToken.Validate()
	require.Error(t, err)
}

func TestRefreshToken_Validate(t *testing.T) {
	// Test valid token
	validToken := domain.RefreshToken("valid-token")
	err := validToken.Validate()
	require.NoError(t, err)

	// Test invalid token(empty string)
	invalidToken := domain.RefreshToken("")
	err = invalidToken.Validate()
	require.Error(t, err)
}
