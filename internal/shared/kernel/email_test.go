package kernel_test

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fixedEmailMaxLen = 254
)

func TestEmail_Validate(t *testing.T) {
	// 正常情况
	err := kernel.Email("abcd@demo.com").Validate()
	require.NoError(t, err)

	// 长度超出限制
	email := kernel.Email("")
	for i := 0; i < fixedEmailMaxLen; i++ {
		email += "a"
	}
	email += "@example.com"
	err = email.Validate()
	require.Equal(t, myErrors.ErrInvalidLength, err)

	// 格式错误情况
	invalidEmails := []string{
		"",                     // 空字符串
		"a@b.c",                // 太短
		"plainaddress",         // 缺少@
		"@missingusername.com", // 缺少用户名
		"username@.com",        // 域名以点开头
		"username@com.",        // 域名以点结尾
		"username@domain..com", // 连续的点
	}

	for _, ie := range invalidEmails {
		err = kernel.Email(ie).Validate()
		require.Equal(t, myErrors.ErrInvalidFormat, err)
	}
}
