package gomail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func NewDialer(config *EmailNotifierConfig) *gomail.Dialer {
	return gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
}

func TestConnection(dialer *gomail.Dialer) error {
	closer, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to dial email server: %w", err)
	}

	if err = closer.Close(); err != nil {
		return fmt.Errorf("failed to close email server connection: %w", err)
	}
	return nil
}
