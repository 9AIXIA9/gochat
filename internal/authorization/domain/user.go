package domain

import (
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

type User struct {
	id                kernel.UserID
	email             kernel.Email
	number            kernel.UserNumber
	passwordEncrypted string
	lastLoggedInAt    time.Time // UTC
	signedUpAt        time.Time // UTC
	eventManager      *event.Manager
}

func NewUser(id kernel.UserID, email kernel.Email, number kernel.UserNumber, passwordEncrypted string, lastLoggedInAt time.Time, signedUpAt time.Time) *User {
	return &User{
		id:                id,
		email:             email,
		number:            number,
		passwordEncrypted: passwordEncrypted,
		lastLoggedInAt:    lastLoggedInAt,
		signedUpAt:        signedUpAt,
		eventManager:      event.NewEventManager(),
	}
}

func (u *User) SignUp(eventIDGenerator event.IDGenerator) error {
	e, err := NewUserCreatedEvent(eventIDGenerator.Generate(), u.id)
	if err != nil {
		return err
	}
	u.eventManager.RecordEvent(e)
	return nil
}

func (u *User) Login() {
	u.lastLoggedInAt = time.Now().UTC()
}

// getter

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) SignedUpAt() time.Time {
	return u.signedUpAt
}

func (u *User) LastLoggedInAt() time.Time {
	return u.lastLoggedInAt
}

func (u *User) PasswordEncrypted() string {
	return u.passwordEncrypted
}
func (u *User) Email() kernel.Email {
	return u.email
}

func (u *User) Number() kernel.UserNumber {
	return u.number
}

func (u *User) GetEvents() []event.Event {
	return u.eventManager.GetEvents()
}
