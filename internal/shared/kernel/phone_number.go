package kernel

import (
	myErrors "gochat/internal/shared/errors"
	"regexp"
	"strings"
)

type PhoneNumber string

func (p PhoneNumber) String() string {
	return string(p)
}

var phoneRegex = regexp.MustCompile(`^1[3456789]\d{9}$`)

func (p PhoneNumber) Validate() error {
	// 去除首尾空格
	s := strings.TrimSpace(string(p))

	// 检查长度
	if len(s) != 11 {
		return myErrors.ErrInvalidLength
	}

	// 使用正则表达式验证格式
	if !phoneRegex.MatchString(s) {
		return myErrors.ErrInvalidFormat
	}

	return nil
}
