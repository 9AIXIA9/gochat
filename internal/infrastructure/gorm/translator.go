package gorm

import (
	"errors"
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"strings"

	"gorm.io/gorm"
)

// TranslateError 统一的 DB 错误转换
func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		return errors.Join(myErrors.ErrDuplicatedKey, err)
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return errors.Join(myErrors.ErrForeignKeyViolated, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errors.Join(myErrors.ErrNotFound, err)

	default:
		return fmt.Errorf("gorm operate failed,err:%w", err)
	}
}
