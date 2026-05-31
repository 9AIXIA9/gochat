package kernel

import (
	"gochat/internal/shared/errors"
	"regexp"
)

type Email string

func (e Email) String() string {
	return string(e)
}

// RFC 5322标准的正则表达式
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validate 方法验证邮箱格式
func (e Email) Validate() error {
	email := e.String()

	if email == "" {
		return errors.ErrInvalidFormat
	}

	// 检查长度 (RFC 3696规定最大254字符)
	if len(email) > 254 {
		return errors.ErrInvalidLength
	}

	// 检查基本结构
	if len(email) < 6 || // 最小长度 a@b.c
		email[0] == '@' || // 不能以@开头
		email[len(email)-1] == '.' || // 不能以点结尾
		email[len(email)-1] == '@' { // 不能以@结尾
		return errors.ErrInvalidFormat
	}

	// 使用正则表达式验证
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidFormat
	}

	return nil
}
