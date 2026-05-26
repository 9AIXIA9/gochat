package errors

import (
	"errors"
)

// 系统相关错误

var (
	ErrEmptyPointer    = errors.New("system: empty pointer")
	ErrHasBeenClosed   = errors.New("system: has been closed")
	ErrChanIsFull      = errors.New("system: chan is full")
	ErrWrongEventTopic = errors.New("system: wrong event topic")
	ErrTimeout         = errors.New("system: operation timed out")
)

// 业务逻辑相关错误
var (
	ErrInvalidLength = errors.New("logic: invalid length")
	ErrInvalidNumber = errors.New("logic: invalid number")
	ErrInvalidFormat = errors.New("logic: invalid format")
	ErrEmptyInput    = errors.New("logic: input is empty")
)

// 数据库相关错误
var (
	ErrDuplicatedKey      = errors.New("database: duplicated key not allowed")
	ErrNotFound           = errors.New("database: entity is not found")
	ErrForeignKeyViolated = errors.New("database: violates foreign key constraint")
)
