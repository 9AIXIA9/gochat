package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gochat/internal/domain"
	"time"
)

var _ domain.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

type RefreshTokenRepository struct {
	rdb *redis.Client
}

func NewRefreshTokenRepository(rdb *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{rdb: rdb}
}

// 生成令牌存储的键
func (r *RefreshTokenRepository) tokenKey(token domain.RefreshToken) string {
	return fmt.Sprintf("%s:refresh_token:%s", KeyPrefix, token)
}

// 生成用户到令牌映射的键
func (r *RefreshTokenRepository) userTokenKey(number domain.UserNumber) string {
	return fmt.Sprintf("%s:user_token:%d", KeyPrefix, number)
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, token domain.RefreshToken, info *domain.RefreshInfo, expireDuration time.Duration) error {
	tokenKey := r.tokenKey(token)
	userTokenKey := r.userTokenKey(info.Auth.UserNumber)

	// 序列化令牌信息
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// 使用事务确保原子性操作
	pipe := r.rdb.TxPipeline()

	// 1. 先删除用户可能存在的旧令牌
	oldToken, err := r.rdb.Get(ctx, userTokenKey).Result()
	if err == nil && oldToken != "" {
		// 用户已有令牌，删除旧令牌
		oldTokenKey := r.tokenKey(domain.RefreshToken(oldToken))
		pipe.Del(ctx, oldTokenKey)
	}

	// 2. 存储新令牌信息
	pipe.Set(ctx, tokenKey, data, expireDuration)

	// 3. 存储用户到令牌的映射（值为令牌字符串）
	pipe.Set(ctx, userTokenKey, string(token), expireDuration)

	// 执行事务
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRepository) FindRefreshToken(ctx context.Context, token domain.RefreshToken) (*domain.RefreshInfo, error) {
	key := r.tokenKey(token)
	val, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // 未找到
	}
	if err != nil {
		return nil, err
	}

	info := new(domain.RefreshInfo)
	if err := json.Unmarshal([]byte(val), info); err != nil {
		return nil, err
	}

	return info, nil
}
