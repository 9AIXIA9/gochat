package domain

import myErrors "gochat/internal/shared/errors"

type (
	Token        string
	RefreshToken Token
	AccessToken  Token
)

func (t Token) Validate() error {
	if len(t) == 0 {
		return myErrors.ErrEmptyInput
	}
	return nil
}

func (t Token) String() string {
	return string(t)
}

func (t AccessToken) Validate() error {
	return Token(t).Validate()
}

func (t AccessToken) String() string {
	return Token(t).String()
}

func (t RefreshToken) Validate() error {
	return Token(t).Validate()
}

func (t RefreshToken) String() string {
	return Token(t).String()
}
