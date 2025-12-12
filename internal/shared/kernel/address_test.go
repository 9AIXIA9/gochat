package kernel_test

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

const fixedMaxAddressLen = 100

func TestAddress_Validate(t *testing.T) {
	// 正常情况
	err := kernel.Address("valid_address").Validate()
	require.NoError(t, err)

	// 空地址
	err = kernel.Address("").Validate()
	require.ErrorIs(t, err, errors.ErrEmptyInput)

	// 超过最大长度
	longAddress := kernel.Address("a" + string(make([]byte, fixedMaxAddressLen)))
	err = longAddress.Validate()
	require.ErrorIs(t, err, errors.ErrInvalidLength)
}
