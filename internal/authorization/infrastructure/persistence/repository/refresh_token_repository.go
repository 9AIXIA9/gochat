package repository

import (
	"context"
	"errors"
	"gochat/internal/authorization/infrastructure/persistence/models"
	myErrors "gochat/internal/shared/errors"
	redisutils "gochat/internal/shared/infrastructure/redis"
	"gochat/internal/shared/kernel"
	"time"

	"github.com/redis/go-redis/v9"
	"gochat/internal/authorization/domain"
)

var _ domain.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

const KeyPrefix = "gochat:authorization"
const KeyRefreshToken = KeyPrefix + ":refresh_token:"

type RefreshTokenRepository struct {
	innerRepository *redisutils.Repository[models.RefreshToken, domain.RefreshToken]
	rdb             *redis.Client
}

func NewRefreshTokenRepository(rdb *redis.Client, converter kernel.GenericModelConverter[*models.RefreshToken, *domain.RefreshToken]) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		innerRepository: redisutils.NewRepository(rdb, converter, KeyRefreshToken),
		rdb:             rdb,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, t *domain.RefreshToken) error {
	// 确保 userID 与 RefreshToken 一一对应
	existing, err := r.innerRepository.Find(ctx, string(t.UserID()))
	if err != nil && !errors.Is(err, myErrors.ErrNotFound) {
		return err
	}
	if err == nil {
		if delErr := r.innerRepository.Delete(ctx, existing.String()); delErr != nil && !errors.Is(delErr, myErrors.ErrNotFound) {
			return delErr
		}
	}

	ttl := time.Until(t.ExpiredAt())

	if err := r.innerRepository.Save(ctx, string(t.UserID()), t, ttl); err != nil {
		return err
	}
	return r.innerRepository.Save(ctx, t.String(), t, ttl)
}

func (r *RefreshTokenRepository) FindByString(ctx context.Context, str string) (*domain.RefreshToken, error) {
	return r.innerRepository.Find(ctx, str)
}
