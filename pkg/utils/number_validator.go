package utils

import (
	"gochat/internal/shared/errors"
)

func ValidateNumber(number string) error {
	for _, digit := range number {
		if digit < '0' || digit > '9' {
			return errors.ErrInvalidNumber
		}
	}
	return nil
}
