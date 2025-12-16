package errors_test

import (
	stderrs "errors"
	"fmt"
	"testing"

	berrors "gochat/internal/shared/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBusiness(t *testing.T) {
	err := berrors.NewBusiness("business error: %s", "invalid input")
	require.Error(t, err)

	var be *berrors.BusinessError
	require.ErrorAs(t, err, &be)
	assert.Equal(t, "business error: invalid input", be.Error())
	assert.Nil(t, stderrs.Unwrap(err))
}

func TestWrapBusiness_Nil(t *testing.T) {
	err := berrors.WrapBusiness(nil, "should be nil")
	assert.NoError(t, err)
}

func TestWrapBusiness_WithCause(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := berrors.WrapBusiness(cause, "wrap: %d", 42)
	require.Error(t, err)

	var be *berrors.BusinessError
	require.ErrorAs(t, err, &be)
	assert.Equal(t, "wrap: 42", be.Error())
	assert.Equal(t, cause, stderrs.Unwrap(err))
}

func TestIsBusinessError_Positive(t *testing.T) {
	err := berrors.NewBusiness("boom")
	require.True(t, berrors.IsBusinessError(err))
}

func TestIsBusinessError_Negative(t *testing.T) {
	require.False(t, berrors.IsBusinessError(fmt.Errorf("not business")))
}
