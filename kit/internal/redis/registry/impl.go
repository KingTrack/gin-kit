package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	"github.com/KingTrack/gin-kit/kit/types/redis/conf"
	"github.com/KingTrack/gin-kit/kit/types/redis/unknown"
	"github.com/redis/go-redis/v9"
)

type Registry struct {
	mu   sync.RWMutex
	rdbs map[tlscontext.ResourceName]*redis.Client
}

func New() *Registry {
	return &Registry{
		rdbs: make(map[tlscontext.ResourceName]*redis.Client, 2),
	}
}

func (r *Registry) Init(ctx context.Context, configs []conf.Config) error {
	for _, v := range configs {
		config := v

		rdb, err := newRedis(ctx, &config)
		if err != nil {
			return err
		}
		r.addOrUpdate(ctx, config.Name, rdb)
	}
	return nil
}

func (r *Registry) addOrUpdate(ctx context.Context, name string, rdb *redis.Client) {
	resourceName := tlscontext.GetResourceName(ctx, name)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.rdbs[resourceName] = rdb
}

func newRedis(ctx context.Context, config *conf.Config) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:           config.Addr, // 单节点，可拓展集群
		Password:       config.Password,
		DB:             config.DB,
		DialTimeout:    time.Duration(config.ConnTimeoutMs) * time.Millisecond,
		ReadTimeout:    time.Duration(config.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout:   time.Duration(config.WriteTimeoutMs) * time.Millisecond,
		MinIdleConns:   config.MinIdleConns,
		MaxIdleConns:   config.MaxIdleConns,
		MaxActiveConns: config.MaxActiveConns,
		PoolSize:       config.MaxActiveConns,
	}

	rdb := redis.NewClient(opts)

	// 可选：ping 测试
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctxTimeout).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %v", err)
	}

	return rdb, nil
}

func (r *Registry) GetRedis(ctx context.Context, name string) *redis.Client {
	resourceName := tlscontext.GetResourceName(ctx, name)

	r.mu.RLock()
	defer r.mu.RUnlock()

	if rdb, ok := r.rdbs[resourceName]; ok {
		return rdb
	}

	return unknown.NewClient()
}
