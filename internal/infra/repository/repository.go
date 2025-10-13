package repository

import (
	"github.com/redis/go-redis/v9"
	"gochat/internal/domain"

	"gorm.io/gorm"
)

var _ domain.Repository = (*Repository)(nil)

type Repository struct {
	db            *gorm.DB
	rdb           *redis.Client
	users         *UserRepository
	rooms         *RoomRepository
	messages      *MessageRepository
	refreshTokens *RefreshTokenRepository
}

func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{
		db:            db,
		rdb:           rdb,
		users:         NewUserRepository(db),
		rooms:         NewRoomRepository(db),
		messages:      NewMessageRepository(db),
		refreshTokens: NewRefreshTokenRepository(rdb),
	}
}

func (r *Repository) Users() domain.UserRepository {
	return r.users
}

func (r *Repository) Rooms() domain.RoomRepository {
	return r.rooms
}

func (r *Repository) Messages() domain.MessageRepository {
	return r.messages
}

func (r *Repository) RefreshTokens() domain.RefreshTokenRepository {
	return r.refreshTokens
}
