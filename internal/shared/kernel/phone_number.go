package kernel

import (
	"gochat/internal/shared/errors"
	"strings"
	"unicode"
)

const (
	maxPhoneNumberLen = 15
	minPhoneNumberLen = 7
)

type PhoneNumber string

func (p PhoneNumber) String() string {
	return string(p)
}

// Validate 验证电话号码格式
func (p PhoneNumber) Validate() error {
	phone := p.String()

	// 移除所有空格、横线、括号等格式字符
	cleanPhone := strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '(' || r == ')' || r == '+' {
			return -1
		}
		return r
	}, phone)

	// 检查长度 (中国手机号: 11位，包含+86是13位)
	if len(cleanPhone) < minPhoneNumberLen || len(cleanPhone) > maxPhoneNumberLen {
		return errors.ErrInvalidLength
	}

	// 检查是否只包含数字（允许开头的+号在clean时已移除）
	for _, r := range cleanPhone {
		if !unicode.IsDigit(r) {
			return errors.ErrInvalidFormat
		}
	}

	return nil
}
