package redis

import (
	"errors"
	myErrors "gochat/internal/shared/errors"

	goRedis "github.com/redis/go-redis/v9"
)

// TranslateError 统一的 Redis 错误转换：仅将 redis.Nil 映射为领域 ErrNotFound
func TranslateError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, goRedis.Nil):
		return errors.Join(myErrors.ErrNotFound, err)
	default:
		return err
	}
}
