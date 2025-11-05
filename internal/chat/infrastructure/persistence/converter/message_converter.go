package converter

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"

	"gorm.io/gorm"
)

var _ gormutils.GenericModelConverter[*model.Message, *domain.Message] = (*MessageConverter)(nil)

type MessageConverter struct {
}

func (c *MessageConverter) ToModel(message *domain.Message) *model.Message {
	domainStates := message.States()

	states := make([]*model.RecipientMessageState, 0, len(domainStates))
	for _, domainState := range domainStates {
		states = append(states, &model.RecipientMessageState{
			MessageID: message.ID(),
			Recipient: domainState.Recipient(),
			State:     domainState.State(),
		})
	}
	return &model.Message{
		Model:   gorm.Model{},
		ID:      message.ID(),
		Sender:  message.Sender(),
		Content: message.Content(),
		SentAt:  message.SentAt(),
		States:  states,
	}
}

func (c *MessageConverter) ToDomain(message *model.Message) *domain.Message {
	domainStates := make([]*domain.RecipientMessageState, 0, len(message.States))
	for _, state := range message.States {
		domainStates = append(domainStates, domain.NewRecipientMessageState(
			state.Recipient,
			state.State,
		))
	}

	return domain.NewMessage(
		message.ID,
		message.Sender,
		message.Content,
		message.SentAt,
		domainStates,
	)
}
