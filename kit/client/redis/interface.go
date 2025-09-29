package client

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IRedis interface {
	Client(ctx context.Context) *redis.Client
}
