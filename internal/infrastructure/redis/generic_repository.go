package redis

import (
	"context"
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository[
	RedisModel any,
	DomainModel any,
] struct {
	client                     *redis.Client
	redisGenericModelConverter GenericModelConverter[*RedisModel, *DomainModel]
	KeyPrefix                  string
}

func NewRepository[
	RedisModel any,
	DomainModel any,
](client *redis.Client, converter GenericModelConverter[*RedisModel, *DomainModel], keyPrefix string) *Repository[RedisModel, DomainModel] {
	return &Repository[RedisModel, DomainModel]{
		client:                     client,
		redisGenericModelConverter: converter,
		KeyPrefix:                  keyPrefix,
	}
}

func (r *Repository[RedisModel, DomainModel]) BuildKey(key string) string {
	return r.KeyPrefix + key
}

func (r *Repository[RedisModel, DomainModel]) Save(ctx context.Context, key string, model *DomainModel, ttl ...time.Duration) error {
	if model == nil {
		return myErrors.ErrEmptyPointer
	}

	redisModel := r.redisGenericModelConverter.ToModel(model)

	data, err := json.Marshal(redisModel)
	if err != nil {
		return err
	}

	fullKey := r.BuildKey(key)
	exp := time.Duration(0)
	if len(ttl) > 0 {
		exp = ttl[0]
	}
	return TranslateError(r.client.Set(ctx, fullKey, data, exp).Err())
}

func (r *Repository[RedisModel, DomainModel]) Find(ctx context.Context, key string) (*DomainModel, error) {
	fullKey := r.BuildKey(key)
	val, err := r.client.Get(ctx, fullKey).Bytes()
	if err != nil {
		return nil, TranslateError(err)
	}

	redisModel := new(RedisModel)
	if err := json.Unmarshal(val, redisModel); err != nil {
		return nil, err
	}
	return r.redisGenericModelConverter.ToDomain(redisModel), nil
}

func (r *Repository[RedisModel, DomainModel]) Finds(ctx context.Context, keys []string) ([]*DomainModel, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	fullKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		fullKeys = append(fullKeys, r.BuildKey(k))
	}

	vals, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		return nil, TranslateError(err)
	}

	domainModels := make([]*DomainModel, 0, len(vals))
	for _, v := range vals {
		if v == nil {
			continue
		}
		byteString, ok := v.(string) // 返回的是JSON bytes
		if !ok {
			continue
		}
		redisModel := new(RedisModel)
		if err := json.Unmarshal([]byte(byteString), redisModel); err != nil {
			return nil, err
		}
		domainModel := r.redisGenericModelConverter.ToDomain(redisModel)

		domainModels = append(domainModels, domainModel)
	}
	return domainModels, nil
}

func (r *Repository[RedisModel, DomainModel]) Delete(ctx context.Context, key string) error {
	fullKey := r.BuildKey(key)
	_, err := r.client.Del(ctx, fullKey).Result()
	return TranslateError(err)
}

func (r *Repository[RedisModel, DomainModel]) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.BuildKey(key)
	n, err := r.client.Exists(ctx, fullKey).Result()
	return n > 0, TranslateError(err)
}

func (r *Repository[RedisModel, DomainModel]) ScanKeys(ctx context.Context, pattern string, count int64) ([]string, error) {
	var (
		cursor      uint64
		matchedKeys []string
	)
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, r.BuildKey(pattern), count).Result()
		if err != nil {
			return nil, TranslateError(err)
		}
		matchedKeys = append(matchedKeys, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return matchedKeys, nil
}

func (r *Repository[RedisModel, DomainModel]) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	fullKey := r.BuildKey(key)
	ok, err := r.client.Expire(ctx, fullKey, ttl).Result()
	return ok, TranslateError(err)
}

func (r *Repository[RedisModel, DomainModel]) TTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := r.BuildKey(key)
	d, err := r.client.TTL(ctx, fullKey).Result()
	return d, TranslateError(err)
}
