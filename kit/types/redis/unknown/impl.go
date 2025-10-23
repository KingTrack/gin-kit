package unknown

import (
	"context"
	"errors"
	"net"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUnknownRedisDB = errors.New("unknown redis db")
)

func New() *redis.Client {
	opts := &redis.Options{
		Addr: "localhost:0",
		DB:   0,
	}

	client := redis.NewClient(opts)

	client.AddHook(&Hook{})

	return redis.NewClient(opts)
}

type Hook struct{}

func (h *Hook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return nil, ErrUnknownRedisDB
	}
}

func (h *Hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		return ErrUnknownRedisDB
	}
}

// BeforeProcessPipeline 拦截 pipeline 执行
func (h *Hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return ErrUnknownRedisDB
	}
}
