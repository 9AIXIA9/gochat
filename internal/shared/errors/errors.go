package errors

import (
	"errors"
)

// 系统相关错误

var (
	ErrEmptyPointer  = errors.New("system: empty pointer")
	ErrHasBeenClosed = errors.New("has been closed")
)

// 业务逻辑相关错误
var (
	ErrInvalidCredential = errors.New("logic: invalid credential")
	ErrInvalidLength     = errors.New("logic: invalid length")
	ErrInvalidNumber     = errors.New("logic: invalid number")
	ErrInvalidFormat     = errors.New("logic: invalid format")
	ErrEmptyInput        = errors.New("logic: input is empty")
	ErrExpired           = errors.New("logic: expired")
	ErrExceedMaxValue    = errors.New("logic: exceed max value")
	//ErrTimeout           = errors.New("logic:timeout")
)

// 数据库相关错误
var (
	ErrDuplicatedKey      = errors.New("database: duplicated key not allowed")
	ErrNotFound           = errors.New("database: entity is not found")
	ErrForeignKeyViolated = errors.New("database: violates foreign key constraint")
)
