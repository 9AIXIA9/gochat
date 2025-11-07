package repository

import (
	"context"
	"gochat/internal/chat/application"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"

	"gorm.io/gorm"
)

var _ application.UserFinder = (*UserRepository)(nil)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number domain.UserNumber) (*domain.User, error) {
	var user model.User

	err := repo.db.WithContext(ctx).
		Preload("MessageLinks").
		Preload("MessageLinks.Message").
		First(&user, "number = ?", number).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	domainMessages := make([]*domain.Message, 0, len(user.MessageLinks))
	for _, link := range user.MessageLinks {
		if link.Message == nil {
			continue
		}
		m := link.Message
		domainMessages = append(domainMessages, domain.NewMessage(
			m.ID,
			link.State, // 状态来自关联表
			m.SenderID,
			m.Content,
			m.CreatedAt,
		))
	}

	return domain.NewUser(user.ID, user.Number, domainMessages), nil
}
