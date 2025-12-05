package repository

import (
	"context"
	"gochat/internal/chat/domain"
	"gochat/internal/chat/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.FriendshipRepository = (*FriendshipRepository)(nil)

type FriendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{
		db: db,
	}
}

func (repo *FriendshipRepository) Save(ctx context.Context, friendship *domain.Friendship) error {
	if err := repo.db.WithContext(ctx).Create(repo.toModel(friendship)).Error; err != nil {
		return gormutils.TranslateError(err)
	}
	return nil
}

func (repo *FriendshipRepository) ExistByUserID(ctx context.Context, userID1, userID2 kernel.UserID) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).
		Model(&model.Friendship{}).
		Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID1, userID2, userID2, userID1).
		Count(&count).Error
	if err != nil {
		return false, gormutils.TranslateError(err)
	}

	return count > 0, nil
}

func (repo *FriendshipRepository) toModel(friendship *domain.Friendship) *model.Friendship {
	return &model.Friendship{
		ID:      friendship.ID(),
		UserID1: friendship.UserID1(),
		UserID2: friendship.UserID2(),
	}
}
