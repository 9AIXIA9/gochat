package repository

import (
	"context"
	"gochat/internal/chat/domain"
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

func (repo *FriendshipRepository) Create(ctx context.Context, friendship *domain.Friendship) error {
	//TODO implement me
	panic("implement me")
}

func (repo *FriendshipRepository) ExistByUserID(ctx context.Context, userID1, userID2 kernel.UserID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
