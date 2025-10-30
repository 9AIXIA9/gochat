package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

type NoticeID kernel.ID

type NoticeTheme string

type Contact string

func (n NoticeTheme) String() string {
	return string(n)
}

func (c Contact) String() string {
	return string(c)
}

func (c Contact) Validate() error {
	if err := kernel.Email(c).Validate(); err == nil {
		return nil
	}
	if err := kernel.PhoneNumber(c).Validate(); err == nil {
		return nil
	}
	return myErrors.ErrInvalidFormat
}

type Notice struct {
	id        NoticeID
	theme     NoticeTheme
	recipient kernel.UserID
	title     string
	content   string
	contact   Contact
}

func NewNotice(
	id NoticeID,
	recipient kernel.UserID,
	theme NoticeTheme,
	title string,
	content string,
	contact Contact,
) *Notice {
	return &Notice{
		id:        id,
		theme:     theme,
		recipient: recipient,
		title:     title,
		content:   content,
		contact:   contact,
	}
}

func (n *Notice) ID() NoticeID {
	return n.id
}

func (n *Notice) Recipient() kernel.UserID {
	return n.recipient
}

func (n *Notice) Theme() NoticeTheme {
	return n.theme
}

func (n *Notice) Title() string {
	return n.title
}

func (n *Notice) Content() string {
	return n.content
}
func (n *Notice) Contact() Contact {
	return n.contact
}
