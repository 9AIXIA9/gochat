package utils

import myErrors "gochat/internal/shared/errors"

func CheckInterfaces(interfaces ...interface{}) error {
	for _, i := range interfaces {
		if i == nil {
			return myErrors.ErrEmptyPointer
		}
	}
	return nil
}
