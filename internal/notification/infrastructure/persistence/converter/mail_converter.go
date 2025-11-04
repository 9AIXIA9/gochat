package converter

import (
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
)

var _ gormutils.GenericModelConverter[*model.Mail, *domain.Mail] = (*MailConverter)(nil)

type MailConverter struct {
}

func (c *MailConverter) ToModel(mail *domain.Mail) *model.Mail {
	return &model.Mail{
		ID:        mail.ID(),
		Theme:     mail.Theme(),
		Recipient: mail.Recipient(),
		Title:     mail.Title(),
		Content:   mail.Content(),
		Email:     mail.Email(),
	}
}

func (c *MailConverter) ToDomain(mail *model.Mail) *domain.Mail {
	return domain.NewMail(mail.ID, mail.Recipient, mail.Theme, mail.Title, mail.Content, mail.Email)
}
