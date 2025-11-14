package gomail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func NewDialer(config *EmailNotifierConfig) (*gomail.Dialer, error) {
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	closer, err := d.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to dial email server: %w", err)
	}
	if err = closer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close email server connection: %w", err)
	}
	return d, nil
}
