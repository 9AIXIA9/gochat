package repository

import (
	"context"
	chatApplication "gochat/internal/chat/application"
	chatDomain "gochat/internal/chat/domain"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/kernel"
	"gochat/internal/social/domain"
	"gochat/internal/social/infrastructure/persistence/model"

	"gorm.io/gorm"
)

var _ chatApplication.RoomMembersFinder = (*RoomRepository)(nil)

type RoomRepository struct {
	innerRepository *gormutils.Repository[model.Room, domain.Room]
	db              *gorm.DB
}

func NewRoomRepository(db *gorm.DB, converter gormutils.GenericModelConverter[*model.Room, *domain.Room]) *RoomRepository {
	return &RoomRepository{innerRepository: gormutils.NewRepository(db, converter), db: db}
}

func (repo *RoomRepository) FindByRoomID(_ context.Context, _ chatDomain.RoomID) ([]kernel.UserID, error) {
	//TODO 还未实现
	return []kernel.UserID{"019a4df2-d7c2-7a5d-8f9c-809077bd8a33", "019a51da-3844-7129-b89c-381faf9b3e88", "019a52d1-f9e0-7702-adff-2346d52745a2"}, nil
}
