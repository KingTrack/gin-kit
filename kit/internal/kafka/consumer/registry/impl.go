package registry

import (
	"context"
	"sync"
	"time"

	"github.com/KingTrack/gin-kit/kit/types/kafka/consumer/conf"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type Registry struct {
	consumers map[string]sarama.ConsumerGroup
	mu        sync.RWMutex
}

func New() *Registry {
	return &Registry{
		consumers: make(map[string]sarama.ConsumerGroup, 2),
	}
}

func (r *Registry) Init(ctx context.Context, configs []conf.Config) error {
	for _, v := range configs {
		config := v

		kafkaConfig := sarama.NewConfig()
		kafkaConfig.Consumer.Offsets.Initial = config.OffsetsInitial
		kafkaConfig.Consumer.Offsets.AutoCommit.Enable = config.OffsetsAutoCommitEnable
		kafkaConfig.Consumer.Offsets.AutoCommit.Interval = time.Duration(config.OffsetsAutoCommitIntervalSec) * time.Second

		consumerGroup, err := sarama.NewConsumerGroup(config.Addrs, config.GroupID, kafkaConfig)
		if err != nil {
			return errors.WithMessagef(err, "kafka consumer registry create consumer group failed, name:%s", config.Name)
		}

		r.addOrUpdate(config.Name, consumerGroup)
	}

	return nil
}

func (r *Registry) addOrUpdate(name string, consumer sarama.ConsumerGroup) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.consumers[name] = consumer
}

func (r *Registry) GetConsumer(ctx context.Context, name string) sarama.ConsumerGroup {
	r.mu.Lock()
	defer r.mu.Unlock()

	if consumer, ok := r.consumers[name]; ok {
		return consumer
	}
	return nil
}
