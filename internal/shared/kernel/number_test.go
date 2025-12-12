package kernel_test

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumber_Validate(t *testing.T) {
	// 正常情况
	err := kernel.Number("1234567890").Validate()
	require.NoError(t, err)

	// 空字符串
	err = kernel.Number("").Validate()
	require.Equal(t, myErrors.ErrEmptyInput, err)

	// 包含非数字字符
	err = kernel.Number("1234abc567").Validate()
	require.Equal(t, myErrors.ErrInvalidNumber, err)

	// 包含特殊字符
	err = kernel.Number("1234!@#567").Validate()
	require.Equal(t, myErrors.ErrInvalidNumber, err)
}
