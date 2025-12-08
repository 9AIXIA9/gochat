package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomMessageRepository = (*RoomMessageRepository)(nil)

type RoomMessageRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomMessageRepository(db *gorm.DB, eventRepo event.Repository) *RoomMessageRepository {
	return &RoomMessageRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomMessageRepository) Create(ctx context.Context, message *domain.RoomMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(repo.toModel(message)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		txCtx := context.WithValue(ctx, "transaction", tx)

		if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, message.GetEvents()); err != nil {
			return gormutils.TranslateError(err)
		}

		return nil
	})
}

func (repo *RoomMessageRepository) FindRoomMessage(ctx context.Context, messageID kernel.MessageID) (*domain.RoomMessage, error) {
	var message model.RoomMessage
	if err := repo.db.WithContext(ctx).
		Preload("States").
		Where("id = ?", messageID).
		First(&message).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&message), nil
}

func (repo *RoomMessageRepository) Updates(ctx context.Context, messages []*domain.RoomMessage) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "transaction", tx)
		for _, message := range messages {
			if err := tx.Model(&model.RoomMessage{}).
				Where("id = ?", message.ID()).
				Updates(map[string]interface{}{
					"content": message.Content(),
				}).Error; err != nil {
				return gormutils.TranslateError(err)
			}

			// Update states
			for userID, state := range message.States() {
				if err := tx.Model(&model.RoomMessageState{}).
					Where("message_id = ? AND user_id = ?", message.ID(), userID).
					Update("state", state).Error; err != nil {
					return gormutils.TranslateError(err)
				}
			}

			if err := repo.eventRepo.CreateUnpublishedEvents(txCtx, message.GetEvents()); err != nil {
				return gormutils.TranslateError(err)
			}
		}
		return nil
	})

}

func (repo *RoomMessageRepository) FindRoomMessagesByRecipientIDAndState(
	ctx context.Context,
	recipientID kernel.UserID,
	state domain.MessageState,
	limit int,
) ([]*domain.RoomMessage, error) {
	var messages []model.RoomMessage
	if err := repo.db.WithContext(ctx).
		Joins("JOIN chat_room_message_states ON chat_room_messages.id = chat_room_message_states.message_id").
		Where("chat_room_message_states.user_id = ? AND chat_room_message_states.state = ?", recipientID, state).
		Preload("States").
		Limit(limit).
		Find(&messages).
		Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(messages), nil
}

func (repo *RoomMessageRepository) UpdatesByUserIDAndRoomID(ctx context.Context, userID kernel.UserID, roomID kernel.RoomID, state domain.MessageState) error {
	// 使用子查询避免在 UPDATE 中直接 JOIN
	sub := repo.db.WithContext(ctx).
		Table("chat_room_messages").
		Select("id").
		Where("room_id = ?", roomID)

	return repo.db.WithContext(ctx).
		Model(&model.RoomMessageState{}).
		Where("chat_room_message_states.user_id = ? AND chat_room_message_states.message_id IN (?)", userID, sub).
		Update("chat_room_message_states.state", state).Error
}

func (repo *RoomMessageRepository) toModel(message *domain.RoomMessage) *model.RoomMessage {
	states := make([]*model.RoomMessageState, 0, len(message.States()))
	for id, state := range message.States() {
		states = append(states, &model.RoomMessageState{
			MessageID: message.ID(),
			UserID:    id,
			State:     state,
		})
	}

	return &model.RoomMessage{
		ID:       message.ID(),
		SenderID: message.SenderID(),
		RoomID:   message.RoomID(),
		Content:  message.Content(),
		SentAt:   message.SentAt(),
		States:   states,
	}
}

func (repo *RoomMessageRepository) toDomain(message *model.RoomMessage) *domain.RoomMessage {
	states := make(map[kernel.UserID]domain.MessageState, len(message.States))
	for _, state := range message.States {
		states[state.UserID] = state.State
	}

	return domain.LoadRoomMessage(
		message.ID,
		message.SenderID,
		states,
		message.RoomID,
		message.Content,
		message.SentAt,
	)
}

func (repo *RoomMessageRepository) toDomains(messages []model.RoomMessage) []*domain.RoomMessage {
	domainMessages := make([]*domain.RoomMessage, 0, len(messages))
	for _, message := range messages {
		domainMessages = append(domainMessages, repo.toDomain(&message))
	}
	return domainMessages
}
