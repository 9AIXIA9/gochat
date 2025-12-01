package repository

import (
	"context"
	"gochat/internal/friendship/domain"
	"gochat/internal/friendship/infrastructure/persistence/model"
	gormutils "gochat/internal/infrastructure/gorm"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db        *gorm.DB
	eventRepo event.Repository
}

func NewUserRepository(db *gorm.DB, eventRepo event.Repository) *UserRepository {
	return &UserRepository{
		db:        db,
		eventRepo: eventRepo,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return gormutils.TranslateError(repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // 冲突的列
		DoNothing: true,
	}).Create(&model.User{
		ID:     user.ID(),
		Number: user.Number(),
	}).Error)
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) error {
	events := user.GetEvents()
	if len(events) == 0 {
		return nil
	}

	return repo.handleUserEvents(ctx, events)
}

func (repo *UserRepository) handleUserEvents(ctx context.Context, evs []event.Event) error {
	//TODO handle different event topics
	return nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id kernel.UserID) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).Model(user).
		Preload("FriendshipsAsUser1").
		Preload("FriendshipsAsUser2").
		Preload("FriendRequestsReceived").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&user), nil
}

func (repo *UserRepository) FindByNumber(ctx context.Context, number kernel.UserNumber) (*domain.User, error) {
	var user model.User
	if err := repo.db.WithContext(ctx).Model(user).
		Preload("FriendshipsAsUser1").
		Preload("FriendshipsAsUser2").
		Preload("FriendRequestsReceived").
		First(&user, "number = ?", number).Error; err != nil {
		return nil, gormutils.TranslateError(err)
	}
	return repo.toDomain(&user), nil
}

func (repo *UserRepository) toDomain(user *model.User) *domain.User {
	friendIDs := make([]kernel.UserID, 0, len(user.FriendshipsAsUser1)+len(user.FriendshipsAsUser2))
	for _, f := range user.FriendshipsAsUser1 {
		friendIDs = append(friendIDs, f.User2ID)
	}
	for _, f := range user.FriendshipsAsUser2 {
		friendIDs = append(friendIDs, f.User1ID)
	}

	requests := make([]*domain.FriendRequest, 0, len(user.FriendRequestsReceived))
	for _, r := range user.FriendRequestsReceived {
		requests = append(requests, domain.LoadFriendRequest(
			r.ID,
			r.From,
			r.To,
			r.Content,
			r.State,
			r.SentAt,
		))
	}

	return domain.LoadUser(user.ID, user.Number, friendIDs, requests)
}
