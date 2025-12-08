package repository

import (
	"context"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/roomship/domain"
	"gochat/internal/roomship/infrastructure/persistence/model"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ domain.RoomRepository = (*RoomRepository)(nil)

type RoomRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewRoomRepository(db *gorm.DB, eventRepo event.Repository) *RoomRepository {
	return &RoomRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo.toModel(room)).Error; err != nil {
			return gormutils.TranslateError(err)
		}

		if err := repo.eventRepo.CreateUnpublishedEvents(gormutils.SetTransaction(ctx, tx), room.GetEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (repo *RoomRepository) FindByID(ctx context.Context, roomID kernel.RoomID) (*domain.Room, error) {
	var room model.Room
	if err := repo.db.WithContext(ctx).First(&room, "id = ?", roomID).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&room), nil
}

func (repo *RoomRepository) FindsByIDs(ctx context.Context, roomIDs []kernel.RoomID) ([]*domain.Room, error) {
	var rooms []model.Room
	if err := repo.db.WithContext(ctx).Where("id IN ?", roomIDs).Find(&rooms).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomains(rooms), nil
}

func (repo *RoomRepository) toModel(room *domain.Room) *model.Room {
	return &model.Room{
		ID:                room.ID(),
		Number:            room.Number(),
		OwnerID:           room.OwnerID(),
		PasswordEncrypted: room.PasswordEncrypted(),
		MaxMemberCount:    room.MaxMemberCount(),
		CreatedAt:         room.CreatedAt(),
	}
}

func (repo *RoomRepository) toDomain(room *model.Room) *domain.Room {
	return domain.LoadRoom(room.ID, room.OwnerID, room.Number, room.PasswordEncrypted, room.MaxMemberCount, room.CreatedAt)
}

func (repo *RoomRepository) toDomains(rooms []model.Room) []*domain.Room {
	domains := make([]*domain.Room, 0, len(rooms))
	for _, room := range rooms {
		domains = append(domains, repo.toDomain(&room))
	}
	return domains
}
