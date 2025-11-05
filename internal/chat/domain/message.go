package domain

import (
	"gochat/internal/shared/kernel"
)

type MessageID kernel.ID

func (i MessageID) String() string {
	return string(i)
}

type State string

const (
	MessageStateCreated   State = "created"
	MessageStateDelivered State = "delivered"
	//MessageStateRead      State = "read"
)

func (s State) String() string {
	return string(s)
}
