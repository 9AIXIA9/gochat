package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/notification/application"
	"gochat/internal/notification/domain"
	"gochat/internal/notification/infrastructure/persistence/model"
	"gochat/internal/shared/kernel"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//TODO 目标：优化数据库操作

var _ application.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SaveMessage(ctx context.Context, recipient kernel.UserID, message *domain.Message) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 持久化消息主体
		msgModel := messageToModel(message)
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(msgModel).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		stateModel := &model.MessageState{
			MessageID: message.ID(),
			UserID:    recipient,
			State:     message.State(),
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(stateModel).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *MessageRepository) FindMessagesByUserID(ctx context.Context, userID kernel.UserID) ([]*domain.Message, error) {
	m := &model.Message{}
	ms := &model.MessageState{}

	// 承载查询结果的轻量结构
	type row struct {
		ID       domain.MessageID    `gorm:"column:id"`
		Content  string              `gorm:"column:content"`
		SenderID kernel.UserID       `gorm:"column:sender_id"`
		SentAt   time.Time           `gorm:"column:created_at"`
		State    domain.MessageState `gorm:"column:state"`
	}

	var rows []row
	err := repo.db.WithContext(ctx).
		Table(m.TableName()+" AS m").
		Select("m.id, m.content, m.sender_id, m.created_at, ms.state").
		Joins("JOIN "+ms.TableName()+" AS ms ON ms.message_id = m.id").
		Where("ms.user_id = ?", userID).
		// 可根据需要调整排序方向
		Order("m.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, gormutils.TranslateError(err)
	}

	result := make([]*domain.Message, 0, len(rows))
	for _, r := range rows {
		result = append(result, domain.NewMessage(
			r.ID,
			r.SenderID,
			r.State,
			r.Content,
			r.SentAt,
		))
	}
	return result, nil
}

func (repo *MessageRepository) UpdateMessageState(ctx context.Context, userID kernel.UserID, messageID domain.MessageID, newState domain.MessageState) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.MessageState{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Update("state", newState).Error)
}

func (repo *MessageRepository) UpdateMessageStates(ctx context.Context, userID kernel.UserID, messageIDs []domain.MessageID, newState domain.MessageState) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Model(&model.MessageState{}).
		Where("user_id = ? AND message_id IN ?", userID, messageIDs).
		Update("state", newState).Error)
}

func messageToModel(message *domain.Message) *model.Message {
	return &model.Message{
		ID:        message.ID(),
		Content:   message.Content(),
		SenderID:  message.Sender(),
		CreatedAt: message.SentAt(),
	}
}
