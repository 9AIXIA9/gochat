package event

import (
	"encoding/json"
	"gochat/internal/shared/kernel"
	"time"
)

var _ Event = (*UserSignedUpEvent)(nil)

type UserSignedUpEvent struct {
	userID kernel.UserID
	email  kernel.Email
	StandardEvent
}

func NewUserSignedUpEvent(id ID, occurredAt time.Time, userID kernel.UserID, email kernel.Email) *UserSignedUpEvent {
	return &UserSignedUpEvent{
		userID:        userID,
		email:         email,
		StandardEvent: NewStandardEvent(id, occurredAt, UserSignedUp),
	}
}

func CreateUserSignedUpEvent(userID kernel.UserID, email kernel.Email, generator kernel.IDGenerator) *UserSignedUpEvent {
	return &UserSignedUpEvent{
		userID:        userID,
		email:         email,
		StandardEvent: CreateStandardEvent(UserSignedUp, generator),
	}
}

func (e *UserSignedUpEvent) UserID() kernel.UserID {
	return e.userID
}
func (e *UserSignedUpEvent) Email() kernel.Email {
	return e.email
}

func (e *UserSignedUpEvent) Marshal() ([]byte, error) {
	type Alias struct {
		Email  kernel.Email
		UserID kernel.UserID
	}
	return json.Marshal(Alias{
		Email:  e.email,
		UserID: e.userID,
	})
}

func (e *UserSignedUpEvent) Unmarshal(data []byte) error {
	type Alias struct {
		UserID kernel.UserID
	}
	var temp Alias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	e.userID = temp.UserID
	return nil
}
