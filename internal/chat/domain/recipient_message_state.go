package domain

import "gochat/internal/shared/kernel"

type RecipientMessageState struct {
	recipient kernel.UserID
	state     MessageState
}

func NewRecipientMessageState(recipient kernel.UserID, state MessageState) *RecipientMessageState {
	return &RecipientMessageState{
		recipient: recipient,
		state:     state,
	}
}

func (m *RecipientMessageState) Recipient() kernel.UserID {
	return m.recipient
}

func (m *RecipientMessageState) State() MessageState {
	return m.state
}
