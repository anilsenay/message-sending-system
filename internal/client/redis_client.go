package client

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) SetJsonWithTTL(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	if len(key) <= 0 {
		return errors.New("key is empty")
	}
	marshal, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, marshal, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) SetJson(ctx context.Context, key string, val interface{}) error {
	return r.SetJsonWithTTL(ctx, key, val, 0)
}
