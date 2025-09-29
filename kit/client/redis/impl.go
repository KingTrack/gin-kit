package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	name string
}

func New(name string) IRedis {
	return &Redis{name: name}
}

func (r *Redis) Client(ctx context.Context) *redis.Client {
	db := runtime.Get().RedisRegistry().GetRedis(ctx, r.name)
	if db == nil {
		return nil
	}
	return db
}
