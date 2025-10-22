package event

import (
	"encoding/json"
	"gochat/internal/shared/kernel"
	"time"
)

var _ Event = (*UserLoggedInEvent)(nil)

type UserLoggedInEvent struct {
	userID kernel.UserID
	StandardEvent
}

func NewUserLoggedInEvent(id ID, occurredAt time.Time, userID kernel.UserID) *UserLoggedInEvent {
	return &UserLoggedInEvent{userID: userID, StandardEvent: NewStandardEvent(id, occurredAt, UserLoggedIn)}
}

func CreateUserLoggedInEvent(userID kernel.UserID, generator kernel.IDGenerator) *UserLoggedInEvent {
	return &UserLoggedInEvent{userID: userID, StandardEvent: CreateStandardEvent(UserLoggedIn, generator)}
}

func (e *UserLoggedInEvent) UserID() kernel.UserID {
	return e.userID
}

func (e *UserLoggedInEvent) Marshal() ([]byte, error) {
	type Alias struct {
		UserID kernel.UserID
	}
	return json.Marshal(Alias{
		UserID: e.userID,
	})
}

func (e *UserLoggedInEvent) Unmarshal(data []byte) error {
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
