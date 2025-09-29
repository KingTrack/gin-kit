package registry

import (
	"context"
	"sync"

	"github.com/KingTrack/gin-kit/kit/internal/httpclient/client"
	"github.com/KingTrack/gin-kit/kit/types/httpclient/conf"
)

type Registry struct {
	clients map[string]*client.Client
	mu      sync.RWMutex
}

func New() *Registry {
	return &Registry{
		clients: make(map[string]*client.Client),
	}
}

func (r *Registry) AddClient(ctx context.Context, config *conf.Config) error {
	return nil
}
