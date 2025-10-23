package client

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/runtime"
	"github.com/KingTrack/gin-kit/kit/types/redis/unknown"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	name string
}

func New(name string) IRedis {
	return &Redis{name: name}
}

func (r *Redis) Client(ctx context.Context) *redis.Client {
	if runtime.Get().RedisRegistry() == nil {
		return unknown.New()
	}
	return runtime.Get().RedisRegistry().GetRedis(ctx, r.name)
}
