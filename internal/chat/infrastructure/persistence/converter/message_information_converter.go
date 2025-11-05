package converter

import (
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
)

var _ gormutils.GenericModelConverter[*model.MessageInformation, *domain.MessageInformation] = (*MessageInformationConverter)(nil)

type MessageInformationConverter struct {
}

func (c *MessageInformationConverter) ToModel(message *domain.MessageInformation) *model.MessageInformation {
	return &model.MessageInformation{
		ID:      message.ID(),
		Type:    message.Type(),
		Sender:  message.Sender(),
		Content: message.Content(),
		SentAt:  message.SentAt(),
	}
}

func (c *MessageInformationConverter) ToDomain(message *model.MessageInformation) *domain.MessageInformation {
	return domain.NewMessageInformation(
		message.ID,
		message.Type,
		message.Sender,
		message.Content,
		message.SentAt,
	)
}
