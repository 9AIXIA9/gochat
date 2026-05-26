package redis

import (
	"context"
	"time"

	kafkamiddleware "gochat/internal/delivery/kafka/middleware"

	"github.com/redis/go-redis/v9"
)

var _ kafkamiddleware.InboxStore = (*inboxStore)(nil)

// inboxStore adapts go-redis to the Kafka inbox middleware interface.
type inboxStore struct {
	client *redis.Client
}

func NewInboxStore(client *redis.Client) kafkamiddleware.InboxStore {
	return &inboxStore{client: client}
}

func (s *inboxStore) Claim(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return s.client.SetNX(ctx, key, "processing", ttl).Result()
}

func (s *inboxStore) Release(ctx context.Context, key string) error {
	_, err := s.client.Del(ctx, key).Result()
	return err
}

func (s *inboxStore) Complete(ctx context.Context, key string, ttl time.Duration) error {
	return s.client.Expire(ctx, key, ttl).Err()
}
