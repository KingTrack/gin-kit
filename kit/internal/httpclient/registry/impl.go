package registry

import (
	"context"
	"sync"

	"github.com/KingTrack/gin-kit/kit/internal/httpclient/client"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
)

type Registry struct {
	clients map[string]*client.HTTPClient
	mu      sync.RWMutex
}

func New() *Registry {
	return &Registry{
		clients: make(map[string]*client.HTTPClient, 2),
	}
}

func (r *Registry) Add(ctx context.Context, config *conf.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return client.New().Init(ctx, config)
}

func (r *Registry) Get(name string) *client.HTTPClient {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.clients[name]
}
