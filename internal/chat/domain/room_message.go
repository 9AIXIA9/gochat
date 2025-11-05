package domain

import (
	"gochat/internal/shared/kernel"
	"time"
)

type RoomMessage struct {
	id     MessageID
	roomID RoomID
	*MessageInformation
	states []*RecipientMessageState
}

func NewRoomMessage(
	id MessageID,
	roomID RoomID,
	sender kernel.UserID,
	content string,
	sentAt time.Time,
	states []*RecipientMessageState,
) *RoomMessage {
	return &RoomMessage{
		id:                 id,
		roomID:             roomID,
		MessageInformation: NewMessageInformation(sender, content, sentAt),
		states:             states,
	}
}

func (r *RoomMessage) ID() MessageID {
	return r.id
}

func (r *RoomMessage) RoomID() RoomID {
	return r.roomID
}

func (r *RoomMessage) States() []*RecipientMessageState {
	return r.states
}

func (r *RoomMessage) Recipients() []kernel.UserID {
	recipients := make([]kernel.UserID, 0, len(r.states))
	for _, state := range r.states {
		recipients = append(recipients, state.recipient)
	}
	return recipients
}
