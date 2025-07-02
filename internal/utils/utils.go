package utils

import (
	"errors"
	"gochat/internal/types"
	"gorm.io/gorm"
	"strings"
)

// NormalizeVarArgs 标准化可变参数，根据参数数量返回不同形式的结果
func NormalizeVarArgs[T any](data ...T) any {
	if l := len(data); l == 1 {
		return data[0]
	} else if l > 1 {
		return data
	} else {
		return nil
	}
}

// GetOption 获取可选的参数
func GetOption[T any](opts ...T) (bool, T) {
	if len(opts) == 0 {
		var zero T
		return false, zero
	} else {
		return true, opts[0]
	}
}

func CheckDuplicateKeyError(err error) error {
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "1062") {
		return types.ErrDuplicateKey
	}
	return err
}

func CheckNotFoundError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return types.ErrNotFound
	}
	return err
}

func IsNotFound(err error) bool {
	return errors.Is(err, types.ErrNotFound)
}

func IsDuplicate(err error) bool {
	return errors.Is(err, types.ErrDuplicateKey)
}
