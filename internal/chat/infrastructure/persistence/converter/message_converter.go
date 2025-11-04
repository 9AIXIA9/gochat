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
	return &model.Message{
		Model:     gorm.Model{},
		ID:        message.ID(),
		Type:      message.Type(),
		Recipient: message.Recipient(),
		State:     message.State(),
		Sender:    message.Sender(),
		Content:   message.Content(),
		SentAt:    message.SentAt(),
	}
}

func (c *MessageConverter) ToDomain(message *model.Message) *domain.Message {
	return domain.NewMessage(
		message.ID,
		message.Type,
		message.Recipient,
		message.State,
		message.Sender,
		message.Content,
		message.SentAt,
	)
}
