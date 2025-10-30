package converter

import (
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ gormutils.GenericModelConverter[*model.Notice, *domain.Notice] = (*NoticeConverter)(nil)

type NoticeConverter struct {
}

func (c *NoticeConverter) ToModel(notice *domain.Notice) *model.Notice {
	return &model.Notice{
		Model:     gorm.Model{},
		ID:        notice.ID(),
		Theme:     notice.Theme(),
		Recipient: notice.Recipient(),
		Title:     notice.Title(),
		Content:   notice.Content(),
		Contact:   notice.Contact().String(),
	}
}

func (c *NoticeConverter) ToDomain(notice *model.Notice) *domain.Notice {
	return domain.NewNotice(notice.ID, notice.Recipient, notice.Theme, notice.Title, notice.Content, domain.Contact(notice.Contact))
}
