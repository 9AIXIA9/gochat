package kernel

import (
	"fmt"
	"gochat/internal/shared/errors"
	"strings"
	"unicode"
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
	if len(cleanPhone) < 7 || len(cleanPhone) > 15 {
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

// Normalize 标准化电话号码格式 (例如: +86-138-0013-8000)
func (p PhoneNumber) Normalize() string {
	phone := p.String()
	cleanPhone := strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '(' || r == ')' {
			return -1
		}
		return r
	}, phone)

	// 添加国际区号前缀（如果是中国手机号且没有+86前缀）
	if len(cleanPhone) == 11 && cleanPhone[0] == '1' {
		return fmt.Sprintf("+86-%s", formatChinesePhone(cleanPhone))
	}

	return phone
}

// IsMobile 判断是否为手机号码
func (p PhoneNumber) IsMobile() bool {
	phone := p.String()
	cleanPhone := strings.Map(func(r rune) rune {
		if r == ' ' || r == '-' || r == '(' || r == ')' || r == '+' {
			return -1
		}
		return r
	}, phone)

	// 中国手机号: 11位，以1开头
	if len(cleanPhone) == 11 && cleanPhone[0] == '1' {
		// 检查第二位是否为3-9
		second := cleanPhone[1]
		return second >= '3' && second <= '9'
	}

	return false
}

// 格式化中国手机号: 138-0013-8000
func formatChinesePhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return fmt.Sprintf("%s-%s-%s",
		phone[0:3],
		phone[3:7],
		phone[7:11])
}
