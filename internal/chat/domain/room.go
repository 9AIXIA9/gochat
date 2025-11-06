package domain

import (
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"time"
)

//TODO 把通过id删除换成通过Number

type RoomID kernel.ID

func (r RoomID) String() string {
	return string(r)
}

type Room struct {
	id       RoomID
	members  []kernel.UserID
	messages []*Message

	eventManager *event.Manager
}

func NewRoom(id RoomID, members []kernel.UserID, messages []*Message) *Room {
	return &Room{
		id:           id,
		members:      members,
		messages:     messages,
		eventManager: event.NewEventManager(),
	}
}

func (r *Room) SendMessage(id MessageID, sender kernel.UserID, content string, generator event.IDGenerator) error {
	if len(r.members) <= 1 {
		return myErrors.ErrInvalidLength
	}

	var ok bool
	states := make([]*RecipientMessageState, 0, len(r.members)-1)
	for _, member := range r.members {
		if member != sender {
			state := NewRecipientMessageState(member, MessageStateCreated)
			states = append(states, state)
		} else {
			ok = true
		}
	}

	if !ok {
		return myErrors.ErrNotBelongTo
	}

	message := NewMessage(id, sender, content, time.Now().UTC(), states)

	r.messages = append(r.messages, message)
	ev, err := NewMessageCreatedEvent(
		generator.Generate(),
		message.id,
		message.sender,
		message.Recipients(),
		message.content,
		message.sentAt,
	)
	if err != nil {
		return err
	}
	r.eventManager.RecordEvent(ev)
	return nil
}

func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) Members() []kernel.UserID {
	return r.members
}

func (r *Room) Messages() []*Message {
	return r.messages
}

func (r *Room) GetEvents() []event.Event {
	return r.eventManager.GetEvents()
}
