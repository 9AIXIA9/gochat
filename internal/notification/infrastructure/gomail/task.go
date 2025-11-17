package gomail

import (
	"fmt"
	"gochat/internal/shared/kernel"

	"gopkg.in/gomail.v2"
)

type taskGenerator struct {
	senderName string
	address    string
}

type task struct {
	message *gomail.Message
	email   kernel.Email
}

func (g *taskGenerator) generateUserCreatedTask(email kernel.Email, number kernel.UserNumber) *task {
	message := gomail.NewMessage()

	// From
	message.SetAddressHeader("From", g.address, g.senderName)

	// To
	message.SetHeader("To", email.String())

	// Subject
	subject := fmt.Sprintf("Welcome to GoChat — Account Created (ID: %s)", number.String())
	message.SetHeader("Subject", subject)

	// Body (HTML + Plain Text Alternative)
	body := fmt.Sprintf(userCreatedEmailFormat, number.String())
	message.SetBody("text/html", body)
	message.AddAlternative("text/plain",
		fmt.Sprintf("Welcome to GoChat!\nYour account has been created.\nUser Number: %s\nIf you did not sign up, ignore this email.", number.String()))

	return &task{
		message: message,
		email:   email,
	}
}
