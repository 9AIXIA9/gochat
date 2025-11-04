package repository

import (
	"context"
	"errors"
	"gochat/internal/authorization/application"
	"gochat/internal/authorization/infrastructure/persistence/model"
	redisutils "gochat/internal/infrastructure/redis"
	myErrors "gochat/internal/shared/errors"
	"time"

	"gochat/internal/authorization/domain"

	"github.com/redis/go-redis/v9"
)

var _ application.RefreshTokenSaver = (*RefreshTokenRepository)(nil)
var _ application.RefreshTokenFinder = (*RefreshTokenRepository)(nil)

const KeyPrefix = "gochat:authorization"
const KeyRefreshToken = KeyPrefix + ":refresh_token:"

type RefreshTokenRepository struct {
	innerRepository *redisutils.Repository[model.RefreshToken, domain.RefreshTokenEntity]
	rdb             *redis.Client
}

func NewRefreshTokenRepository(rdb *redis.Client, converter redisutils.GenericModelConverter[*model.RefreshToken, *domain.RefreshTokenEntity]) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		innerRepository: redisutils.NewRepository(rdb, converter, KeyRefreshToken),
		rdb:             rdb,
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, t *domain.RefreshTokenEntity) error {
	// 确保 userID 与 Token 一一对应
	existing, err := r.innerRepository.Find(ctx, t.UserID().String())
	if err != nil && !errors.Is(err, myErrors.ErrNotFound) {
		return err
	}
	if err == nil {
		if delErr := r.innerRepository.Delete(ctx, existing.Token().String()); delErr != nil && !errors.Is(delErr, myErrors.ErrNotFound) {
			return delErr
		}
	}

	ttl := time.Until(t.ExpiredAt())

	if err := r.innerRepository.Save(ctx, t.UserID().String(), t, ttl); err != nil {
		return err
	}
	return r.innerRepository.Save(ctx, t.Token().String(), t, ttl)
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshTokenEntity, error) {
	return r.innerRepository.Find(ctx, token.String())
}
