package kernel_test

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGender_Validate(t *testing.T) {
	// 正常情况
	for i := 0; i < 3; i++ {
		err := kernel.Gender(i).Validate()
		require.NoError(t, err)
	}

	// 异常情况
	err := kernel.Gender(99).Validate()
	require.ErrorIs(t, err, myErrors.ErrInvalidFormat)
}
