package kernel

import myErrors "gochat/internal/shared/errors"

type Gender int

const (
	UnknownGender Gender = 0
	MaleGender    Gender = 1
	FemaleGender  Gender = 2
)

func (g Gender) String() string {
	switch g {
	case MaleGender:
		return "male"
	case FemaleGender:
		return "female"
	default:
		return "unknown"
	}
}

func (g Gender) Validate() error {
	switch g {
	case MaleGender, FemaleGender, UnknownGender:
		return nil
	default:
		return myErrors.ErrInvalidFormat
	}
}
