package validate

import myErrors "gochat/internal/shared/errors"

// NotNil checks for nil values and returns ErrEmptyPointer when any is nil.
func NotNil(values ...interface{}) error {
	for _, v := range values {
		if v == nil {
			return myErrors.ErrEmptyPointer
		}
	}
	return nil
}
