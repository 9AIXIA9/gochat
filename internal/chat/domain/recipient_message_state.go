package domain

import "gochat/internal/shared/kernel"

type RecipientMessageState struct {
	recipient kernel.UserID
	state     State
}

func NewRecipientMessageState(recipient kernel.UserID, state State) *RecipientMessageState {
	return &RecipientMessageState{
		recipient: recipient,
		state:     state,
	}
}

func (m *RecipientMessageState) Recipient() kernel.UserID {
	return m.recipient
}

func (m *RecipientMessageState) State() State {
	return m.state
}
