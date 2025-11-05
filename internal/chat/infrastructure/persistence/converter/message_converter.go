package converter

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
)

type MessageConverter struct {
}

func (c *MessageConverter) ToInformation(id domain.MessageID, information *domain.MessageInformation) *model.MessageInformation {
	return &model.MessageInformation{
		ID:      id,
		Sender:  information.Sender(),
		Content: information.Content(),
		SentAt:  information.SentAt(),
	}
}

func (c *MessageConverter) ToRecipientState(id domain.MessageID, recipientMessageState *domain.RecipientMessageState) *model.RecipientMessageState {
	return &model.RecipientMessageState{
		MessageID: id,
		Recipient: recipientMessageState.Recipient(),
		State:     recipientMessageState.State(),
	}
}

func (c *MessageConverter) ToRecipientStates(id domain.MessageID, recipientMessageStates []*domain.RecipientMessageState) []*model.RecipientMessageState {
	states := make([]*model.RecipientMessageState, 0, len(recipientMessageStates))
	for _, state := range recipientMessageStates {
		states = append(states, c.ToRecipientState(id, state))
	}
	return states
}
