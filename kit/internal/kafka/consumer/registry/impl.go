package registry

import (
	"sync"

	"github.com/Shopify/sarama"
)

type Registry struct {
	producers map[string]sarama.Consumer
	mu        sync.RWMutex
}

func New() *Registry {
	return &Registry{
		producers: make(map[string]sarama.Consumer, 2),
	}
}
