package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"homework-3/internal/pkg/repository"

	"github.com/redis/go-redis/v9"
)

var RedisCache *Redis

type Redis struct {
	client *redis.Client
}

func NewRedisPvzCache(opt *redis.Options) {
	RedisCache = &Redis{redis.NewClient(opt)}
}

func (rc *Redis) GetPvz(ctx context.Context, id int64) (*repository.Pvz, error) {
	bytes, err := rc.client.Get(ctx, hashId(id)).Result()
	if err == redis.Nil {
		return nil, repository.ErrObjectNotFound
	} else if err != nil {
		return nil, err
	}

	var pvz repository.Pvz
	err = json.Unmarshal([]byte(bytes), &pvz)
	if err != nil {
		return nil, err
	}

	return &pvz, nil
}

func (rc *Redis) SetPvz(ctx context.Context, id int64, pvz repository.Pvz, expiration time.Duration) error {
	data, err := json.Marshal(pvz)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, hashId(id), data, expiration).Err()
}

func hashId(id int64) string {
	return fmt.Sprintf("%d", id)
}
