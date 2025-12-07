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
	passwordEncrypted PasswordEncrypted
	signedUpAt        time.Time // UTC
	eventManager      *event.Manager
}

func LoadUser(
	id kernel.UserID,
	email kernel.Email,
	number kernel.UserNumber,
	passwordEncrypted PasswordEncrypted,
	signedUpAt time.Time,
) *User {
	return &User{
		id:                id,
		email:             email,
		number:            number,
		passwordEncrypted: passwordEncrypted,
		signedUpAt:        signedUpAt,
		eventManager:      event.NewEventManager(),
	}
}

func CreateUser(
	email kernel.Email,
	passwordEncrypted PasswordEncrypted,
	userIDGenerator UserIDGenerator,
	numberGenerator UserNumberGenerator,
	eventIDGenerator event.IDGenerator,
) (*User, error) {
	u := &User{
		id:                userIDGenerator.Generate(),
		email:             email,
		number:            numberGenerator.Generate(),
		passwordEncrypted: passwordEncrypted,
		signedUpAt:        time.Now().UTC(),
		eventManager:      event.NewEventManager(),
	}

	ev, err := NewUserCreatedEvent(u.id, eventIDGenerator)
	if err != nil {
		return nil, err
	}

	u.eventManager.RecordEvent(ev)
	return u, nil
}

// getter

func (u *User) ID() kernel.UserID {
	return u.id
}

func (u *User) SignedUpAt() time.Time {
	return u.signedUpAt
}

func (u *User) PasswordEncrypted() PasswordEncrypted {
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
