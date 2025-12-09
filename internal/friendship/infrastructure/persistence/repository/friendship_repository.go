package repository

import (
	"context"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.FriendshipRepository = (*FriendshipRepository)(nil)

type FriendshipRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewFriendshipRepository(db *gorm.DB, eventRepo event.Repository) *FriendshipRepository {
	return &FriendshipRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *FriendshipRepository) Create(ctx context.Context, friendship *domain.Friendship) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(friendship)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), friendship.GetEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (repo *FriendshipRepository) FindByID(ctx context.Context, friendshipID domain.FriendshipID) (*domain.Friendship, error) {
	var friendship model.Friendship
	if err := repo.db.WithContext(ctx).First(&friendship, "id = ?", friendshipID).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&friendship), nil
}

func (repo *FriendshipRepository) ExistByUserID(ctx context.Context, userID1, userID2 kernel.UserID) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&model.Friendship{}).
		Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID1, userID2, userID2, userID1).
		Count(&count).Error
	if err != nil {
		return false, gormutils.TranslateError(err)
	}
	return count > 0, nil
}

func (repo *FriendshipRepository) FindsByUserID(ctx context.Context, userID kernel.UserID, limit int, baseID domain.FriendshipID) ([]*domain.Friendship, error) {
	var friendships []model.Friendship
	query := repo.db.WithContext(ctx).Model(&model.Friendship{}).
		Where("user_id1 = ? OR user_id2 = ?", userID, userID).
		Order("id DESC").
		Limit(limit)
	if baseID != "" {
		query = query.Where("id < ?", baseID)
	}
	if err := query.Find(&friendships).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(friendships), nil
}

func (repo *FriendshipRepository) toModel(friendship *domain.Friendship) *model.Friendship {
	return &model.Friendship{
		ID:        friendship.ID(),
		UserID1:   friendship.UserID1(),
		UserID2:   friendship.UserID2(),
		CreatedAt: friendship.CreatedAt(),
	}
}

func (repo *FriendshipRepository) toDomain(friendship *model.Friendship) *domain.Friendship {
	return domain.LoadFriendship(friendship.ID, friendship.UserID1, friendship.UserID2, friendship.CreatedAt)
}

func (repo *FriendshipRepository) toDomains(friendships []model.Friendship) []*domain.Friendship {
	domains := make([]*domain.Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		domains = append(domains, repo.toDomain(&friendship))
	}
	return domains
}
