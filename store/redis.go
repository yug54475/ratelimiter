package store

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//go:embed ratelimit.lua
var luaScript string

type RedisStore struct {
	client *redis.Client
	script *redis.Script
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
		script: redis.NewScript(luaScript),
	}
}

func (r *RedisStore) Allow(ctx context.Context, clientID string, limit int, windowSeconds int) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", clientID)
	result, err := r.script.Run(ctx, r.client, []string{key}, limit, windowSeconds).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
