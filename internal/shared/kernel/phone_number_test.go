package kernel_test

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhoneNumber_Validate(t *testing.T) {
	// 正常情况
	err := kernel.PhoneNumber("+86 138-0013-8000").Validate()
	require.NoError(t, err)

	// 长度过短
	err = kernel.PhoneNumber("123").Validate()
	require.Equal(t, myErrors.ErrInvalidLength, err)

	// 长度过长
	longNumber := "1234567890123456" // 16位
	err = kernel.PhoneNumber(longNumber).Validate()
	require.Equal(t, myErrors.ErrInvalidLength, err)

	// 包含非法字符
	err = kernel.PhoneNumber("138-0013-800A").Validate()
	require.Equal(t, myErrors.ErrInvalidFormat, err)
}
