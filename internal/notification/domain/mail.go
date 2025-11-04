package domain

import (
	"gochat/internal/shared/kernel"
)

type MailID kernel.ID

func (i MailID) String() string {
	return kernel.ID(i).String()
}

type MailTheme string

type Mail struct {
	id        MailID
	theme     MailTheme
	recipient kernel.UserID
	title     string
	content   string
	email     kernel.Email
}

func NewMail(
	id MailID,
	recipient kernel.UserID,
	theme MailTheme,
	title string,
	content string,
	email kernel.Email,
) *Mail {
	return &Mail{
		id:        id,
		theme:     theme,
		recipient: recipient,
		title:     title,
		content:   content,
		email:     email,
	}
}

func (n *Mail) ID() MailID {
	return n.id
}

func (n *Mail) Recipient() kernel.UserID {
	return n.recipient
}

func (n *Mail) Theme() MailTheme {
	return n.theme
}

func (n *Mail) Title() string {
	return n.title
}

func (n *Mail) Content() string {
	return n.content
}
func (n *Mail) Email() kernel.Email {
	return n.email
}
