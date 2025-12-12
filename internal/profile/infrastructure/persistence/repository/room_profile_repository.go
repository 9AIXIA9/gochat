package repository

import (
	"context"
	"gochat/internal/profile/domain"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomProfileRepository = (*RoomProfileRepository)(nil)

type RoomProfileRepository struct {
	db *gorm.DB
}

func NewRoomProfileRepository(db *gorm.DB) *RoomProfileRepository {
	return &RoomProfileRepository{db: db}
}

func (repo *RoomProfileRepository) Create(ctx context.Context, profile *domain.RoomProfile) error {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomProfileRepository) FindByID(ctx context.Context, id kernel.RoomID) (*domain.RoomProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *RoomProfileRepository) Update(ctx context.Context, profile *domain.RoomProfile) error {
	//TODO implement me
	panic("implement me")
}
