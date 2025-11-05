package converter

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
)

type MessageConverter struct {
}

func (c *MessageConverter) PrivateMessageToState(message *domain.PrivateMessage) *model.MessageState {
	return &model.MessageState{
		MessageID: message.ID(),
		Recipient: message.Recipient(),
		State:     message.State(),
	}
}

func (c *MessageConverter) ToPrivateMessage(state *model.MessageState, information *model.MessageInformation) *domain.PrivateMessage {
	return domain.NewPrivateMessage(information.ID, information.Sender, information.Content, information.SentAt, state.State, state.Recipient)
}
