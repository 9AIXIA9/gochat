package domain

import (
	"encoding/json"
	"gochat/internal/shared/contract"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type Notification struct {
	id          kernel.MessageID
	recipientID kernel.UserID
	state       NotificationState
	rawPayload  json.RawMessage

	eventManager *event.Manager
}

func LoadNotification(
	id kernel.MessageID,
	recipientID kernel.UserID,
	state NotificationState,
	rawPayload json.RawMessage,
) *Notification {
	return &Notification{
		id:           id,
		recipientID:  recipientID,
		state:        state,
		rawPayload:   rawPayload,
		eventManager: event.NewEventManager(),
	}
}

func CreateNotification(
	id kernel.MessageID,
	recipientID kernel.UserID,
	rawPayload json.RawMessage,
	eventIDGenerator event.IDGenerator,
) (*Notification, error) {
	if len(recipientID) == 0 {
		return nil, ErrNotificationNoRecipient
	}

	notification := &Notification{
		id:           id,
		recipientID:  recipientID,
		state:        StateUndelivered,
		rawPayload:   rawPayload,
		eventManager: event.NewEventManager(),
	}

	ev, err := contract.NewPushRequestedEvent(id, recipientID, rawPayload, eventIDGenerator)
	if err != nil {
		return nil, err
	}

	notification.eventManager.RecordEvent(ev)
	return notification, nil
}

func (n *Notification) ID() kernel.MessageID {
	return n.id
}

func (n *Notification) RecipientID() kernel.UserID {
	return n.recipientID
}

func (n *Notification) State() NotificationState {
	return n.state
}

func (n *Notification) RawPayload() json.RawMessage {
	return n.rawPayload
}

func (n *Notification) GetEvents() []event.Event {
	return n.eventManager.GetEvents()
}
