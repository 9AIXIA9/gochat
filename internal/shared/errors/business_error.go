package errors

import (
	"errors"
	"fmt"
)

type BusinessError struct {
	Message string
	Cause   error
}

func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) Unwrap() error {
	return e.Cause
}

func NewBusiness(format string, args ...interface{}) error {
	return &BusinessError{Message: fmt.Sprintf(format, args...)}
}

func WrapBusiness(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &BusinessError{
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

func IsBusinessError(err error) bool {
	var e *BusinessError
	if errors.As(err, &e) {
		return true
	}
	return false
}
