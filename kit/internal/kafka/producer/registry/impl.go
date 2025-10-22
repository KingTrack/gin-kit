package registry

import (
	"context"
	"sync"

	"github.com/KingTrack/gin-kit/kit/types/kafka/producer/conf"
	"github.com/KingTrack/gin-kit/kit/types/kafka/producer/unknown"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type Registry struct {
	producers map[string]sarama.SyncProducer
	mu        sync.RWMutex
}

func New() *Registry {
	return &Registry{
		producers: make(map[string]sarama.SyncProducer, 4),
	}
}

func (r *Registry) Init(ctx context.Context, configs []conf.Config) error {
	for _, v := range configs {
		config := v

		kafkaConfig := sarama.NewConfig()
		kafkaConfig.Producer.RequiredAcks = config.RequiredAcks
		kafkaConfig.Producer.Retry.Max = config.RetryMax
		kafkaConfig.Producer.Return.Successes = true

		syncProducer, err := sarama.NewSyncProducer(config.Addrs, kafkaConfig)
		if err != nil {
			return errors.WithMessagef(err, "kafka syncProducer registry create sync syncProducer failed, name:%s", config.Name)
		}

		r.addOrUpdate(config.Name, syncProducer)
	}

	return nil
}

func (r *Registry) addOrUpdate(name string, producer sarama.SyncProducer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.producers[name] = producer
}

func (r *Registry) GetProducer(ctx context.Context, name string) sarama.SyncProducer {
	r.mu.Lock()
	defer r.mu.Unlock()
	if producer, ok := r.producers[name]; ok {
		return producer
	}
	return unknown.New()
}
